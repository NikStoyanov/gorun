/* Modified from:
https://github.com/googlemaps/google-maps-services-go/blob/master/polyline.go

Performs encoding from gps coodinates (GpxData) to a polyline string using
the algorithm described in:
https://developers.google.com/maps/documentation/utilities/polylinealgorithm
*/

package main

import (
	"bytes"
	"io"
)

// PolyEncode returns a new encoded Polyline from a given path of GpxData points.
func PolyEncode(path []GpxData) string {
	var prevLat, prevLng int64

	encodePoly := new(bytes.Buffer)
	encodePoly.Grow(len(path) * 4)

	for _, point := range path {
		// Rounding precision
		lat := int64(point.Latitude * 1e5)
		lng := int64(point.Longitude * 1e5)

		encodeInt(lat-prevLat, encodePoly)
		encodeInt(lng-prevLng, encodePoly)

		prevLat, prevLng = lat, lng
	}

	return encodePoly.String()
}

// encodeInt writes an encoded int64 to the passed io.ByteWriter.
func encodeInt(v int64, w io.ByteWriter) {
	// Left bitwise shift to leave space for a sign bit as the right most bit.
	if v < 0 {
		// Invert a negative coordinate using two's complement.
		v = ^(v << 1)
	} else {
		v <<= 1
	}

	// Add a continuation bit at the LHS for non-last chunks using OR 0x20
	// (0x20 = 100000)
	for v >= 0x20 {
		// Get the last 5 bits (0x1f) and Add 63 to
		// get "better" looking polyline characters in ASCII.
		w.WriteByte((0x20 | (byte(v) & 0x1f)) + 63)
		v >>= 5
	}

	// Modify the last chunk.
	w.WriteByte(byte(v) + 63)
}
