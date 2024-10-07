package main

import (
	"github.com/linbaozhong/gentity/example/model/db"
	"github.com/linbaozhong/gentity/pkg/log"
)

func main() {
	b := db.NewCompany()
	b.FullCorpName = "test1"
	a := *b

	log.Info(a)
}
