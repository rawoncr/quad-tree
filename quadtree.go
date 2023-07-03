package quad_tree

import (
	"github.com/golang/geo/r2"
	"math"
)

type QuadTree struct {
	root *Node
}

type Point struct {
	ID       int32
	location r2.Point
}

type Node struct {
	box      Box
	points   []r2.Point
	children [4]*Node
	depth    int
}

type NearestPointResult struct {
	Distance float64
	Point    r2.Point
}

func FromPoints(points []r2.Point) QuadTree {
	index := QuadTree{root: &Node{
		box:   NewBox(points),
		depth: 0,
	}}
	for _, p := range points {
		index.Insert(p)
	}
	return index
}

func (tree QuadTree) NearestPoint(needle r2.Point) NearestPointResult {
	best := math.MaxFloat64
	return tree.root.NearestPoint(needle, best)
}

func (n Node) NearestPoint(needle r2.Point, best float64) NearestPointResult {
	result := NearestPointResult{Distance: math.MaxFloat64}
	if n.isLeaf() {
		for _, p := range n.points {
			d := distance(p, needle)
			if d <= best {
				result.Distance = d
				result.Point = p
				best = d
			}
		}
	} else {
		for _, c := range n.children {
			if c != nil && c.box.ShouldInspect(needle, best) {
				nearest := c.NearestPoint(needle, best)
				if nearest.Distance < result.Distance {
					result = nearest
				}
			}
		}
	}
	return result
}

func (n Node) isLeaf() bool {
	return n.children[0] == nil && n.children[1] == nil && n.children[2] == nil && n.children[3] == nil
}

func (n Node) hasCapacity() bool {
	return n.points == nil || len(n.points) < 8
}

func distance(p1, p2 r2.Point) float64 {
	return math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2)
}

func (tree *QuadTree) Insert(p r2.Point) {
	tree.root.insert(p)
}

func (n *Node) insert(p r2.Point) bool {
	if !n.box.ContainsPoint(p) || n.depth == 40 {
		return false
	}
	if n.isLeaf() && n.hasCapacity() {
		n.points = append(n.points, p)
		return true
	}
	n.rebalance(p)
	return true
}

func (n *Node) rebalance(p r2.Point) {
	for _, po := range append(n.points, p) {
		for id, b := range n.box.Children() {
			if b.ContainsPoint(po) {
				n.FindOrCreateChildFromBox(id, b).insert(po)
			}
		}
	}
	n.points = nil
}

func (n *Node) FindOrCreateChildFromBox(id int, b Box) *Node {
	if n.children[id] == nil {
		n.children[id] = &Node{box: b, depth: n.depth + 1}
	}
	return n.children[id]
}
