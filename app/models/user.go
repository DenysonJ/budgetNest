package models

import (
	contractsorm "github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/database/orm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	orm.Model
	Id       int
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	orm.SoftDeletes
}

func (u *User) DispatchesEvents() map[contractsorm.EventType]func(contractsorm.Event) error {
	return map[contractsorm.EventType]func(contractsorm.Event) error{
		contractsorm.EventCreating: func(event contractsorm.Event) error {
			userPassword := event.GetAttribute("Password").(string)
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			event.SetAttribute("Password", string(hashedPassword))
			return nil
		},
	}
}
