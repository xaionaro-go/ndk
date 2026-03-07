// Simulates c-for-go output for Android Font (android/font.h).
// This file is parsed at AST level only; it does not compile.
package font

// Opaque handle types.
type AFont C.AFont
type AFontMatcher C.AFontMatcher

// Integer typedefs.
type Font_weight_t int32

// Font weight constants.
const (
	AFONT_WEIGHT_THIN   Font_weight_t = 100
	AFONT_WEIGHT_LIGHT  Font_weight_t = 300
	AFONT_WEIGHT_NORMAL Font_weight_t = 400
	AFONT_WEIGHT_MEDIUM Font_weight_t = 500
	AFONT_WEIGHT_BOLD   Font_weight_t = 700
	AFONT_WEIGHT_BLACK  Font_weight_t = 900
)

// --- Font functions ---
func AFont_close(font *AFont)            {}
func AFont_getWeight(font *AFont) uint16 { return 0 }
func AFont_isItalic(font *AFont) bool    { return false }

// --- FontMatcher functions ---
func AFontMatcher_create() *AFontMatcher                                      { return nil }
func AFontMatcher_destroy(matcher *AFontMatcher)                              {}
func AFontMatcher_setStyle(matcher *AFontMatcher, weight uint16, italic bool) {}
func AFontMatcher_match(matcher *AFontMatcher, familyName *byte, text *uint16, length uint32, runLengthOut *uint32) *AFont {
	return nil
}
