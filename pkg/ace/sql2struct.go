package ace

import (
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/sqlparser"
	"strings"
)

type Driverer interface {
	GetTables(db *DB, dbName string) ([]*sqlparser.Table, error)
}

type mysql struct {
}

var Mysql mysql

func getTables(db *DB, dbName string) ([]*sqlparser.Table, error) {
	// 表名,表注释
	rows, err := db.Query(`SELECT table_name,table_comment FROM information_schema.tables WHERE table_schema = ? and table_type = 'BASE TABLE'`, dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ts := make([]*sqlparser.Table, 0)
	for rows.Next() {
		var tableName string
		var comment string
		err = rows.Scan(&tableName, &comment)
		if err != nil {
			return nil, err
		}
		ts = append(ts, &sqlparser.Table{Name: tableName, Comment: comment})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ts, nil
}

func getColumns(db *DB, dbName string) (map[string][]*sqlparser.Column, error) {
	// 表字段信息
	rows, err := db.Query(`SELECT table_name,column_name,column_default,data_type,column_type,ifnull(character_maximum_length,0),ifnull(numeric_precision,0),ifnull(numeric_scale,0),column_key,extra,column_comment FROM information_schema.COLUMNS WHERE table_schema = ?`, dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms := make(map[string][]*sqlparser.Column)
	for rows.Next() {
		var tableName string
		col := &sqlparser.Column{}
		err = rows.Scan(&tableName, &col.Name, &col.Default, &col.Type, &col.ColumnType, &col.Size, &col.Precision, &col.Scale, &col.Key, &col.Extra, &col.Comment)
		if strings.ToUpper(col.Extra) == dialect.AutoInc {
			col.AutoIncr = true
		}
		if strings.Contains(col.ColumnType, "unsigned") {
			col.Unsigned = true
		}
		if cols, ok := ms[tableName]; ok {
			ms[tableName] = append(cols, col)
		} else {
			ms[tableName] = []*sqlparser.Column{col}
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ms, nil
}

func (mysql) GetTables(db *DB, dbName string) ([]*sqlparser.Table, error) {
	// 表名,表注释
	ts, err := getTables(db, dbName)
	if err != nil {
		return nil, err
	}

	// 表字段信息
	ms, err := getColumns(db, dbName)
	if err != nil {
		return nil, err
	}

	for _, t := range ts {
		t.ColumnsX = ms[t.Name]
	}
	return ts, nil
}

//
// func (mysql) GetTables(db *DB, dbName string) (map[string][]*sqlparser.Column, error) {
// 	rows,err:=db.Query(`SELECT table_name,table_comment FROM information_schema.tables WHERE table_schema = 'dispatch' and table_type = 'BASE TABLE'`)
// 	if err!= nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
//
// 	rows, err = db.Query(`SELECT table_name,column_name,column_default,data_type,column_type,ifnull(character_maximum_length,0),ifnull(numeric_precision,0),ifnull(numeric_scale,0),column_key,extra,column_comment FROM information_schema.COLUMNS WHERE table_schema = ?`, dbName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	ms := make(map[string][]*sqlparser.Column)
// 	for rows.Next() {
// 		var tableName string
// 		col := &sqlparser.Column{}
// 		err = rows.Scan(&tableName, &col.Name, &col.Default, &col.Type, &col.ColumnType, &col.Size, &col.Precision, &col.Scale, &col.Key, &col.Extra, &col.Comment)
// 		if strings.ToUpper(col.Extra) == dialect.AutoInc {
// 			col.AutoIncr = true
// 		}
// 		if strings.Contains(col.ColumnType, "unsigned") {
// 			col.Unsigned = true
// 		}
// 		if cols, ok := ms[tableName]; ok {
// 			ms[tableName] = append(cols, col)
// 		} else {
// 			ms[tableName] = []*sqlparser.Column{col}
// 		}
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return ms, nil
// }
