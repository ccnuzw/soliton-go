package user

import (
	"errors"
	"reflect"
)

// EmailAddress 是领域值对象。
type EmailAddress struct {
	Value string `json:"value"`
}

// NewEmailAddress 创建一个新的 EmailAddress。
func NewEmailAddress(value string) (EmailAddress, error) {
	vo := EmailAddress{
		Value: value,
	}
	if err := vo.Validate(); err != nil {
		return EmailAddress{}, err
	}
	return vo, nil
}

// Validate 执行值对象的领域校验规则。
func (v EmailAddress) Validate() error {
	// TODO: 在此添加校验逻辑
	return errors.New("not implemented")
}

// Equals 比较两个值对象是否相等。
func (v EmailAddress) Equals(other EmailAddress) bool {
	return reflect.DeepEqual(v, other)
}
