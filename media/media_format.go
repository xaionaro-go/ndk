package media

// #cgo LDFLAGS: -lmediandk
// #include "sys/types.h"
// #include "media/NdkMediaFormat.h"
import "C"

type AMediaFormat C.AMediaFormat

func (mf *AMediaFormat) CPointer() *C.AMediaFormat {
	return (*C.AMediaFormat)(mf)
}
