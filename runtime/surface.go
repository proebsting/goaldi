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
	Width      int     // width in pixels
	Height     int     // height in pixels
	PixPerPt   float64 // density in pixels/point
	draw.Image         // underlying image
}

//  configure minimum backing store for anti-aliasing coarse-grained screens
const MinPPP = 3 // minimum PixPerPt acceptable

//  app configuration (valid after app initialization)
var cfg app.Config             // current app window configuration
var pixPerPt float64           // our actual PPP value w/ anti-aliasing
var gli *glutil.Image          // GLutil image currently displayed on screen
var once sync.Once             // initialization interlock
var appGo = make(chan bool)    // signal for starting app loop
var appReady = make(chan bool) // signal when initialization is complete

//  MemSurface creates a new off-line Surface with the given characteristics.
func MemSurface(w int, h int, ppp float64) *Surface {
	return newSurface(image.NewRGBA(image.Rect(0, 0, w, h)), ppp)
}

//  AppSurface creates a Surface for use in a golang/x/mobile/app.
func AppSurface() *Surface {
	once.Do(appInit)
	return newSurface(gli, pixPerPt)
}

//  newSurface initializes and returns a new App or Mem surface.
func newSurface(im draw.Image, ppp float64) *Surface {
	w := im.Bounds().Max.X
	h := im.Bounds().Max.Y
	draw.Draw(im, im.Bounds(), image.White, image.Point{}, draw.Src) // erase
	return &Surface{w, h, ppp, im}
}

//  appRepaint is called 60x/second to draw the current Surface on the screen
func appRepaint() {
	gli.Upload()
	gli.Draw(
		geom.Point{0, 0},
		geom.Point{cfg.Width, 0},
		geom.Point{0, cfg.Height},
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
	app.Register(app.Callbacks{
		Start:  appStart,
		Config: appConfig,
		Stop:   appStop,
		Touch:  appTouch,
		Draw:   appRepaint,
	})
	app.Run(app.Callbacks{}) // n.b. argument deprecated
}

//  appStart does the actual initialization once the app driver has started
func appStart() {
	cfg = app.GetConfig()
	if geom.PixelsPerPt >= MinPPP {
		pixPerPt = float64(geom.PixelsPerPt)
	} else {
		pixPerPt = MinPPP
	}
	w := int(math.Ceil(float64(cfg.Width) * float64(pixPerPt)))
	h := int(math.Ceil(float64(cfg.Height) * float64(pixPerPt)))
	gli = glutil.NewImage(w, h)
	draw.Draw(gli, gli.Bounds(), image.White, image.Point{}, draw.Src) // erase
	appReady <- true
}

//  appConfig responds to a resizing of the application window
func appConfig(new, old app.Config) {
	cfg = new
	//#%#%#%# DO SOMETHING MORE...
	//#%#%#%# SEND TO GOALDI PROGRAM...
}

//  appTouch responds to a mouse (or finger) event
func appTouch(e event.Touch) {
	//#%#%#%# SEND TO GOALDI PROGRAM...
}

//  appStop responds to an app "stop" call (#%#% whatever that means...)
func appStop() {
	//#%#%#%# SEND TO GOALDI PROGRAM ?????
	fmt.Fprint(os.Stderr, "Shutdown by window system")
	Shutdown(0)
}
