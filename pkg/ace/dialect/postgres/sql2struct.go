package postgres

import (
	"database/sql"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/sqlparser"
	"strings"
)

func Tables(db *sql.DB, dbName string) ([]*sqlparser.Table, error) {
	// PostgreSQL 获取表名和注释
	rows, err := db.Query(`
		SELECT 
			t.table_name,
			COALESCE(d.description, '') as table_comment
		FROM information_schema.tables t
		LEFT JOIN pg_description d ON d.objoid = (t.table_schema || '.' || t.table_name)::regclass::oid
		WHERE t.table_schema = $1 
		  AND t.table_type = 'BASE TABLE'
	`, dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ts := make([]*sqlparser.Table, 0)
	for rows.Next() {
		var tableName, comment string
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

func Columns(db *sql.DB, dbName string) (map[string][]*sqlparser.Column, error) {
	// PostgreSQL 获取列信息
	rows, err := db.Query(`
		SELECT 
			c.table_name,
			c.column_name,
			c.column_default,
			c.data_type,
			c.udt_name,
			COALESCE(c.character_maximum_length, 0),
			COALESCE(c.numeric_precision, 0),
			COALESCE(c.numeric_scale, 0),
			c.is_nullable,
			c.is_identity,
			COALESCE(col_description((c.table_schema || '.' || c.table_name)::regclass::oid, c.ordinal_position), '')
		FROM information_schema.columns c
		WHERE c.table_schema = $1
		ORDER BY c.table_name, c.ordinal_position
	`, dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms := make(map[string][]*sqlparser.Column)
	for rows.Next() {
		var tableName string
		col := &sqlparser.Column{}
		var isNullable, isIdentity, comment string
		err = rows.Scan(
			&tableName,
			&col.Name,
			&col.Default,
			&col.Type,
			&col.ColumnType,
			&col.Size,
			&col.Precision,
			&col.Scale,
			&isNullable,
			&isIdentity,
			&comment,
		)
		if err != nil {
			return nil, err
		}

		col.Comment = comment
		col.Nullable = isNullable == "YES"

		// 判断自增：PostgreSQL 10+ 用 is_identity，旧版本用 column_default
		if isIdentity == "YES" {
			col.AutoIncr = true
			col.Extra = "IDENTITY"
		} else if col.Default != nil && strings.Contains(fmt.Sprint(col.Default), "nextval") {
			col.AutoIncr = true
			col.Extra = "SERIAL"
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

	// 获取主键信息
	if err = fillPrimaryKeys(db, dbName, ms); err != nil {
		return nil, err
	}

	// 获取唯一键信息
	if err = fillUniqueKeys(db, dbName, ms); err != nil {
		return nil, err
	}

	return ms, nil
}

// fillPrimaryKeys 填充主键信息
func fillPrimaryKeys(db *sql.DB, dbName string, ms map[string][]*sqlparser.Column) error {
	rows, err := db.Query(`
		SELECT 
			tc.table_name,
			kcu.column_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu 
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		WHERE tc.table_schema = $1
		  AND tc.constraint_type = 'PRIMARY KEY'
	`, dbName)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, columnName string
		err = rows.Scan(&tableName, &columnName)
		if err != nil {
			return err
		}
		if cols, ok := ms[tableName]; ok {
			for _, col := range cols {
				if col.Name == columnName {
					col.Key = "PRI"
					break
				}
			}
		}
	}
	return rows.Err()
}

// fillUniqueKeys 填充唯一键信息
func fillUniqueKeys(db *sql.DB, dbName string, ms map[string][]*sqlparser.Column) error {
	rows, err := db.Query(`
		SELECT 
			tc.table_name,
			kcu.column_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu 
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		WHERE tc.table_schema = $1
		  AND tc.constraint_type = 'UNIQUE'
	`, dbName)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, columnName string
		err = rows.Scan(&tableName, &columnName)
		if err != nil {
			return err
		}
		if cols, ok := ms[tableName]; ok {
			for _, col := range cols {
				if col.Name == columnName {
					col.Key = "UNI"
					break
				}
			}
		}
	}
	return rows.Err()
}

func (m *PostgreSQL) GetTables(db *sql.DB, dbName string) ([]*sqlparser.Table, error) {
	// 表名,表注释
	ts, err := Tables(db, dbName)
	if err != nil {
		return nil, err
	}

	// 表字段信息
	ms, err := Columns(db, dbName)
	if err != nil {
		return nil, err
	}

	for _, t := range ts {
		t.ColumnsX = ms[t.Name]
	}
	return ts, nil
}
