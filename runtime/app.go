//  app.go -- graphics code specific to application windows

package runtime

import (
	"fmt"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
	"os"
	"sync"
	"time"
)

//  Minimum acceptable pixel density, in pixels per point.
//  Used for anti-aliasing coarse-grained screens.
//  (An irrational value to try and avoid Moire effects.)
const MinPPP = 2.7183 // minimum PixPerPt acceptable

//  Size of the event buffer.
const EVBUFSIZE = 100

//  Shutdown allowance
const SHUTDOWN = 200 * time.Millisecond

//  An App struct holds the application window configuration information.
type App struct {
	*Canvas                 // associated canvas
	event.Config            // current app window configuration
	Events       chan Event // window event channel
	PixPerPt     float64    // our actual PPP value w/ anti-aliasing
	TL, TR, BL   geom.Point // position for rendering
}

//  App.String() produces a printable representation of the App struct.
func (a *App) String() string {
	return fmt.Sprintf("App(%.2fx%.2f+%.2f+%.2f,%.2f,%v)",
		a.TR.X-a.TL.X, a.BL.Y-a.TL.Y, a.TL.X, a.TL.Y, a.PixPerPt,
		a.Canvas)
}

var OneApp App // data for the one app window

//  An Event is an action in a window.
type Event struct {
	ID     int64   // touch sequence ID (event.TouchSequenceID)
	Action string  // "touch" | "drag" | "release"
	X, Y   float64 // location in user coordinates
}

//  AppSize returns the current size for an application canvas.
//  On the first call, it starts up the application main loop.
func AppSize() (w int, h int, d float64) {
	appOnce.Do(func() { // on first call only:
		appGo <- true // start the application loop
		<-appGo       // wait for signal from the start callback
	})
	d = float64(OneApp.Config.PixelsPerPt)
	if d < MinPPP {
		d = MinPPP
	}
	w = int(d*float64(OneApp.Config.Width) + 0.5)
	h = int(d*float64(OneApp.Config.Height) + 0.5)
	return w, h, d
}

//  AppCanvas(c) installs canvas c as the application canvas.
func AppCanvas(c *Canvas) {
	c.App = &OneApp
	if OneApp.Canvas != nil {
		OneApp.Canvas.App = nil // disconnect old app canvas
	}
	OneApp.Canvas = c
}

//  startup synchronization
var appOnce sync.Once       // one-time initialization flag
var appGo = make(chan bool) // thread handoff synchronization

//  AppMain, when signalled, runs the main mobile application loop.
//  The Go library requires that this be run in the main thread.
//  One or more Config events will precede the Start event.
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

//  evtStart signals that the app is ready to go.
func evtStart() {
	appGo <- true
}

//  evtConfig responds to configuration (init or resize) of the app window.
func evtConfig(new, old event.Config) {
	// save for use in drawing the canvas
	OneApp.Config = new
	// send to Goaldi program event channel
	OneApp.Events <- Event{0, "config", float64(new.Width), float64(new.Height)}
}

//  evtTouch responds to a mouse (or finger) event
func evtTouch(e event.Touch, g event.Config) {
	// convert to user coordinates
	//#%#%#% assumes that the origin is at the center of the canvas
	m := OneApp.PixPerPt / OneApp.Canvas.PixPerPt
	x := m * (float64(e.Loc.X - OneApp.Config.Width/2))
	y := m * (float64(e.Loc.Y - OneApp.Config.Height/2))
	// send to the channel
	var s string
	switch e.Change {
	case event.ChangeOn:
		s = "touch"
	case event.ChangeNone:
		s = "drag"
	case event.ChangeOff:
		s = "release"
	default:
		panic(fmt.Sprintf("Unexpected event type: %v", e))
	}
	OneApp.Events <- Event{int64(e.ID), s, x, y}
}

//  evtStop responds to an app "stop" call
func evtStop() {
	OneApp.Events <- Event{0, "stop", 0, 0}              // send to program
	time.Sleep(SHUTDOWN)                                 // allow to shutdown
	fmt.Fprint(os.Stderr, "Shutdown by window system\n") // force kill
	Shutdown(0)
}

//  evtRepaint is called 60x/second to draw the current Canvas on the screen
func evtRepaint(g event.Config) {
	gl.ClearColor(.5, .5, .5, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	if OneApp.Canvas == nil { // if canvas not set yet
		return
	}
	OneApp.SetConfig(OneApp.Config) //#%#%#%#% recalc every time?
	gli := OneApp.Image
	gli.Upload()
	gli.Draw(g, OneApp.TL, OneApp.TR, OneApp.BL, gli.Bounds())
}

//  App.SetConfig updates the App struct for a new window configuration.
//  Mostly this means figuring out where to draw the OneApp Canvas image
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
