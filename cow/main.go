package main

import (
	"math/rand"
	"time"
	"runtime"
	"fmt"
)


type Cell struct {
	state string //Either alive or dead grass.

	cowpresent bool //is the cow here at this iteration?

	probability float64 //p or 2p, depending on Cow Type

}

type Matrix [][]Cell //a 2 dimensional slice of Cells.

type Neighborhood struct { //just organizing, so using pointers.

	neighbors []*Cell //center, then all 9 adjacent cells
	center    *Cell
}


type OrderedPair struct{
	x,y int
}

func main() {

	//seeding PRNG
	rand.Seed(time.Now().UTC().UnixNano())

	// p := 1.0

	numGens:= 144 //144 20-minute rounds per simulation.
	numProcs := runtime.NumCPU() // use all cores available
	numSimulations := 10000

	fmt.Println("I have", runtime.NumCPU(), "cores available.")

	fmt.Println("Simulating Clarabelle 10,000 times beginning at (0,0) on an 11x11 board of uneaten grass.")

fmt.Println("Clarabelle ate", 100.0*SimulateClarabelleParallel(numSimulations, numProcs, numGens),"% of the grass given 48 hours.")


fmt.Println("Simulating Bernadette 10,000 times beginning at (0,0) on an 11x11 board of uneaten grass.")

minsTo50:=20.0*SimulateBernadetteParallel(numSimulations, numProcs, numGens)

hrsTo50:=minsTo50/60.0

fmt.Println("To reach 50% of grass eaten, over 10,000 simulations, Bernadette took:", hrsTo50,"hours.")


}




//want to see number of trials when Bernadette reaches 50% on average as well as average percentage eaten by 144*20mins.
//so, we will terminate simulation when half the grass is eaten on any board.

func SimulateBernadetteParallel(numSimulations, numProcs,numGens int) float64{ //want to see number of trials when Bernadette reaches 50% on average.


allTrialsToHalf := []float64{} //create an empty slice of float64s

c := make(chan float64, numProcs) //each processor will get a somewhat equal number of simulations.
// channel returns the percentage of eaten grass.
	for k := 0; k < numProcs; k++ {
			go SimulateBernadetteSingleProc(numGens, numSimulations, numProcs, c)
}

//grabbing values from channel, then averaging.
for k:=1;k<=numProcs; k++{
  currTrials := <- c
	allTrialsToHalf = append(allTrialsToHalf, currTrials) //will be one percent grass eaten value for each processor.
	// fmt.Println(len(percentages)) useful for debugging.
}

total := 0.0
for _, val := range allTrialsToHalf{
	total += val
}
return total/float64(numProcs)
}


func SimulateBernadetteSingleProc(numGens, numSimulations, numProcs int, c chan float64){

sumAvg:=0.0
countSims := 0.0

for i:=0;i<numSimulations/numProcs; i++ { //run this singple proc process a number of times
	//equal to total simulations divided by available processors.

trialsTaken := (SimulateBernadette(numGens)) //this is working.
//percent of grass eaten this 144-increment trial.

sumAvg += trialsTaken //20 minutes is one trial.
countSims ++

}

procSimAvg := sumAvg/countSims
// fmt.Println(procSimAvg) //working

c <- procSimAvg
}







//takes in numgens, and simulates Bernadette's 144 movements (48 hours) once.
//we can then use this function to count total living squares upon the farmer's return, then average this amount over a monte-carlo simulated
// number of trials (and use Go's concurrency to run these trials!).
func SimulateBernadette(numGens int) float64 {

	boardSlice:=make([]Matrix,numGens)
	p := 1.0 //arbitrarily set p = 1.0.

	cowType := "Bernadette"
	board := InitializeBoard(cowType, p) //all living grass set to 2.0. Once eaten, grass will be set to p = 1.0

	board[0][0].cowpresent = true //we will put the cow at (0,0) as a starting position (assumed).

	for i := 0; i < numGens; i++ {
		boardSlice[i] = MoveBernadetteOnce(board)

		if TakeStockEaten(boardSlice[i])>=0.50{ //if more than half the grass is eaten, return the number of trials.

		return float64(i) //return the current generation if eaten more than half of grass.
	}
	}

	return float64(numGens) //return 144 trials
}

//move this cow once and she feeds or relaxes.
func MoveBernadetteOnce(board Matrix) Matrix {
	for i := range board {
		for j := range board[i] {
			if board[i][j].cowpresent == true { //detect where cow is and
				// eat grass if applicable.

				board[i][j].state = "dead" //eat the grass
				board[i][j].probability = 1.0 //set movement probability to 1.0

				board[i][j].cowpresent = false
				//move the cow to a random neighbor with probability p.

				var x, y = RollBernadette(board, i, j) //get new coordinates to move cow to.

				for InField(board, x, y) == false {
					x, y = RollBernadette(board, i, j)
					//while roll puts the cow's movement out of the board, re-roll and reset x and y.
				}
				board[x][y].cowpresent = true

			}
		}
	}
	return board
}

//this function selects a new square with probabilities in mind.
func RollBernadette(board Matrix, i, j int) (int, int) {

pile := []OrderedPair{} //create a slice of strings representing neighbors.

		for x := i-1; x <= i+1; x++ {
			for y := j-1; y <= j + 1; y ++ {
				if InField(board,x,y) == true {

					if board[x][y].probability == 2.0{
						var currPair OrderedPair
						currPair.x = x
						currPair.y = y

						pile = append(pile,currPair)
						pile = append(pile,currPair) //append the current coordinates to the pile
						//TWICE if probability is 2.0*p = 2.0*1.0

					} else if board[x][y].probability == 1.0{
						var currPair OrderedPair
						currPair.x = x
						currPair.y = y
						pile = append(pile,currPair) //just append once if the grass is dead.

					}
				}

			}
		}
// pile is now a slice of ordered pairs for which we can randomly choose from.

n:=len(pile)

roll:=rand.Intn(n)

toX := pile[roll].x
toY := pile[roll].y

	return toX, toY
}




// **************



func SimulateClarabelleParallel(numSimulations, numProcs,numGens int) float64{

percentages := []float64{} //create an empty slice of float64s

c := make(chan float64, numProcs) //each processor will get a somewhat equal number of simulations.
// channel returns the percentage of eaten grass.
	for k := 0; k < numProcs; k++ {
			go SimulateClarabelleSingleProc(numGens, numSimulations, numProcs, c)
}

//grabbing values from channel, then averaging.
for k:=1;k<=numProcs; k++{
  currPercentage := <- c
	percentages = append(percentages, currPercentage) //will be one percent grass eaten value for each processor.
	// fmt.Println(len(percentages)) useful for debugging.
}

total := 0.0
for _, val := range percentages{
	total += val
}
return total/float64(numProcs)
}


func SimulateClarabelleSingleProc(numGens, numSimulations, numProcs int, c chan float64){

sumAvg:=0.0
countSims := 0.0

for i:=0;i<numSimulations/numProcs; i++ { //run this singple proc process a number of times
	//equal to total simulations divided by available processors.

percentEaten := TakeStockEaten(SimulateClarabelle(numGens)) //working.
//percent of grass eaten this 144-increment trial.

sumAvg += percentEaten
countSims ++

}

procSimAvg := sumAvg/countSims
// fmt.Println(procSimAvg) //working

c <- procSimAvg
}




//takes in numgens, and simulates Clarabelle's 144 movements (48 hours) once.
//we can then use this function to count total living squares upon the farmer's return, then average this amount over a monte-carlo simulated
// number of trials (and use Go's concurrency to run these trials!).
func SimulateClarabelle(numGens int) Matrix {

	boardSlice:=make([]Matrix,numGens)
	p := 1.0

	cowType := "Clarabelle"
	board := InitializeBoard(cowType, p)

	board[0][0].cowpresent = true //we will put the cow at (0,0) as a starting assumption.

	for i := 0; i < numGens; i++ {
		boardSlice[i] = MoveClarabelleOnce(board)
	}

	return boardSlice[numGens-1] //return final board.
}

//move this cow once and she feeds or relaxes.
func MoveClarabelleOnce(board Matrix) Matrix {
	for i := range board {
		for j := range board[i] {
			if board[i][j].cowpresent == true { //detect where cow is and
				// eat grass if applicable.

				board[i][j].state = "dead" //eat the grass

				board[i][j].cowpresent = false
				//move the cow to a random neighbor with probability p.

				var x, y = Roll(i, j) //get new coordinates to move cow to.

				for InField(board, x, y) == false {
					x, y = Roll(i, j)
					//while roll puts the cow's movement out of the board, re-roll and reset x and y.
				}
				board[x][y].cowpresent = true

			}
		}
	}
	return board
}




//returns percentage eaten of final board in slice of boards.
func TakeStockEaten(finalBoard Matrix) float64 {

totalSq := 11.0*11.0

countEaten := 0.0
	for i:=range finalBoard {
		for j:=range finalBoard[i]{
			if finalBoard[i][j].state == "dead"{
				countEaten ++ //increment if grass has been eaten.
			}
		}
	}
return countEaten/totalSq

}



func InitializeBoard(cowType string, p float64) Matrix {

	startBoard := make(Matrix, 11)
	for r := range startBoard {
		// range over rows and make each column
		startBoard[r] = make([]Cell, 11)
	}



	for i := range startBoard {
		for j := range startBoard[i] {
			startBoard[i][j].state = "alive"    // all grass starts at alive.
			startBoard[i][j].cowpresent = false //we will put cow at (0,0) as a starting assumption.
			// setting probabilities
			if cowType == "Bernadette" {
				startBoard[i][j].probability = 2.0 * p // all living grass has 2.0 probability.
			} else if cowType == "Clarabelle" {
				startBoard[i][j].probability = p
			}

		}
	}
	return startBoard
}


func Roll(i, j int) (int, int) {

	x := i
	y := j

	roll := rand.Intn(8) //return a random integer on [0,8), representing clockwise neighbors (top,top/right,right,...,top/left)
	//infield statements
	switch roll { //depending on the roll, a new square will get the cow.
	case 0:
		x--
	case 1:
		x--
		y++
	case 2:
		y++
	case 3:
		x++
		y++
	case 4:
		x++
	case 5:
		x++
		y--
	case 6:
		y--
	case 7:
		x--
		y--
	}

	return x, y
}

//InField takes a GameBoard and i/j indices.  It returns true if (i,j) is a valid entry
//of the board.
func InField(currBoard Matrix, i, j int) bool {
	numRows := len(currBoard) //should be 11 in this instance
	numCols := len(currBoard[0])
	if i < 0 || j < 0 || i >= numRows || j >= numCols {
		return false
	}
	// if we make it here, we are in the field.
	return true
}
