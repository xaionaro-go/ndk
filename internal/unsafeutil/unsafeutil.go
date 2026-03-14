package unsafeutil

import "unsafe"

// SlicePointer returns an unsafe.Pointer to the first element of a slice.
// Returns nil if the slice is empty.
func SlicePointer[T any](s []T) unsafe.Pointer {
	if len(s) == 0 {
		return nil
	}
	return unsafe.Pointer(&s[0])
}
