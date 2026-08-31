// Package migrations is the schema history: one file per migration, named
// after the migration it registers, applied in lexical order.
package migrations

import "github.com/go-rio/migrate"

// collection is filled by each migration file's init.
var collection = migrate.NewCollection()

// Collection returns the full, ordered migration set.
func Collection() *migrate.Collection { return collection }
