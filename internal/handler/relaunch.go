// Copyright © 2023 Linbaozhong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"fmt"
	"github.com/linbaozhong/gentity/internal/base"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	version = "0.4.8"
)

var (
	command  string // 命令: api,dao,db,sql,check
	path     string // struct文件路径
	fullpath string // struct文件全路径
	parent   string // struct文件父级目录

	tablePath string // table文件全路径
	daoPath   string // dao文件全路径

	driver string // 数据库驱动
	dns    string // 数据库连接字符串

	launch = &cobra.Command{
		Use:   `gentity command [Struct路径] ["SQL文件路径" | "数据库驱动" "数据库连接字符串"]`,
		Short: "ORM 代码生成工具. ver." + version,
		Example: `	gentity api project_name
	gentity dao
	gentity dao .\do
	gentity db .\do mysql "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	gentity sql .\do mysql .\database.sql
	gentity check .\dto
	gentity version`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err     error
				pkgPath string
			)
			// struct全路径
			fullpath, err = filepath.Abs(path)
			if err != nil {
				showError(err)
				return
			}
			//
			switch command {
			case "api": // 如果command="api"，则初始化api模板
				if len(args) > 1 {
					generateApi(path)
				} else {
					showError("The project name is not entered")
				}
				return
			case "dao", "sql", "db":
				// 上一级目录
				parent = filepath.Dir(fullpath)
				pos := strings.LastIndex(fullpath, string(os.PathSeparator))
				if pos > 0 {
					parent = fullpath[:pos]
				}
				//
				definePath := filepath.Join(parent, "define")
				// 包名
				packageName := fullpath[pos+1:]

				if len(driver) > 0 && len(dns) > 0 {
					// 生成结构体
					// has, _ := regexp.MatchString("@.*:|:.*@", dns)
					if command == "db" {
						err = db2struct(driver, dns, fullpath, packageName)
					} else if command == "sql" {
						err = sql2struct(driver, dns, fullpath, packageName)
					}
					if err != nil {
						showError(err)
					}
				}
				tablePath = filepath.Join(definePath, "table")
				// 创建生成dao层代码的目录
				err = os.MkdirAll(tablePath, os.ModePerm)
				if err != nil {
					showError(err)
				}
				daoPath = filepath.Join(definePath, "dao")
				err = os.MkdirAll(daoPath, os.ModePerm)
				if err != nil {
					showError(err)
				}
				// 写入daoBase
				err = writeDaoBase(daoPath)
				if err != nil {
					showError(err)
				}
			case "check":
			case "version":
				fmt.Println("ver." + version)
				return
			default:
				showError("The command is not entered")
				return
			}

			// 包目录
			pkgPath, err = base.PkgPath(nil, path)
			if err != nil {
				showError(err)
				return
			}
			// 遍历结构体目录中的文件，生成dao层代码
			dirs, err := os.ReadDir(path)
			if err != nil {
				showError(err)
			}

			if command == "check" {
				defer dtoFile.Close()
			}

			for _, dir := range dirs {
				if dir.IsDir() {
					continue
				}
				var filename = dir.Name() // struct文件名
				if filepath.Ext(filename) != ".go" {
					continue
				}

				if command == "check" {
					err = generateCheck(filename)
				} else {
					err = parseFile(filename, pkgPath)
				}
				if err != nil {
					showError(err)
				}
			}
		},
	}
)

func showError(msg any) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Println("Error:", msg, file, line)
	os.Exit(1)
}

func Execute() {
	lens := len(os.Args)
	path = "."
	driver = "mysql"
	if lens > 1 {
		command = os.Args[1]
		if lens > 2 {
			path = os.Args[2]
			if lens > 3 {
				driver = os.Args[3]
				if lens > 4 {
					dns = os.Args[4]
				}
			}
		}
	}

	var err error

	err = launch.Execute()
	if err != nil {
		showError(err)
	}
}

/*
gentity api project_name
gentity dao
gentity dao .\do
gentity db .\do mysql "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
gentity sql .\do mysql .\database.sql
gentity check .\dto

*/
