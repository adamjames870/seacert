package migrations

import "embed"

// FS is the filesystem containing all SQL migrations.
//
//go:embed *.sql
var FS embed.FS
