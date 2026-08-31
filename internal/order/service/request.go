package service

type OrderID struct {
	ID uint `uri:"id" validate:"required && number"`
}

type OrderCreate struct {
	// the user's existence is checked in the usecase, not by a DB constraint
	UserID uint  `json:"user_id" form:"user_id" validate:"required && number"`
	Amount int64 `json:"amount" form:"amount" validate:"required && number && min:1"`
}
