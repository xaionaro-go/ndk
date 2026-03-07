// Simulates c-for-go output for Android Media NDK.
// This file is parsed at AST level only; it does not compile.
package media

import "unsafe"

// Opaque handle types.
type AMediaCodec C.AMediaCodec
type AMediaFormat C.AMediaFormat
type AMediaExtractor C.AMediaExtractor
type AMediaMuxer C.AMediaMuxer
type AMediaDrm C.AMediaDrm
type AMediaCrypto C.AMediaCrypto

// Buffer info (opaque for CGo bridge).
type AMediaCodecBufferInfo C.AMediaCodecBufferInfo

// Integer typedefs.
type Media_status_t int32

// Status codes.
const (
	AMEDIA_OK                       Media_status_t = 0
	AMEDIA_ERROR_UNKNOWN            Media_status_t = -10000
	AMEDIA_ERROR_MALFORMED          Media_status_t = -10001
	AMEDIA_ERROR_UNSUPPORTED        Media_status_t = -10002
	AMEDIA_ERROR_INVALID_OBJECT     Media_status_t = -10003
	AMEDIA_ERROR_INVALID_PARAMETER  Media_status_t = -10004
	AMEDIA_ERROR_INVALID_OPERATION  Media_status_t = -10005
	AMEDIA_ERROR_END_OF_STREAM      Media_status_t = -10006
	AMEDIA_ERROR_IO                 Media_status_t = -10007
	AMEDIA_ERROR_WOULD_BLOCK        Media_status_t = -10008
)

// --- MediaCodec functions ---
func AMediaCodec_createDecoderByType(mimeType *byte) *AMediaCodec { return nil }
func AMediaCodec_createEncoderByType(mimeType *byte) *AMediaCodec { return nil }
func AMediaCodec_delete(codec *AMediaCodec) Media_status_t        { return 0 }
func AMediaCodec_configure(codec *AMediaCodec, format *AMediaFormat, window unsafe.Pointer, crypto *AMediaCrypto, flags uint32) Media_status_t {
	return 0
}
func AMediaCodec_start(codec *AMediaCodec) Media_status_t { return 0 }
func AMediaCodec_stop(codec *AMediaCodec) Media_status_t  { return 0 }
func AMediaCodec_flush(codec *AMediaCodec) Media_status_t { return 0 }
func AMediaCodec_dequeueInputBuffer(codec *AMediaCodec, timeoutUs int64) int32 {
	return 0
}
func AMediaCodec_dequeueOutputBuffer(codec *AMediaCodec, info *AMediaCodecBufferInfo, timeoutUs int64) int32 {
	return 0
}
func AMediaCodec_getOutputFormat(codec *AMediaCodec) *AMediaFormat { return nil }
func AMediaCodec_queueInputBuffer(codec *AMediaCodec, index int32, offset int32, size int32, time uint64, flags uint32) Media_status_t {
	return 0
}
func AMediaCodec_releaseOutputBuffer(codec *AMediaCodec, index int32, render bool) Media_status_t {
	return 0
}

// --- MediaFormat functions ---
func AMediaFormat_new() *AMediaFormat                                        { return nil }
func AMediaFormat_delete(format *AMediaFormat) Media_status_t                { return 0 }
func AMediaFormat_setString(format *AMediaFormat, name *byte, value *byte)   {}
func AMediaFormat_setInt32(format *AMediaFormat, name *byte, value int32)    {}
func AMediaFormat_getInt32(format *AMediaFormat, name *byte, out *int32) bool { return false }
func AMediaFormat_getString(format *AMediaFormat, name *byte, out **byte) bool { return false }

// --- MediaExtractor functions ---
func AMediaExtractor_new() *AMediaExtractor                        { return nil }
func AMediaExtractor_delete(ex *AMediaExtractor) Media_status_t    { return 0 }
func AMediaExtractor_setDataSourceFd(ex *AMediaExtractor, fd int32, offset int64, length int64) Media_status_t {
	return 0
}
func AMediaExtractor_getTrackCount(ex *AMediaExtractor) int32          { return 0 }
func AMediaExtractor_getTrackFormat(ex *AMediaExtractor, idx int32) *AMediaFormat {
	return nil
}
func AMediaExtractor_selectTrack(ex *AMediaExtractor, idx int32) Media_status_t { return 0 }
func AMediaExtractor_advance(ex *AMediaExtractor) bool                          { return false }
func AMediaExtractor_readSampleData(ex *AMediaExtractor, buffer unsafe.Pointer, capacity int32) int32 {
	return 0
}
func AMediaExtractor_getSampleTime(ex *AMediaExtractor) int64 { return 0 }

// --- MediaMuxer functions ---
func AMediaMuxer_new(fd int32, format int32) *AMediaMuxer               { return nil }
func AMediaMuxer_delete(muxer *AMediaMuxer) Media_status_t              { return 0 }
func AMediaMuxer_addTrack(muxer *AMediaMuxer, format *AMediaFormat) int32 { return 0 }
func AMediaMuxer_start(muxer *AMediaMuxer) Media_status_t               { return 0 }
func AMediaMuxer_stop(muxer *AMediaMuxer) Media_status_t                { return 0 }

var _ = unsafe.Pointer(nil)
