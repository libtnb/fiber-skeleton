package data

import "github.com/go-rio/migrate"

// user_id is a plain column, not a foreign key: modules do not share tables,
// the reference is validated in biz.Users.
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
