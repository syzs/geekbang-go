package service

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"wrap_error/dao"
)

func TestLogin(t *testing.T) {
	s := NewSrevice()

	// user not exists
	err := s.Login("Abby", "pwd")
	log.Print(err)
	assert.Equal(t, true, errors.Is(err, dao.SqlErrNoRows))

	// wrong pwd
	err = s.Login("Eva", "wrong pwd")
	log.Print(err)
	assert.Equal(t, false, errors.Is(err, dao.SqlErrNoRows))

}