package dialect

type Column struct {
	Name    string `db:"name"`
	Type    string `db:"type"`
	Default any    `db:"default"`
	Size    *int   `db:"size"`
	Key     string `db:"column_key"`
	Extra   string `db:"extra"`
	Comment string `db:"comment"`
}
type Table struct {
	Name    string
	Comment string
	Columns []Column
}
