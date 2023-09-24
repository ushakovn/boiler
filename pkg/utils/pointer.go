package utils

func PtrVal[T any](ptr *T) T {
	if ptr == nil {
		return *new(T)
	}
	return *ptr
}
