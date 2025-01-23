package main

import (
	"fmt"
	"math"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const refreshInterval = 25 * time.Millisecond

var (
	view *tview.Box
	app  *tview.Application
)

var ballX, ballY = 1, 1
var dx, dy = 1, 1

func main() {
	app = tview.NewApplication()
	view = tview.NewBox().
		SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle("[green:white] [::b]Bouncing Ball ")
	view.SetDrawFunc(bounce)

	go refresh()

	if err := app.SetRoot(view, true).Run(); err != nil {
		panic(err)
	}
}

func bounce(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	// Ensure the ball remains in the box.
	ballX = int(math.Max(float64(math.Min(float64(ballX), float64(width-2))), 1))
	ballY = int(math.Max(float64(math.Min(float64(ballY), float64(height-2))), 1))

	// Display coordinates.
	msg := fmt.Sprintf("x=%d, y=%d - width=%d, height=%d", ballX, ballY, width, height)
	tview.Print(screen, msg, x, height/2, width, tview.AlignCenter, tcell.ColorLime)

	// Draw the ball.
	tview.Print(screen, "o", ballX, ballY, 1, tview.AlignCenter, tcell.ColorWhite)

	// Move the ball.
	ballX += dx
	ballY += dy

	// Bounce the ball off the walls.
	if ballX >= width-2 || ballX <= 1 {
		dx = -dx
	}
	if ballY >= height-2 || ballY <= 1 {
		dy = -dy
	}

	return 0, 0, 0, 0
}

func refresh() {
	tick := time.NewTicker(refreshInterval)
	for range tick.C {
		app.Draw()
	}
}
