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

package base

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"os"
	"path"
	"path/filepath"
)

// DefaultConfig for loading Go base.
var DefaultConfig = &packages.Config{Mode: packages.NeedName}

// PkgPath returns the Go package name for given target path.
// Even if the existing path is not exist yet in the filesystem.
//
// If base.Config is nil, DefaultConfig will be used to load base.
func PkgPath(config *packages.Config, target string) (string, error) {
	if config == nil {
		config = DefaultConfig
	}
	pathCheck, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	var parts []string
	if _, err := os.Stat(pathCheck); os.IsNotExist(err) {
		parts = append(parts, filepath.Base(pathCheck))
		pathCheck = filepath.Dir(pathCheck)
	}
	// Try maximum 2 directories above the given
	// target to find the root package or module.
	var (
		pkgs []*packages.Package
	)

	for i := 0; i < 2; i++ {
		pkgs, err = packages.Load(config, pathCheck)
		if err != nil {
			return "", fmt.Errorf("load package info: %w", err)
		}
		if len(pkgs) == 0 || len(pkgs[0].Errors) != 0 {
			parts = append(parts, filepath.Base(pathCheck))
			pathCheck = filepath.Dir(pathCheck)
			continue
		}
		pkgPath := pkgs[0].PkgPath
		for j := len(parts) - 1; j >= 0; j-- {
			pkgPath = path.Join(pkgPath, parts[j])
		}
		return filepath.ToSlash(filepath.Dir(pkgPath)), nil
	}
	if len(pkgs) > 0 {
		return pkgs[0].PkgPath, nil
	}
	return "", fmt.Errorf("root package or module was not found for: %s", target)
}
