package dao

import (
	goerror "errors"
	"fmt"
	"github.com/pkg/errors"
)

var (
	// the custom sentinel error
	SqlErrNoRows = goerror.New("sql: no rows in result set")
	InternalErr  = goerror.New("sql: internal error")

	NotFound = goerror.New("not found")
)

type User struct {
	Name string
	Pwd  string
}

type MockDBConnection struct {
}

func (db *MockDBConnection) FindByUserName(userName string) (u *User, err error) {
	user, err := db.findByUserName(userName)
	// wrap sentinel error
	// 避免返回预定义的error，如果error在向上返回时，有中间逻辑层对其进行了封装，破坏原始错误信息，将导致上层的判等或断言失效
	//fmt.Errorf("query user from db failed: %v", err)
	if err == SqlErrNoRows {
		return nil, errors.Wrap(NotFound, "user not found")
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("query user:%s failed", userName))
	}
	return user, nil
}

func (db *MockDBConnection) findByUserName(userName string) (u *User, err error) {
	if userName == "Alice" {
		return nil, SqlErrNoRows
	} else if userName == "Eva" {
		u = &User{"Eva", "pwd"}
		return u, nil
	}
	return nil, InternalErr
}
