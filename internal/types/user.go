package types

import "time"

type UserLoginCreds struct {
	Username string `json:"username" example:"hornero3dx"`
	Password string `json:"password" example:"hornero3dx"`
}
type TokenResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;size:120;not null"`
	Password  string    `gorm:"size:255;not null"` // hash
	Role      string    `gorm:"size:16;default:user"`
	Status    string    `gorm:"size:16;default:active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserCreateCreds struct {
	Username string `json:"username" example:"usertest"`
	Password string `json:"password" example:"usertest"`
	Role     string `json:"role" example:"user"`
}
type UserUpdateCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type UserDeleteCreds struct {
	Username string `json:"username"`
}
