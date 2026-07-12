package service

type OrderID struct {
	ID uint `uri:"id" validate:"required && number"`
}

type OrderCreate struct {
	// UserID is validated for shape here; that the user actually exists is a
	// business rule checked in the usecase through the user module, not a
	// database rule that would reach across the module boundary.
	UserID uint  `json:"user_id" form:"user_id" validate:"required && number"`
	Amount int64 `json:"amount" form:"amount" validate:"required && number && min:1"`
}
