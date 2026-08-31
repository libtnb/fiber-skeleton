package migrations

import "github.com/go-rio/migrate"

func init() {
	collection.Add("20260101000000_create_users_table", func(s *migrate.Schema) {
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
