package handler

import (
	"fmt"
	"github.com/linbaozhong/gentity/internal/schema"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"os"
	"path/filepath"
	"regexp"
)

func sql2struct(driver, sqlPath, outputPath, packageName string) error {
	dialect.Register(driver)

	sqlPath, err := filepath.Abs(sqlPath)
	if err != nil {
		return err
	}
	_, err = os.Stat(sqlPath)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(outputPath, "gentity_model.go"), os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	buf, err := schema.SqlFile2Struct(sqlPath, packageName)
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func db2struct(driver, dns, outputPath, packageName string) error {
	db, err := ace.Connect(driver, dns)
	if err != nil {
		return err
	}
	defer db.Close()

	re := regexp.MustCompile(`/([^/]+)\?`)
	match := re.FindStringSubmatch(dns)
	if match == nil || len(match) == 0 {
		return fmt.Errorf("Could not parse database name from the connection string.")
	}
	dr, err := getDriver(driver)
	if err != nil {
		return err
	}
	// match[1] 存储的就是 dbname 的值
	tables, err := dr.GetTables(db, match[1])
	if err != nil {
		return err
	}

	buf, err := schema.DB2Struct(tables, packageName)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(outputPath, "gentity_model.go"), os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(buf)
	if err != nil {
		return err
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
