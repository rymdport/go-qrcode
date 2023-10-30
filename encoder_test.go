// go-qrcode
// Copyright 2014 Tom Harwood

package qrcode

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/rymdport/go-qrcode/bitset"
)

func TestClassifyDataMode(t *testing.T) {
	tests := []struct {
		data   string
		actual []segment
	}{
		{
			"\x30",
			[]segment{
				{
					dataModeNumeric,
					"\x30",
				},
			},
		},
		{
			"\x30\x41\x42\x43\x20\x00\xf0\xf1\xf2\x31",
			[]segment{
				{
					dataModeNumeric,
					"\x30",
				},
				{
					dataModeAlphanumeric,
					"\x41\x42\x43\x20",
				},
				{
					dataModeByte,
					"\x00\xf0\xf1\xf2",
				},
				{
					dataModeNumeric,
					"\x31",
				},
			},
		},
	}

	for _, test := range tests {
		encoder := newDataEncoder(dataEncoderType1To9)
		encoder.encode(test.data)

		if !reflect.DeepEqual(test.actual, encoder.actual) {
			t.Errorf("Got %v, expected %v", encoder.actual, test.actual)
		}
	}
}

func TestByteModeLengthCalculations(t *testing.T) {
	tests := []struct {
		dataEncoderType dataEncoderType
		dataMode        dataMode
		numSymbols      int
		expectedLength  int
	}{}

	for i, test := range tests {
		encoder := newDataEncoder(test.dataEncoderType)
		var resultLength int

		resultLength, err := encoder.encodedLength(test.dataMode, test.numSymbols)

		if test.expectedLength == -1 {
			if err == nil {
				t.Errorf("Test %d: got length %d, expected error", i, resultLength)
			}
		} else if resultLength != test.expectedLength {
			t.Errorf("Test %d: got length %d, expected length %d", i, resultLength,
				test.expectedLength)
		}
	}
}

func TestSingleModeEncodings(t *testing.T) {
	tests := []struct {
		dataEncoderType dataEncoderType
		dataMode        dataMode
		data            string
		expected        *bitset.Bitset
	}{
		{
			dataEncoderType1To9,
			dataModeNumeric,
			"01234567",
			bitset.NewFromBase2String("0001 0000001000 0000001100 0101011001 1000011"),
		},
		{
			dataEncoderType1To9,
			dataModeAlphanumeric,
			"AC-42",
			bitset.NewFromBase2String("0010 000000101 00111001110 11100111001 000010"),
		},
		{
			dataEncoderType1To9,
			dataModeByte,
			"123",
			bitset.NewFromBase2String("0100 00000011 00110001 00110010 00110011"),
		},
		{
			dataEncoderType10To26,
			dataModeByte,
			"123",
			bitset.NewFromBase2String("0100 00000000 00000011 00110001 00110010 00110011"),
		},
		{
			dataEncoderType27To40,
			dataModeByte,
			"123",
			bitset.NewFromBase2String("0100 00000000 00000011 00110001 00110010 00110011"),
		},
	}

	for _, test := range tests {
		encoder := newDataEncoder(test.dataEncoderType)
		encoded := bitset.New()

		encoder.encodeDataRaw(test.data, test.dataMode, encoded)

		if !test.expected.Equals(encoded) {
			t.Errorf("For %s got %s, expected %s", test.data, encoded.String(),
				test.expected.String())
		}
	}
}

type testModeSegment struct {
	dataMode dataMode
	numChars int
}

func TestOptimiseEncoding(t *testing.T) {
	tests := []struct {
		dataEncoderType dataEncoderType
		actual          []testModeSegment
		optimised       []testModeSegment
	}{
		// Coalescing multiple segments.
		{
			dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 1}, // length = 4 + 9 + 6 = 19 bits
				{dataModeNumeric, 1},      // length = 4 + 10 + 4 = 18 bits
				{dataModeAlphanumeric, 1}, // 19 bits.
				{dataModeNumeric, 1},      // 18 bits.
				{dataModeAlphanumeric, 1}, // 19 bits.
				// total = 93 bits.
			},
			[]testModeSegment{
				{dataModeAlphanumeric, 5}, // length = 4 + 9 + 22 + 6 = 41.
			},
		},
		// Coalesing not necessary.
		{
			dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 1},
				{dataModeNumeric, 20},
			},
			[]testModeSegment{
				{dataModeAlphanumeric, 1},
				{dataModeNumeric, 20},
			},
		},
		// Switch to more general dataMode.
		{
			dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 100},
				{dataModeByte, 1},
				{dataModeNumeric, 1},
			},
			[]testModeSegment{
				{dataModeAlphanumeric, 100},
				{dataModeByte, 2},
			},
		},
		// Sometimes encoding everything as bytes is more efficient.
		{
			dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 1},
				{dataModeByte, 1},
				{dataModeNumeric, 1},
			},
			[]testModeSegment{
				{dataModeByte, 3},
			},
		},
		// https://www.google.com/123456789012345678901234567890
		// BBBBBAAABBBABBBBBBABBBANNNNNNNNNNNNNNNNNNNNNNNNNNNNNN
		{
			dataEncoderType1To9,
			[]testModeSegment{
				{dataModeByte, 5},
				{dataModeAlphanumeric, 3},
				{dataModeByte, 3},
				{dataModeAlphanumeric, 1},
				{dataModeByte, 6},
				{dataModeAlphanumeric, 1},
				{dataModeAlphanumeric, 4},
				{dataModeNumeric, 30},
			},
			[]testModeSegment{
				{dataModeByte, 23},
				{dataModeNumeric, 30},
			},
		},
		// https://www.google.com/123
		// BBBBBAAABBBABBBBBBABBBANNN
		// Small segments are inefficient because of additional metadata.
		{
			dataEncoderType1To9,
			[]testModeSegment{
				{dataModeByte, 5},
				{dataModeAlphanumeric, 3},
				{dataModeByte, 3},
				{dataModeAlphanumeric, 1},
				{dataModeByte, 6},
				{dataModeAlphanumeric, 1},
				{dataModeAlphanumeric, 4},
				{dataModeNumeric, 3},
			},
			[]testModeSegment{
				{dataModeByte, 26},
			},
		},
		// HTTPS://WWW.GOOGLE.COM/123
		// AAAAAAAAAAAAAAAAAAAAAAANNN
		{
			dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 23},
				{dataModeNumeric, 3},
			},
			[]testModeSegment{
				{dataModeAlphanumeric, 26},
			},
		},
		{
			dataEncoderType27To40,
			[]testModeSegment{
				{dataModeByte, 1},
				{dataModeNumeric, 1},
				{dataModeByte, 1},
				{dataModeNumeric, 1},
				{dataModeByte, 1},
				{dataModeNumeric, 1},
				{dataModeByte, 1},
				{dataModeNumeric, 1},
			},
			[]testModeSegment{
				{dataModeByte, 8},
			},
		},
		// HTTPS://ABC.DE/Q/393AABB6998877XYZ0518AUQCRVJN25
		// AAAAAAAAAAAAAAAAANNNAAAANNNNNNNAAANNNNAAAAAAAANN
		// different to below---------^--------------------
		{
			dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 17},
				{dataModeNumeric, 3},
				{dataModeAlphanumeric, 4},
				{dataModeNumeric, 7},
				{dataModeAlphanumeric, 3},
				{dataModeNumeric, 4},
				{dataModeAlphanumeric, 8},
				{dataModeNumeric, 2},
			},
			[]testModeSegment{
				{dataModeAlphanumeric, 48},
			},
		},
		// HTTPS://ABC.DE/Q/393AABB699E877XYZ0518AUQCRVJN25
		// AAAAAAAAAAAAAAAAANNNAAAANNNANNNAAANNNNAAAAAAAANN
		// different to above---------^--------------------
		{
			dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 17},
				{dataModeNumeric, 3},
				{dataModeAlphanumeric, 4},
				{dataModeNumeric, 3},
				{dataModeAlphanumeric, 1},
				{dataModeNumeric, 3},
				{dataModeAlphanumeric, 3},
				{dataModeNumeric, 4},
				{dataModeAlphanumeric, 8},
				{dataModeNumeric, 2},
			},
			[]testModeSegment{
				{dataModeAlphanumeric, 48},
			},
		},
		// 0123456789
		// NNNNNNNNNN
		{
			dataEncoderType1To9,
			[]testModeSegment{
				{dataModeNumeric, 10},
			},
			[]testModeSegment{
				{dataModeNumeric, 10},
			},
		},
	}

	for _, test := range tests {
		numTotalChars := 0
		for _, v := range test.actual {
			numTotalChars += v.numChars
		}

		var builder strings.Builder
		builder.Grow(numTotalChars)

		for _, v := range test.actual {
			for j := 0; j < v.numChars; j++ {
				switch v.dataMode {
				case dataModeNumeric:
					builder.WriteByte('1')
				case dataModeAlphanumeric:
					builder.WriteByte('A')
				case dataModeByte:
					builder.WriteByte('#')
				default:
					t.Fatal("Unrecognised data mode")
				}
			}
		}

		encoder := newDataEncoder(test.dataEncoderType)

		_, err := encoder.encode(builder.String())

		if err != nil {
			t.Errorf("Got %s, expected valid encoding", err.Error())
		} else {
			ok := true

			if len(encoder.optimised) != len(test.optimised) {
				ok = false
			} else {
				for i, s := range test.optimised {
					if encoder.optimised[i].dataMode != s.dataMode ||
						len(encoder.optimised[i].data) != s.numChars {
						ok = false
						break
					}
				}
			}

			if !ok {
				t.Errorf("got %s, expected %s", segmentsString(encoder.optimised),
					testModeSegmentsString(test.optimised))
			}
		}
	}
}

func testModeSegmentsString(segments []testModeSegment) string {
	result := "["

	for i, segment := range segments {
		if i > 0 {
			result += ", "
		}

		result += fmt.Sprintf("%d*%s", segment.numChars,
			dataModeString(segment.dataMode))
	}

	result += "]"

	return result
}

func segmentsString(segments []segment) string {
	result := "["

	for i, segment := range segments {
		if i > 0 {
			result += ", "
		}

		result += fmt.Sprintf("%d*%s", len(segment.data),
			dataModeString(segment.dataMode))
	}

	result += "]"

	return result
}
