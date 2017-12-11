package testdecoding

import "regexp"

var (
	opTable  = regexp.MustCompile(`^table\s\w.+\.(.*?)\:\s.*$`)
	opType   = regexp.MustCompile(`\:\s(.*?)\:`)
	opValues = regexp.MustCompile(`table\s\w.+\.\w.+\:\s\w.+\:\s(.*)`)
)

// Operation the operation unit
type Operation struct {
	Value string
	Table string
	Type  string
}

// ExtractOpTable extracts the database table
func ExtractOpTable(s string) string {
	return opTable.FindStringSubmatch(s)[1]
}

// ExtractOpType extracts the database table operation
func ExtractOpType(s string) string {
	return opType.FindStringSubmatch(s)[1]
}

// ExtractOpValues extracts the database table operation values
func ExtractOpValues(s string) string {
	return opValues.FindStringSubmatch(s)[1]
}

// NewOperation returns the operation struct
func NewOperation(s string) Operation {
	return Operation{
		Value: ExtractOpValues(s),
		Table: ExtractOpTable(s),
		Type:  ExtractOpType(s),
	}
}
