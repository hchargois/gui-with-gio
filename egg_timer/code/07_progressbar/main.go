package main

import (
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Timer struct {
	elapsed time.Duration
	target  time.Duration
	running bool
}

func NewTimer(target time.Duration) *Timer {
	return &Timer{
		target: target,
	}
}

func (t *Timer) IsRunning() bool {
	return t.running
}

func (t *Timer) Start() {
	t.running = true
}

func (t *Timer) Stop() {
	t.running = false
}

func (t *Timer) Reset() {
	t.elapsed = 0
	t.running = false
}

func (t *Timer) IsFinished() bool {
	return t.elapsed >= t.target
}

func (t *Timer) Progress() float32 {
	return float32(t.elapsed) / float32(t.target)
}

func (t *Timer) Advance(dt time.Duration) {
	if t.running {
		t.elapsed += dt
		if t.elapsed >= t.target {
			t.elapsed = t.target
			t.running = false
		}
	}
}

func main() {
	go func() {
		// create new window
		w := new(app.Window)
		w.Option(app.Title("Egg timer"))
		w.Option(app.Size(unit.Dp(400), unit.Dp(600)))
		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

type C = layout.Context
type D = layout.Dimensions

func draw(w *app.Window) error {
	// ops are the operations from the UI
	var ops op.Ops

	var startButton widget.Clickable
	var stopButton widget.Clickable
	var resetButton widget.Clickable

	var lastFrameTime time.Time

	// th defines the material design style
	th := material.NewTheme()

	timer := NewTimer(3 * time.Second)

	for {
		// listen for events in the window
		switch e := w.Event().(type) {

		// this is sent when the application should re-render
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			timer.Advance(gtx.Now.Sub(lastFrameTime))

			if startButton.Clicked(gtx) {
				timer.Start()
			}
			if stopButton.Clicked(gtx) {
				timer.Stop()
			}
			if resetButton.Clicked(gtx) {
				timer.Reset()
			}

			// Let's try out the flexbox layout concept
			layout.Flex{
				// Vertical alignment, from top to bottom
				Axis: layout.Vertical,
				// Empty space is left at the start, i.e. at the top
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx C) D {
						bar := material.ProgressBar(th, timer.Progress())
						return bar.Layout(gtx)
					},
				),
				layout.Rigid(
					func(gtx C) D {
						// We start by defining a set of margins
						margins := layout.Inset{
							Top:    unit.Dp(25),
							Bottom: unit.Dp(25),
							Right:  unit.Dp(35),
							Left:   unit.Dp(35),
						}
						// Then we lay out within those margins ...
						return margins.Layout(gtx,
							// ...the same function we earlier used to create a button
							func(gtx C) D {
								if timer.IsRunning() {
									return material.Button(th, &stopButton, "stop").Layout(gtx)
								}
								if timer.IsFinished() {
									return material.Button(th, &resetButton, "reset").Layout(gtx)
								}
								return material.Button(th, &startButton, "start").Layout(gtx)
							},
						)
					},
				),
			)
			if timer.IsRunning() {
				gtx.Execute(op.InvalidateCmd{})
			}
			lastFrameTime = gtx.Now
			e.Frame(gtx.Ops)

		// this is sent when the application is closed
		case app.DestroyEvent:
			return e.Err
		}

	}
}
