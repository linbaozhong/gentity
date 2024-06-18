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
	"github.com/mitchellh/go-homedir"
	"log"

	"github.com/spf13/cobra"
)

var (
	home   string
	launch = &cobra.Command{
		Use:   "Gentity",
		Short: "DAL代码生成工具",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

func initConfig() {
	if home == "" {
		var err error
		home, err = homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
	}
}
