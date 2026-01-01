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

	// 1. Preparar contraseñas hasheadas
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPassword := string(hashedBytes)

	// 2. Crear usuarios de prueba
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

	// 4. Casos de prueba
	t.Run("Login exitoso con usuario activo", func(t *testing.T) {
		creds := types.UserLoginCreds{
			Username: "active_user",
			Password: "password123",
		}

		token, err := AuthUseCase(creds)

		if err != nil {
			t.Errorf("No debería haber error en login activo: %v", err)
		}
		if token == nil || token.Token == "" {
			t.Error("Se esperaba un token válido")
		}
	})

	t.Run("Login fallido con usuario deshabilitado (Status Check)", func(t *testing.T) {
		creds := types.UserLoginCreds{
			Username: "banned_user",
			Password: "password123",
		}

		token, err := AuthUseCase(creds)

		if token != nil {
			t.Error("No se debería generar un token para un usuario deshabilitado")
		}
		if err == nil || err.Error() != "this account is disabled, please contact support" {
			t.Errorf("Se esperaba error de cuenta deshabilitada, se obtuvo: %v", err)
		}
	})

	t.Run("Login fallido por contraseña incorrecta", func(t *testing.T) {
		creds := types.UserLoginCreds{
			Username: "active_user",
			Password: "wrong_password",
		}

		_, err := AuthUseCase(creds)

		if err == nil || err.Error() != "Invalid username or password" {
			t.Error("Se esperaba error de credenciales inválidas")
		}
	})
}
