// Package migration holds the schema migrations: one file per migration,
// registered at init time and compiled into the binary. Rollbacks are derived
// automatically from the declarations.
package migration

import "github.com/libtnb/migrate"

func init() {
	migrate.Add("20260101000000_create_users", func(s *migrate.Schema) {
		s.Create("users", func(t *migrate.Table) {
			t.ID()
			t.String("name")
			t.Timestamps()
			t.SoftDeletes()
		})
	})
}
