package user

import (
	"errors"
	"fmt"
)

//params {"username":"alisher","password":"lovelove"}
//token {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoiYWxpc2hlciIsImlkIjoiNWRkZTI4YjU0OWMxMTVlNGFmMDIyMzhiIn0sImlhdCI6MTU3NDg0MDUwMSwiZXhwIjoxNTc1NDQ1MzAxfQ.kSUyOCd4SVl4ja7EJGhMYUp61gVnelh3m5H7hFpc_Zs"}
type User struct {
	ID       uint32
	Login    string
	password string
}

type UserRepo struct {
	data map[string]*User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		data: make(map[string]*User),
	}
}

func (r *UserRepo) GetData(username string) (*User, bool) {
	u, ok := r.data[username]
	return u, ok
}

var (
	ErrNoUser              = errors.New("user not found")
	ErrBadPass             = errors.New("invald password")
	ErrAlreadyExist        = errors.New("already exists")
	autoID          uint32 = 1
)

func (repo *UserRepo) Authorize(login, pass string) (*User, error) {
	u, ok := repo.data[login]
	if !ok {
		return nil, ErrNoUser
	}

	if u.password != pass {
		return nil, ErrBadPass
	}

	return u, nil
}

func (repo *UserRepo) Register(login, pass string) (*User, error) {
	if _, ok := repo.data[login]; ok {
		return nil, ErrAlreadyExist
	}
	u := &User{
		ID:       autoID,
		Login:    login,
		password: pass,
	}
	repo.data[login] = u
	for k, v := range repo.data {
		fmt.Println("REPO:\n\t", k, v)
	}
	autoID++
	return u, nil
}
