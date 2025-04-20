//go:build !android
// +build !android

package callfn

func CallFn(fn uintptr) {
	panic("this platform is not supported")
}
