package gock

import "github.com/chocokacang/gock/database"

type ORM struct {
	// The ORM perform single create, update, delete operations in transactions by default to ensure database data integrity.
	// You can disable it by setting `SkipDefaultTransaction` to true
}

type DB struct {
	//
	Dialector database.Dialector
	srv       *Server
}
