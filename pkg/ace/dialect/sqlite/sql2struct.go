package sqlite

import (
	"database/sql"
	"github.com/linbaozhong/gentity/pkg/sqlparser"
	"strconv"
	"strings"
)

// type Driverer interface {
// 	GetTables(db *ace.DB, dbName string) ([]*sqlparser.Table, error)
// }

func GetTables(db *sql.DB, dbName string) ([]*sqlparser.Table, error) {
	// SQLite 获取所有表名（排除 sqlite_ 系统表）
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ts := make([]*sqlparser.Table, 0)
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		// SQLite 不原生支持表注释
		ts = append(ts, &sqlparser.Table{Name: tableName, Comment: ""})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ts, nil
}

func GetColumns(db *sql.DB, dbName string) (map[string][]*sqlparser.Column, error) {
	// 先获取所有表名
	tables, err := GetTables(db, dbName)
	if err != nil {
		return nil, err
	}

	ms := make(map[string][]*sqlparser.Column)

	for _, t := range tables {
		cols, err := getTableColumns(db, t.Name)
		if err != nil {
			return nil, err
		}
		ms[t.Name] = cols
	}

	return ms, nil
}

// getTableColumns 获取单个表的列信息
func getTableColumns(db *sql.DB, tableName string) ([]*sqlparser.Column, error) {
	// PRAGMA table_info 返回: cid, name, type, notnull, dflt_value, pk
	rows, err := db.Query("PRAGMA table_info(" + tableName + ")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []*sqlparser.Column
	for rows.Next() {
		var cid int
		var notnull, pk int
		var dfltValue any
		col := &sqlparser.Column{}

		err = rows.Scan(&cid, &col.Name, &col.ColumnType, &notnull, &dfltValue, &pk)
		if err != nil {
			return nil, err
		}

		col.Type = parseSQLiteType(col.ColumnType)
		col.Nullable = notnull == 0 // SQLite: notnull=1 表示 NOT NULL
		col.Default = dfltValue
		col.Index = cid

		// 主键判断
		if pk == 1 {
			col.Key = "PRI"
			// SQLite: INTEGER PRIMARY KEY 是自增的 (ROWID 别名)
			if strings.ToUpper(col.ColumnType) == "INTEGER" {
				col.AutoIncr = true
				col.Extra = "AUTOINCREMENT"
			}
		}

		// 解析长度/精度（如 VARCHAR(255)）
		col.Size, col.Precision, col.Scale = parseTypeInfo(col.ColumnType)

		// SQLite 不原生支持列注释
		col.Comment = ""

		cols = append(cols, col)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// 获取唯一键信息
	if err = fillUniqueKeys(db, tableName, cols); err != nil {
		return nil, err
	}

	return cols, nil
}

// fillUniqueKeys 填充唯一键信息
func fillUniqueKeys(db *sql.DB, tableName string, cols []*sqlparser.Column) error {
	// PRAGMA index_list 返回: seq, name, unique, origin, partial
	rows, err := db.Query("PRAGMA index_list(" + tableName + ")")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var seq int
		var indexName string
		var unique int
		var origin, partial string

		err = rows.Scan(&seq, &indexName, &unique, &origin, &partial)
		if err != nil {
			return err
		}

		// 只处理唯一索引，排除主键索引
		if unique == 1 && origin != "pk" {
			// PRAGMA index_info 返回: seqno, cid, name
			idxRows, err := db.Query("PRAGMA index_info(" + indexName + ")")
			if err != nil {
				return err
			}

			for idxRows.Next() {
				var seqno, cid int
				var colName string
				err = idxRows.Scan(&seqno, &cid, &colName)
				if err != nil {
					idxRows.Close()
					return err
				}
				// 标记唯一键列
				for _, col := range cols {
					if col.Name == colName && col.Key != "PRI" {
						col.Key = "UNI"
						break
					}
				}
			}
			idxRows.Close()
			if err = idxRows.Err(); err != nil {
				return err
			}
		}
	}
	return rows.Err()
}

// parseSQLiteType 将 SQLite 类型映射到通用类型
func parseSQLiteType(sqliteType string) string {
	upper := strings.ToUpper(sqliteType)

	// INTEGER 类型
	if strings.Contains(upper, "INT") {
		return "INTEGER"
	}

	// REAL / FLOAT / DOUBLE 类型
	if strings.Contains(upper, "REAL") ||
		strings.Contains(upper, "FLOA") ||
		strings.Contains(upper, "DOUB") {
		return "REAL"
	}

	// TEXT / CHAR / VARCHAR / CLOB 类型
	if strings.Contains(upper, "CHAR") ||
		strings.Contains(upper, "TEXT") ||
		strings.Contains(upper, "CLOB") {
		return "TEXT"
	}

	// BLOB 类型
	if strings.Contains(upper, "BLOB") {
		return "BLOB"
	}

	// NUMERIC / DECIMAL / BOOLEAN / DATE / DATETIME
	if strings.Contains(upper, "NUMERIC") ||
		strings.Contains(upper, "DECIMAL") ||
		strings.Contains(upper, "BOOL") ||
		strings.Contains(upper, "DATE") ||
		strings.Contains(upper, "TIME") {
		return "NUMERIC"
	}

	// 默认
	return "TEXT"
}

// parseTypeInfo 解析类型中的长度和精度信息
// 例如: VARCHAR(255) -> size=255
//
//	DECIMAL(10,2) -> precision=10, scale=2
func parseTypeInfo(columnType string) (size, precision, scale int) {
	if columnType == "" {
		return 0, 0, 0
	}

	// 提取括号中的内容
	start := strings.Index(columnType, "(")
	end := strings.Index(columnType, ")")
	if start == -1 || end == -1 || end <= start {
		return 0, 0, 0
	}

	content := columnType[start+1 : end]
	parts := strings.Split(content, ",")

	if len(parts) >= 1 {
		size, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
		precision = size
	}
	if len(parts) >= 2 {
		scale, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
	}

	return size, precision, scale
}

// func GetTables(db *ace.DB, dbName string) ([]*sqlparser.Table, error) {
// 	// 表名,表注释
// 	ts, err := getTables(db, dbName)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// 表字段信息
// 	ms, err := getColumns(db, dbName)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	for _, t := range ts {
// 		t.ColumnsX = ms[t.Name]
// 	}
// 	return ts, nil
// }
