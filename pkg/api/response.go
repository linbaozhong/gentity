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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/types"
	"github.com/linbaozhong/gentity/pkg/util"
	"io"
	"net/http"
)

func Fail(c Context, e error, args ...any) error {
	j := types.NewResult()
	defer j.Free()

	if er, ok := e.(*types.Error); ok {
		j.Code = er.Code
		j.Message = er.Error()
		j.Info = er.Info
	} else if len(args) == 0 {
		j.Code = UnKnown.Code
		j.Message = e.Error()
	} else {
		j.Info = fmt.Sprintf("%s", args[0])
	}

	return c.JSON(j)
}

func Ok(c Context, args ...any) error {
	j := types.NewResult()
	defer j.Free()

	if len(args) > 0 {
		j.Data = args[0]
	}
	// 缓存
	if c.Method() == http.MethodGet {
		key := c.Values().Get(hasCacheKey)
		if _key, ok := key.(cacheKey); ok {
			buf, e := json.Marshal(j)
			if e == nil {
				setCache(context.Background(), _key, buf)
			}
		}
	}
	return c.JSON(j)
}

// SendLocalFile 发送本地文件
// path: 文件路径
// name: 文件名
func SendLocalFile(c Context, path, name string) error {
	return c.SendFile(path, name)
}

// SendUrlFile 发送url文件
// url: 文件url
// name: 文件名
// 注意: 该函数会将文件下载到本地, 然后再发送到客户端, 所以需要确保文件url是可访问的, 否则会出现错误
func SendUrlFile(c Context, url, name string) error {
	if url == "" {
		return errors.New("url is empty")
	}
	url, _ = util.IsUrl(url, "https")
	// 发送 HTTP 请求获取网络文件
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch file from URL: %s, status code: %d", url, resp.StatusCode)
	}

	// 设置响应头
	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", name))
	c.Header("Content-Length", resp.Header.Get("Content-Length"))
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("Content-Transfer-Encoding", "binary")

	// 将文件内容复制到 HTTP 响应体
	_, err = io.Copy(c.ResponseWriter(), resp.Body)
	if err != nil {
		// 复制文件内容失败，设置 500 状态码
		c.StatusCode(http.StatusInternalServerError)
	}
	return err
}
