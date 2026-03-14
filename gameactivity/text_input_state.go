//go:build ignore

package gameactivity

// TextInputState represents the state of a text input field.
type TextInputState struct {
	// Text is the current text content (UTF-8).
	Text string

	// SelectionStart is the start offset of the text selection.
	SelectionStart int32

	// SelectionEnd is the end offset of the text selection.
	SelectionEnd int32

	// ComposingStart is the start offset of the composing region (IME).
	ComposingStart int32

	// ComposingEnd is the end offset of the composing region (IME).
	ComposingEnd int32
}
