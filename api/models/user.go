package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower((action)) {
	case "update":
		if u.Username == "" {
			return errors.New("Username can not empty!")
		}
		if u.Password == "" {
			return errors.New("Password can not empty!")
		}
		if u.Email == "" {
			return errors.New("Email can not empty")
		}
		if u.Email != "" {
			if err := checkmail.ValidateFormat(u.Email); err != nil {
				return errors.New("Invalid email")
			}
		}
		return nil
	case "login":
		if u.Username == "" {
			return errors.New("Username can not empty")
		}
		if u.Password == "" {
			return errors.New("Password can not empty")
		}
		return nil
	case "forgotpassword":
		if u.Email == "" {
			return errors.New("Email can not empty")
		}
		if u.Email != "" {
			return errors.New("Invalid email")
		}
		return nil
	default:
		if u.Username == "" {
			return errors.New("Username can not empty!")
		}
		if u.Password == "" {
			return errors.New("Password can not empty!")
		}
		if u.Email == "" {
			return errors.New("Email can not empty")
		}
		if u.Email != "" {
			if err := checkmail.ValidateFormat(u.Email); err != nil {
				return errors.New("Invalid email")
			}
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) UpdateUser(db *gorm.DB, uid uint32) (*User, error) {
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"username":  u.Username,
			"password":  u.Password,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}

	err = db.Model(&User{}).Where("id = ?", uid).Take(&User{}).Error
	if err != nil {
		return &User{}, db.Error
	}
	return u, nil
}

func (u *User) DeleteUser(db *gorm.DB, uid uint32) (int64, error) {
	result := db.Where("id = ?", uid).Take(&User{}).Delete(&User{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (u *User) UpdatePassword(db *gorm.DB) error {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&User{}).Where("email = ?", u.Email).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
