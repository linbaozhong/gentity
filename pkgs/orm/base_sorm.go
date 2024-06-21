package orm

import "fmt"

const (
	Quote_Char = "`"
)

type TableField struct {
	Name  string
	Json  string
	Table string
}

func (f TableField) IF(v1, v2 interface{}) string {
	return fmt.Sprintf("IF(%s,%v,%v) AS %s", f.Quote(), v1, v2, f.Name)
}

// AsName
func (f TableField) AsName(s string) string {
	return f.Quote() + " AS " + s
}

func (f TableField) PureQuote() string {
	return Quote_Char + f.Name + Quote_Char
}

func (f TableField) Quote() string {
	if f.Table == "" {
		return f.PureQuote()
	}
	return f.Table + "." + f.PureQuote()
}

func (f TableField) generate(op string) string {
	return f.Quote() + op + "?"
}
