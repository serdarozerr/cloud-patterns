package validators

import (
	v "github.com/serdarozerr/request-reply/pkg"
)

type CreateUser struct {
	Name     string
	Email    string
	Password string
	Age      int
}

func (c CreateUser) Validate() map[string]string {
	errors := make(map[string]string)
	if v.ValidateEmptyField(c.Name) {
		errors["Name"] = "Name cannot be empty"
	}
	if v.ValidateEmptyField(c.Email) {
		errors["Email"] = "Email cannot be empty"
	}
	if v.ValidateEmptyField(c.Password) {
		errors["Password"] = "Password cannot be empty"
	}
	if c.Age == 0 {
		errors["Age"] = "Age cannot be zero"
	}
	return errors
}
