package user

// ActiveUserSpec 表示一个领域规格（Specification）。
type ActiveUserSpec struct{}

// IsSatisfiedBy 判断目标对象是否满足规格。
func (s ActiveUserSpec) IsSatisfiedBy(target *User) bool {
	// TODO: 实现规格校验逻辑
	return true
}
