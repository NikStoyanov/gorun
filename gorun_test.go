package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

// generateTestFile returns a gpx file
func generateTestFile() *bytes.Buffer {
	r := bytes.NewBufferString(`
	<?xml version="1.0" encoding="UTF-8"?>
	<gpx creator="StravaGPX" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd" version="1.1" xmlns="http://www.topografix.com/GPX/1/1" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3">
	 <metadata>
	  <time>2019-04-20T08:59:54Z</time>
	 </metadata>
	 <trk>
	  <name>Morning Hike</name>
	  <type>4</type>
	  <trkseg>
	   <trkpt lat="53.3647850" lon="-1.8160720">
		<ele>244.5</ele>
		<time>2019-04-20T08:59:54Z</time>
		<extensions>
		 <gpxtpx:TrackPointExtension>
		  <gpxtpx:cad>53</gpxtpx:cad>
		 </gpxtpx:TrackPointExtension>
		</extensions>
	   </trkpt>
	   <trkpt lat="53.3647720" lon="-1.8160070">
		<ele>244.4</ele>
		<time>2019-04-20T08:59:57Z</time>
		<extensions>
		 <gpxtpx:TrackPointExtension>
		  <gpxtpx:cad>53</gpxtpx:cad>
		 </gpxtpx:TrackPointExtension>
		</extensions>
	   </trkpt>
	  </trkseg>
	 </trk>
	</gpx>
	`)

	return r
}

func TestReadGPXFile(t *testing.T) {
	x := NewGPX()
	testString := generateTestFile()

	content := testString.Bytes()

	tmpfile, _ := ioutil.TempFile(os.TempDir(), "testfile.gpx")

	tmpfile.Write(content)

	x.readGPXFile(tmpfile.Name())

	tmpfile.Close()
}
