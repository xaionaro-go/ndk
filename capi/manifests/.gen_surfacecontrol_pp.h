#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wc23-extensions"
#pragma clang diagnostic pop
int __system_property_get(const char* __name, char* __value);
int atoi(const char* __s) ;
static __inline__ int android_get_device_api_level() {
  char value[92] = { 0 };
  if (__system_property_get("ro.build.version.sdk", value) < 1) return -1;
  int api_level = atoi(value);
  return (api_level > 0) ? api_level : -1;
}

typedef long int ptrdiff_t;
typedef long unsigned int size_t;
typedef int wchar_t;
typedef struct {
  long long __clang_max_align_nonce1
                                                          ;
  long double __clang_max_align_nonce2
                                                            ;
} max_align_t;

typedef signed char __int8_t;
typedef unsigned char __uint8_t;
typedef short __int16_t;
typedef unsigned short __uint16_t;
typedef int __int32_t;
typedef unsigned int __uint32_t;
typedef long __int64_t;
typedef unsigned long __uint64_t;
typedef long __intptr_t;
typedef unsigned long __uintptr_t;
typedef __int8_t int8_t;
typedef __uint8_t uint8_t;
typedef __int16_t int16_t;
typedef __uint16_t uint16_t;
typedef __int32_t int32_t;
typedef __uint32_t uint32_t;
typedef __int64_t int64_t;
typedef __uint64_t uint64_t;
typedef __intptr_t intptr_t;
typedef __uintptr_t uintptr_t;
typedef int8_t int_least8_t;
typedef uint8_t uint_least8_t;
typedef int16_t int_least16_t;
typedef uint16_t uint_least16_t;
typedef int32_t int_least32_t;
typedef uint32_t uint_least32_t;
typedef int64_t int_least64_t;
typedef uint64_t uint_least64_t;
typedef int8_t int_fast8_t;
typedef uint8_t uint_fast8_t;
typedef int64_t int_fast64_t;
typedef uint64_t uint_fast64_t;
typedef int64_t int_fast16_t;
typedef uint64_t uint_fast16_t;
typedef int64_t int_fast32_t;
typedef uint64_t uint_fast32_t;
typedef uint64_t uintmax_t;
typedef int64_t intmax_t;
typedef struct {
 intmax_t quot;
 intmax_t rem;
} imaxdiv_t;
intmax_t imaxabs(intmax_t __i) ;
imaxdiv_t imaxdiv(intmax_t __numerator, intmax_t __denominator) ;
intmax_t strtoimax(const char* __s, char* * __end_ptr, int __base);
uintmax_t strtoumax(const char* __s, char* * __end_ptr, int __base);
intmax_t wcstoimax(const wchar_t* __s, wchar_t* * __end_ptr, int __base);
uintmax_t wcstoumax(const wchar_t* __s, wchar_t* * __end_ptr, int __base);
enum ADataSpace {
    ADATASPACE_UNKNOWN = 0,
    ADATASPACE_STANDARD_MASK = 63 << 16,
    ADATASPACE_STANDARD_UNSPECIFIED = 0 << 16,
    ADATASPACE_STANDARD_BT709 = 1 << 16,
    ADATASPACE_STANDARD_BT601_625 = 2 << 16,
    ADATASPACE_STANDARD_BT601_625_UNADJUSTED = 3 << 16,
    ADATASPACE_STANDARD_BT601_525 = 4 << 16,
    ADATASPACE_STANDARD_BT601_525_UNADJUSTED = 5 << 16,
    ADATASPACE_STANDARD_BT2020 = 6 << 16,
    ADATASPACE_STANDARD_BT2020_CONSTANT_LUMINANCE = 7 << 16,
    ADATASPACE_STANDARD_BT470M = 8 << 16,
    ADATASPACE_STANDARD_FILM = 9 << 16,
    ADATASPACE_STANDARD_DCI_P3 = 10 << 16,
    ADATASPACE_STANDARD_ADOBE_RGB = 11 << 16,
    ADATASPACE_TRANSFER_MASK = 31 << 22,
    ADATASPACE_TRANSFER_UNSPECIFIED = 0 << 22,
    ADATASPACE_TRANSFER_LINEAR = 1 << 22,
    ADATASPACE_TRANSFER_SRGB = 2 << 22,
    ADATASPACE_TRANSFER_SMPTE_170M = 3 << 22,
    ADATASPACE_TRANSFER_GAMMA2_2 = 4 << 22,
    ADATASPACE_TRANSFER_GAMMA2_6 = 5 << 22,
    ADATASPACE_TRANSFER_GAMMA2_8 = 6 << 22,
    ADATASPACE_TRANSFER_ST2084 = 7 << 22,
    ADATASPACE_TRANSFER_HLG = 8 << 22,
    ADATASPACE_RANGE_MASK = 7 << 27,
    ADATASPACE_RANGE_UNSPECIFIED = 0 << 27,
    ADATASPACE_RANGE_FULL = 1 << 27,
    ADATASPACE_RANGE_LIMITED = 2 << 27,
    ADATASPACE_RANGE_EXTENDED = 3 << 27,
    ADATASPACE_SCRGB_LINEAR = 406913024,
    ADATASPACE_SRGB = 142671872,
    ADATASPACE_SCRGB = 411107328,
    ADATASPACE_DISPLAY_P3 = 143261696,
    ADATASPACE_BT2020_PQ = 163971072,
    ADATASPACE_BT2020_ITU_PQ = 298188800,
    ADATASPACE_ADOBE_RGB = 151715840,
    ADATASPACE_JFIF = 146931712,
    ADATASPACE_BT601_625 = 281149440,
    ADATASPACE_BT601_525 = 281280512,
    ADATASPACE_BT2020 = 147193856,
    ADATASPACE_BT709 = 281083904,
    ADATASPACE_DCI_P3 = 155844608,
    ADATASPACE_SRGB_LINEAR = 138477568,
    ADATASPACE_BT2020_HLG = 168165376,
    ADATASPACE_BT2020_ITU_HLG = 302383104,
    ADATASPACE_DEPTH = 4096,
    ADATASPACE_DYNAMIC_DEPTH = 4098,
    STANDARD_MASK = ADATASPACE_STANDARD_MASK,
    STANDARD_UNSPECIFIED = ADATASPACE_STANDARD_UNSPECIFIED,
    STANDARD_BT709 = ADATASPACE_STANDARD_BT709,
    STANDARD_BT601_625 = ADATASPACE_STANDARD_BT601_625,
    STANDARD_BT601_625_UNADJUSTED = ADATASPACE_STANDARD_BT601_625_UNADJUSTED,
    STANDARD_BT601_525 = ADATASPACE_STANDARD_BT601_525,
    STANDARD_BT601_525_UNADJUSTED = ADATASPACE_STANDARD_BT601_525_UNADJUSTED,
    STANDARD_BT470M = ADATASPACE_STANDARD_BT470M,
    STANDARD_BT2020 = ADATASPACE_STANDARD_BT2020,
    STANDARD_FILM = ADATASPACE_STANDARD_FILM,
    STANDARD_DCI_P3 = ADATASPACE_STANDARD_DCI_P3,
    STANDARD_ADOBE_RGB = ADATASPACE_STANDARD_ADOBE_RGB,
    TRANSFER_MASK = ADATASPACE_TRANSFER_MASK,
    TRANSFER_UNSPECIFIED = ADATASPACE_TRANSFER_UNSPECIFIED,
    TRANSFER_LINEAR = ADATASPACE_TRANSFER_LINEAR,
    TRANSFER_SMPTE_170M = ADATASPACE_TRANSFER_SMPTE_170M,
    TRANSFER_GAMMA2_2 = ADATASPACE_TRANSFER_GAMMA2_2,
    TRANSFER_GAMMA2_6 = ADATASPACE_TRANSFER_GAMMA2_6,
    TRANSFER_GAMMA2_8 = ADATASPACE_TRANSFER_GAMMA2_8,
    TRANSFER_SRGB = ADATASPACE_TRANSFER_SRGB,
    TRANSFER_ST2084 = ADATASPACE_TRANSFER_ST2084,
    TRANSFER_HLG = ADATASPACE_TRANSFER_HLG,
    RANGE_MASK = ADATASPACE_RANGE_MASK,
    RANGE_UNSPECIFIED = ADATASPACE_RANGE_UNSPECIFIED,
    RANGE_FULL = ADATASPACE_RANGE_FULL,
    RANGE_LIMITED = ADATASPACE_RANGE_LIMITED,
    RANGE_EXTENDED = ADATASPACE_RANGE_EXTENDED,
};
struct AChoreographer;
typedef struct AChoreographer AChoreographer;
typedef int64_t AVsyncId;
struct AChoreographerFrameCallbackData;
typedef struct AChoreographerFrameCallbackData AChoreographerFrameCallbackData;
typedef void (*AChoreographer_frameCallback)(long frameTimeNanos, void* data);
typedef void (*AChoreographer_frameCallback64)(int64_t frameTimeNanos, void* data);
typedef void (*AChoreographer_vsyncCallback)(
        const AChoreographerFrameCallbackData* callbackData, void* data);
typedef void (*AChoreographer_refreshRateCallback)(int64_t vsyncPeriodNanos, void* data);
AChoreographer* AChoreographer_getInstance() ;
void AChoreographer_postFrameCallback(AChoreographer* choreographer,
                                      AChoreographer_frameCallback callback, void* data)
                                                                                                 ;
void AChoreographer_postFrameCallbackDelayed(AChoreographer* choreographer,
                                             AChoreographer_frameCallback callback, void* data,
                                             long delayMillis)
                                                                                    ;
void AChoreographer_postFrameCallback64(AChoreographer* choreographer,
                                        AChoreographer_frameCallback64 callback, void* data)
                           ;
void AChoreographer_postFrameCallbackDelayed64(AChoreographer* choreographer,
                                               AChoreographer_frameCallback64 callback, void* data,
                                               uint32_t delayMillis) ;
void AChoreographer_postVsyncCallback(AChoreographer* choreographer,
                                        AChoreographer_vsyncCallback callback, void* data)
                           ;
void AChoreographer_registerRefreshRateCallback(AChoreographer* choreographer,
                                                AChoreographer_refreshRateCallback, void* data)
                           ;
void AChoreographer_unregisterRefreshRateCallback(AChoreographer* choreographer,
                                                  AChoreographer_refreshRateCallback, void* data)
                           ;
int64_t AChoreographerFrameCallbackData_getFrameTimeNanos(
        const AChoreographerFrameCallbackData* data) ;
size_t AChoreographerFrameCallbackData_getFrameTimelinesLength(
        const AChoreographerFrameCallbackData* data) ;
size_t AChoreographerFrameCallbackData_getPreferredFrameTimelineIndex(
        const AChoreographerFrameCallbackData* data) ;
AVsyncId AChoreographerFrameCallbackData_getFrameTimelineVsyncId(
        const AChoreographerFrameCallbackData* data, size_t index) ;
int64_t AChoreographerFrameCallbackData_getFrameTimelineExpectedPresentationTimeNanos(
        const AChoreographerFrameCallbackData* data, size_t index) ;
int64_t AChoreographerFrameCallbackData_getFrameTimelineDeadlineNanos(
        const AChoreographerFrameCallbackData* data, size_t index) ;

typedef struct ARect {
    int32_t left;
    int32_t top;
    int32_t right;
    int32_t bottom;
} ARect;
enum AHardwareBuffer_Format {
    AHARDWAREBUFFER_FORMAT_R8G8B8A8_UNORM = 1,
    AHARDWAREBUFFER_FORMAT_R8G8B8X8_UNORM = 2,
    AHARDWAREBUFFER_FORMAT_R8G8B8_UNORM = 3,
    AHARDWAREBUFFER_FORMAT_R5G6B5_UNORM = 4,
    AHARDWAREBUFFER_FORMAT_R16G16B16A16_FLOAT = 0x16,
    AHARDWAREBUFFER_FORMAT_R10G10B10A2_UNORM = 0x2b,
    AHARDWAREBUFFER_FORMAT_BLOB = 0x21,
    AHARDWAREBUFFER_FORMAT_D16_UNORM = 0x30,
    AHARDWAREBUFFER_FORMAT_D24_UNORM = 0x31,
    AHARDWAREBUFFER_FORMAT_D24_UNORM_S8_UINT = 0x32,
    AHARDWAREBUFFER_FORMAT_D32_FLOAT = 0x33,
    AHARDWAREBUFFER_FORMAT_D32_FLOAT_S8_UINT = 0x34,
    AHARDWAREBUFFER_FORMAT_S8_UINT = 0x35,
    AHARDWAREBUFFER_FORMAT_Y8Cb8Cr8_420 = 0x23,
    AHARDWAREBUFFER_FORMAT_YCbCr_P010 = 0x36,
    AHARDWAREBUFFER_FORMAT_YCbCr_P210 = 0x3c,
    AHARDWAREBUFFER_FORMAT_R8_UNORM = 0x38,
    AHARDWAREBUFFER_FORMAT_R16_UINT = 0x39,
    AHARDWAREBUFFER_FORMAT_R16G16_UINT = 0x3a,
    AHARDWAREBUFFER_FORMAT_R10G10B10A10_UNORM = 0x3b,
};
enum AHardwareBuffer_UsageFlags {
    AHARDWAREBUFFER_USAGE_CPU_READ_NEVER = 0UL,
    AHARDWAREBUFFER_USAGE_CPU_READ_RARELY = 2UL,
    AHARDWAREBUFFER_USAGE_CPU_READ_OFTEN = 3UL,
    AHARDWAREBUFFER_USAGE_CPU_READ_MASK = 0xFUL,
    AHARDWAREBUFFER_USAGE_CPU_WRITE_NEVER = 0UL << 4,
    AHARDWAREBUFFER_USAGE_CPU_WRITE_RARELY = 2UL << 4,
    AHARDWAREBUFFER_USAGE_CPU_WRITE_OFTEN = 3UL << 4,
    AHARDWAREBUFFER_USAGE_CPU_WRITE_MASK = 0xFUL << 4,
    AHARDWAREBUFFER_USAGE_GPU_SAMPLED_IMAGE = 1UL << 8,
    AHARDWAREBUFFER_USAGE_GPU_FRAMEBUFFER = 1UL << 9,
    AHARDWAREBUFFER_USAGE_GPU_COLOR_OUTPUT = AHARDWAREBUFFER_USAGE_GPU_FRAMEBUFFER,
    AHARDWAREBUFFER_USAGE_COMPOSER_OVERLAY = 1ULL << 11,
    AHARDWAREBUFFER_USAGE_PROTECTED_CONTENT = 1UL << 14,
    AHARDWAREBUFFER_USAGE_VIDEO_ENCODE = 1UL << 16,
    AHARDWAREBUFFER_USAGE_SENSOR_DIRECT_DATA = 1UL << 23,
    AHARDWAREBUFFER_USAGE_GPU_DATA_BUFFER = 1UL << 24,
    AHARDWAREBUFFER_USAGE_GPU_CUBE_MAP = 1UL << 25,
    AHARDWAREBUFFER_USAGE_GPU_MIPMAP_COMPLETE = 1UL << 26,
    AHARDWAREBUFFER_USAGE_FRONT_BUFFER = 1ULL << 32,
    AHARDWAREBUFFER_USAGE_VENDOR_0 = 1ULL << 28,
    AHARDWAREBUFFER_USAGE_VENDOR_1 = 1ULL << 29,
    AHARDWAREBUFFER_USAGE_VENDOR_2 = 1ULL << 30,
    AHARDWAREBUFFER_USAGE_VENDOR_3 = 1ULL << 31,
    AHARDWAREBUFFER_USAGE_VENDOR_4 = 1ULL << 48,
    AHARDWAREBUFFER_USAGE_VENDOR_5 = 1ULL << 49,
    AHARDWAREBUFFER_USAGE_VENDOR_6 = 1ULL << 50,
    AHARDWAREBUFFER_USAGE_VENDOR_7 = 1ULL << 51,
    AHARDWAREBUFFER_USAGE_VENDOR_8 = 1ULL << 52,
    AHARDWAREBUFFER_USAGE_VENDOR_9 = 1ULL << 53,
    AHARDWAREBUFFER_USAGE_VENDOR_10 = 1ULL << 54,
    AHARDWAREBUFFER_USAGE_VENDOR_11 = 1ULL << 55,
    AHARDWAREBUFFER_USAGE_VENDOR_12 = 1ULL << 56,
    AHARDWAREBUFFER_USAGE_VENDOR_13 = 1ULL << 57,
    AHARDWAREBUFFER_USAGE_VENDOR_14 = 1ULL << 58,
    AHARDWAREBUFFER_USAGE_VENDOR_15 = 1ULL << 59,
    AHARDWAREBUFFER_USAGE_VENDOR_16 = 1ULL << 60,
    AHARDWAREBUFFER_USAGE_VENDOR_17 = 1ULL << 61,
    AHARDWAREBUFFER_USAGE_VENDOR_18 = 1ULL << 62,
    AHARDWAREBUFFER_USAGE_VENDOR_19 = 1ULL << 63,
};
typedef struct AHardwareBuffer_Desc {
    uint32_t width;
    uint32_t height;
    uint32_t layers;
    uint32_t format;
    uint64_t usage;
    uint32_t stride;
    uint32_t rfu0;
    uint64_t rfu1;
} AHardwareBuffer_Desc;
typedef struct AHardwareBuffer_Plane {
    void* data;
    uint32_t pixelStride;
    uint32_t rowStride;
} AHardwareBuffer_Plane;
typedef struct AHardwareBuffer_Planes {
    uint32_t planeCount;
    AHardwareBuffer_Plane planes[4];
} AHardwareBuffer_Planes;
typedef struct AHardwareBuffer AHardwareBuffer;
int AHardwareBuffer_allocate(const AHardwareBuffer_Desc* desc,
                             AHardwareBuffer* * outBuffer) ;
void AHardwareBuffer_acquire(AHardwareBuffer* buffer) ;
void AHardwareBuffer_release(AHardwareBuffer* buffer) ;
void AHardwareBuffer_describe(const AHardwareBuffer* buffer,
                              AHardwareBuffer_Desc* outDesc) ;
int AHardwareBuffer_lock(AHardwareBuffer* buffer, uint64_t usage, int32_t fence,
                         const ARect* rect, void* * outVirtualAddress)
                           ;
int AHardwareBuffer_unlock(AHardwareBuffer* buffer, int32_t* fence)
                           ;
int AHardwareBuffer_sendHandleToUnixSocket(const AHardwareBuffer* buffer, int socketFd)
                           ;
int AHardwareBuffer_recvHandleFromUnixSocket(int socketFd,
                                             AHardwareBuffer* * outBuffer)
                           ;
int AHardwareBuffer_lockPlanes(AHardwareBuffer* buffer, uint64_t usage, int32_t fence,
                               const ARect* rect,
                               AHardwareBuffer_Planes* outPlanes) ;
int AHardwareBuffer_isSupported(const AHardwareBuffer_Desc* desc) ;
int AHardwareBuffer_lockAndGetInfo(AHardwareBuffer* buffer, uint64_t usage, int32_t fence,
                                   const ARect* rect,
                                   void* * outVirtualAddress,
                                   int32_t* outBytesPerPixel,
                                   int32_t* outBytesPerStride) ;
int AHardwareBuffer_getId(const AHardwareBuffer* buffer, uint64_t* outId)
                           ;
enum AHdrMetadataType {
    HDR10_SMPTE2086 = 1,
    HDR10_CTA861_3 = 2,
    HDR10PLUS_SEI = 3,
};
struct AColor_xy {
    float x;
    float y;
};
struct AHdrMetadata_smpte2086 {
    struct AColor_xy displayPrimaryRed;
    struct AColor_xy displayPrimaryGreen;
    struct AColor_xy displayPrimaryBlue;
    struct AColor_xy whitePoint;
    float maxLuminance;
    float minLuminance;
};
struct AHdrMetadata_cta861_3 {
    float maxContentLightLevel;
    float maxFrameAverageLightLevel;
};
enum ANativeWindow_LegacyFormat {
    WINDOW_FORMAT_RGBA_8888 = AHARDWAREBUFFER_FORMAT_R8G8B8A8_UNORM,
    WINDOW_FORMAT_RGBX_8888 = AHARDWAREBUFFER_FORMAT_R8G8B8X8_UNORM,
    WINDOW_FORMAT_RGB_565 = AHARDWAREBUFFER_FORMAT_R5G6B5_UNORM,
};
enum ANativeWindowTransform {
    ANATIVEWINDOW_TRANSFORM_IDENTITY = 0x00,
    ANATIVEWINDOW_TRANSFORM_MIRROR_HORIZONTAL = 0x01,
    ANATIVEWINDOW_TRANSFORM_MIRROR_VERTICAL = 0x02,
    ANATIVEWINDOW_TRANSFORM_ROTATE_90 = 0x04,
    ANATIVEWINDOW_TRANSFORM_ROTATE_180 = ANATIVEWINDOW_TRANSFORM_MIRROR_HORIZONTAL |
                                                  ANATIVEWINDOW_TRANSFORM_MIRROR_VERTICAL,
    ANATIVEWINDOW_TRANSFORM_ROTATE_270 = ANATIVEWINDOW_TRANSFORM_ROTATE_180 |
                                                  ANATIVEWINDOW_TRANSFORM_ROTATE_90,
};
struct ANativeWindow;
typedef struct ANativeWindow ANativeWindow;
typedef struct ANativeWindow_Buffer {
    int32_t width;
    int32_t height;
    int32_t stride;
    int32_t format;
    void* bits;
    uint32_t reserved[6];
} ANativeWindow_Buffer;
void ANativeWindow_acquire(ANativeWindow* window);
void ANativeWindow_release(ANativeWindow* window);
int32_t ANativeWindow_getWidth(ANativeWindow* window);
int32_t ANativeWindow_getHeight(ANativeWindow* window);
int32_t ANativeWindow_getFormat(ANativeWindow* window);
int32_t ANativeWindow_setBuffersGeometry(ANativeWindow* window,
        int32_t width, int32_t height, int32_t format);
int32_t ANativeWindow_lock(ANativeWindow* window, ANativeWindow_Buffer* outBuffer,
        ARect* inOutDirtyBounds);
int32_t ANativeWindow_unlockAndPost(ANativeWindow* window);
int32_t ANativeWindow_setBuffersTransform(ANativeWindow* window, int32_t transform) ;
int32_t ANativeWindow_setBuffersDataSpace(ANativeWindow* window, int32_t dataSpace) ;
int32_t ANativeWindow_getBuffersDataSpace(ANativeWindow* window) ;
int32_t ANativeWindow_getBuffersDefaultDataSpace(ANativeWindow* window) ;
enum ANativeWindow_FrameRateCompatibility {
    ANATIVEWINDOW_FRAME_RATE_COMPATIBILITY_DEFAULT = 0,
    ANATIVEWINDOW_FRAME_RATE_COMPATIBILITY_FIXED_SOURCE = 1
};
int32_t ANativeWindow_setFrameRate(ANativeWindow* window, float frameRate, int8_t compatibility)
                           ;
void ANativeWindow_tryAllocateBuffers(ANativeWindow* window) ;
enum ANativeWindow_ChangeFrameRateStrategy {
    ANATIVEWINDOW_CHANGE_FRAME_RATE_ONLY_IF_SEAMLESS = 0,
    ANATIVEWINDOW_CHANGE_FRAME_RATE_ALWAYS = 1
} ;
int32_t ANativeWindow_setFrameRateWithChangeStrategy(ANativeWindow* window, float frameRate,
        int8_t compatibility, int8_t changeFrameRateStrategy)
                           ;
inline int32_t ANativeWindow_clearFrameRate(ANativeWindow* window) {
    return ANativeWindow_setFrameRateWithChangeStrategy(window, 0,
            ANATIVEWINDOW_FRAME_RATE_COMPATIBILITY_DEFAULT,
            ANATIVEWINDOW_CHANGE_FRAME_RATE_ONLY_IF_SEAMLESS);
}
struct ASurfaceControl;
typedef struct ASurfaceControl ASurfaceControl;
ASurfaceControl* ASurfaceControl_createFromWindow(ANativeWindow* parent,
                                                            const char* debug_name)
                           ;
ASurfaceControl* ASurfaceControl_create(ASurfaceControl* parent,
                                                  const char* debug_name)
                           ;
void ASurfaceControl_acquire(ASurfaceControl* surface_control) ;
void ASurfaceControl_release(ASurfaceControl* surface_control) ;
struct ASurfaceTransaction;
typedef struct ASurfaceTransaction ASurfaceTransaction;
ASurfaceTransaction* ASurfaceTransaction_create() ;
void ASurfaceTransaction_delete(ASurfaceTransaction* transaction) ;
void ASurfaceTransaction_apply(ASurfaceTransaction* transaction) ;
typedef struct ASurfaceTransactionStats ASurfaceTransactionStats;
typedef void (*ASurfaceTransaction_OnComplete)(void* context,
                                               ASurfaceTransactionStats* stats);
typedef void (*ASurfaceTransaction_OnCommit)(void* context,
                                             ASurfaceTransactionStats* stats);
typedef void (*ASurfaceTransaction_OnBufferRelease)(void* context,
                                                    int release_fence_fd);
int64_t ASurfaceTransactionStats_getLatchTime(
        ASurfaceTransactionStats* surface_transaction_stats) ;
int ASurfaceTransactionStats_getPresentFenceFd(
        ASurfaceTransactionStats* surface_transaction_stats) ;
void ASurfaceTransactionStats_getASurfaceControls(
        ASurfaceTransactionStats* surface_transaction_stats,
        ASurfaceControl* * * outASurfaceControls,
        size_t* outASurfaceControlsSize) ;
void ASurfaceTransactionStats_releaseASurfaceControls(
        ASurfaceControl* * surface_controls) ;
int64_t ASurfaceTransactionStats_getAcquireTime(
        ASurfaceTransactionStats* surface_transaction_stats,
        ASurfaceControl* surface_control) ;
int ASurfaceTransactionStats_getPreviousReleaseFenceFd(
        ASurfaceTransactionStats* surface_transaction_stats,
        ASurfaceControl* surface_control) ;
void ASurfaceTransaction_setOnComplete(ASurfaceTransaction* transaction,
                                       void* context,
                                       ASurfaceTransaction_OnComplete func)
                           ;
void ASurfaceTransaction_setOnCommit(ASurfaceTransaction* transaction,
                                     void* context,
                                     ASurfaceTransaction_OnCommit func)
                           ;
void ASurfaceTransaction_reparent(ASurfaceTransaction* transaction,
                                  ASurfaceControl* surface_control,
                                  ASurfaceControl* new_parent) ;
enum ASurfaceTransactionVisibility {
    ASURFACE_TRANSACTION_VISIBILITY_HIDE = 0,
    ASURFACE_TRANSACTION_VISIBILITY_SHOW = 1,
};
void ASurfaceTransaction_setVisibility(ASurfaceTransaction* transaction,
                                       ASurfaceControl* surface_control,
                                       enum ASurfaceTransactionVisibility visibility)
                           ;
void ASurfaceTransaction_setZOrder(ASurfaceTransaction* transaction,
                                   ASurfaceControl* surface_control, int32_t z_order)
                           ;
void ASurfaceTransaction_setBuffer(ASurfaceTransaction* transaction,
                                   ASurfaceControl* surface_control,
                                   AHardwareBuffer* buffer, int acquire_fence_fd)
                           ;
void ASurfaceTransaction_setBufferWithRelease(ASurfaceTransaction* transaction,
                                              ASurfaceControl* surface_control,
                                              AHardwareBuffer* buffer,
                                              int acquire_fence_fd, void* context,
                                              ASurfaceTransaction_OnBufferRelease func)
                           ;
void ASurfaceTransaction_setColor(ASurfaceTransaction* transaction,
                                  ASurfaceControl* surface_control, float r, float g,
                                  float b, float alpha, enum ADataSpace dataspace)
                           ;
void ASurfaceTransaction_setGeometry(ASurfaceTransaction* transaction,
                                     ASurfaceControl* surface_control,
                                     const ARect* source,
                                     const ARect* destination,
                                     int32_t transform) ;
void ASurfaceTransaction_setCrop(ASurfaceTransaction* transaction,
                                 ASurfaceControl* surface_control,
                                 const ARect* crop)
                           ;
void ASurfaceTransaction_setPosition(ASurfaceTransaction* transaction,
                                     ASurfaceControl* surface_control, int32_t x,
                                     int32_t y) ;
void ASurfaceTransaction_setBufferTransform(ASurfaceTransaction* transaction,
                                            ASurfaceControl* surface_control,
                                            int32_t transform) ;
void ASurfaceTransaction_setScale(ASurfaceTransaction* transaction,
                                  ASurfaceControl* surface_control, float xScale,
                                  float yScale) ;
enum ASurfaceTransactionTransparency {
    ASURFACE_TRANSACTION_TRANSPARENCY_TRANSPARENT = 0,
    ASURFACE_TRANSACTION_TRANSPARENCY_TRANSLUCENT = 1,
    ASURFACE_TRANSACTION_TRANSPARENCY_OPAQUE = 2,
};
void ASurfaceTransaction_setBufferTransparency(ASurfaceTransaction* transaction,
                                               ASurfaceControl* surface_control,
                                               enum ASurfaceTransactionTransparency transparency)
                                                                  ;
void ASurfaceTransaction_setDamageRegion(ASurfaceTransaction* transaction,
                                         ASurfaceControl* surface_control,
                                         const ARect* rects, uint32_t count)
                                                            ;
void ASurfaceTransaction_setDesiredPresentTime(ASurfaceTransaction* transaction,
                                               int64_t desiredPresentTime) ;
void ASurfaceTransaction_setBufferAlpha(ASurfaceTransaction* transaction,
                                        ASurfaceControl* surface_control, float alpha)
                                                           ;
void ASurfaceTransaction_setBufferDataSpace(ASurfaceTransaction* transaction,
                                            ASurfaceControl* surface_control,
                                            enum ADataSpace data_space) ;
void ASurfaceTransaction_setHdrMetadata_smpte2086(ASurfaceTransaction* transaction,
                                                  ASurfaceControl* surface_control,
                                                  struct AHdrMetadata_smpte2086* metadata)
                                                                     ;
void ASurfaceTransaction_setHdrMetadata_cta861_3(ASurfaceTransaction* transaction,
                                                 ASurfaceControl* surface_control,
                                                 struct AHdrMetadata_cta861_3* metadata)
                                                                    ;
void ASurfaceTransaction_setExtendedRangeBrightness(ASurfaceTransaction* transaction,
                                                    ASurfaceControl* surface_control,
                                                    float currentBufferRatio, float desiredRatio)
                                                                                      ;
void ASurfaceTransaction_setDesiredHdrHeadroom(ASurfaceTransaction* transaction,
                                               ASurfaceControl* surface_control,
                                               float desiredHeadroom)
                                          ;
void ASurfaceTransaction_setFrameRate(ASurfaceTransaction* transaction,
                                      ASurfaceControl* surface_control, float frameRate,
                                      int8_t compatibility) ;
void ASurfaceTransaction_setFrameRateWithChangeStrategy(ASurfaceTransaction* transaction,
                                                        ASurfaceControl* surface_control,
                                                        float frameRate, int8_t compatibility,
                                                        int8_t changeFrameRateStrategy)
                                                                           ;
void ASurfaceTransaction_clearFrameRate(ASurfaceTransaction* transaction,
                                        ASurfaceControl* surface_control)
                                                                          ;
void ASurfaceTransaction_setEnableBackPressure(ASurfaceTransaction* transaction,
                                               ASurfaceControl* surface_control,
                                               _Bool enableBackPressure) ;
void ASurfaceTransaction_setFrameTimeline(ASurfaceTransaction* transaction,
                                          AVsyncId vsyncId) ;

