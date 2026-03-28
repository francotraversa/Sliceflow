package services

import (
	"os"
	"testing"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Couldn't connect to DB: %v", err)
	}
	db.AutoMigrate(&types.User{})

	storage.OverrideDatabaseInstance(db)
	return db
}
func TestAuthUseCase(t *testing.T) {
	db := setupTest(t)

	// 1. Prepare hashed passwords
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPassword := string(hashedBytes)

	// 2. Create test users
	activeUser := types.User{
		Username: "active_user",
		Password: hashedPassword,
		Status:   "active",
		Role:     "user",
	}
	disabledUser := types.User{
		Username: "banned_user",
		Password: hashedPassword,
		Status:   "disabled",
		Role:     "user",
	}

	db.Create(&activeUser)
	db.Create(&disabledUser)

	os.Setenv("JWT_SECRET", "test_secret")
	os.Setenv("TTL", "60")

	// 4. Test cases
	t.Run("Successful login with active user", func(t *testing.T) {
		creds := types.UserLoginCreds{
			Username: "active_user",
			Password: "password123",
		}

		token, err := AuthUseCase(creds)

		if err != nil {
			t.Errorf("Should not have errored on active login: %v", err)
		}
		if token == nil || token.Token == "" {
			t.Error("Expected a valid token")
		}
	})

	t.Run("Failed login with disabled user (Status Check)", func(t *testing.T) {
		creds := types.UserLoginCreds{
			Username: "banned_user",
			Password: "password123",
		}

		token, err := AuthUseCase(creds)

		if token != nil {
			t.Error("Should not generate a token for a disabled user")
		}
		if err == nil || err.Error() != "this account is disabled, please contact support" {
			t.Errorf("Expected disabled account error, got: %v", err)
		}
	})

	t.Run("Failed login due to wrong password", func(t *testing.T) {
		creds := types.UserLoginCreds{
			Username: "active_user",
			Password: "wrong_password",
		}

		_, err := AuthUseCase(creds)

		if err == nil || err.Error() != "Invalid username or password" {
			t.Error("Expected invalid credentials error")
		}
	})
}
