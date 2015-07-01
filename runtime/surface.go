//  surface.go -- the image type underlying a Goaldi canvas

package runtime

import (
	"fmt"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"image"
	"image/draw"
	"math"
	"os"
	"sync"
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

//  Surface.String() produces a printable representation of a Surface struct.
func (s Surface) String() string {
	a := "-"
	if s.App != nil {
		a = "A"
	}
	return fmt.Sprintf("Surface(%s,%dx%dx%.2f)",
		a, s.Width, s.Height, s.PixPerPt)
}

//  An App struct holds the application window configuration information.
//  Only one Surface can have such a window.
type App struct {
	*glutil.Image            // GLutil image currently displayed on screen
	*Surface                 // associated surface
	event.Config             // current app window configuration
	Events        chan Event // window events
	PixPerPt      float64    // our actual PPP value w/ anti-aliasing
	TL, TR, BL    geom.Point // position for rendering
}

//  App.String() produces a printable representation of the App struct.
func (a App) String() string {
	return fmt.Sprintf("App(%.2fx%.2f+%.2f+%.2f,%.2f,%v)",
		a.TR.X-a.TL.X, a.BL.Y-a.TL.Y, a.TL.X, a.TL.Y, a.PixPerPt,
		a.Surface)
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

//  newSurface initializes and returns a new App or Mem surface.
func newSurface(app *App, im draw.Image, ppp float64) *Surface {
	w := im.Bounds().Max.X
	h := im.Bounds().Max.Y
	s := &Surface{app, w, h, ppp, im}
	if app != nil {
		app.Surface = s
	}
	draw.Draw(im, im.Bounds(), image.White, image.Point{}, draw.Src) // erase
	return s
}

//  MemSurface creates a new off-line Surface with the given characteristics.
func MemSurface(w int, h int, ppp float64) *Surface {
	return newSurface(nil, image.NewRGBA(image.Rect(0, 0, w, h)), ppp)
}

//  AppSurface creates a Surface for use in a golang/x/mobile/app.
func AppSurface() *Surface {
	appOnce.Do(func() { // on first call only:
		appGo <- true // start initialization in main thread
		<-appGo       // wait for it to complete
	})
	return OneApp.Surface
}

//  startup synchronization
var appOnce sync.Once       // initialization interlock
var appGo = make(chan bool) // thread handoff synchronization

//  evtRepaint is called 60x/second to draw the current Surface on the screen
func evtRepaint(g event.Config) {
	gli := OneApp.Image
	gli.Upload()
	gli.Draw(g, OneApp.TL, OneApp.TR, OneApp.BL, gli.Bounds())
}

//  AppMain, when signalled, starts up the main mobile application loop.
//  The Go library requires that this be run in the main thread.
//  #%#%#%#% Is that still a requirement?
func AppMain() {
	<-appGo // block until the first canvas call
	OneApp.Events = make(chan Event, EVBUFSIZE)
	app.Run(app.Callbacks{
		Start:  evtStart,
		Config: evtConfig,
		Stop:   evtStop,
		Touch:  evtTouch,
		Draw:   evtRepaint,
	})
	panic("app.Run() returned")
}

//  evtStart now does nothing.
//  Actual initialization occurs in response to the first Config event.
func evtStart() {
}

//  evtInit initilizes the app in response to the first config event.
func evtInit(cfg event.Config) {
	if cfg.PixelsPerPt >= MinPPP {
		OneApp.PixPerPt = float64(cfg.PixelsPerPt)
	} else {
		OneApp.PixPerPt = MinPPP
	}
	w := int(math.Ceil(float64(cfg.Width) * OneApp.PixPerPt))
	h := int(math.Ceil(float64(cfg.Height) * OneApp.PixPerPt))
	gli := glutil.NewImage(w, h)
	draw.Draw(gli, gli.Bounds(), image.White, image.Point{}, draw.Src) // erase
	OneApp.Image = gli
	OneApp.SetConfig(cfg) // do before setting Ready
	OneApp.Surface = newSurface(&OneApp, OneApp.Image, OneApp.PixPerPt)
	appGo <- true
}

//  evtConfig responds to configuration (init or resize) of the app window.func
func evtConfig(new, old event.Config) {
	if OneApp.PixPerPt == 0 { // if not initialized
		evtInit(new) // then do so
	}
	OneApp.SetConfig(new)
	//#%#%#% DO SOMETHING MORE...
	//#%#%#% SEND TO GOALDI PROGRAM...
	//#%#%#% IMPLICATIONS ON MEANINGS OF CANVAS PPP VALUES & SCALING?
}

//  evtTouch responds to a mouse (or finger) event
func evtTouch(e event.Touch, g event.Config) {
	// convert to user coordinates
	//#%#%#% assumes that the origin is at the center of the canvas
	m := OneApp.PixPerPt / OneApp.Surface.PixPerPt
	x := m * (float64(e.Loc.X - OneApp.Config.Width/2))
	y := m * (float64(e.Loc.Y - OneApp.Config.Height/2))
	// send to the channel
	OneApp.Events <- Event{int64(e.ID), int(e.Type), x, y}
}

//  evtStop responds to an app "stop" call
func evtStop() {
	//#%#%#%# SEND TO GOALDI PROGRAM ?????
	fmt.Fprint(os.Stderr, "Shutdown by window system")
	Shutdown(0)
}

//  App.SetConfig updates the App struct for a new window configuration.
//  Mostly this means figuring out where to draw the OneApp Surface image
//  in the reconfigured window.
func (a *App) SetConfig(g event.Config) {
	a.Config = g
	rwidth := float64(a.Image.Rect.Max.X)  // raster width in pixels
	rheight := float64(a.Image.Rect.Max.X) // raster height in pixels
	raspr := rwidth / rheight              // raster aspect ratio
	daspr := float64(g.Width / g.Height)   // display aspect ratio
	if daspr > raspr {
		// sidebar configuration
		a.PixPerPt = rheight / float64(g.Height)
		rwpts := geom.Pt(raspr) * g.Height // raster width in pts
		dx := (g.Width - rwpts) / 2
		a.TL = geom.Point{dx, 0}
		a.TR = geom.Point{dx + rwpts, 0}
		a.BL = geom.Point{dx, g.Height}
	} else {
		// letterbox configuration
		a.PixPerPt = rwidth / float64(g.Width)
		rhpts := g.Width / geom.Pt(raspr) // raster height in pts
		dy := (g.Height - rhpts) / 2
		a.TL = geom.Point{0, dy}
		a.TR = geom.Point{g.Width, dy}
		a.BL = geom.Point{0, dy + rhpts}
	}
}
