package media

// #cgo LDFLAGS: -lmediandk
// #include "sys/types.h"
// #include "media/NdkMediaFormat.h"
import "C"

type MediaFormat C.AMediaFormat

func (mf *MediaFormat) CPointer() *C.AMediaFormat {
	return (*C.AMediaFormat)(mf)
}
