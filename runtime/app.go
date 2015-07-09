//  app.go -- graphics code specific to application windows

package runtime

import (
	"fmt"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/exp/f32"
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
}

//  App.String() produces a printable representation of the App struct.
func (a *App) String() string {
	return fmt.Sprintf("App(%v,%.2f)", a.Canvas, a.PixPerPt)
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
	m := &f32.Affine{}
	OneApp.SetMatrix(m)
	OneApp.Display(OneApp.Canvas, m)
}

//  App.SetMatrix(m) initializes a transformation matrix for the base canvas.
func (a *App) SetMatrix(m *f32.Affine) {
	rwidth := float64(a.Image.Rect.Max.X)  // raster width in pixels
	rheight := float64(a.Image.Rect.Max.Y) // raster height in pixels
	raspr := rwidth / rheight              // raster aspect ratio
	g := a.Config
	daspr := float64(g.Width / g.Height) // display aspect ratio
	dx := float32(0)
	dy := float32(0)
	if daspr > raspr {
		// sidebar configuration
		a.PixPerPt = rheight / float64(g.Height)
		rwpts := geom.Pt(raspr) * g.Height // raster width in pts
		dx = float32(g.Width-rwpts) / 2
	} else {
		// letterbox configuration
		a.PixPerPt = rwidth / float64(g.Width)
		rhpts := g.Width / geom.Pt(raspr) // raster height in pts
		dy = float32(g.Height-rhpts) / 2
	}
	sc := float32(1 / a.PixPerPt)
	m.Identity()
	m.Translate(m, dx, dy)
	m.Scale(m, sc, sc)
}

//  App.Display(canvas,xform) displays a canvas on the app screen.
func (a *App) Display(c *Canvas, m *f32.Affine) {
	w := float32(c.Image.Rect.Max.X)
	h := float32(c.Image.Rect.Max.Y)
	tl := pj(m, 0, 0)
	tr := pj(m, w, 0)
	bl := pj(m, 0, h)
	c.Image.Upload()
	c.Image.Draw(OneApp.Config, tl, tr, bl, c.Image.Bounds())
}

//  pj(xform, x, y) -- project a point using an affine transform
func pj(m *f32.Affine, x float32, y float32) geom.Point {
	return geom.Point{
		geom.Pt(m[0][0]*x + m[0][1]*y + m[0][2]),
		geom.Pt(m[1][0]*x + m[1][1]*y + m[1][2]),
	}
}
