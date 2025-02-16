package handler

import (
	"fmt"
	"github.com/linbaozhong/gentity/internal/schema"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
)

func sql2struct(driver, sqlPath, outputPath, packageName string) error {
	dialect.Register(driver)

	_sqlPath, e := filepath.Abs(sqlPath)
	if e != nil {
		return e
	}
	_, e = os.Stat(_sqlPath)
	if e != nil {
		return e
	}
	_f, e := os.OpenFile(filepath.Join(outputPath, "gentity_model.go"), os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		return e
	}
	defer _f.Close()

	_buf, e := schema.SqlFile2Struct(_sqlPath, packageName)
	if e != nil {
		return e
	}

	_formatted, e := format.Source(_buf)
	if e != nil {
		return e
	}

	_, e = _f.Write(_formatted)
	if e != nil {
		return e
	}
	return nil
}

func db2struct(driver, dns, outputPath, packageName string) error {
	_db, e := ace.Connect(driver, dns)
	if e != nil {
		return e
	}
	defer _db.Close()

	_re := regexp.MustCompile(`/([^/]+)\?`)
	match := _re.FindStringSubmatch(dns)
	if match == nil || len(match) == 0 {
		return fmt.Errorf("Could not parse database name from the connection string.")
	}
	_dr, e := getDriver(driver)
	if e != nil {
		return e
	}
	// match[1] 存储的就是 dbname 的值
	tables, e := _dr.GetTables(_db, match[1])
	if e != nil {
		return e
	}

	_buf, e := schema.DB2Struct(tables, packageName)
	if e != nil {
		return e
	}

	formatted, e := format.Source(_buf)
	if e != nil {
		return e
	}

	_f, e := os.OpenFile(filepath.Join(outputPath, "gentity_model.go"), os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		return e
	}
	defer _f.Close()

	_, e = _f.Write(formatted)
	if e != nil {
		return e
	}
	return nil
}

func getDriver(driver string) (ace.Driverer, error) {
	switch driver {
	case "mysql":
		return ace.Mysql, nil
	}
	return nil, fmt.Errorf("Unsupported driver %s", driver)
}
