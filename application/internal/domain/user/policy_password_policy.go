package user

import "errors"

// PasswordPolicy 表示一个领域策略（Policy）。
type PasswordPolicy struct{}

// Validate 执行策略校验，返回错误表示不满足策略。
func (p PasswordPolicy) Validate(target *User) error {
	// TODO: 实现策略校验逻辑
	return errors.New("not implemented")
}
