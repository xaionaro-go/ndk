//go:build ignore

package gameactivity

/*
#include "include/game-activity/GameActivity.h"
*/
import "C"

// Window flags for SetWindowFlags.
const (
	FlagAllowLockWhileScreenOn = C.GAMEACTIVITY_FLAG_ALLOW_LOCK_WHILE_SCREEN_ON
	FlagDimBehind              = C.GAMEACTIVITY_FLAG_DIM_BEHIND
	FlagNotFocusable           = C.GAMEACTIVITY_FLAG_NOT_FOCUSABLE
	FlagNotTouchable           = C.GAMEACTIVITY_FLAG_NOT_TOUCHABLE
	FlagFullscreen             = C.GAMEACTIVITY_FLAG_FULLSCREEN
	FlagLayoutInScreen         = C.GAMEACTIVITY_FLAG_LAYOUT_IN_SCREEN
	FlagKeepScreenOn           = C.GAMEACTIVITY_FLAG_KEEP_SCREEN_ON
	FlagSecure                 = C.GAMEACTIVITY_FLAG_SECURE
	FlagShowWhenLocked         = C.GAMEACTIVITY_FLAG_SHOW_WHEN_LOCKED
	FlagDismissKeyguard        = C.GAMEACTIVITY_FLAG_DISMISS_KEYGUARD
)

// ShowSoftInputFlags controls how the soft keyboard is shown.
type ShowSoftInputFlags uint32

const (
	ShowSoftInputImplicit ShowSoftInputFlags = C.GAMEACTIVITY_SHOW_SOFT_INPUT_IMPLICIT
	ShowSoftInputForced   ShowSoftInputFlags = C.GAMEACTIVITY_SHOW_SOFT_INPUT_FORCED
)

// HideSoftInputFlags controls how the soft keyboard is hidden.
type HideSoftInputFlags uint32

const (
	HideSoftInputImplicitOnly HideSoftInputFlags = C.GAMEACTIVITY_HIDE_SOFT_INPUT_IMPLICIT_ONLY
	HideSoftInputNotAlways    HideSoftInputFlags = C.GAMEACTIVITY_HIDE_SOFT_INPUT_NOT_ALWAYS
)

// InsetsType identifies the type of window insets to query.
type InsetsType int32
