// Copyright (c) 2012 The gocql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package marshal

import (
	"errors"
	"fmt"
	"math/big"
)

var (
	bigOne = big.NewInt(1)
	// ErrorUDTUnavailable is returned when easycql is used on protocol older than version 3.
	ErrorUDTUnavailable = errors.New("UDT are not available on protocols less than 3, please update config")
)

func DecInt(x []byte) int32 {
	if len(x) != 4 {
		return 0
	}
	return int32(x[0])<<24 | int32(x[1])<<16 | int32(x[2])<<8 | int32(x[3])
}

func DecShort(p []byte) int16 {
	if len(p) != 2 {
		return 0
	}
	return int16(p[0])<<8 | int16(p[1])
}

func DecTiny(p []byte) int8 {
	if len(p) != 1 {
		return 0
	}
	return int8(p[0])
}

func bytesToInt64(data []byte) (ret int64) {
	for i := range data {
		ret |= int64(data[i]) << (8 * uint(len(data)-i-1))
	}
	return ret
}

// VarIntToInt64 decodes a varint and returns a flag indicating whether the value fits into int64.
func VarIntToInt64(data []byte) (int64, bool) {
	if len(data) > 8 {
		return 0, false
	}

	int64Val := bytesToInt64(data)
	if len(data) > 0 && len(data) < 8 && data[0]&0x80 > 0 {
		int64Val -= 1 << uint(len(data)*8)
	}
	return int64Val, true
}

func DecBigInt(data []byte) int64 {
	if len(data) != 8 {
		return 0
	}
	return int64(data[0])<<56 | int64(data[1])<<48 |
		int64(data[2])<<40 | int64(data[3])<<32 |
		int64(data[4])<<24 | int64(data[5])<<16 |
		int64(data[6])<<8 | int64(data[7])
}

func DecBool(v []byte) bool {
	if len(v) == 0 {
		return false
	}
	return v[0] != 0
}

// DecBigInt2C sets the value of n to the big-endian two's complement
// value stored in the given data. If data[0]&80 != 0, the number
// is negative. If data is empty, the result will be 0.
func DecBigInt2C(data []byte, n *big.Int) {
	n.SetBytes(data)
	if len(data) > 0 && data[0]&0x80 > 0 {
		n.Sub(n, new(big.Int).Lsh(bigOne, uint(len(data))*8))
	}
}

// ReadBytes decodes bytes from p and returns rest of p.
// It swallows whole p if there is not enough bytes for the data.
// https://github.com/apache/cassandra/blob/698078fecf3914b2a5f9d2ea344868f677f5afb2/doc/native_protocol_v4.spec#L227-L228
// Deprecated: Use ReadBytes2 instead.
func ReadBytes(p []byte) (bytes, rest []byte) {
	bytes, rest, err := ReadBytes2(p)
	if err != nil {
		return p, nil
	}
	return bytes, rest
}

// ReadBytes2 decodes bytes from p and returns rest of p.
// It returns an error if there is not enough bytes to read the data.
// https://github.com/apache/cassandra/blob/698078fecf3914b2a5f9d2ea344868f677f5afb2/doc/native_protocol_v4.spec#L227-L228
func ReadBytes2(p []byte) (bytes, rest []byte, err error) {
	size := readInt(p)
	p = p[4:]
	if size < 0 {
		return nil, p, nil
	}
	if len(p) < int(size) {
		return nil, nil, fmt.Errorf("read bytes: expecting %d bytes, got %d", size, len(p))
	}
	return p[:size], p[size:], nil
}

func readInt(p []byte) int32 {
	return int32(p[0])<<24 | int32(p[1])<<16 | int32(p[2])<<8 | int32(p[3])
}
