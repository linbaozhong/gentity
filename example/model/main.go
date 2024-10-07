package main

import (
	"github.com/linbaozhong/gentity/example/model/db"
	"github.com/linbaozhong/gentity/pkg/log"
)

// noCopy 是一个空结构体，用于防止嵌入它的结构体被复制。
// 它实现了Locker接口的Lock和Unlock方法，但这些方法什么都不做。
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// MyStruct 是一个包含 noCopy 的结构体，因此它不应该被复制。
type MyStruct struct {
	noCopy noCopy // 嵌入 noCopy 来防止复制
	data   interface{}
}

// NewMyStruct 是 MyStruct 的构造函数，用于创建一个新的实例。
func NewMyStruct(data interface{}) *MyStruct {
	return &MyStruct{
		data: data,
	}
}

func main() {
	b := db.NewCompany()
	b.FullCorpName = "test1"

	a := *b

	log.Info(a)
}
