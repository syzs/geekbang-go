package service

import (
	"log"
	"wrap_error/dao"

	"github.com/pkg/errors"
)

var (
	LoginErrorPrefix = "Login Failed: "
)

type Service struct {
	DB *dao.MockDBConnection
}

func NewSrevice() *Service {
	return &Service{&dao.MockDBConnection{}}
}

func (s *Service) Login(userName, pwd string) (err error) {
	defer func() {
		if err != nil {
			// %+v: 打印堆栈信息
			log.Printf("original error: %T, %v\n", errors.Cause(err), errors.Cause(err))
			log.Printf("stack trace:\n%+v\n", err)
		}
	}()
	user, err := s.DB.FindByUserName(userName)
	// 日志记录错误 或 对 调试 有帮助的信息，否则即为噪音日志
	if err != nil {
		// wrap 携带堆栈信息，后续通过 WithMessage 添加错误描述
		// tips: 多次使用 wrap 会添加多次堆栈信息，需要注意下层返回的错误是都是 wrap 过的
		errors.WithMessage(err, LoginErrorPrefix)
		return err
	}
	if user.Pwd != pwd {
		err = errors.New(LoginErrorPrefix + " wrong password")
		return errors.WithMessage(err, LoginErrorPrefix)
	}
	return nil
}
