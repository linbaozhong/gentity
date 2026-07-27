package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ColumnInfo 列信息（用于 DDL 生成）
type ColumnInfo struct {
	Name    string // 列名
	Type    string // Go 类型
	Size    string // 长度（如 "20"）
	IsPK    bool   // 是否主键
	IsAuto  bool   // 是否自增
	NotNil  bool   // 是否 NOT NULL
	Comment string // 列注释
}

// Go 类型到 SQL 类型的映射表
type typeMapEntry struct {
	MySQL     string
	Postgres  string
	SQLite    string
	SQLServer string
}

var goTypeMap = map[string]typeMapEntry{
	"types.BigInt":  {"BIGINT", "BIGINT", "INTEGER", "BIGINT"},
	"types.Int":     {"INT", "INTEGER", "INTEGER", "INT"},
	"types.Int8":    {"TINYINT", "SMALLINT", "INTEGER", "TINYINT"},
	"types.Int16":   {"SMALLINT", "SMALLINT", "INTEGER", "SMALLINT"},
	"types.Int32":   {"INT", "INTEGER", "INTEGER", "INT"},
	"types.Int64":   {"BIGINT", "BIGINT", "INTEGER", "BIGINT"},
	"types.Uint":    {"INT UNSIGNED", "INTEGER", "INTEGER", "INT"},
	"types.Uint8":   {"TINYINT UNSIGNED", "SMALLINT", "INTEGER", "TINYINT"},
	"types.Uint16":  {"SMALLINT UNSIGNED", "SMALLINT", "INTEGER", "SMALLINT"},
	"types.Uint32":  {"INT UNSIGNED", "INTEGER", "INTEGER", "INT"},
	"types.Uint64":  {"BIGINT UNSIGNED", "BIGINT", "INTEGER", "BIGINT"},
	"types.String":  {"VARCHAR", "VARCHAR", "TEXT", "NVARCHAR"},
	"types.Float32": {"FLOAT", "REAL", "REAL", "FLOAT"},
	"types.Float64": {"DOUBLE", "DOUBLE PRECISION", "REAL", "FLOAT"},
	"types.Bool":    {"TINYINT(1)", "BOOLEAN", "INTEGER", "BIT"},
	"types.Time":    {"DATETIME", "TIMESTAMP", "TEXT", "DATETIME2"},
	"types.Money":   {"BIGINT UNSIGNED", "BIGINT", "INTEGER", "BIGINT"},
	"time.Time":     {"DATETIME", "TIMESTAMP", "TEXT", "DATETIME2"},
	"string":        {"VARCHAR", "VARCHAR", "TEXT", "NVARCHAR"},
	"int":           {"INT", "INTEGER", "INTEGER", "INT"},
	"int8":          {"TINYINT", "SMALLINT", "INTEGER", "TINYINT"},
	"int16":         {"SMALLINT", "SMALLINT", "INTEGER", "SMALLINT"},
	"int32":         {"INT", "INTEGER", "INTEGER", "INT"},
	"int64":         {"BIGINT", "BIGINT", "INTEGER", "BIGINT"},
	"uint":          {"INT UNSIGNED", "INTEGER", "INTEGER", "INT"},
	"uint8":         {"TINYINT UNSIGNED", "SMALLINT", "INTEGER", "TINYINT"},
	"uint16":        {"SMALLINT UNSIGNED", "SMALLINT", "INTEGER", "SMALLINT"},
	"uint32":        {"INT UNSIGNED", "INTEGER", "INTEGER", "INT"},
	"uint64":        {"BIGINT UNSIGNED", "BIGINT", "INTEGER", "BIGINT"},
	"float32":       {"FLOAT", "REAL", "REAL", "FLOAT"},
	"float64":       {"DOUBLE", "DOUBLE PRECISION", "REAL", "FLOAT"},
	"bool":          {"TINYINT(1)", "BOOLEAN", "INTEGER", "BIT"},
}

func getSQLType(driver, goType, size string) string {
	entry, ok := goTypeMap[goType]
	if !ok {
		return "VARCHAR(255)"
	}
	var sqlType string
	switch driver {
	case "mysql":
		sqlType = entry.MySQL
	case "postgres":
		sqlType = entry.Postgres
	case "sqlite":
		sqlType = entry.SQLite
	case "sqlserver":
		sqlType = entry.SQLServer
	default:
		sqlType = entry.MySQL
	}
	if size != "" && needsSize(sqlType) {
		if strings.Contains(sqlType, "(") {
			sqlType = strings.Split(sqlType, "(")[0] + "(" + size + ")"
		} else {
			sqlType = sqlType + "(" + size + ")"
		}
	}
	return sqlType
}

func needsSize(sqlType string) bool {
	return strings.Contains(sqlType, "VARCHAR") ||
		strings.Contains(sqlType, "NVARCHAR") ||
		(strings.Contains(sqlType, "INT") && !strings.Contains(sqlType, "BIGINT") &&
			!strings.Contains(sqlType, "TINYINT") && !strings.Contains(sqlType, "SMALLINT") &&
			!strings.Contains(sqlType, "UNSIGNED"))
}

// generateDDL 根据结构体生成 DDL 文件到指定目录
func generateDDL(tds []TempData, driver, outDir string) error {
	dir, _ := filepath.Abs(outDir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, td := range tds {
		tableName := td.TableName
		if tableName == "" {
			tableName = getFieldName(td.StructName)
		}

		var columns []ColumnInfo
		for _, f := range td.Columns {
			col := ColumnInfo{
				Name:    f.Col,
				Type:    f.Type,
				Size:    f.Size,
				IsPK:    td.HasPrimaryKey && td.PrimaryKey.Col == f.Col,
				IsAuto:  f.Auto,
				Comment: f.Comment,
			}
			if strings.HasPrefix(f.Type, "*") {
				col.NotNil = false
			} else {
				col.NotNil = true
			}
			columns = append(columns, col)
		}

		ddl := generateCreateTable(driver, tableName, td.TableComment, columns)

		outFile := filepath.Join(dir, tableName+".sql")
		if err := os.WriteFile(outFile, []byte(ddl), 0644); err != nil {
			return err
		}
		fmt.Println("Generated:", outFile)
	}
	return nil
}

func generateCreateTable(driver, tableName, tableComment string, columns []ColumnInfo) string {
	var sb strings.Builder
	q := getQuoter(driver)

	sb.WriteString("-- " + driver + " DDL for " + tableName + "\n")
	sb.WriteString("CREATE TABLE IF NOT EXISTS ")
	sb.WriteString(q(tableName))
	sb.WriteString(" (\n")

	var pkCol string
	for i, col := range columns {
		sb.WriteString("  ")
		sb.WriteString(q(col.Name))
		sb.WriteString(" ")
		sb.WriteString(getSQLType(driver, col.Type, col.Size))

		if col.IsAuto {
			sb.WriteString(" " + getAutoIncrement(driver))
			sb.WriteString(" NOT NULL")
		} else if col.NotNil {
			sb.WriteString(" NOT NULL")
		} else {
			sb.WriteString(" NULL")
		}

		if col.Comment != "" && (driver == "mysql" || driver == "postgres" || driver == "sqlserver") {
			sb.WriteString(" COMMENT " + sqlQuote(col.Comment))
		}

		if col.IsPK {
			pkCol = col.Name
		}

		if i < len(columns)-1 || col.IsPK {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}

	if pkCol != "" {
		sb.WriteString("  PRIMARY KEY (" + q(pkCol) + ")\n")
	}

	sb.WriteString(")")

	if tableComment != "" && (driver == "mysql" || driver == "postgres" || driver == "sqlserver") {
		sb.WriteString(" COMMENT=" + sqlQuote(tableComment))
	}

	sb.WriteString(getTableOptions(driver))
	sb.WriteString(";\n")

	return sb.String()
}

func sqlQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func getQuoter(driver string) func(string) string {
	switch driver {
	case "postgres", "sqlite":
		return func(name string) string { return `"` + name + `"` }
	case "sqlserver":
		return func(name string) string { return "[" + name + "]" }
	default:
		return func(name string) string { return "`" + name + "`" }
	}
}

func getAutoIncrement(driver string) string {
	switch driver {
	case "postgres":
		return "AUTO_INCREMENT"
	case "sqlserver":
		return "IDENTITY(1,1)"
	case "sqlite":
		return "AUTOINCREMENT"
	default:
		return "AUTO_INCREMENT"
	}
}

func getTableOptions(driver string) string {
	switch driver {
	case "mysql":
		return " ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	default:
		return ""
	}
}
