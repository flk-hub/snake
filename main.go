package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Status int

const (
	GameOver Status = iota
	Running
	Paused
)

const (
	ScreenOffset int32 = 75
	ScreenWidth  int32 = 800
	ScreenHeight int32 = 800
	Fps          int32 = 10

	ObjectsSize    int32 = 20
	BlocksInRow    int32 = (ScreenWidth - 2*ScreenOffset) / ObjectsSize
	BlocksInColumn int32 = (ScreenHeight - 2*ScreenOffset) / ObjectsSize
)

type Grid [BlocksInRow][BlocksInColumn]bool

type Position struct {
	X, Y int32
}

type Snake struct {
	head      Position
	size      int
	direction int
	body      [BlocksInRow * BlocksInColumn]Position
}

var GameGrid Grid

func (grid *Grid) Draw() {
	for x := int32(0); x < BlocksInRow; x++ {
		for y := int32(0); y < BlocksInColumn; y++ {

			if grid[x][y] {
				rl.DrawRectangle((x*ObjectsSize)+ScreenOffset, (y*ObjectsSize)+ScreenOffset, ObjectsSize, ObjectsSize, rl.DarkGray)
				continue
			}
			//rl.DrawRectangleLines((x*ObjectsSize)+ScreenOffset, (y*ObjectsSize)+ScreenOffset, ObjectsSize, ObjectsSize, rl.DarkGray)
		}
	}
	// game box
	rl.DrawRectangleLines(ScreenOffset, ScreenOffset, BlocksInRow*ObjectsSize, BlocksInColumn*ObjectsSize, rl.DarkGray)
}

func genRandomPosition() Position {
	return Position{
		X: rl.GetRandomValue(0, BlocksInRow-1),
		Y: rl.GetRandomValue(0, BlocksInColumn-1),
	}
}

func genFood() (Position, error) {
	var tries uint

	for {
		food := genRandomPosition()
		// prevent to generate food on top of snake
		if !GameGrid.CheckCollision(&food) {
			return food, nil
		}
		if GameGrid.CheckCollision(&food) {
			log.Println("food on top of snake")
			tries += 1

		}
		if tries > 5 {
			log.Println("tried to generate food 5 times, giving up")
			return Position{}, errors.New("failed to generate food")
		}

	}

}

func NewSnake() *Snake {
	return &Snake{
		head:      genRandomPosition(),
		direction: rl.KeyRight,
		size:      1,
	}
}

func (s *Snake) MoveHead() {

	switch s.direction {
	case rl.KeyLeft:
		s.head.X -= 1
	case rl.KeyRight:
		s.head.X += 1
	case rl.KeyUp:
		s.head.Y -= 1
	case rl.KeyDown:
		s.head.Y += 1
	}

	// reset position

	if s.head.X < 0 {
		// limit the x axis
		s.head.X = BlocksInRow - 1
	}
	if s.head.X >= BlocksInRow {
		s.head.X = 0
	}
	if s.head.Y < 0 {
		// limit the y axis
		s.head.Y = BlocksInColumn - 1
	}
	if s.head.Y >= BlocksInColumn {
		s.head.Y = 0
	}

}
func (s *Grid) CheckCollision(p *Position) bool {

	// if position is true on the grid thn there is a collision
	return GameGrid[p.X][p.Y]

}

func (s *Snake) InFoodRange(food *Position) bool {
	return s.head.X == food.X &&
		s.head.Y == food.Y

}

func (s *Snake) Move() error {

	// reset last position
	tailX := s.body[s.size-1].X
	tailY := s.body[s.size-1].Y
	GameGrid[tailX][tailY] = false
	for i := s.size - 1; i >= 0; i-- {
		// prevent overflow if 20 at the 20 pos should be 19
		if i == 0 {
			// pass the body to the head position
			s.body[i].X = s.head.X
			s.body[i].Y = s.head.Y
		} else {
			s.body[i].X = s.body[i-1].X
			s.body[i].Y = s.body[i-1].Y
		}

		GameGrid[s.body[i].X][s.body[i].Y] = true

		log.Printf("Body Coordinates Index:%v X, Y: %v \n", i, s.body[i])

	}

	s.MoveHead()
	if GameGrid.CheckCollision(&s.head) {
		log.Print("There is a collision with the head and body")
		return errors.New("collision with the body is not allowed")
	}
	return nil

}

func main() {
	var err error
	var score uint
	rl.InitWindow(ScreenWidth, ScreenHeight, "Snake!!!")
	defer rl.CloseWindow()

	rl.SetTargetFPS(Fps)
	gameStatus := Running
	snake, food := NewSnake(), genRandomPosition()
	// game loop
	for !rl.WindowShouldClose() {

		// handle the pause mode
		if rl.IsKeyPressed(rl.KeySpace) {
			if gameStatus == Paused {
				gameStatus = Running
			} else {
				gameStatus = Paused
			}
		}

		if gameStatus == Running {
			// save snake direction
			if rl.IsKeyDown(rl.KeyLeft) {
				snake.direction = rl.KeyLeft
			}
			if rl.IsKeyDown(rl.KeyRight) {
				snake.direction = rl.KeyRight
			}
			if rl.IsKeyDown(rl.KeyUp) {
				snake.direction = rl.KeyUp
			}
			if rl.IsKeyDown(rl.KeyDown) {
				snake.direction = rl.KeyDown
			}

			if err = snake.Move(); err != nil {
				gameStatus = GameOver
			}
			//check if eat the food after moving
			if snake.InFoodRange(&food) {
				snake.size += 1
				score += 10
				food, err = genFood()
				if err != nil {
					gameStatus = GameOver
				}

			}

		}

		rl.BeginDrawing()
		if gameStatus != GameOver {
			rl.DrawText(fmt.Sprintf("SCORE: %d", score), ScreenOffset, 55, 20, rl.Maroon)
			GameGrid.Draw()
			log.Printf("Head Coordinates X, Y:%v\n", snake.head)
			rl.DrawRectangle(ScreenOffset+snake.head.X*ObjectsSize, ScreenOffset+snake.head.Y*ObjectsSize, ObjectsSize, ObjectsSize, rl.SkyBlue)
			rl.DrawRectangle(ScreenOffset+food.X*ObjectsSize, ScreenOffset+food.Y*ObjectsSize, ObjectsSize, ObjectsSize, rl.Red)
			if gameStatus == Paused {
				rl.DrawText("PAUSED", ScreenWidth/4, ScreenHeight/4, 30, rl.Green)
			}
		} else {
			rl.DrawText("GAME OVER", ScreenWidth/4, ScreenHeight/4, 30, rl.Maroon)
			rl.DrawText(fmt.Sprintf("Reason: %v", strings.ToUpper(err.Error())), ScreenHeight/4, ScreenHeight/3, 10, rl.Maroon)
		}
		rl.ClearBackground(rl.RayWhite)
		rl.EndDrawing()
	}
}
