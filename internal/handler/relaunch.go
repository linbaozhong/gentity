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

var (
	path   string
	launch = &cobra.Command{
		Use:   "Gentity",
		Short: "DAL代码生成工具",
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

			pkgPath, err := base.PkgPath(nil, path)
			if err != nil {
				showError(err)
			}
			pos = strings.LastIndex(pkgPath, "/")
			if pos > 0 {
				pkgPath = pkgPath[:pos]
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
	launch.PersistentFlags().StringVarP(&path, "path", "p", ".", "model路径")

	err := launch.Execute()
	if err != nil {
		panic(err)
	}
}
