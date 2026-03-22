package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserLoginCreds struct {
	Username string `json:"username" example:"hornero3dx"`
	Password string `json:"password" example:"hornero3dx"`
}
type TokenResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

type JwtCustomClaims struct {
	UserId    uint   `json:"user_id"`
	CompanyId uint   `json:"company_id"`
	Role      string `json:"role"`
	User      string `json:"user"`
	jwt.RegisteredClaims
}

type User struct {
	IdUser    uint      `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;size:120;not null"`
	Password  string    `gorm:"size:255;not null"` // hash
	Role      string    `gorm:"size:16;default:user"`
	Status    string    `gorm:"size:16;default:active"`
	IdCompany uint      `gorm:"not null" json:"id_company"`
	Company   *Company  `gorm:"foreignKey:IdCompany;references:IdCompany"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Company struct {
	IdCompany uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"uniqueIndex;size:120;not null"`
	Status    string    `gorm:"size:16;default:active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `gorm:"index"`
}

type CompanyCreateDTO struct {
	Name string `json:"name" example:"companytest"`
}

type UserCreateCreds struct {
	Username  string `json:"username" example:"usertest"`
	Password  string `json:"password" example:"usertest"`
	Role      string `json:"role" example:"user"`
	IdCompany uint   `json:"id_company" example:"1"`
}
type UserUpdateCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type UserDeleteCreds struct {
	IdUser uint `json:"id_user"`
}

type UserIDActivate struct {
	IdUser uint `json:"id_user"`
}
