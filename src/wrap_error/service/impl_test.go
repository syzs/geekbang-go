package service

import (
	"fmt"

	"testing"
	"wrap_error/dao"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	s := NewSrevice()

	// user not exists
	err := s.Login("Abby", "pwd")
	// errors.Is(err, dao.SqlErrNoRows) ==> err == dao.SqlErrNoRows
	assert.Equal(t, true, errors.Is(err, dao.SqlErrNoRows))
	assert.Equal(t, "sql: no rows in result set", fmt.Sprintf("%v", errors.Cause(err)))

	// wrong pwd
	err = s.Login("Eva", "wrong pwd")
	assert.Equal(t, false, errors.Is(err, dao.SqlErrNoRows))
	assert.Equal(t, "wrong password", fmt.Sprintf("%v", errors.Cause(err)))
	assert.Equal(t, "Login Failed: wrong password", err.Error())

	// login success
	err = s.Login("Eva", "pwd")
	assert.Equal(t, nil, err)
}
