package store

import (
	"crypto/rand"
	// "errors"
	// "fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/rs/zerolog/log"
)

type User struct {
	ID 		       int    `pg:"-"`
	Username       string `pg:"username,unique" binding:"required,min=5,max=30"`
	Password       string `pg:"-" binding:"required,min=7,max=32"`
	HashedPassword []byte `json:"-"`
	Salt           []byte `json:"-"`
	CreatedAt      time.Time
	ModifiedAt     time.Time
}

var Users []*User

func AddUser(user *User) error {
	salt, err := GenerateSalt()
	if err != nil {
		return err
	}
	toHash := append([]byte(user.Password), salt...)
	hashedPassword, err := bcrypt.GenerateFromPassword(toHash, bcrypt.DefaultCost)
	
	user.Salt = salt
	user.HashedPassword = hashedPassword
	
	_, errInsertion := db.Model(user).Returning("*").Insert()
	if err != nil {
		return errInsertion
	}

	return nil
}

func Authenticate(username string, password string) (*User, error) {
	user := new(User)
	if err := db.Model(user).Where(
		"username = ?", username).Select(); err != nil {
		return nil, err
	}
	salted := append([]byte(password), user.Salt...)
	if errCheck := bcrypt.CompareHashAndPassword(user.HashedPassword, salted); 
	errCheck != nil {
		return nil, errCheck
	}
	return user, nil
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		log.Error().Err(err).Msg("Unable to create salt")
		return nil, err
	}
	return salt, nil
}