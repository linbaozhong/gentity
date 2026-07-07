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

package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/linbaozhong/gentity/pkg/types"
	"github.com/linbaozhong/gentity/pkg/util"
)

var (
	Param_InvalidErr = types.NewError(-620, "参数无效")
	UnKnownErr       = types.NewError(-610, "未知错误")
)

// Fail writes an error JSON response via the Context interface.
func Fail(c Context, e error, args ...any) error {
	j := types.NewResult()
	defer j.Free()

	if er, ok := e.(*types.Error); ok {
		j.Code = er.Code
		j.Message = er.Message
	} else if len(args) == 0 {
		j.Code = UnKnownErr.Code
		j.Message = e.Error()
	} else {
		j.Message = fmt.Sprintf("%s", args[0])
	}

	go logError(c, e)

	if err := c.JSON(j); err != nil {
		return err
	}
	return nil
}

// Ok writes a success JSON response with optional caching for GET requests.
func Ok(c Context, args ...any) error {
	j := types.NewResult()
	defer j.Free()

	if len(args) > 0 {
		j.Data = args[0]
	}

	// Cache GET responses
	if c.Method() == http.MethodGet {
		key := c.Values().Get(hasCacheKey)
		if ck, ok := key.(cacheKey); ok {
			buf, e := json.Marshal(j)
			if e == nil {
				setCache(context.Background(), ck, buf)
			}
		}
	}

	if err := c.JSON(j); err != nil {
		return Fail(c, err)
	}
	return nil
}

// SendLocalFile sends a local file as an attachment.
func SendLocalFile(c Context, path, name string) error {
	return c.SendFile(path, name)
}

// SendUrlFile downloads a file from URL and sends it as attachment.
func SendUrlFile(c Context, urlStr, name string) error {
	if urlStr == "" {
		return errors.New("url is empty")
	}
	urlStr, _ = util.IsUrl(urlStr, "https")

	resp, err := http.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch file from URL: %s, status code: %d", urlStr, resp.StatusCode)
	}

	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	c.Header("Content-Length", resp.Header.Get("Content-Length"))
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("Content-Transfer-Encoding", "binary")

	_, err = io.Copy(c.ResponseWriter(), resp.Body)
	if err != nil {
		c.StatusCode(http.StatusInternalServerError)
	}
	return err
}
