package media

// #cgo LDFLAGS: -lmediandk
// #include "sys/types.h"
// #include "media/NdkMediaCodec.h"
// #include "media/NdkMediaFormat.h"
import "C"

type MediaCodec C.AMediaCodec

func (mc *MediaCodec) CPointer() *C.AMediaCodec {
	return (*C.AMediaCodec)(mc)
}

func (mc *MediaCodec) SetParameters(params *MediaFormat) error {
	status := C.AMediaCodec_setParameters(mc.CPointer(), params.CPointer())
	return CMediaStatusToError(status)
}
