package dao

import (
	goerror "errors"
	"fmt"
	"github.com/pkg/errors"
)

var (
	// the custom sentinel error
	SqlErrNoRows = goerror.New("sql: no rows in result set")
)

type User struct {
	Name string
	Pwd  string
}

type MockDBConnection struct {
}

func (db *MockDBConnection) FindByUserName(userName string) (u *User, err error) {
	if userName != "Eva" {
		err = SqlErrNoRows
		// wrap sentinel error
		// 避免返回预定义的error，如果error在向上返回时，有中间逻辑层对其进行了封装，破坏原始错误信息，将导致上层的判等或断言失效
		//fmt.Errorf("query user from db failed: %v", err)
		err = errors.Wrap(err, fmt.Sprintf("user: %s not exists", userName))
		return
	}
	u = &User{"Eva", "pwd"}
	return
}
