// Package migrations holds the versioned SQL migrations that define the
// database schema. They are the single source of truth for the schema and are
// embedded into the API binary so every deployment applies the same set.
package migrations

import "embed"

// FS contains every NNNNNN_description.{up,down}.sql file in this directory.
//
//go:embed *.sql
var FS embed.FS
