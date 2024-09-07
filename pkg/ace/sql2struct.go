package ace

import (
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

type Driverer interface {
	GetTables(db *DB, dbName string) (map[string][]dialect.Column, error)
}

type mysql struct {
}

var Mysql mysql

func (mysql) GetTables(db *DB, dbName string) (map[string][]dialect.Column, error) {
	rows, err := db.Query(`SELECT table_name,column_name,column_default,data_type,ifnull(character_maximum_length,0),column_key,extra,column_comment FROM information_schema.COLUMNS WHERE table_schema = ?`, dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms := make(map[string][]dialect.Column)
	for rows.Next() {
		var tableName string
		col := dialect.Column{}
		err = rows.Scan(&tableName, &col.Name, &col.Default, &col.Type, &col.Size, &col.Key, &col.Extra, &col.Comment)

		if cols, ok := ms[tableName]; ok {
			ms[tableName] = append(cols, col)
		} else {
			ms[tableName] = []dialect.Column{col}
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ms, nil
}
