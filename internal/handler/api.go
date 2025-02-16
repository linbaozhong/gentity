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
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func generateApi(name string) error {
	_dir, e := os.Getwd()
	if e != nil {
		return e
	}
	_dir = filepath.Join(_dir, name)
	_, e = os.Stat(_dir)
	if e != nil {
		if os.IsNotExist(e) {
			e = os.Mkdir(_dir, os.ModePerm)
			if e != nil {
				showError(e)
				return e
			}
			log.Printf("The application directory %s is created successfully. \n", name)
		} else {
			showError(e)
			return e
		}
	} else {
		showError("_dir already exists")
	}
	// 进入目录
	e = os.Chdir(_dir)
	if e != nil {
		showError(e)
		return e
	}

	// 初始化模块
	e = apiInitModule(name)
	if e != nil {
		showError(e)
		return e
	}

	// 生成其他文件
	e = apiCmd(name)
	if e != nil {
		showError(e)
		return e
	}

	e = apiHandler(name)
	if e != nil {
		showError(e)
		return e
	}

	e = apiRouter(name)
	if e != nil {
		showError(e)
		return e
	}

	e = apiModel(name)
	if e != nil {
		showError(e)
		return e
	}

	e = apiConstant(name)
	if e != nil {
		showError(e)
		return e
	}

	e = apiServiceInit(name)
	if e != nil {
		showError(e)
		return e
	}

	e = apiService(name)
	if e != nil {
		showError(e)
		return e
	}

	log.Printf("go mod tidying ... \n")
	return exec.Command("go", "mod", "tidy").Run()
}

func apiInitModule(name string) error {
	log.Printf("Initialize the new go.mod file. \n")
	return exec.Command("go", "mod", "init", name).Run()
}

func apiCmd(name string) error {
	log.Printf("Creating the main.go application entry file. \n")
	e := os.Mkdir("cmd", os.ModePerm)
	if e != nil {
		return e
	}
	_, e = os.Stat("cmd/main.go")
	if e != nil {
		if os.IsNotExist(e) {
			_f, e := os.OpenFile("cmd/main.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if e != nil {
				return e
			}
			defer _f.Close()

			_tmpl := template.New("")
			_, e = _tmpl.ParseFS(resources.TemplatesFS, "templates/api_cmd_main.tmpl")
			if e != nil {
				showError(e)
			}
			return _tmpl.ExecuteTemplate(_f, "api_cmd_main.tmpl", struct {
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
	log.Printf("Creating new api handler file. \n")
	e := os.MkdirAll("internal/handler", os.ModePerm)
	if e != nil {
		return e
	}
	_, e = os.Stat("internal/handler/user_handler.go")
	if e != nil {
		if os.IsNotExist(e) {
			_f, e := os.OpenFile("internal/handler/user_handler.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if e != nil {
				return e
			}
			defer _f.Close()

			_tmpl := template.New("")
			_, e = _tmpl.ParseFS(resources.TemplatesFS, "templates/api_internal_handler.tmpl")
			return _tmpl.ExecuteTemplate(_f, "api_internal_handler.tmpl", struct {
				ModulePath string
			}{
				ModulePath: name,
			})
		}
	} else {
		showError("internal/handler/user_handler.go already exists")
	}
	return nil
}

func apiRouter(name string) error {
	log.Printf("Creating new api route file. \n")
	e := os.MkdirAll("internal/router", os.ModePerm)
	if e != nil {
		return e
	}
	_, e = os.Stat("internal/router/router.go")
	if e != nil {
		if os.IsNotExist(e) {
			_f, e := os.OpenFile("internal/router/router.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if e != nil {
				return e
			}
			defer _f.Close()

			_tmpl := template.New("")
			_, e = _tmpl.ParseFS(resources.TemplatesFS, "templates/api_internal_router.tmpl")
			return _tmpl.ExecuteTemplate(_f, "api_internal_router.tmpl", struct {
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
	log.Printf("Creating new api service file. \n")
	const parent = "internal/service"
	e := os.MkdirAll(parent, os.ModePerm)
	if e != nil {
		return e
	}
	_, e = os.Stat(parent + "/init.go")
	if e != nil {
		if os.IsNotExist(e) {
			_f, e := os.OpenFile(parent+"/init.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if e != nil {
				return e
			}
			defer _f.Close()

			_tmpl := template.New("")
			_, e = _tmpl.ParseFS(resources.TemplatesFS, "templates/api_internal_service_init.tmpl")
			return _tmpl.ExecuteTemplate(_f, "api_internal_service_init.tmpl", struct {
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
	const parent = "internal/service/user"
	e := os.MkdirAll(parent, os.ModePerm)
	if e != nil {
		return e
	}
	_, e = os.Stat(parent + "/user_service.go")
	if e != nil {
		if os.IsNotExist(e) {
			_f, e := os.OpenFile(parent+"/user_service.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if e != nil {
				return e
			}
			defer _f.Close()

			_tmpl := template.New("")
			_, e = _tmpl.ParseFS(resources.TemplatesFS, "templates/api_internal_service.tmpl")
			return _tmpl.ExecuteTemplate(_f, "api_internal_service.tmpl", struct {
				ModulePath string
			}{
				ModulePath: name,
			})
		}
	} else {
		showError("internal/service/user/user_service.go already exists")
	}
	return nil
}

func apiModel(name string) error {
	log.Printf("Creating new api model file. \n")
	const parent = "internal/model/dto"
	e := os.MkdirAll(parent, os.ModePerm)
	if e != nil {
		return e
	}
	_, e = os.Stat(parent + "/user.go")
	if e != nil {
		if os.IsNotExist(e) {
			_f, e := os.OpenFile(parent+"/user.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if e != nil {
				return e
			}
			defer _f.Close()

			_tmpl := template.New("")
			_, e = _tmpl.ParseFS(resources.TemplatesFS, "templates/api_internal_model_dto.tmpl")
			return _tmpl.ExecuteTemplate(_f, "api_internal_model_dto.tmpl", struct {
				ModulePath string
			}{
				ModulePath: name,
			})
		}
	} else {
		showError(parent + "/user.go already exists")
	}
	return nil
}

func apiConstant(name string) error {
	log.Printf("Creating new api constant file. \n")
	const parent = "internal/constant"
	e := os.MkdirAll(parent, os.ModePerm)
	if e != nil {
		return e
	}
	_, e = os.Stat(parent + "/error.go")
	if e != nil {
		if os.IsNotExist(e) {
			_f, e := os.OpenFile(parent+"/error.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if e != nil {
				return e
			}
			defer _f.Close()

			_tmpl := template.New("")
			_, e = _tmpl.ParseFS(resources.TemplatesFS, "templates/api_internal_constant.tmpl")
			return _tmpl.ExecuteTemplate(_f, "api_internal_constant.tmpl", struct {
				ModulePath string
			}{
				ModulePath: name,
			})
		}
	} else {
		showError(parent + "/error.go already exists")
	}
	return nil
}
