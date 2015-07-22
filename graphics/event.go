//  event.go -- window event handling

package graphics

import (
	"fmt"
	g "goaldi/runtime"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/gl"
	"os"
	"time"
)

//  An Event is an action in a window.
type Event struct {
	ID        int64   // touch sequence ID (event.TouchSequenceID)
	Action    string  // "touch" | "drag" | "release"
	Lookahead string  // following action, if already pending in queue
	X, Y      float64 // location in user coordinates
}

//  Event.String() produces a printable representation of an Event.
func (e *Event) String() string {
	return fmt.Sprintf("Event(%d,%s,%s,%.2f,%.2f)",
		e.ID, e.Action, e.Lookahead, e.X, e.Y)
}

//  eventQueuer runs as a goroutine to buffer window events.
//  The point of this is to set the event.Lookahead field so that apps
//  can collapse multiple consecutive "drag" or "config" events.
func eventQueuer(inq, outq chan *Event) {

	pending := make([]*Event, 0)
	for { // loop forever handling queue events

		// at this point the queue is empty, so block awaiting an event
		pending = append(pending, <-inq)

		// when the buffer is non-empty, check for both input and output
		for len(pending) > 0 {
			select {
			case outq <- pending[0]:
				// successfully sent an event
				// remove it from the queue
				pending = pending[1:]
			case ev := <-inq:
				// got a new event
				// set the lookahead field of its queued predecessor
				pending[len(pending)-1].Lookahead = ev.Action
				// and add the new event to the queue
				pending = append(pending, ev)
			}
		}
	}
}

//  evtStart signals that the app is ready to go.
func evtStart() {
	appGo <- true
}

//  evtConfig responds to configuration (initial or resize) of the app window.
func evtConfig(new, old event.Config) {
	// save for use in drawing the canvas
	OneApp.Config = new
	// send to Goaldi program event channel
	OneApp.ToEvtQ <- &Event{0, "config", "",
		float64(new.Width), float64(new.Height)}
}

//  evtTouch responds to a mouse (or finger) event
func evtTouch(e event.Touch, f event.Config) {
	// convert to user coordinates
	//#%#%#% assumes that the origin is at the center of the canvas
	m := OneApp.CvScale / OneApp.Canvas.PixPerPt
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
	OneApp.ToEvtQ <- &Event{int64(e.ID), s, "", x, y}
}

//  evtStop responds to an app "stop" call
func evtStop() {
	OneApp.ToEvtQ <- &Event{0, "stop", "", 0, 0} // send to event queue
	go func() {
		// allow program a chance to shut down -- then kill it
		time.Sleep(SHUTDOWN)
		fmt.Fprint(os.Stderr, "Shutdown by window system\n")
		g.Shutdown(0)
	}()
}

//  evtRepaint is called 60x/second to draw the current Canvas on the screen
func evtRepaint(f event.Config) {
	gl.ClearColor(.5, .5, .5, 1)  // color for margins
	gl.Clear(gl.COLOR_BUFFER_BIT) // clear area behind base canvas
	if OneApp.Canvas == nil {     // if canvas not set yet
		return
	}
	OneApp.ConfigDisplay()                   // #%#%# recalculate this every time???
	OneApp.ShowTree(IDENTITY, OneApp.Sprite) // render canvas and sprites
}
