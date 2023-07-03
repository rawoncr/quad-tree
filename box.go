package quad_tree

import (
	"github.com/golang/geo/r1"
	"github.com/golang/geo/r2"
)

type Box struct {
	rectangle r2.Rect
}

func NewBox(points []r2.Point) Box {
	rectangle := r2.EmptyRect()
	for _, p := range points {
		rectangle = rectangle.AddPoint(p)
	}
	return Box{rectangle: rectangle}
}

func (b Box) Children() [4]Box {
	// ul represents the upper left rectangle.
	ul := Box{
		rectangle: r2.Rect{
			X: r1.Interval{Lo: b.rectangle.Lo().X, Hi: b.rectangle.Center().X},
			Y: r1.Interval{Lo: b.rectangle.Center().Y, Hi: b.rectangle.Hi().Y},
		},
	}
	// ur represents the upper right rectangle.
	ur := Box{
		rectangle: r2.Rect{
			X: r1.Interval{Lo: b.rectangle.Center().X, Hi: b.rectangle.Hi().X},
			Y: r1.Interval{Lo: b.rectangle.Center().Y, Hi: b.rectangle.Hi().Y},
		},
	}
	// ll represents the lower left rectangle.
	ll := Box{
		rectangle: r2.Rect{
			X: r1.Interval{Lo: b.rectangle.Lo().X, Hi: b.rectangle.Center().X},
			Y: r1.Interval{Lo: b.rectangle.Lo().Y, Hi: b.rectangle.Center().Y},
		},
	}
	// lr represents the lower right rectangle.
	lr := Box{
		rectangle: r2.Rect{
			X: r1.Interval{Lo: b.rectangle.Center().X, Hi: b.rectangle.Hi().X},
			Y: r1.Interval{Lo: b.rectangle.Lo().Y, Hi: b.rectangle.Center().Y},
		},
	}
	return [4]Box{ul, ur, lr, ll}
}

func (b Box) ContainsPoint(p r2.Point) bool {
	return b.rectangle.ContainsPoint(p)
}

func (b Box) ShouldInspect(needle r2.Point, radius float64) bool {
	closestPoint := b.rectangle.ClampPoint(b.rectangle.Center())
	if distance(closestPoint, needle) <= radius {
		return true
	}
	return false
}
