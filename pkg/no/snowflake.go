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

package no

import (
	"errors"
	"github.com/bwmarrin/snowflake"
	"golang.org/x/sync/singleflight"
	"strconv"
)

var (
	generateNodes = make(map[int64]*snowflake.Node)
	sg            singleflight.Group
)

func getGenerateNode(node int64) *snowflake.Node {
	if nd, ok := generateNodes[node]; ok && nd != nil {
		return nd
	}
	r, e, _ := sg.Do(strconv.FormatInt(node, 10), func() (r interface{}, e error) {
		r, e = snowflake.NewNode(node)
		if e == nil {
			generateNodes[node] = r.(*snowflake.Node)
		} else {
			generateNodes[node] = nil
		}
		return
	})
	if e != nil {
		panic(e)
	}
	return r.(*snowflake.Node)
}

// GetId 获取id,node 为节点号,取值范围为[0,1023]
func GetId(node uint) (snowflake.ID, error) {
	if node > 1023 {
		return 0, errors.New("node is too large")
	}
	return getGenerateNode(int64(node)).Generate(), nil
}
