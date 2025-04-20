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
	version = "0.9"
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
				e       error
				pkgPath string
			)
			// struct全路径
			fullpath, e = filepath.Abs(path)
			if e != nil {
				showError(e)
				return
			}
			// 上一级目录
			parent = filepath.Dir(fullpath)
			pos := strings.LastIndex(fullpath, string(os.PathSeparator))
			if pos > 0 {
				parent = fullpath[:pos]
			}
			// 包名
			_packageName := fullpath[pos+1:]
			//
			switch command {
			case "api": // 如果command="api"，则初始化api模板
				if len(args) > 1 {
					generateApi(path)
				} else {
					showError("The project name is not entered")
				}
				return
			case "sql", "db": // 如果command="db"或"sql"，则根据数据表信息生成结构体
				if len(driver) == 0 || len(dns) == 0 {
					showError("The database driver and connection string must be entered")
					return
				}
				// 生成结构体
				if command == "db" { // 从dns获取表结构生成结构体
					e = db2struct(driver, dns, fullpath, _packageName)
				} else if command == "sql" { // 根据sql建表文件生成结构体
					e = sql2struct(driver, dns, fullpath, _packageName)
				}
				if e != nil {
					showError(e)
				}
				fallthrough
			case "dao": // 根据struct生成dao层代码和序列化器
				_definePath := filepath.Join(parent, "define")
				tablePath = filepath.Join(_definePath, "table")
				// 创建生成dao层代码的目录
				e = os.MkdirAll(tablePath, os.ModePerm)
				if e != nil {
					showError(e)
				}
				daoPath = filepath.Join(_definePath, "dao")
				e = os.MkdirAll(daoPath, os.ModePerm)
				if e != nil {
					showError(e)
				}
				//// 写入daoBase
				//e = writeDaoBase(daoPath)
				//if e != nil {
				//	showError(e)
				//}
			case "check", "dto":
			case "swag":
			case "version":
				fmt.Println("ver." + version)
				return
			default:
				showError("The command is not entered")
				return
			}

			// 包目录
			pkgPath, e = base.PkgPath(nil, path)
			if e != nil {
				showError(e)
				return
			}
			// 遍历结构体目录中的文件，生成dao层代码
			_dirs, e := os.ReadDir(path)
			if e != nil {
				showError(e)
			}

			if command == "check" || command == "dto" {
				os.Remove(filepath.Join(fullpath, dentityDTO))
				defer dtoFile.Close()
			}
			// // 初始化进度条
			// _bar := progressbar.NewOptions(len(_dirs),
			//	progressbar.OptionEnableColorCodes(true),
			//	progressbar.OptionSetTheme(progressbar.Theme{
			//		Saucer:        "[red]=[reset]",
			//		SaucerHead:    "[red]>[reset]",
			//		SaucerPadding: " ",
			//		BarStart:      "[",
			//		BarEnd:        "]",
			//	}),
			//	progressbar.OptionShowCount())
			_allTds := make([]TempData, 0)
			for _, dir := range _dirs {
				// _bar.Add(1)
				if dir.IsDir() {
					continue
				}
				var filename = dir.Name() // struct文件名
				if filepath.Ext(filename) != ".go" || strings.HasSuffix(filename, ".gen.go") {
					continue
				}
				// 解析文件
				if command == "swag" {
					parseSwagger(filename, pkgPath)
					continue
				} else if command == "check" || command == "dto" {
					if filename == dentityDTO {
						continue
					}
					_tds, e := parseFile(filename, pkgPath, "checker", "request", "response")
					if e != nil {
						showError(e)
					}
					e = generateDTO(_tds, filename)
				} else {
					_tds, e := parseFile(filename, pkgPath, "tablename")
					if e != nil {
						showError(e)
					}
					e = generateDao(_tds, filename)
					// 将_tds中的元素追加到_allTds中
					_allTds = append(_allTds, _tds...)
				}
				if e != nil {
					showError(e)
				}
			}
			// 写入daoBase
			e = writeDaoBase(daoPath, _allTds)
			if e != nil {
				showError(e)
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
	_lens := len(os.Args)
	path = "."
	driver = "mysql"
	if _lens > 1 {
		command = os.Args[1]
		if _lens > 2 {
			path = os.Args[2]
			if _lens > 3 {
				driver = os.Args[3]
				if _lens > 4 {
					dns = os.Args[4]
				}
			}
		}
	}

	var e error

	e = launch.Execute()
	if e != nil {
		showError(e)
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
