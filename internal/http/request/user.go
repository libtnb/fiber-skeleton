package request

type UserID struct {
	ID uint `uri:"id" validate:"required|number"`
}

type UserAdd struct {
	Name string `json:"name" form:"name" validate:"required|min_len:3|max_len:255"`
}

type UserUpdate struct {
	ID   uint   `uri:"id" validate:"required|number"`
	Name string `json:"name" form:"name" validate:"required|min_len:3|max_len:255"`
}
