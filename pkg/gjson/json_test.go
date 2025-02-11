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

package gjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-json-experiment/json/jsontext"
	"testing"
)

var (
	jsonString = `{
"name":"linbaozhong",
"age":18,
"address":{
"city":"beijing",
"street":"chaoyang"
},
"hobbies":[
"football",
"bask \"et\" ball"
]
}`
)

type User struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Address Address  `json:"address"`
	Hobbies []string `json:"hobbies"`
}
type Address struct {
	City   string `json:"city"`
	Street string `json:"street"`
}
type Custom struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Address Address  `json:"address"`
	Hobbies []string `json:"hobbies"`
}

func TestEncoder(t *testing.T) {
	user := User{
		Name: "linbaozhong",
		Age:  18,
		Address: Address{
			City:   "beijing",
			Street: "chaoyang",
		},
		Hobbies: []string{"football", "bask \"et\" ball"},
	}
	buf, e := json.Marshal(&user)
	if e != nil {
		t.Log(e)
	}
	t.Log(string(buf))
}
func TestDecoder(t *testing.T) {
	// start := time.Now()
	// for i := 0; i < 500; i++ {
	// 	var u User
	// 	e := json.Unmarshal([]byte(jsonString), &u)
	// 	if e != nil {
	// 		t.Log(e)
	// 	}
	// }
	// t.Log(time.Since(start).Nanoseconds())
	//
	// start = time.Now()
	// for i := 0; i < 500; i++ {
	var c Custom
	e := json.Unmarshal([]byte(jsonString), &c)
	if e != nil {
		t.Log(e)
	}

	// }
	// t.Log(time.Since(start).Nanoseconds())
}
func (u *User) MarshalJSON() ([]byte, error) {
	fmt.Println("---------------")
	var buf = bytes.NewBuffer(nil)
	enc := jsontext.NewEncoder(buf, jsontext.Multiline(true))

	enc.WriteToken(jsontext.ObjectStart)
	enc.WriteToken(jsontext.String("name"))
	enc.WriteToken(jsontext.String(u.Name))
	enc.WriteToken(jsontext.String("age"))
	enc.WriteToken(jsontext.Int(int64(u.Age)))
	enc.WriteToken(jsontext.String("address"))
	enc.WriteToken(jsontext.ObjectStart)
	enc.WriteToken(jsontext.String("city"))
	enc.WriteToken(jsontext.String(u.Address.City))
	enc.WriteToken(jsontext.String("street"))
	enc.WriteToken(jsontext.String(u.Address.Street))
	enc.WriteToken(jsontext.ObjectEnd)
	enc.WriteToken(jsontext.String("hobbies"))
	enc.WriteToken(jsontext.ArrayStart)
	for _, hobby := range u.Hobbies {
		enc.WriteToken(jsontext.String(hobby))
	}
	enc.WriteToken(jsontext.ArrayEnd)
	enc.WriteToken(jsontext.ObjectEnd)

	return buf.Bytes(), nil
}

func (u *User) UnmarshalJSON(data []byte) error {
	reader := bytes.NewReader(data)
	dec := jsontext.NewDecoder(reader)

	var (
		readUser    func(*User)
		readAddress func(*Address)
		readHobbies func(*[]string)
	)
	readUser = func(user *User) {
		for {
			token, err := dec.ReadToken()
			if err != nil {
				break
			}

			switch token.String() {
			case "{":
				break
			case "}":
				return
			case "name":
				token, err := dec.ReadToken()
				if err != nil {
					break
				}
				user.Name = token.String()
			case "age":
				token, err := dec.ReadToken()
				if err != nil {
					break
				}
				user.Age = int(token.Int())
			case "address":
				readAddress(&user.Address)
			case "hobbies":
				readHobbies(&user.Hobbies)
			}
		}

	}
	readAddress = func(address *Address) {
		for {
			token, err := dec.ReadToken()
			if err != nil {
				break
			}
			switch token.String() {
			case "{":
				break
			case "}":
				return
			case "city":
				token, err := dec.ReadToken()
				if err != nil {
					break
				}
				address.City = token.String()
			case "street":
				token, err := dec.ReadToken()
				if err != nil {
					break
				}
				address.Street = token.String()
			}
		}
	}
	readHobbies = func(hobbies *[]string) {
		for {
			token, err := dec.ReadToken()
			if err != nil {
				fmt.Println(err)
				break
			}
			switch token.String() {
			case "[":
				break
			case "]":
				return
			default:
				*hobbies = append(*hobbies, token.String())
			}
		}
	}
	readUser(u)

	return nil
}

func (u *Custom) UnmarshalJSON(data []byte) error {
	ok := ValidBytes(data)
	if !ok {
		return errors.New("invalid json")
	}
	result := ParseBytes(data)
	result.ForEach(func(key, value Result) bool {
		fmt.Println("--------", value.Raw)
		switch key.Str {
		case "name":
			u.Name = value.Str
		case "age":
			u.Age = int(value.Int())
		case "address":
			value.ForEach(func(key, value Result) bool {
				switch key.Str {
				case "city":
					u.Address.City = value.Str
				case "street":
					u.Address.Street = value.Str
				}
				return true
			})
		case "hobbies":
			for _, val := range value.Array() {
				u.Hobbies = append(u.Hobbies, val.Str)
			}
		}
		return true
	})
	return nil
}
