package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const maxTailLength = 10
const maxBalls = 10

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

var (
	view           *tview.Box
	app            *tview.Application
	ballSpeed      = 15
	ballCharacters = []string{"〇", "◯", "○", "○", "◌", "◌", "◌", "◌", "◌", "◌", "⋅"}
	ballColors     = []tcell.Color{tcell.ColorRed, tcell.ColorGreen, tcell.ColorBlue, tcell.ColorYellow, tcell.ColorOrange, tcell.ColorPurple, tcell.ColorAqua, tcell.ColorOlive, tcell.ColorSilver, tcell.ColorMaroon}
	balls          = [maxBalls]BallWithTail{}
	ballCount      = 0
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

func addBall() {
	if ballCount < maxBalls {
		_, _, width, height := view.GetRect()
		balls[ballCount] = BallWithTail{
			head:       Ball{x: rand.Intn(width), y: rand.Intn(height)},
			directionX: rand.Intn(2)*2 - 1,
			directionY: rand.Intn(2)*2 - 1,
			color:      ballColors[ballCount],
			tail:       [maxTailLength]Ball{},
			tailLength: 0,
		}
		ballCount++
	}
}

func removeBall() {
	if ballCount > 1 {
		ballCount--
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
	case tcell.KeyRune:
		if event.Rune() == '+' {
			addBall()
		} else if event.Rune() == '-' {
			removeBall()
		}
	}
	return event
}

func bounce(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	if ballCount == 0 {
		addBall()
	}

	for ball := 0; ball < ballCount; ball++ {
		keepBallInBounds(&balls[ball], width, height)
		drawBall(screen, &balls[ball])
		moveBall(&balls[ball], width, height)
	}

	showStatusAndInstructions(width, height, screen, x)

	return 0, 0, 0, 0
}

func keepBallInBounds(ball *BallWithTail, width int, height int) {
	// Keep the ball in bounds.
	ball.head.x = max(min(ball.head.x, width-2), 1)
	ball.head.y = max(min(ball.head.y, height-2), 1)

	// Keep the tail in bounds.
	for _, tail := range ball.tail {
		tail.x = max(min(tail.x, width-2), 1)
		tail.y = max(min(tail.y, height-2), 1)
	}
}

func drawBall(screen tcell.Screen, ball *BallWithTail) {
	tview.Print(screen, ballCharacters[0], ball.head.x, ball.head.y, 1, tview.AlignCenter, ball.color)

	for i := 0; i < ball.tailLength; i++ {
		tview.Print(screen, ballCharacters[i+1], ball.tail[i].x, ball.tail[i].y, 1, tview.AlignCenter, dimColor(ball.color, int32((maxTailLength-i)*10)))
	}
}

func dimColor(color tcell.Color, percentage int32) tcell.Color {
	r, g, b := color.RGB()
	return tcell.NewRGBColor(r*percentage/100, g*percentage/100, b*percentage/100)
}

func updateTail(ball *BallWithTail) {
	ball.tailLength = min(ball.tailLength+1, maxTailLength)

	for i := ball.tailLength - 1; i > 0; i-- {
		ball.tail[i] = ball.tail[i-1]
	}
	ball.tail[0] = ball.head
}

func moveBall(ball *BallWithTail, width int, height int) {
	updateTail(ball)
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
	msg := fmt.Sprintf("[width=%d, height=%d, Speed=%d, Balls=%d]", width, height, ballSpeed, ballCount)
	tview.Print(screen, msg, x, height/2, width, tview.AlignCenter, tcell.ColorLime)
	tview.Print(screen, "[ESC[] to exit, [UP[]/[Down[] to change speed, [+[]/[-[] Add/Remove Ball", x, height/2+1, width, tview.AlignCenter, tcell.ColorDarkGoldenrod)
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
