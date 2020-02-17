package main

import (
  "image"
)

func DrawGoLBoard(board Matrix, cellWidth int) Canvas {
  numRows := len(board)
  numCols := len(board[0])

  w := cellWidth*numCols
  h := cellWidth*numRows

  c := CreateNewCanvas(w,h)


  black := MakeColor(0,0,0)
  c.SetFillColor(black)
  c.Clear()

  white := MakeColor(255,255,255)
  c.SetStrokeColor(white)

  //Grid lines
  DrawGridLines(c, cellWidth)

  c.SetFillColor(white)
  for row := range board {
    for col := range board[row] {
      if board[row][col].state == "alive" {
        DrawSquare(c, row, col, cellWidth)
      }
    }
  }

  return c
}

func DrawGridLines(c Canvas, cellWidth int) {
  w, h := c.width, c.height

  for i := 1; i < w/cellWidth; i++ {
    x := i*cellWidth
    c.MoveTo(float64(x), 0)
    c.LineTo(float64(x), float64(h))
  }

  for j := 1; j < h/cellWidth; j++ {
    y := j*cellWidth
    c.MoveTo(0,float64(y))
    c.LineTo(float64(w),float64(y))
  }

  c.Stroke()
}

func DrawSquare(c Canvas, row, col, cellWidth int) {
  x1 := col*cellWidth
  y1 := row*cellWidth

  x2 := (col+1)*cellWidth
  y2 := (row+1)*cellWidth

  c.ClearRect(x1, y1, x2, y2)
}

func ConvertGameBoards(boards []Matrix, cellWidth int) []image.Image{
	numGenerations := len(boards)
	imageList := make([]image.Image, numGenerations)
	for i:=range boards {
		imageList[i] = DrawGoLBoard(boards[i], cellWidth).img
	}
	return imageList
}
