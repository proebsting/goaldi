//  surface.go -- the image type underlying a Goaldi canvas

package runtime

import (
	//"code.google.com/p/freetype-go/freetype"
	"fmt"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event"
	//"golang.org/x/mobile/font"
	"golang.org/x/mobile/geom"
	//"golang.org/x/mobile/gl"
	"golang.org/x/mobile/gl/glutil"
	"image"
	//"image/color"
	"image/draw"
	//"log"
	"math"
	"os"
	"sync"
	//"time"
)

//  A Surface is the actual writing area for a canvas.
//  It can be written to a file and/or displayed on the screen.
type Surface struct {
	*App               // app configuration, or nil
	Width      int     // width in pixels
	Height     int     // height in pixels
	PixPerPt   float64 // density in pixels/point
	draw.Image         // underlying image
}

//  An App struct holds the application window configuration information.
//  Only one Surface can have such a window.
type App struct {
	*glutil.Image            // GLutil image currently displayed on screen
	app.Config               // current app window configuration
	Events        chan Event // window events
	pixPerPt      float64    // our actual PPP value w/ anti-aliasing
}

var OneApp App // data for the one app

//  An Event is an action in a window.
type Event struct {
	ID     int64   // touch sequence ID (event.TouchSequenceID)
	Action int     // 0=begin, 1=move, 2=release (event.TouchType)
	X, Y   float64 // location in user coordinates
}

//  minimum backing store for anti-aliasing coarse-grained screens
const MinPPP = 3 // minimum PixPerPt acceptable

//  size of the event buffer
const EVBUFSIZE = 1000

//  startup synchronization
var appOnce sync.Once          // initialization interlock
var appGo = make(chan bool)    // signal for starting app loop
var appReady = make(chan bool) // signal when initialization is complete

//  MemSurface creates a new off-line Surface with the given characteristics.
func MemSurface(w int, h int, ppp float64) *Surface {
	return newSurface(nil, image.NewRGBA(image.Rect(0, 0, w, h)), ppp)
}

//  AppSurface creates a Surface for use in a golang/x/mobile/app.
func AppSurface() *Surface {
	appOnce.Do(appInit)
	return newSurface(&OneApp, OneApp.Image, OneApp.pixPerPt)
}

//  newSurface initializes and returns a new App or Mem surface.
func newSurface(app *App, im draw.Image, ppp float64) *Surface {
	w := im.Bounds().Max.X
	h := im.Bounds().Max.Y
	draw.Draw(im, im.Bounds(), image.White, image.Point{}, draw.Src) // erase
	return &Surface{app, w, h, ppp, im}
}

//  evtRepaint is called 60x/second to draw the current Surface on the screen
func evtRepaint() {
	gli := OneApp.Image
	gli.Upload()
	gli.Draw(
		geom.Point{0, 0},
		geom.Point{OneApp.Config.Width, 0},
		geom.Point{0, OneApp.Config.Height},
		gli.Bounds(),
	)
}

//  appInit starts the main loop and waits for its initialization to finish
func appInit() {
	appGo <- true
	<-appReady
}

//  AppMain, when signalled, starts up the main mobile application loop.
//  The Go library requires that this be run in the main thread.
func AppMain() {
	<-appGo
	OneApp.Events = make(chan Event, EVBUFSIZE)
	app.Register(app.Callbacks{
		Start:  evtStart,
		Config: evtConfig,
		Stop:   evtStop,
		Touch:  evtTouch,
		Draw:   evtRepaint,
	})
	app.Run(app.Callbacks{}) // n.b. argument deprecated
}

//  evtStart does the actual initialization once the app driver has started
func evtStart() {
	OneApp.Config = app.GetConfig()
	if geom.PixelsPerPt >= MinPPP {
		OneApp.pixPerPt = float64(geom.PixelsPerPt)
	} else {
		OneApp.pixPerPt = MinPPP
	}
	w := int(math.Ceil(float64(OneApp.Config.Width) * float64(OneApp.pixPerPt)))
	h := int(math.Ceil(float64(OneApp.Config.Height) * float64(OneApp.pixPerPt)))
	gli := glutil.NewImage(w, h)
	draw.Draw(gli, gli.Bounds(), image.White, image.Point{}, draw.Src) // erase
	OneApp.Image = gli
	appReady <- true
}

//  evtConfig responds to a resizing of the application window
func evtConfig(new, old app.Config) {
	OneApp.Config = new
	//#%#%#%# DO SOMETHING MORE...
	//#%#%#%# SEND TO GOALDI PROGRAM...
}

//  evtTouch responds to a mouse (or finger) event
func evtTouch(e event.Touch) {
	// convert to user coordinates
	//#%#%#% assumes the window has not been resized
	//#%#%#% and the origin is still at the center
	x := float64(e.Loc.X - OneApp.Config.Width/2)
	y := float64(e.Loc.Y - OneApp.Config.Height/2)
	// send to the channel
	OneApp.Events <- Event{int64(e.ID), int(e.Type), x, y}
}

//  evtStop responds to an app "stop" call (#%#% whatever that means...)
func evtStop() {
	//#%#%#%# SEND TO GOALDI PROGRAM ?????
	fmt.Fprint(os.Stderr, "Shutdown by window system")
	Shutdown(0)
}
