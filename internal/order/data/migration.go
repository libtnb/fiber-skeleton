package data

import "github.com/go-rio/migrate"

// The order module owns its schema. user_id is a plain column, not a foreign
// key into the user module's table: modules do not share tables, so the
// reference is validated in the business layer (biz.Users), which keeps the
// two modules independently splittable into services.
func init() {
	migrate.Add("20260102000000_create_orders", func(s *migrate.Schema) {
		s.Create("orders", func(t *migrate.Table) {
			t.ID()
			t.BigInteger("user_id")
			t.BigInteger("amount")
			t.Timestamps()
			t.SoftDeletes()
			t.Index("user_id")
		})
	})
}
