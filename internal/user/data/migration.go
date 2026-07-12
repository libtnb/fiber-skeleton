package data

import "github.com/go-rio/migrate"

// The user module owns its schema: the migration is registered at init time,
// compiled into the binary and run by the migrator like every other module's.
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
