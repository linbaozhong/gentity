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
	"github.com/linbaozhong/gentity/pkg/schema"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	path    string // model路径
	driver  string // 数据库驱动
	dns     string // 数据库连接字符串
	sqlPath string // sql文件路径

	launch = &cobra.Command{
		Use:   `gentity model路径 ["SQL文件路径" | "数据库驱动" "数据库连接字符串"]`,
		Short: "ORM 代码生成工具",
		Example: `	gentity
	gentity .\db
	gentity .\db mysql "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	gentity . .\database.sql`,
		Run: func(cmd *cobra.Command, args []string) {
			fullpath, err := filepath.Abs(path)
			if err != nil {
				showError(err)
			}

			parent := filepath.Dir(fullpath)
			pos := strings.LastIndex(fullpath, string(os.PathSeparator))
			if pos > 0 {
				parent = fullpath[:pos]
			}
			// package name
			packageName := fullpath[pos+1:]

			pkgPath, err := base.PkgPath(nil, path)
			if err != nil {
				showError(err)
			}
			pos = strings.LastIndex(pkgPath, "/")
			if pos > 0 {
				pkgPath = pkgPath[:pos]
			}
			//
			if len(driver) == 0 {
				sqlPath, err = filepath.Abs(sqlPath)
				if err != nil {
					showError(err)
				}
				_, err = os.Stat(sqlPath)
				if err != nil {
					showError(err)
				}
				f, e := os.OpenFile(filepath.Join(fullpath, "gentity_model.go"), os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
				if e != nil {
					showError(e.Error())
				}
				defer f.Close()

				buf, err := schema.SqlFile2Struct(sqlPath, packageName)
				if err != nil {
					showError(err)
				}
				_, e = f.Write(buf)
				if e != nil {
					showError(e)
				}
			}

			dirs, err := os.ReadDir(path)
			if err != nil {
				showError(err)
			}
			for _, dir := range dirs {
				if dir.IsDir() {
					continue
				}
				var filename = dir.Name()
				fmt.Println(filename)
				if filepath.Ext(filename) != ".go" {
					continue
				}

				err = parseFile(parent, filepath.Join(fullpath, filename), pkgPath)
				if err != nil {
					showError(err)
				}
			}
		},
	}
)

func showError(msg interface{}) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Println("Error:", msg, file, line)
	os.Exit(1)
}

func Execute() {
	lens := len(os.Args)
	if lens > 1 {
		path = os.Args[1]
		if lens > 3 {
			driver = os.Args[2]
			dns = os.Args[3]
		} else if lens > 2 {
			sqlPath = os.Args[2]
		}
	} else {
		path = "."
	}

	err := launch.Execute()
	if err != nil {
		panic(err)
	}
}
