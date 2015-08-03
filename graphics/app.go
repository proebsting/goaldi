//  app.go -- graphics code specific to application windows

package graphics

import (
	"fmt"
	g "goaldi/runtime"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/config"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/geom"
	"os"
	"sync"
	"time"
)

//  Minimum acceptable pixel density, in pixels per point.
//  This provides some anti-aliasing if a coarse-grained screen is zoomed.
//  An irrational value attempts to avoid Moire effects.
const MinPPP = 2.7183 // minimum PixPerPt acceptable

//  Shutdown allowance granted to user program before we force quit.
const SHUTDOWN = 200 * time.Millisecond

//  Identity transform (should be treated as a constant)
var IDENTITY = &f32.Affine{{1, 0, 0}, {0, 1, 0}}

//  An App struct holds the application window configuration information.
type App struct {
	*Canvas               // associated canvas
	CvScale float64       // canvas scaling
	Config  config.Event  // current app window configuration
	ToEvtQ  chan<- *Event // channel for sending events
	Events  <-chan *Event // channel for getting events
}

//  OneApp is the actual data for the single application window.
var OneApp App

//  App.String() produces a printable representation of the App struct.
func (a *App) String() string {
	return fmt.Sprintf("App(%v,%.2f)", a.Canvas, a.CvScale)
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
	w = int(d*float64(OneApp.Config.WidthPt) + 0.5)
	h = int(d*float64(OneApp.Config.HeightPt) + 0.5)
	return w, h, d
}

//  AppCanvas(c) installs canvas c as the application canvas.
func AppCanvas(c *Canvas) {
	c.App = &OneApp
	if OneApp.Canvas != nil {
		OneApp.Canvas.App = nil // disconnect previous app canvas
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
	toqueue := make(chan *Event)
	fmqueue := make(chan *Event)
	go eventQueuer(toqueue, fmqueue)
	OneApp.ToEvtQ = toqueue
	OneApp.Events = fmqueue
	app.Main(func(a app.App) {
		for e := range a.Events() {
			switch e := app.Filter(e).(type) {
			case config.Event:
				evtConfig(e)
			case touch.Event:
				evtTouch(e)
			case paint.Event:
				evtRepaint()
				a.EndPaint(e)
			case lifecycle.Event:
				if e.Crosses(lifecycle.StageVisible) == lifecycle.CrossOn {
					appGo <- true // now alive and visible
				} else if e.Crosses(lifecycle.StageVisible) == lifecycle.CrossOff {
					evtStop() // lost our window
				} // else something else happened, e.g. gained/lost focus
			}
		}
	})
	// allow program a chance to shut down -- then kill it
	time.Sleep(SHUTDOWN)
	fmt.Fprint(os.Stderr, "Shutdown by window system\n")
	g.Shutdown(0)
}

//  App.ConfigDisplay configures the transformation matrix
//  for displaying the underlying canvas.
func (a *App) ConfigDisplay() {
	rwidth := float64(a.Image.Bounds().Max.X)  // raster width in pixels
	rheight := float64(a.Image.Bounds().Max.Y) // raster height in pixels
	raspr := rwidth / rheight                  // raster aspect ratio
	f := a.Config
	daspr := float64(f.WidthPt / f.HeightPt) // display aspect ratio
	dx := float32(0)
	dy := float32(0)
	if daspr > raspr {
		// sidebar configuration
		a.CvScale = rheight / float64(f.HeightPt)
		rwpts := geom.Pt(raspr) * f.HeightPt // raster width in pts
		dx = float32(f.WidthPt-rwpts) / 2
	} else {
		// letterbox configuration
		a.CvScale = rwidth / float64(f.WidthPt)
		rhpts := f.WidthPt / geom.Pt(raspr) // raster height in pts
		dy = float32(f.HeightPt-rhpts) / 2
	}
	m := &a.Sprite.Xform
	m.Translate(IDENTITY, dx, dy)
	m.Scale(m, float32(1/a.CvScale), float32(1/a.CvScale))
}

//  App.ShowTree(xform, sprite) renders the tree of sprites on the screen.
func (a *App) ShowTree(m0 *f32.Affine, e *Sprite) {
	var m f32.Affine
	m.Mul(m0, &e.Xform)
	a.Display(e.Source, &m)
	for _, c := range e.Children {
		a.ShowTree(&m, c)
	}
}

//  App.Display(canvas,xform) displays a canvas on the app screen.
func (a *App) Display(c *Canvas, m *f32.Affine) {
	w := float32(c.Image.Bounds().Max.X)
	h := float32(c.Image.Bounds().Max.Y)
	tl := pj(m, 0, 0)
	tr := pj(m, w, 0)
	bl := pj(m, 0, h)
	c.GLI.Upload()
	c.GLI.Draw(OneApp.Config, tl, tr, bl, c.Image.Bounds())
}

//  pj(xform, x, y) -- project a point using an affine transform
func pj(m *f32.Affine, x float32, y float32) geom.Point {
	return geom.Point{
		X: geom.Pt(m[0][0]*x + m[0][1]*y + m[0][2]),
		Y: geom.Pt(m[1][0]*x + m[1][1]*y + m[1][2]),
	}
}
