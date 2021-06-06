package service

import (
	"wrap_error/dao"

	"github.com/pkg/errors"
)

type Service struct {
	DB *dao.DBConnection
}

func NewSrevice() *Service {
	return &Service{&dao.DBConnection{}}
}

func (s *Service) Login(userName, pwd string) error {
	user, err := s.DB.FindByUserName(userName)
	if err != nil {
		return err
	}
	if user.Pwd != pwd {
		return errors.New("wrong password")
	}
	return nil
}
