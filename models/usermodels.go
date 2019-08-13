package models

import "golang.org/x/crypto/bcrypt"

type UserLogin struct {
	UserName string `gorm:"UNIQUE;NOT NULL" json:"username"`
	Password string `gorm:"NOT NULL" json:"password"`
}

type User struct {
	ID   int    `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	Name string `gorm:"NOT NULL" json:"name"`
	UserLogin
}

type UserDetails struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"username"`
}

type UserUpdate struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u *User) BeforeSave() (err error) {
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	u.Password = string(hash)
	return
}
