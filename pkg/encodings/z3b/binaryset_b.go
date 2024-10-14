// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.
//go:build disabled
// +build disabled

package z3b

// Binary Sets: Each set maps to a unique range of byte values.
// The sets are designed to optimize encoding of random data by
// maximizing the likelihood that any random byte is present in the current set.
// Each byte value appears in at least two sets, and some in three, distributing overlaps evenly.
var (
	binarySet1 = [setSize]byte{
		0x00, 0x09, 0x12, 0x1B, 0x24, 0x2D, 0x36, 0x3F,
		0x48, 0x51, 0x5A, 0x63, 0x6C, 0x75, 0x7E, 0x87,
		0x90, 0x99, 0xA2, 0xAB, 0xB4, 0xBD, 0xC6, 0xCF,
		0xD8, 0xE1, 0xEA, 0xF3, 0xFC, 0x05, 0x0E, 0x17,
		0x20, 0x29, 0x32, 0x3B, 0x44, 0x4D, 0x56, 0x5F,
		0x68, 0x71, 0x7A, 0x83, 0x8C, 0x95, 0x9E, 0xA7,
		0xB0, 0xB9, 0xC2, 0xCB, 0xD4, 0xDD, 0xE6, 0xEF,
		0xF8, 0x01, 0x0A, 0x13, 0x1C, 0x25, 0x2E, 0x37,
		0x40, 0x49, 0x52, 0x5B, 0x64, 0x6D, 0x76, 0x7F,
		0x88, 0x91, 0x9A, 0xA3, 0xAC, 0xB5, 0xBE, 0xC7,
		0xD0, 0xD9, 0xE2, 0xEB, 0xF4, 0xFD,
	}

	binarySet2 = [setSize]byte{
		0x01, 0x0A, 0x13, 0x1C, 0x25, 0x2E, 0x37, 0x40,
		0x49, 0x52, 0x5B, 0x64, 0x6D, 0x76, 0x7F, 0x88,
		0x91, 0x9A, 0xA3, 0xAC, 0xB5, 0xBE, 0xC7, 0xD0,
		0xD9, 0xE2, 0xEB, 0xF4, 0xFD, 0x06, 0x0F, 0x18,
		0x21, 0x2A, 0x33, 0x3C, 0x45, 0x4E, 0x57, 0x60,
		0x69, 0x72, 0x7B, 0x84, 0x8D, 0x96, 0x9F, 0xA8,
		0xB1, 0xBA, 0xC3, 0xCC, 0xD5, 0xDE, 0xE7, 0xF0,
		0xF9, 0x02, 0x0B, 0x14, 0x1D, 0x26, 0x2F, 0x38,
		0x41, 0x4A, 0x53, 0x5C, 0x65, 0x6E, 0x77, 0x80,
		0x89, 0x92, 0x9B, 0xA4, 0xAD, 0xB6, 0xBF, 0xC8,
		0xD1, 0xDA, 0xE3, 0xEC, 0xF5, 0xFE,
	}

	binarySet3 = [setSize]byte{
		0x02, 0x0B, 0x14, 0x1D, 0x26, 0x2F, 0x38, 0x41,
		0x4A, 0x53, 0x5C, 0x65, 0x6E, 0x77, 0x80, 0x89,
		0x92, 0x9B, 0xA4, 0xAD, 0xB6, 0xBF, 0xC8, 0xD1,
		0xDA, 0xE3, 0xEC, 0xF5, 0xFE, 0x07, 0x10, 0x19,
		0x22, 0x2B, 0x34, 0x3D, 0x46, 0x4F, 0x58, 0x61,
		0x6A, 0x73, 0x7C, 0x85, 0x8E, 0x97, 0xA0, 0xA9,
		0xB2, 0xBB, 0xC4, 0xCD, 0xD6, 0xDF, 0xE8, 0xF1,
		0xFA, 0x03, 0x0C, 0x15, 0x1E, 0x27, 0x30, 0x39,
		0x42, 0x4B, 0x54, 0x5D, 0x66, 0x6F, 0x78, 0x81,
		0x8A, 0x93, 0x9C, 0xA5, 0xAE, 0xB7, 0xC0, 0xC9,
		0xD2, 0xDB, 0xE4, 0xED, 0xF6, 0xFF,
	}

	binarySet4 = [setSize]byte{
		0x03, 0x0C, 0x15, 0x1E, 0x27, 0x30, 0x39, 0x42,
		0x4B, 0x54, 0x5D, 0x66, 0x6F, 0x78, 0x81, 0x8A,
		0x93, 0x9C, 0xA5, 0xAE, 0xB7, 0xC0, 0xC9, 0xD2,
		0xDB, 0xE4, 0xED, 0xF6, 0xFF, 0x08, 0x11, 0x1A,
		0x23, 0x2C, 0x35, 0x3E, 0x47, 0x50, 0x59, 0x62,
		0x6B, 0x74, 0x7D, 0x86, 0x8F, 0x98, 0xA1, 0xAA,
		0xB3, 0xBC, 0xC5, 0xCE, 0xD7, 0xE0, 0xE9, 0xF2,
		0xFB, 0x04, 0x0D, 0x16, 0x1F, 0x28, 0x31, 0x3A,
		0x43, 0x4C, 0x55, 0x5E, 0x67, 0x70, 0x79, 0x82,
		0x8B, 0x94, 0x9D, 0xA6, 0xAF, 0xB8, 0xC1, 0xCA,
		0xD3, 0xDC, 0xE5, 0xEE, 0xF7, 0x00,
	}

	binarySet5 = [setSize]byte{
		0x04, 0x0D, 0x16, 0x1F, 0x28, 0x31, 0x3A, 0x43,
		0x4C, 0x55, 0x5E, 0x67, 0x70, 0x79, 0x82, 0x8B,
		0x94, 0x9D, 0xA6, 0xAF, 0xB8, 0xC1, 0xCA, 0xD3,
		0xDC, 0xE5, 0xEE, 0xF7, 0x00, 0x09, 0x12, 0x1B,
		0x24, 0x2D, 0x36, 0x3F, 0x48, 0x51, 0x5A, 0x63,
		0x6C, 0x75, 0x7E, 0x87, 0x90, 0x99, 0xA2, 0xAB,
		0xB4, 0xBD, 0xC6, 0xCF, 0xD8, 0xE1, 0xEA, 0xF3,
		0xFC, 0x05, 0x0E, 0x17, 0x20, 0x29, 0x32, 0x3B,
		0x44, 0x4D, 0x56, 0x5F, 0x68, 0x71, 0x7A, 0x83,
		0x8C, 0x95, 0x9E, 0xA7, 0xB0, 0xB9, 0xC2, 0xCB,
		0xD4, 0xDD, 0xE6, 0xEF, 0xF8, 0x01,
	}

	binarySet6 = [setSize]byte{
		0x05, 0x0E, 0x17, 0x20, 0x29, 0x32, 0x3B, 0x44,
		0x4D, 0x56, 0x5F, 0x68, 0x71, 0x7A, 0x83, 0x8C,
		0x95, 0x9E, 0xA7, 0xB0, 0xB9, 0xC2, 0xCB, 0xD4,
		0xDD, 0xE6, 0xEF, 0xF8, 0x01, 0x0A, 0x13, 0x1C,
		0x25, 0x2E, 0x37, 0x40, 0x49, 0x52, 0x5B, 0x64,
		0x6D, 0x76, 0x7F, 0x88, 0x91, 0x9A, 0xA3, 0xAC,
		0xB5, 0xBE, 0xC7, 0xD0, 0xD9, 0xE2, 0xEB, 0xF4,
		0xFD, 0x06, 0x0F, 0x18, 0x21, 0x2A, 0x33, 0x3C,
		0x45, 0x4E, 0x57, 0x60, 0x69, 0x72, 0x7B, 0x84,
		0x8D, 0x96, 0x9F, 0xA8, 0xB1, 0xBA, 0xC3, 0xCC,
		0xD5, 0xDE, 0xE7, 0xF0, 0xF9, 0x02,
	}

	binarySet7 = [setSize]byte{
		0x06, 0x0F, 0x18, 0x21, 0x2A, 0x33, 0x3C, 0x45,
		0x4E, 0x57, 0x60, 0x69, 0x72, 0x7B, 0x84, 0x8D,
		0x96, 0x9F, 0xA8, 0xB1, 0xBA, 0xC3, 0xCC, 0xD5,
		0xDE, 0xE7, 0xF0, 0xF9, 0x02, 0x0B, 0x14, 0x1D,
		0x26, 0x2F, 0x38, 0x41, 0x4A, 0x53, 0x5C, 0x65,
		0x6E, 0x77, 0x80, 0x89, 0x92, 0x9B, 0xA4, 0xAD,
		0xB6, 0xBF, 0xC8, 0xD1, 0xDA, 0xE3, 0xEC, 0xF5,
		0xFE, 0x07, 0x10, 0x19, 0x22, 0x2B, 0x34, 0x3D,
		0x46, 0x4F, 0x58, 0x61, 0x6A, 0x73, 0x7C, 0x85,
		0x8E, 0x97, 0xA0, 0xA9, 0xB2, 0xBB, 0xC4, 0xCD,
		0xD6, 0xDF, 0xE8, 0xF1, 0xFA, 0x03,
	}
)