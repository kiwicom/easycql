// Copyright (c) 2012 The gocql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package marshal

func AppendBytes(p, d []byte) []byte {
	if d == nil {
		return appendInt(p, -1)
	}
	p = appendInt(p, int32(len(d)))
	p = append(p, d...)
	return p
}

func appendInt(p []byte, n int32) []byte {
	return append(p, byte(n>>24),
		byte(n>>16),
		byte(n>>8),
		byte(n))
}
