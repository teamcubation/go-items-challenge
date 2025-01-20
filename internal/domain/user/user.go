package user

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type Credentials struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}
