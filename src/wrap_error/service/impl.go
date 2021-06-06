package service

import (
	"log"
	"wrap_error/dao"

	"github.com/pkg/errors"
)

var (
	LoginErrorPrefix = "Login Failed"
)

type Service struct {
	DB *dao.MockDBConnection
}

func NewSrevice() *Service {
	return &Service{&dao.MockDBConnection{}}
}

func (s *Service) Login(userName, pwd string) (err error) {
	defer func() {
		// 日志记录错误 或 对 调试 有帮助的信息，否则即为噪音日志
		if err != nil {
			// Cause: 获取 root error
			log.Printf("original error: %T, %v\n", errors.Cause(err), errors.Cause(err))
			log.Printf("%s\n", err.Error())
			// %+v: 打印堆栈信息
			log.Printf("stack trace:\n%+v\n", err)
		}
	}()

	user, err := s.DB.FindByUserName(userName)
	// 调用其他包内的函数返回的error，不进行降级处理的话，则直接抛出
	// 注意：降级处理后的 error 不能再抛出
	if err != nil {
		return err
	}
	if user.Pwd != pwd {
		// wrap 携带堆栈信息，后续通过 WithMessage 添加错误描述
		// tips: 多次使用 wrap 会添加多次堆栈信息，需要注意 调用的第三方的库 或 其他的函数 返回的错误是否是 wrap 过的
		err = errors.New( "wrong password")
		return errors.WithMessage(err, LoginErrorPrefix)
	}
	return nil
}
