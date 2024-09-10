package schema

type Column struct {
	Name     string
	Type     string
	Default  any
	Size     int
	Unsigned bool
	Key      string
	Extra    string
	AutoIncr bool
	Comment  string
}

// type Table struct {
//	Name    string
//	Comment string
//	Columns []Column
// }
