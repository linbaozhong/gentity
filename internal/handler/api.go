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
	"github.com/linbaozhong/gentity/internal/resources"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func generateApi(name string) error {
	dir := filepath.Join(filepath.Dir(os.Args[0]), name)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, os.ModePerm)
			if err != nil {
				showError(err)
				return err
			}
		} else {
			showError(err)
			return err
		}
	} else {
		showError("dir already exists")
	}
	// 进入目录
	err = os.Chdir(dir)
	if err != nil {
		showError(err)
		return err
	}

	// 初始化模块
	err = apiInitModule(name)
	if err != nil {
		showError(err)
		return err
	}

	// 生成其他文件
	err = apiCmd(name)
	if err != nil {
		showError(err)
		return err
	}

	err = apiHandler(name)
	if err != nil {
		showError(err)
		return err
	}

	err = apiRouter(name)
	if err != nil {
		showError(err)
		return err
	}

	err = apiServiceInit(name)
	if err != nil {
		showError(err)
		return err
	}

	err = apiService(name)
	if err != nil {
		showError(err)
		return err
	}

	return exec.Command("go", "mod", "tidy").Run()
}

func apiInitModule(name string) error {
	return exec.Command("go", "mod", "init", name).Run()
}

func apiCmd(name string) error {
	err := os.Mkdir("cmd", os.ModePerm)
	if err != nil {
		return err
	}
	_, err = os.Stat("cmd/main.go")
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.OpenFile("cmd/main.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if err != nil {
				return err
			}
			defer f.Close()

			tmpl := template.New("")
			_, err = tmpl.ParseFS(resources.TemplatesFS, "templates/api_cmd_main.tmpl")
			if err != nil {
				showError(err)
			}
			return tmpl.ExecuteTemplate(f, "api_cmd_main.tmpl", struct {
				ModulePath string
			}{
				ModulePath: name,
			})
		}
	} else {
		showError("cmd/main.go already exists")
	}
	return nil
}

func apiHandler(name string) error {
	err := os.MkdirAll("internal/handler", os.ModePerm)
	if err != nil {
		return err
	}
	_, err = os.Stat("internal/handler/app_handler.go")
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.OpenFile("internal/handler/app_handler.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if err != nil {
				return err
			}
			defer f.Close()

			tmpl := template.New("")
			_, err = tmpl.ParseFS(resources.TemplatesFS, "templates/api_internal_handler.tmpl")
			return tmpl.ExecuteTemplate(f, "api_internal_handler.tmpl", struct {
				ModulePath string
			}{
				ModulePath: name,
			})
		}
	} else {
		showError("internal/handler/app_handler.go already exists")
	}
	return nil
}

func apiRouter(name string) error {
	err := os.MkdirAll("internal/router", os.ModePerm)
	if err != nil {
		return err
	}
	_, err = os.Stat("internal/router/router.go")
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.OpenFile("internal/router/router.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if err != nil {
				return err
			}
			defer f.Close()

			tmpl := template.New("")
			_, err = tmpl.ParseFS(resources.TemplatesFS, "templates/api_internal_router.tmpl")
			return tmpl.ExecuteTemplate(f, "api_internal_router.tmpl", struct {
				ModulePath string
			}{
				ModulePath: name,
			})
		}
	} else {
		showError("internal/router/router.go already exists")
	}
	return nil
}
func apiServiceInit(name string) error {
	const parent = "internal/service"
	err := os.MkdirAll(parent, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = os.Stat(parent + "/init.go")
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.OpenFile(parent+"/init.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if err != nil {
				return err
			}
			defer f.Close()

			tmpl := template.New("")
			_, err = tmpl.ParseFS(resources.TemplatesFS, "templates/api_internal_service_init.tmpl")
			return tmpl.ExecuteTemplate(f, "api_internal_service_init.tmpl", struct {
				ModulePath string
			}{
				ModulePath: name,
			})
		}
	} else {
		showError("internal/service/init.go already exists")
	}
	return nil
}

func apiService(name string) error {
	const parent = "internal/service/app"
	err := os.MkdirAll(parent, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = os.Stat(parent + "/app_service.go")
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.OpenFile(parent+"/app_service.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if err != nil {
				return err
			}
			defer f.Close()

			tmpl := template.New("")
			_, err = tmpl.ParseFS(resources.TemplatesFS, "templates/api_internal_service.tmpl")
			return tmpl.ExecuteTemplate(f, "api_internal_service.tmpl", struct {
				ModulePath string
			}{
				ModulePath: name,
			})
		}
	} else {
		showError("internal/service/app/app_service.go already exists")
	}
	return nil
}
