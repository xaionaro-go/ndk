// Copyright 2018-2024 The gooid Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package storage

import (
	"github.com/xaionaro-go/ndk/ndk"
)

type AssetManager = ndk.AssetManager
type AssetDir = ndk.AssetDir
type Asset = ndk.Asset

/* Available modes for opening assets */
const (
	ASSET_MODE_UNKNOWN   = ndk.ASSET_MODE_UNKNOWN
	ASSET_MODE_RANDOM    = ndk.ASSET_MODE_RANDOM
	ASSET_MODE_STREAMING = ndk.ASSET_MODE_STREAMING
	ASSET_MODE_BUFFER    = ndk.ASSET_MODE_BUFFER
)
