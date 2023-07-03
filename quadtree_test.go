package quad_tree

import (
	"github.com/golang/geo/r2"
	"log"
	"math"
	"math/rand"
	"testing"
)

func TestIndex_Search(t *testing.T) {

	points := makeRandomPoints()
	index := FromPoints(points)

	log.Println("Testing index search with", len(points), "points")

	for i := 0; i < len(points); i++ {
		// any random point.
		needle := r2.Point{X: rand.Float64() * 1000, Y: rand.Float64() * 1000}
		// determine nearest point by brute force.
		want := nearestPointBruteForce(needle, points)
		// use the index to find the nearest point.
		got := index.NearestPoint(needle)
		if want.Distance != got.Distance || want.Point != got.Point {
			t.Fatalf("test: %d, expected %f and got %f", i, want.Distance, got.Distance)
		}
	}

}

// nearestPoint determine the nearest point by brute force.
func nearestPointBruteForce(needle r2.Point, points []r2.Point) NearestPointResult {
	best := math.MaxFloat64
	result := NearestPointResult{}
	for _, p := range points {
		d := distance(needle, p)
		if d < best {
			best = d
			result.Point = p
			result.Distance = d
		}
	}
	return result
}

func makeRandomPoints() []r2.Point {
	points := make([]r2.Point, rand.Intn(100))
	for i, _ := range points {
		points[i].X = rand.Float64() * 1000
		points[i].Y = rand.Float64() * 1000
	}

	return points
}
