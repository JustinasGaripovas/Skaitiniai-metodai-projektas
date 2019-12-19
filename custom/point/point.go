package point

import (
	"math"
	"math/rand"
	"time"
)

type Point struct {
	X float64
	Y float64
}

func New(X float64, Y float64) Point {
	point := Point{X, Y}
	return point
}

func Random() Point {
	rand.NewSource(time.Now().UnixNano())
	min := -10
	max := 10
	randomX := rand.Intn(max - min + 1) + min
	randomY := rand.Intn(max - min + 1) + min

	point := Point{float64(randomX), float64(randomY)}
	return point
}

func RandomArray(amount int) []Point{
	var staticPoints []Point

	for i := 0; i < amount; i++  {
		staticPoints = append(staticPoints, Random())
	}

	return staticPoints
}

func (myPoint Point) DistanceTo(point Point) float64 {
	return math.Sqrt(math.Pow(myPoint.X-point.X, 2) + math.Pow(myPoint.Y-point.Y, 2));
}

func CopyArray(points []Point) []Point{
	var newArray []Point
	for i := 0; i < len(points); i++  {
		newArray = append(newArray, points[i])
	}
	return newArray
}

func ArrayDistanceDifferenceSum(points []Point) float64{

	average := ArrayDistanceAverage(points)
	difference := 0.0
	length := len(points)

	for i := 0; i < length - 1 ; i++  {
		for j := i +1; j < length ; j++  {
			difference += math.Abs(average - points[i].DistanceTo(points[j]))
		}
	}
	return difference
}

func ArrayDistanceAverage(points []Point) float64 {
	sum := 0.0
	length := len(points)
	for i := 0; i < length - 1 ; i++  {
		for j := i + 1; j < length ; j++  {
			sum += points[i].DistanceTo(points[j])
		}
	}
	distanceCount := float64(length) * (float64(length)-1)

	return sum / distanceCount
}



