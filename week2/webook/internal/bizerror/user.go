package bizerror

import (
	"fmt"
	"strings"
)

const (
	IncorrectUserNameOrPassword = iota + 10001
)

// ErrorIdentify 扩展error接口，带唯一标识。
type ErrorIdentify interface {
	// 接口继承
	error
	// ID 返回error的唯一标识
	ID() int32
}

// BizError 再次扩展ErrorIdentify，扩展其他能力
type BizError interface {
	// ErrorIdentify 接口继承
	ErrorIdentify

	// Name error名称
	Name() string
	// Message error简单描述
	Message() string
	// Unwrap 返回包装的error
	Unwrap() error
	// Errorf 提供重写error错误信息的能力，返回一个新error（避免改变原本的error）
	Errorf(format string, a ...interface{}) BizError
}

// bizError 实现IError接口
type bizError struct {
	ErrID   int32
	ErrName string
	Msg     string
	wrapErr error
	// ErrorMsgList 用于包装error的多个错误描述，有些场景用户需要重写一个error的描述，多次重写的描述通过这个字段承载
	ErrorMsgList []string
}

func New(id int32, name string, msg string) *bizError {
	return &bizError{
		ErrID:   id,
		ErrName: name,
		Msg:     msg,
	}
}

func (b *bizError) ID() int32 {
	return b.ErrID
}

func (b *bizError) Name() string {
	return b.ErrName
}

func (b *bizError) Message() string {
	return b.Msg
}

func (b *bizError) Unwrap() error {
	return b.wrapErr
}

func (b *bizError) Errorf(format string, a ...interface{}) BizError {
	// 之所以clone，是为了避免改变原本的error
	newBizError := b.clone()
	err0 := fmt.Errorf(format, a...)
	if newBizError.wrapErr == nil {
		newBizError.wrapErr = err0
	}

	// append重写的error信息
	newBizError.ErrorMsgList = append(newBizError.ErrorMsgList, err0.Error())
	return newBizError
}

func (b *bizError) clone() *bizError {
	return &bizError{
		ErrID:        b.ErrID,
		ErrName:      b.ErrName,
		Msg:          b.Msg,
		wrapErr:      b.wrapErr,
		ErrorMsgList: b.ErrorMsgList,
	}
}

// Error 返回error信息
func (b *bizError) Error() string {
	return fmt.Sprintf("id=%d, name=%s, msg=%s, wrapErr=[%v], errorMsgList=[%v]",
		b.ID(), b.Name(), b.Message(), b.Unwrap(), strings.Join(b.ErrorMsgList, "|"))
}

// Is 比较error是否相等
func (b *bizError) Is(err error) bool {
	err0, ok := err.(*bizError)
	if ok {
		// 通过ID比较相等
		return err0.ErrID == b.ErrID
	}
	return false
}
