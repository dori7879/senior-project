package model

import (
	"errors"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Users []*User

type User struct {
	ID           uint `gorm:"primaryKey"`
	FirstName    string
	LastName     string
	Email        string `gorm:"unique"`
	PasswordHash []byte
	DateJoined   time.Time
	LastLogin    time.Time
	IsActive     bool
}

type UserDtos []*UserDto

type UserDto struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	DateJoined string `json:"date_joined"`
	LastLogin  string `json:"last_login"`
	IsActive   string `json:"is_active"`
}

type UserForm struct {
	FirstName string `json:"first_name" form:"alpha_space,max=255"`
	LastName  string `json:"last_name" form:"alpha_space,max=255"`
	Email     string `json:"email" form:"required,email,max=255"`
	Password  string `json:"password" form:"required,max=255"`
	Role      string `json:"role" form:"alpha_space,max=16"`
}

func (hw User) ToDto() *UserDto {
	return &UserDto{
		ID:         hw.ID,
		FirstName:  hw.FirstName,
		LastName:   hw.LastName,
		Email:      hw.Email,
		DateJoined: hw.DateJoined.Format(time.RFC3339),
		LastLogin:  hw.LastLogin.Format(time.RFC3339),
		IsActive:   strconv.FormatBool(hw.IsActive),
	}
}

func (hws Users) ToDto() UserDtos {
	dtos := make([]*UserDto, len(hws))
	for i, hw := range hws {
		dtos[i] = hw.ToDto()
	}
	return dtos
}

func (f *UserForm) ToModel() (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(f.Password), 12)
	if err != nil {
		return nil, err
	}

	if f.FirstName == "" || f.LastName == "" {
		return nil, errors.New("UserForm does not have either first name or last name")
	}

	return &User{
		FirstName:    f.FirstName,
		LastName:     f.LastName,
		Email:        f.Email,
		PasswordHash: hashedPassword,
		DateJoined:   time.Now(),
		IsActive:     true,
	}, nil
}
