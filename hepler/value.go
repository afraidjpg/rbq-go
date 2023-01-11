package hepler

// GetPoint 获取变量的指针
func GetPoint[T any](v T) *T {
	return &v
}
