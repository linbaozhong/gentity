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

// Mode 项目模式
type Mode string

const (
	FrontEnd Mode = "frontend" //前端项目
	BackEnd  Mode = "backend"  //后端项目
)

func (m Mode) String() string {
	return string(m)
}

func parseFolder(name string) (string, error) {
	_dir, e := os.Getwd()
	if e != nil {
		return "", e
	}
	return filepath.Join(_dir, name), nil
}

func generate(name string, modes ...Mode) error {
	var mode = BackEnd
	if len(modes) > 0 {
		mode = modes[0]
	}

	_dir, e := parseFolder(filepath.Join(name, mode.String()))
	if e != nil {
		return e
	}
	_, e = os.Stat(_dir)
	if e != nil {
		if os.IsNotExist(e) {
			e = os.MkdirAll(_dir, os.ModePerm)
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

	switch mode {
	case FrontEnd:
		e = frontend(name)
	case BackEnd:
		e = backend(name)
	}
	return e
}

// initApi 将指定文件夹初始化为后端 api 项目框架
func initApi(dir string) error {
	name, e := parseFolder(dir)
	if e != nil {
		return e
	}
	return backend(filepath.Base(name))
}

// backend 生成后端 api 项目框架
func backend(name string) error {
	// 初始化模块
	e := apiInitModule(name)
	if e != nil {
		showError(e)
		return e
	}

	// 生成其他文件
	if e = apiCmd(name); e != nil {
		showError(e)
		return e
	}

	if e = apiHandler(name); e != nil {
		showError(e)
		return e
	}

	if e = apiRouter(name); e != nil {
		showError(e)
		return e
	}

	if e = apiModel(name); e != nil {
		showError(e)
		return e
	}

	if e = apiConstant(name); e != nil {
		showError(e)
		return e
	}

	if e = apiLib(name); e != nil {
		showError(e)
		return e
	}

	if e = apiServiceInit(name); e != nil {
		showError(e)
		return e
	}

	if e = apiService(name); e != nil {
		showError(e)
		return e
	}

	// 生成dto和dao
	if e = os.Chdir("internal/model"); e != nil {
		showError(e)
		return e
	}

	log.Printf("The DTO code is being generated ... \n")
	if e = exec.Command("gentity", "dao", "./do").Run(); e != nil {
		showError(e)
		return e
	}

	log.Printf("The DAO code is being generated ... \n")
	if e = exec.Command("gentity", "dto", "./dto").Run(); e != nil {
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
	const parent = "internal/service"
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
		showError("internal/service/user_service.go already exists")
	}
	return nil
}

func apiModel(name string) error {
	log.Printf("Creating new api dto file. \n")
	var parent = "internal/model/dto"
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
			if e != nil {
				return e
			}
			e = _tmpl.ExecuteTemplate(_f, "api_internal_model_dto.tmpl", struct {
				ModulePath string
			}{
				ModulePath: name,
			})
			if e != nil {
				return e
			}
		}
	} else {
		showError(parent + "/user.go already exists")
	}
	// 创建 do文件
	log.Printf("Creating new api do file. \n")
	parent = "internal/model/do"
	e = os.MkdirAll(parent, os.ModePerm)
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
			_, e = _tmpl.ParseFS(resources.TemplatesFS, "templates/api_internal_model_do.tmpl")
			if e != nil {
				return e
			}
			e = _tmpl.ExecuteTemplate(_f, "api_internal_model_do.tmpl", struct {
				ModulePath string
			}{
				ModulePath: name,
			})
			if e != nil {
				return e
			}
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
	_, e = os.Stat(parent + "/consts.go")
	if e != nil {
		if os.IsNotExist(e) {
			_f, e := os.OpenFile(parent+"/consts.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
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
		showError(parent + "/consts.go already exists")
	}
	return nil
}

func apiLib(name string) error {
	log.Printf("Creating new api lib file. \n")
	const parent = "internal/lib"
	e := os.MkdirAll(parent, os.ModePerm)
	if e != nil {
		return e
	}
	_, e = os.Stat(parent + "/middleware.go")
	if e != nil {
		if os.IsNotExist(e) {
			_f, e := os.OpenFile(parent+"/middleware.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
			if e != nil {
				return e
			}
			defer _f.Close()

			_tmpl := template.New("")
			_, e = _tmpl.ParseFS(resources.TemplatesFS, "templates/api_internal_lib_middleware.tmpl")
			return _tmpl.ExecuteTemplate(_f, "api_internal_lib_middleware.tmpl", struct {
				ModulePath string
			}{
				ModulePath: name,
			})
		}
	} else {
		showError(parent + "/middleware.go already exists")
	}
	return nil
}
