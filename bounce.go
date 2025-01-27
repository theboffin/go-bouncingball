package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const maxTailLength = 10

// Represents a ball, its position, and its direction.
type Ball struct {
	x int
	y int
}

type BallWithTail struct {
	head       Ball
	directionX int
	directionY int
	color      tcell.Color
	tail       [maxTailLength]Ball
	tailLength int
}

var ball BallWithTail = BallWithTail{
	head: Ball{x: 1,
		y: 1},
	directionX: 1,
	directionY: 1,
	color:      tcell.ColorRed,
	tail:       [maxTailLength]Ball{},
	tailLength: 0,
}

var (
	view      *tview.Box
	app       *tview.Application
	ballSpeed = 15
)

func main() {
	app = tview.NewApplication().SetInputCapture(handleKeyboard)
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

func handleKeyboard(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		app.Stop()
	case tcell.KeyUp:
		ballSpeed = max(ballSpeed-1, 1)
	case tcell.KeyDown:
		ballSpeed = min(ballSpeed+1, 100)
	}
	return event
}

func bounce(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	keepBallInBounds(width, height)
	showStatusAndInstructions(width, height, screen, x)
	drawBall(screen)
	moveBall(width, height)

	return 0, 0, 0, 0
}

func keepBallInBounds(width int, height int) {
	// Keep the ball in bounds.
	ball.head.x = max(min(ball.head.x, width-2), 1)
	ball.head.y = max(min(ball.head.y, height-2), 1)

	// Keep the tail in bounds.
	for _, tail := range ball.tail {
		tail.x = max(min(tail.x, width-2), 1)
		tail.y = max(min(tail.y, height-2), 1)
	}
}

func drawBall(screen tcell.Screen) {
	tview.Print(screen, "[::b]O", ball.head.x, ball.head.y, 1, tview.AlignCenter, ball.color)

	for i := 0; i < ball.tailLength; i++ {
		tview.Print(screen, "*", ball.tail[i].x, ball.tail[i].y, 1, tview.AlignCenter, dimColor(ball.color, int32((maxTailLength-i)*10)))
	}
}

func dimColor(color tcell.Color, percentage int32) tcell.Color {
	r, g, b := color.RGB()
	return tcell.NewRGBColor(r*percentage/100, g*percentage/100, b*percentage/100)
}

func updateTail() {
	ball.tailLength = min(ball.tailLength+1, maxTailLength)

	for i := ball.tailLength - 1; i > 0; i-- {
		ball.tail[i] = ball.tail[i-1]
	}
	ball.tail[0] = ball.head
}

func moveBall(width int, height int) {
	updateTail()
	ball.head.x += ball.directionX
	ball.head.y += ball.directionY

	// Bounce the ball off the walls.
	if ball.head.x >= width-2 || ball.head.x <= 1 {
		ball.directionX = -ball.directionX
	}
	if ball.head.y >= height-2 || ball.head.y <= 1 {
		ball.directionY = -ball.directionY
	}
}

func showStatusAndInstructions(width int, height int, screen tcell.Screen, x int) {
	msg := fmt.Sprintf("x=%d, y=%d - [width=%d, height=%d, Speed=%d]", ball.head.x, ball.head.y, width, height, ballSpeed)
	tview.Print(screen, msg, x, height/2, width, tview.AlignCenter, tcell.ColorLime)
	tview.Print(screen, "Press ESC to exit, Cursor UP/Down to change speed", x, height/2+1, width, tview.AlignCenter, tcell.ColorDarkGoldenrod)
}

func refresh() {
	tick := time.NewTicker(time.Duration(ballSpeed) * time.Millisecond)
	for range tick.C {
		// Refresh the screen.
		app.Draw()
		// Reset the timer to the current speed.
		tick.Reset(time.Duration(ballSpeed) * time.Millisecond)
	}
}
