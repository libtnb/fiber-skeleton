package data

import "github.com/go-rio/migrate"

// The module owns its schema; init registers it with the migrator.
func init() {
	migrate.Add("20260101000000_create_users", func(s *migrate.Schema) {
		s.Create("users", func(t *migrate.Table) {
			t.ID()
			t.String("name")
			t.Timestamps()
			t.SoftDeletes()
			// live rows cannot share a name; a soft-deleted row releases it
			t.Unique("name").Where("deleted_at IS NULL")
		})
	})
}
