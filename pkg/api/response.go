// Copyright © 2023 SnowIM. All rights reserved.
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

package api

import (
	"fmt"
	"github.com/linbaozhong/gentity/pkg/types"

	"github.com/kataras/iris/v12"
)

func Fail(c Context, e error, args ...interface{}) {
	j := types.NewResult()
	defer j.Free()

	if er, ok := e.(types.Error); ok {
		j.Code = er.Code
		j.Message = er.Error()
		j.Info = er.Info
	} else if len(args) == 0 {
		j.Code = UnKnown.Code
		j.Message = e.Error()
	} else {
		j.Info = fmt.Sprintf("%s", args[0])
	}
	c.JSON(j)
}
func FailWithParam(c Context, key string, message string) {
	j := types.NewResult()
	defer j.Free()

	j.Code = Param_Invalid.Code
	j.Data = iris.Map{
		"key":     key,
		"message": message,
	}
	c.JSON(j)
}

func Ok(c Context, args ...interface{}) {
	j := types.NewResult()
	defer j.Free()

	if len(args) > 0 {
		j.Data = args[0]
	}

	c.JSON(j)
}
