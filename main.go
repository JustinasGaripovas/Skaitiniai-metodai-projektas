package main

import (
	"custom/point"
	"custom/vector"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

const delta = 0.0001
const numberOfTimesToRepeat = 4

type SpecificPointDirection struct {
	gradient vector.Vector
	id       int
}

func main() {
	writer, _ := os.Create("result.csv")
	writer.WriteString(" 	Point count ,	Thread count,	Time\n")

	//Points
	for i := 2; i <= 20; i += 2 {
		//Threads
		for thread := 1; thread <= 20; thread++ {

			//Average time

			var avgTime = 0.0

			for j := 0; j < numberOfTimesToRepeat; j++ {
				startTime := time.Now()
				calculationInitialization(i, i*2, thread)
				endTime := time.Now()
				elapsedTime := endTime.Sub(startTime).Seconds() * 1000
				avgTime += elapsedTime
			}

			avgTime /= numberOfTimesToRepeat

		//	fmt.Println(strconv.Itoa(i*3) + "," + strconv.Itoa(thread) + "," + strconv.FormatFloat(avgTime, 'f', 0, 64) + "\n")
			writer.WriteString(strconv.Itoa(i*3) + "," + strconv.Itoa(thread) + "," + strconv.FormatFloat(avgTime, 'f', 0, 64) + "\n")
		}
		writer.WriteString("\n")
	}
}

func calculationInitialization(staticPointsAmount int, dynamicPointsAmount int, threadCount int) {

	staticPoints := point.RandomArray(staticPointsAmount)
	dynamicPoints := point.RandomArray(dynamicPointsAmount)
	allPoints := append(dynamicPoints, staticPoints...)

	gradientChannel := make(chan SpecificPointDirection, threadCount)

	var temporaryDynamicPoints []int

	for i := 0; i < dynamicPointsAmount; i++ {
		temporaryDynamicPoints = append(temporaryDynamicPoints, i)
	}

	var threadLoad [][]int

	loadSize := dynamicPointsAmount / threadCount

	if len(temporaryDynamicPoints) < threadCount {
		for i := 0; i < len(temporaryDynamicPoints); i++ {
			threadLoad = append(threadLoad, []int{temporaryDynamicPoints[i]})
		}
	} else {

		for i := 0; i < len(dynamicPoints); i += loadSize {
			end := i + loadSize
			if end > len(temporaryDynamicPoints) {
				end = len(temporaryDynamicPoints)
			}
			threadLoad = append(threadLoad, temporaryDynamicPoints[i:end])
		}
	}
	fmt.Printf("ThreadCount %d  Dynamic points %d Static points %d Load size %d \n",threadCount,dynamicPointsAmount,staticPointsAmount, loadSize)
	fmt.Println(threadLoad)

	for {
		previousValue := point.ArrayDistanceDifferenceSum(allPoints)
		
		for i := 0; i < len(threadLoad); i++ {
			go calculationGoroutine(allPoints, threadLoad[i], gradientChannel)
		}

		var gradients = make([]SpecificPointDirection, dynamicPointsAmount)

		for i := 0; i < dynamicPointsAmount; i++ {
			gradients[i] = <-gradientChannel
		}

		for i := 0; i < dynamicPointsAmount; i++ {
			gradient := gradients[i]
			allPoints[ gradient.id ].X -= gradient.gradient.X
			allPoints[ gradient.id ].Y -= gradient.gradient.Y
		}

		if previousValue < point.ArrayDistanceDifferenceSum(allPoints) {
			break
		}
	}
}

func calculationGoroutine(points []point.Point, threadLoad []int, resultChannel chan SpecificPointDirection) {
	for i := 0; i < len(threadLoad); i++ {

		grad1 := distanceGradient(points, threadLoad[i])
		grad1 = grad1.Normalized()

		grad2 := functionGradient(points[threadLoad[i]])
		grad2 = grad2.Normalized()

		gradient := vector.New(grad1.X+grad2.X, grad1.Y+grad2.Y)
		gradient = gradient.Normalized()

		resultChannel <- SpecificPointDirection{gradient, threadLoad[i]}
	}
}

func distanceGradient(points []point.Point, id int) vector.Vector {
	xPoints := point.CopyArray(points)
	xPoints[id].X += delta

	yPoints := point.CopyArray(points)
	yPoints[id].Y += delta

	gradientX := (point.ArrayDistanceDifferenceSum(xPoints) - point.ArrayDistanceDifferenceSum(points)) / delta
	gradientY := (point.ArrayDistanceDifferenceSum(yPoints) - point.ArrayDistanceDifferenceSum(points)) / delta

	return vector.New(gradientX, gradientY)
}

func functionGradient(currentPoint point.Point) vector.Vector {
	gradientX := (function(point.New(currentPoint.X+delta, currentPoint.Y)) - function(currentPoint)) / delta
	gradientY := (function(point.New(currentPoint.X, currentPoint.Y+delta)) - function(currentPoint)) / delta
	return vector.New(gradientX, gradientY)
}

func function(p point.Point) float64 {
	return p.X*math.Pow(math.E, -((p.X*p.X+p.Y*p.Y)/10)) + 1.5
}
