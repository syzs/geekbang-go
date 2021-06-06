package dao

import (
	goerror "errors"

	"github.com/pkg/errors"
)

var (
	SqlErrNoRows = goerror.New("ErrNoRows")
)

type User struct {
	Name string
	Pwd  string
}

type DBConnection struct {
}

func (db *DBConnection) FindByUserName(userName string) (u *User, err error) {
	if userName != "Eva" {
		err = SqlErrNoRows
		err = errors.Wrap(err, "user not exists")
		return
	}
	u = &User{"Eva", "pwd"}
	return
}
