package re

import "regexp"

const (
	// BeginRegexp regexp to match with Begin
	BeginRegexp = `^BEGIN\s(?P<id>\d+)$`
	// CommitRegexp regexp to match with Commit
	CommitRegexp = `^COMMIT\s(?P<id>\d+)$`
	// OperationRegexp regexp to match with Operation
	OperationRegexp = `^table\s` +
		`(?P<schema>\w+)` + `\.` +
		`(?P<table>\w+)` + `\:\s` +
		`(?P<operation>INSERT|UPDATE|DELETE)` + `\:\s` +
		`(?P<values>.+)` + `$`
	// RowRegexp regexp to match with the row statement
	RowRegexp = `` +
		`(?P<column>\w+)` + `\[` +
		`(?P<datatype>[\w\s]+)` + `\]\:` +
		`(?P<value>\d{1,}\.\d{1,}|'.*?'|\w+)`
)

var (
	// Begin regexp object for BeginRegexp
	Begin = regexp.MustCompile(BeginRegexp)
	// Commit regexp object for CommitRegexp
	Commit = regexp.MustCompile(CommitRegexp)
	// Operation regexp object for OperationRegexp
	Operation = regexp.MustCompile(OperationRegexp)
	// Row regexp object for RowRegexp
	Row = regexp.MustCompile(RowRegexp)
)
