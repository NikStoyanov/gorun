package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/NikStoyanov/gorun/util"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/joho/godotenv"
	"github.com/tkrajina/gpxgo/gpx"
)

// GoRun an app to display and edit GPX data using Google maps
type GoRun struct {
	Img          mapsImage
	gpxDataArray []util.GpxData

	image   *canvas.Image
	iDEntry *widget.Entry
}

// mapSize is the size of the google maps image
type mapSize struct {
	Horiz int32
	Vert  int32
}

// mapMarker is used to draw markers along the polyline
type mapMarker struct {
	Color string
	Label string
	Coord util.GpxData
}

// mapsImage is to draw image from the Google maps API
type mapsImage struct {
	URL        string
	Base       string
	MapType    string
	Polyline   string
	APIKey     string
	Size       mapSize
	Scale      int16
	MarkStart  mapMarker
	MarkFinish mapMarker
}

// NewGPX returns a new activity app
func NewGPX() *GoRun {
	return &GoRun{}
}

// getGoogleAPI generates the hyper ref to Google Maps
func (x *GoRun) getGoogleAPI() {
	x.Img.URL = fmt.Sprintf("%smaptype=%s&path=enc:%s&key=%s&size=%dx%d&scale=%d&markers=color:%s|label:%s|%f,%f&markers=color:%s|label:%s|%f,%f",
		x.Img.Base, x.Img.MapType, x.Img.Polyline,
		x.Img.APIKey, x.Img.Size.Horiz, x.Img.Size.Vert, x.Img.Scale,
		x.Img.MarkStart.Color, x.Img.MarkStart.Label,
		x.Img.MarkStart.Coord.Latitude, x.Img.MarkStart.Coord.Longitude,
		x.Img.MarkFinish.Color, x.Img.MarkFinish.Label,
		x.Img.MarkFinish.Coord.Latitude, x.Img.MarkFinish.Coord.Longitude)
}

// Submit will lookup the xkcd cartoon and do something useful with it
func (x *GoRun) Submit() {
	filename := x.iDEntry.Text
	x.readGPXFile(filename)
	x.getGoogleAPI()

	// Get the google maps api url
	r := x.Img.URL
	resp, err := http.Get(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		data, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(data, x)
		x.DataToScreen()
	} else {
		fmt.Println("Error from Google Maps API ", resp.Status, resp.StatusCode)
	}
}

// readGPXFile will open the GPX file and return and Latitude-Longitude array
func (x *GoRun) readGPXFile(fileName string) {
	r, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	gpxFile, err := gpx.ParseBytes(r)
	if err != nil {
		panic(err)
	}

	// Extract Latitude-Longitude points
	for _, track := range gpxFile.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				x.gpxDataArray = append(x.gpxDataArray, util.GpxData{point.Latitude, point.Longitude})
			}
		}
	}

	x.Img.Polyline = util.PolyEncode(x.gpxDataArray)
	x.Img.MarkStart.Coord = x.gpxDataArray[0]
	x.Img.MarkFinish.Coord = x.gpxDataArray[len(x.gpxDataArray)-1]
}

// downloadImage gets the generated image from Google Maps
func (x *GoRun) downloadImage() {
	response, e := http.Get(x.Img.URL)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	file, err := ioutil.TempFile(os.TempDir(), "runmap.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	x.image.File = file.Name()
	canvas.Refresh(x.image)
}

// DataToScreen displays the downloaded image
func (x *GoRun) DataToScreen() {
	go x.downloadImage()
}

// NewForm generates a new display form
func (x *GoRun) NewForm(w fyne.Window) fyne.Widget {
	form := &widget.Form{}

	entry := widget.NewEntry()
	x.iDEntry = entry
	form.Append("File name", entry)

	return form
}

func (x *GoRun) initialize() {
	x.Img.Base = "https://maps.googleapis.com/maps/api/staticmap?"
	x.Img.MapType = "roadmap"
	x.Img.Size.Horiz = 1600
	x.Img.Size.Vert = 1600
	x.Img.Scale = 2

	x.Img.MarkStart.Color = "yellow"
	x.Img.MarkFinish.Color = "green"
	x.Img.MarkStart.Label = "S"
	x.Img.MarkFinish.Label = "F"
	x.Img.APIKey = os.Getenv("GOOGLE_STATIC_MAPS_API_KEY")
}

// Show starts a new xkcd widget
func Show(app fyne.App) {
	x := NewGPX()
	x.initialize()
	w := app.NewWindow("GPX Viewer")

	form := x.NewForm(w)
	submitButton := widget.NewButton("Submit", func() {
		x.Submit()
	})
	cancelButton := widget.NewButton("Cancel", func() {
		w.Close()
	})
	submitButton.Style = widget.PrimaryButton

	buttons := widget.NewHBox(
		cancelButton,
		submitButton)

	x.image = &canvas.Image{FillMode: canvas.ImageFillOriginal}
	w.SetContent(fyne.NewContainerWithLayout(
		layout.NewBorderLayout(form, buttons, nil, nil),
		form, buttons, x.image))
	w.Show()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := app.New()

	Show(app)
	app.Run()
}
