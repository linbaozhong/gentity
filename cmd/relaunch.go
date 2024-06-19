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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	path   string
	launch = &cobra.Command{
		Use:   "Gentity",
		Short: "DAL代码生成工具",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("path:%s\n", path)

			fullpath, err := filepath.Abs(path)
			if err != nil {
				showError(err)
			}
			log.Printf("fullpath:%s\n", fullpath)

			parent := fullpath
			pos := strings.LastIndex(fullpath, string(os.PathSeparator))
			if pos > 0 {
				parent = fullpath[:pos]
			}
			log.Printf("parent:%s\n", parent)

			// 创建目录
			parent = filepath.Join(parent, "table")
			err = os.MkdirAll(parent, os.ModePerm)
			if err != nil {
				showError(err)
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
				if filepath.Ext(filename) != ".go" {
					continue
				}

				err = parseFile(parent, filepath.Join(fullpath, filename))
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
	launch.PersistentFlags().StringVarP(&path, "path", "p", ".", "model路径")

	err := launch.Execute()
	if err != nil {
		panic(err)
	}
}
