package main

import (
  "fmt"
)

func PrintBoards(boards []Matrix) {
  for _, board := range boards {
    PrintBoard(board)
  }
}

func PrintBoard(board Matrix) {
  for _, row := range board {
    PrintRow(row)
  }
  fmt.Println()
}

func PrintRow(row []Cell) {
  for c := range row {
    if row[c].state == "alive" {//alive
      fmt.Print("⬜")
    } else {
      fmt.Print("⬛")
    }
  }
  fmt.Println()
}
