package store

import (
	"context"
	"crypto/rand"
	// "errors"
	// "fmt"
	"time"

	// "github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID 		       int    `pg:"-"`
	Username       string `pg:"username,unique" binding:"required,min=5,max=30"`
	Password       string `pg:"-" binding:"required,min=7,max=32"`
	HashedPassword []byte `json:"-"`
	Salt           []byte `json:"-"`
	CreatedAt      time.Time
	ModifiedAt     time.Time
	Posts          []*Post `json:"-" pg:"fk:user_id,rel:has-many,on_delete:CASCADE"`
}

var Users []*User
var _ pg.AfterSelectHook = (*User)(nil)

func (user *User) AfterSelect(ctx context.Context) error {
	if user.Posts == nil {
		user.Posts = []*Post{}
	}
	return nil
}

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

func FetchUser(id int) (*User, error) {
	user := new(User)
	user.ID = id
	err := db.Model(user).Returning("*").WherePK().Select()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching user")
		return nil, err
	}
	return user, nil
}
