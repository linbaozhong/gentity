package schema

type Column struct {
	Name     string `db:"name"`
	Type     string `db:"type"`
	Default  any    `db:"default"`
	Size     int    `db:"size"`
	Key      string `db:"column_key"`
	Extra    string `db:"extra"`
	AutoIncr bool   `db:"auto_increment"`
	Comment  string `db:"comment"`
}

//type Table struct {
//	Name    string
//	Comment string
//	Columns []Column
//}
