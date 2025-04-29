// Package crypto provides encryption and decryption utilities for OpenKore.
// This file implements the Crypton encryption algorithm.
package crypto

import (
	"encoding/binary"
)

// Crypton represents a Crypton encryption instance with the necessary state.
type Crypton struct {
	eKey []uint32
}

// Global S-box and S-tab tables
var (
	sBox [2][256]byte
	sTab [4][256]uint32
	// tableGenerated tracks if the tables have been generated
	tableGenerated bool
)

// New creates a new Crypton instance with the given key and key length.
func New(key []byte, keyLen int) *Crypton {
	// Generate tables if not already done
	if !tableGenerated {
		generateTable()
		tableGenerated = true
	}

	c := &Crypton{
		eKey: make([]uint32, 52), // 52 is the size needed for the encryption key
	}

	// Set up the key
	setKey(key, keyLen, c.eKey)
	return c
}

// Encrypt encrypts a 16-byte block using the Crypton algorithm.
func (c *Crypton) Encrypt(block []byte) []byte {
	if len(block) != 16 {
		panic("Crypton.Encrypt: block must be 16 bytes")
	}

	// Unpack the block into 4 uint32 values
	var b0, b1 [4]uint32
	for i := 0; i < 4; i++ {
		b1[i] = binary.LittleEndian.Uint32(block[i*4 : (i+1)*4])
	}

	// Initial key addition
	b0[0] = b1[0] ^ c.eKey[0]
	b0[1] = b1[1] ^ c.eKey[1]
	b0[2] = b1[2] ^ c.eKey[2]
	b0[3] = b1[3] ^ c.eKey[3]

	// Rounds
	f0Rnd(c.eKey, 4, &b0, &b1)
	f1Rnd(c.eKey, 8, &b0, &b1)
	f0Rnd(c.eKey, 12, &b0, &b1)
	f1Rnd(c.eKey, 16, &b0, &b1)
	f0Rnd(c.eKey, 20, &b0, &b1)
	f1Rnd(c.eKey, 24, &b0, &b1)
	f0Rnd(c.eKey, 28, &b0, &b1)
	f1Rnd(c.eKey, 32, &b0, &b1)
	f0Rnd(c.eKey, 36, &b0, &b1)
	f1Rnd(c.eKey, 40, &b0, &b1)
	f0Rnd(c.eKey, 44, &b0, &b1)

	// Final transformation
	b0[0] = gammaTau(&b1, 0, 1, 0) ^ c.eKey[48]
	b0[1] = gammaTau(&b1, 1, 0, 1) ^ c.eKey[49]
	b0[2] = gammaTau(&b1, 2, 1, 0) ^ c.eKey[50]
	b0[3] = gammaTau(&b1, 3, 0, 1) ^ c.eKey[51]

	// Pack the result back into a byte array
	result := make([]byte, 16)
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(result[i*4:(i+1)*4], b0[i])
	}

	return result
}

// byte extracts a specific byte from a uint32 value.
func byteAt(data uint32, bytePos int) byte {
	if bytePos > 0 && bytePos < 4 {
		data = data >> (8 * bytePos)
	}
	return byte(data & 0xff)
}

// gammaTau applies the gamma and tau transformations.
func gammaTau(rB *[4]uint32, m int, p, q int) uint32 {
	return uint32(sBox[p][byteAt(rB[0], m)]) |
		(uint32(sBox[q][byteAt(rB[1], m)]) << 8) |
		(uint32(sBox[p][byteAt(rB[2], m)]) << 16) |
		(uint32(sBox[q][byteAt(rB[3], m)]) << 24)
}

// rotl rotates a uint32 value left by the specified number of bits.
func rotl(x uint32, bit int) uint32 {
	if bit > 0 && bit < 32 {
		return (x << bit) | (x >> (32 - bit))
	}
	return x
}

// pi applies the pi transformation to a slice of uint32 values.
func pi(rB []uint32, posB, n0, n1, n2, n3 int) uint32 {
	ma := [4]uint32{0x3fcff3fc, 0xfc3fcff3, 0xf3fc3fcf, 0xcff3fc3f}
	return (rB[posB+0] & ma[n0]) ^
		(rB[posB+1] & ma[n1]) ^
		(rB[posB+2] & ma[n2]) ^
		(rB[posB+3] & ma[n3])
}

// phiN applies the phi_n transformation.
func phiN(x uint32, n0, n1, n2, n3 int) uint32 {
	mb := [4]uint32{0xcffccffc, 0xf33ff33f, 0xfccffccf, 0x3ff33ff3}
	return (x & mb[n0]) ^
		(rotl(x, 8) & mb[n1]) ^
		(rotl(x, 16) & mb[n2]) ^
		(rotl(x, 24) & mb[n3])
}

// generateTable generates the S-box and S-tab tables.
func generateTable() {
	pBox := [3][16]byte{
		{15, 9, 6, 8, 9, 9, 4, 12, 6, 2, 6, 10, 1, 3, 5, 15},
		{10, 15, 4, 7, 5, 2, 14, 6, 9, 3, 12, 8, 13, 1, 11, 0},
		{0, 4, 8, 4, 2, 15, 8, 13, 1, 1, 15, 7, 2, 11, 14, 15},
	}

	for i := 0; i < 256; i++ {
		xl := byte((i >> 4) & 0x0f)
		xr := byte(i & 0x0f)

		yr := xr ^ pBox[1][xl^pBox[0][xr]]
		yl := xl ^ pBox[0][xr] ^ pBox[2][yr]

		yr |= (yl << 4)

		sBox[0][i] = yr
		sBox[1][yr] = byte(i)

		xr32 := uint32(yr) * 0x01010101
		xl32 := uint32(i) * 0x01010101

		sTab[0][i] = xr32 & 0x3fcff3fc
		sTab[1][yr] = xl32 & 0xfc3fcff3
		sTab[2][i] = xr32 & 0xf3fc3fcf
		sTab[3][yr] = xl32 & 0xcff3fc3f
	}
}

// h0Block applies the h0 transformation to the key.
func h0Block(rEKey []uint32, n, r0, r1 int, rc uint32) {
	rEKey[4*n+8] = rotl(rEKey[4*n+0], r0)
	rEKey[4*n+9] = rc ^ rEKey[4*n+1]
	rEKey[4*n+10] = rotl(rEKey[4*n+2], r1)
	rEKey[4*n+11] = rc ^ rEKey[4*n+3]
}

// h1Block applies the h1 transformation to the key.
func h1Block(rEKey []uint32, n, r0, r1 int, rc uint32) {
	rEKey[4*n+8] = rc ^ rEKey[4*n+0]
	rEKey[4*n+9] = rotl(rEKey[4*n+1], r0)
	rEKey[4*n+10] = rc ^ rEKey[4*n+2]
	rEKey[4*n+11] = rotl(rEKey[4*n+3], r1)
}

// setKey sets up the encryption key.
func setKey(inKey []byte, keyLen int, rEKey []uint32) {
	kp := [4]uint32{0xbb67ae85, 0x3c6ef372, 0xa54ff53a, 0x510e527f}
	kq := [4]uint32{0x9b05688c, 0x1f83d9ab, 0x5be0cd19, 0xcbbb9d5d}

	// Convert byte array to uint32 array
	key := make([]uint32, 8)
	for i := 0; i < len(inKey)/4 && i < 8; i++ {
		key[i] = binary.LittleEndian.Uint32(inKey[i*4:])
	}

	// Initialize key
	rEKey[2] = 0
	rEKey[3] = 0
	rEKey[6] = 0
	rEKey[7] = 0

	i := (keyLen + 63) / 64
	if i == 4 {
		rEKey[3] = key[6]
		rEKey[7] = key[7]
		rEKey[2] = key[4]
		rEKey[6] = key[5]
		rEKey[0] = key[0]
		rEKey[4] = key[1]
		rEKey[1] = key[2]
		rEKey[5] = key[3]
	} else if i == 3 {
		rEKey[2] = key[4]
		rEKey[6] = key[5]
		rEKey[0] = key[0]
		rEKey[4] = key[1]
		rEKey[1] = key[2]
		rEKey[5] = key[3]
	} else if i == 2 {
		rEKey[0] = key[0]
		rEKey[4] = key[1]
		rEKey[1] = key[2]
		rEKey[5] = key[3]
	}

	// Key setup
	var tmp [4]uint32
	tmp[0] = pi(rEKey, 0, 0, 1, 2, 3) ^ kp[0]
	tmp[1] = pi(rEKey, 0, 1, 2, 3, 0) ^ kp[1]
	tmp[2] = pi(rEKey, 0, 2, 3, 0, 1) ^ kp[2]
	tmp[3] = pi(rEKey, 0, 3, 0, 1, 2) ^ kp[3]

	rEKey[0] = gammaTau(&tmp, 0, 0, 1)
	rEKey[1] = gammaTau(&tmp, 1, 1, 0)
	rEKey[2] = gammaTau(&tmp, 2, 0, 1)
	rEKey[3] = gammaTau(&tmp, 3, 1, 0)

	tmp[0] = pi(rEKey, 4, 1, 2, 3, 0) ^ kq[0]
	tmp[1] = pi(rEKey, 4, 2, 3, 0, 1) ^ kq[1]
	tmp[2] = pi(rEKey, 4, 3, 0, 1, 2) ^ kq[2]
	tmp[3] = pi(rEKey, 4, 0, 1, 2, 3) ^ kq[3]

	rEKey[4] = gammaTau(&tmp, 0, 1, 0)
	rEKey[5] = gammaTau(&tmp, 1, 0, 1)
	rEKey[6] = gammaTau(&tmp, 2, 1, 0)
	rEKey[7] = gammaTau(&tmp, 3, 0, 1)

	t0 := rEKey[0] ^ rEKey[1] ^ rEKey[2] ^ rEKey[3]
	t1 := rEKey[4] ^ rEKey[5] ^ rEKey[6] ^ rEKey[7]

	rEKey[0] ^= t1
	rEKey[1] ^= t1
	rEKey[2] ^= t1
	rEKey[3] ^= t1
	rEKey[4] ^= t0
	rEKey[5] ^= t0
	rEKey[6] ^= t0
	rEKey[7] ^= t0

	h0Block(rEKey, 0, 8, 16, 0x01010101)
	h1Block(rEKey, 1, 16, 24, 0x01010101)

	h1Block(rEKey, 2, 24, 8, 0x02020202)
	h0Block(rEKey, 3, 8, 16, 0x02020202)

	h0Block(rEKey, 4, 16, 24, 0x04040404)
	h1Block(rEKey, 5, 24, 8, 0x04040404)

	h1Block(rEKey, 6, 8, 16, 0x08080808)
	h0Block(rEKey, 7, 16, 24, 0x08080808)

	h0Block(rEKey, 8, 24, 8, 0x10101010)
	h1Block(rEKey, 9, 8, 16, 0x10101010)

	h1Block(rEKey, 10, 16, 24, 0x20202020)

	rEKey[48] = phiN(rEKey[48], 3, 0, 1, 2)
	rEKey[49] = phiN(rEKey[49], 2, 3, 0, 1)
	rEKey[50] = phiN(rEKey[50], 1, 2, 3, 0)
	rEKey[51] = phiN(rEKey[51], 0, 1, 2, 3)
}

// f0Rnd applies the f0 round function.
func f0Rnd(rKp []uint32, posKp int, rB0, rB1 *[4]uint32) {
	rB1[0] = sTab[0][byteAt(rB0[0], 0)] ^
		sTab[1][byteAt(rB0[1], 0)] ^
		sTab[2][byteAt(rB0[2], 0)] ^
		sTab[3][byteAt(rB0[3], 0)] ^
		rKp[posKp+0]

	rB1[1] = sTab[1][byteAt(rB0[0], 1)] ^
		sTab[2][byteAt(rB0[1], 1)] ^
		sTab[3][byteAt(rB0[2], 1)] ^
		sTab[0][byteAt(rB0[3], 1)] ^
		rKp[posKp+1]

	rB1[2] = sTab[2][byteAt(rB0[0], 2)] ^
		sTab[3][byteAt(rB0[1], 2)] ^
		sTab[0][byteAt(rB0[2], 2)] ^
		sTab[1][byteAt(rB0[3], 2)] ^
		rKp[posKp+2]

	rB1[3] = sTab[3][byteAt(rB0[0], 3)] ^
		sTab[0][byteAt(rB0[1], 3)] ^
		sTab[1][byteAt(rB0[2], 3)] ^
		sTab[2][byteAt(rB0[3], 3)] ^
		rKp[posKp+3]
}

// f1Rnd applies the f1 round function.
func f1Rnd(rKp []uint32, posKp int, rB0, rB1 *[4]uint32) {
	rB0[0] = sTab[1][byteAt(rB1[0], 0)] ^
		sTab[2][byteAt(rB1[1], 0)] ^
		sTab[3][byteAt(rB1[2], 0)] ^
		sTab[0][byteAt(rB1[3], 0)] ^
		rKp[posKp+0]

	rB0[1] = sTab[2][byteAt(rB1[0], 1)] ^
		sTab[3][byteAt(rB1[1], 1)] ^
		sTab[0][byteAt(rB1[2], 1)] ^
		sTab[1][byteAt(rB1[3], 1)] ^
		rKp[posKp+1]

	rB0[2] = sTab[3][byteAt(rB1[0], 2)] ^
		sTab[0][byteAt(rB1[1], 2)] ^
		sTab[1][byteAt(rB1[2], 2)] ^
		sTab[2][byteAt(rB1[3], 2)] ^
		rKp[posKp+2]

	rB0[3] = sTab[0][byteAt(rB1[0], 3)] ^
		sTab[1][byteAt(rB1[1], 3)] ^
		sTab[2][byteAt(rB1[2], 3)] ^
		sTab[3][byteAt(rB1[3], 3)] ^
		rKp[posKp+3]
}
