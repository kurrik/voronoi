// Copyright 2012 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// An implementation of Fortune's algorithm to get Voronoi edges for a set of
// points.
package voronoi

import (
	"container/heap"
	"fmt"
	"math"
)

type Point struct {
	X float64
	Y float64
}

func Pt(x float64, y float64) *Point {
	return &Point{X: x, Y: y}
}

type Vertices []*Point

type Edge struct {
	Start     *Point
	End       *Point
	Direction *Point
	Left      *Point
	Right     *Point
	F         float64
	G         float64
	Neighbor  *Edge
}

func Ed(x1 float64, y1 float64, x2 float64, y2 float64) *Edge {
	return &Edge{Start: Pt(x1, y1), End: Pt(x2, y2)}
}

func NewEdge(s *Point, a *Point, b *Point) *Edge {
	e := &Edge{
		Start:    s,
		Left:     a,
		Right:    b,
		Neighbor: nil,
		End:      nil,
	}
	e.F = (b.X - a.X) / (a.Y - b.Y)
	e.G = s.Y - e.F*s.X

	if math.IsInf(float64(e.F), -1) {
		e.Direction = Pt(1, 0)
	} else if math.IsInf(float64(e.F), 1) {
		e.Direction = Pt(-1, 0)
	} else {
		e.Direction = Pt(b.Y-a.Y, -(b.X - a.X))
	}
	return e
}

type Edges []*Edge

type Event struct {
	Point   *Point
	IsPlace bool
	Y       float64
	Arch    *Parabola
}

func NewEvent(p *Point, place bool) *Event {
	e := &Event{
		Point:   p,
		IsPlace: place,
		Y:       p.Y,
		Arch:    nil,
	}
	return e
}

type EventList []*Event

func (l EventList) Find(e *Event) int {
	for i := 0; i < len(l); i++ {
		if l[i] == e {
			return i
		}
	}
	return -1
}

func (l EventList) Last() *Event {
	if len(l) == 0 {
		return nil
	}
	return l[len(l)-1]
}

func (l *EventList) Remove(e *Event) bool {
	i := l.Find(e)
	if i == -1 {
		return false
	}
	a := *l
	*l = append(a[:i], a[i+1:]...)
	return true
}

type EventQueue []*Event

func (q EventQueue) Len() int {
	return len(q)
}

func (q EventQueue) Less(i int, j int) bool {
	// Sorted by Y descending.
	return q[i].Y > q[j].Y
}

func (q EventQueue) Swap(i int, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *EventQueue) Push(x interface{}) {
	a := *q
	n := len(a)
	if n+1 > cap(a) {
		c := make(EventQueue, len(a), 2*cap(a)+1)
		copy(c, a)
		a = c
	}
	a = a[0 : n+1]
	event := x.(*Event)
	a[n] = event
	*q = a
}

func (q *EventQueue) Pop() interface{} {
	a := *q
	n := len(a)
	event := a[n-1]
	*q = a[0 : n-1]
	return event
}

type Parabola struct {
	IsLeaf bool
	Site   *Point
	Edge   *Edge
	Event  *Event
	Parent *Parabola
	left   *Parabola
	right  *Parabola
}

func NewParabola() *Parabola {
	return &Parabola{
		Site:   nil,
		IsLeaf: false,
		Edge:   nil,
		Event:  nil,
		Parent: nil,
		left:   nil,
		right:  nil,
	}
}

func NewLeafParabola(s *Point) *Parabola {
	p := NewParabola()
	p.Site = s
	p.IsLeaf = true
	return p
}

func (p *Parabola) Left() *Parabola {
	return p.left
}

func (p *Parabola) Right() *Parabola {
	return p.right
}

func (p *Parabola) SetLeft(c *Parabola) {
	p.left = c
	c.Parent = p
}

func (p *Parabola) SetRight(c *Parabola) {
	p.right = c
	c.Parent = p
}

func (p *Parabola) GetLeft() *Parabola {
	return p.GetLeftParent().GetLeftChild()
}

func (p *Parabola) GetRight() *Parabola {
	return p.GetRightParent().GetRightChild()
}

func (p *Parabola) Print() {
	fmt.Printf("Parabola: %p\n", &p)
	fmt.Printf("  Site:   %v\n", p.Site)
	fmt.Printf("  IsLeaf: %v\n", p.IsLeaf)
	fmt.Printf("  Event:  %v\n", p.Event)
	fmt.Printf("  Edge:   %v\n", p.Edge)
	fmt.Printf("  Parent: %v\n", p.Parent)
	fmt.Printf("  Left:   %v\n", p.left)
	fmt.Printf("  Right:  %v\n", p.right)
}

func (p *Parabola) GetLeftParent() *Parabola {
	par := p.Parent
	plast := p
	for par.Left() == plast {
		if par.Parent == nil {
			return nil
		}
		plast = par
		par = par.Parent
	}
	return par
}

func (p *Parabola) GetRightParent() *Parabola {
	par := p.Parent
	plast := p
	for par.Right() == plast {
		if par.Parent == nil {
			return nil
		}
		plast = par
		par = par.Parent
	}
	return par
}

func (p *Parabola) GetLeftChild() *Parabola {
	if p == nil {
		return nil
	}
	par := p.Left()
	for !par.IsLeaf {
		par = par.Right()
	}
	return par
}

func (p *Parabola) GetRightChild() *Parabola {
	if p == nil {
		return nil
	}
	par := p.Right()
	for !par.IsLeaf {
		par = par.Left()
	}
	return par
}

type Voronoi struct {
	Edges    Edges
	Vertices Vertices
	Places   *Vertices
	Width    float64
	Height   float64
	Root     *Parabola
	Y        float64
	del      EventList
	points   Vertices
	queue    EventQueue
}

func (v *Voronoi) GetEdges(places *Vertices, w float64, h float64) Edges {
	v.Places = places
	v.Width = w
	v.Height = h
	v.Root = nil
	v.Edges = make(Edges, 0, 0)
	v.points = make(Vertices, 0, 0)

	v.queue = make(EventQueue, 0, len(*places)+1)
	for _, p := range *places {
		heap.Push(&v.queue, NewEvent(p, true))
	}

	v.del = make(EventList, 0, 0)
	var e *Event
	for v.queue.Len() > 0 {
		e = heap.Pop(&v.queue).(*Event)
		v.Y = e.Point.Y
		if i := v.del.Find(e); i != -1 {
			v.del.Remove(e)
			continue
		}
		if e.IsPlace {
			v.insertParabola(e.Point)
		} else {
			v.removeParabola(e)
		}
	}

	v.finishEdge(v.Root)

	for _, e := range v.Edges {
		if e.Neighbor != nil {
			e.Start = e.Neighbor.End
			e.Neighbor = nil
		}
	}
	return v.Edges
}

func (v *Voronoi) insertParabola(p *Point) {
	if v.Root == nil {
		v.Root = NewLeafParabola(p)
		return
	}
	if v.Root.IsLeaf && v.Root.Site.Y-p.Y < 1 {
		fp := v.Root.Site
		v.Root.IsLeaf = false
		v.Root.SetLeft(NewLeafParabola(fp))
		v.Root.SetRight(NewLeafParabola(p))
		s := Pt((p.X+fp.X)/2.0, v.Height)
		v.points = append(v.points, s)
		if p.X > fp.X {
			v.Root.Edge = NewEdge(s, fp, p)
		} else {
			v.Root.Edge = NewEdge(s, p, fp)
		}
		v.Edges = append(v.Edges, v.Root.Edge)
		return
	}

	par := v.getParabolaByX(p.X)
	if par.Event != nil {
		v.del = append(v.del, par.Event)
		par.Event = nil
	}

	start := Pt(p.X, v.getY(par.Site, p.X))
	v.points = append(v.points, start)

	el := NewEdge(start, par.Site, p)
	er := NewEdge(start, p, par.Site)

	el.Neighbor = er
	v.Edges = append(v.Edges, el)

	par.Edge = er
	par.IsLeaf = false

	p0 := NewLeafParabola(par.Site)
	p1 := NewLeafParabola(p)
	p2 := NewLeafParabola(par.Site)

	par.SetRight(p2)
	par.SetLeft(NewParabola())
	par.Left().Edge = el
	par.Left().SetLeft(p0)
	par.Left().SetRight(p1)

	v.checkCircle(p0)
	v.checkCircle(p2)
}

func (v *Voronoi) removeParabola(e *Event) {
	var (
		p1 = e.Arch
		xl = p1.GetLeftParent()
		xr = p1.GetRightParent()
		p0 = xl.GetLeftChild()
		p2 = xr.GetRightChild()
	)
	if p0.Event != nil {
		v.del = append(EventList{p0.Event}, v.del...)
		p0.Event = nil
	}
	if p2.Event != nil {
		v.del = append(EventList{p2.Event}, v.del...)
		p2.Event = nil
	}

	p := Pt(e.Point.X, v.getY(p1.Site, e.Point.X))
	v.points = append(v.points, p)

	xl.Edge.End = p
	xr.Edge.End = p

	var (
		higher *Parabola
		par    *Parabola = p1
	)
	for par != v.Root {
		par = par.Parent
		if par == xl {
			higher = xl
		}
		if par == xr {
			higher = xr
		}
	}

	higher.Edge = NewEdge(p, p0.Site, p2.Site)
	v.Edges = append(v.Edges, higher.Edge)

	gparent := p1.Parent.Parent
	if p1.Parent.Left() == p1 {
		if gparent.Left() == p1.Parent {
			gparent.SetLeft(p1.Parent.Right())
		}
		if gparent.Right() == p1.Parent {
			gparent.SetRight(p1.Parent.Right())
		}
	} else {
		if gparent.Left() == p1.Parent {
			gparent.SetLeft(p1.Parent.Left())
		}
		if gparent.Right() == p1.Parent {
			gparent.SetRight(p1.Parent.Left())
		}
	}

	p1.Parent = nil

	v.checkCircle(p0)
	v.checkCircle(p2)
}

func (v *Voronoi) getEdgeIntersection(a *Edge, b *Edge) *Point {
	var (
		x = (b.G - a.G) / (a.F - b.F)
		y = a.F*x + a.G
	)

	if math.IsInf(float64(b.F), 0) {
		x = b.Start.X
		y = a.F*x + a.G
	}
	if math.IsInf(float64(a.F), 0) {
		x = a.Start.X
		y = b.F*x + b.G
	}

	if (x-a.Start.X)/a.Direction.X < 0 {
		return nil
	}
	if (y-a.Start.Y)/a.Direction.Y < 0 {
		return nil
	}
	if (x-b.Start.X)/b.Direction.X < 0 {
		return nil
	}
	if (y-b.Start.Y)/b.Direction.Y < 0 {
		return nil
	}
	p := Pt(x, y)
	v.points = append(v.points, p)
	return p
}

func (v *Voronoi) checkCircle(b *Parabola) {
	var (
		lp = b.GetLeftParent()
		rp = b.GetRightParent()
		a  = lp.GetLeftChild()
		c  = rp.GetRightChild()
	)
	if a == nil || c == nil || a.Site == c.Site {
		return
	}
	s := v.getEdgeIntersection(lp.Edge, rp.Edge)
	if s == nil {
		return
	}
	var (
		dx = a.Site.X - s.X
		dy = a.Site.Y - s.Y
		d  = float64(math.Sqrt(float64((dx * dx) + (dy * dy))))
	)
	if s.Y-d >= v.Y {
		return
	}
	e := NewEvent(Pt(s.X, s.Y-d), false)
	v.points = append(v.points, e.Point)
	b.Event = e
	e.Arch = b
	heap.Push(&v.queue, e)
}

func (v *Voronoi) getParabolaByX(xx float64) *Parabola {
	par := v.Root
	var x float64 = 0.0
	for !par.IsLeaf {
		x = v.getXOfEdge(par, v.Y)
		if x > xx {
			par = par.Left()
		} else {
			par = par.Right()
		}
	}
	return par
}

func (v *Voronoi) getY(p *Point, x float64) float64 {
	var (
		dp = 2 * (p.Y - v.Y)
		a1 = 1 / dp
		b1 = -2 * p.X / dp
		c1 = v.Y + dp/4 + p.X*p.X/dp
	)
	return a1*x*x + b1*x + c1
}

func (v *Voronoi) finishEdge(n *Parabola) {
	if n.IsLeaf {
		return
	}
	var mx float64
	if n.Edge.Direction.X > 0.0 {
		if v.Width > n.Edge.Start.X+10 {
			mx = v.Width
		} else {
			mx = n.Edge.Start.X + 10
		}
	} else {
		if 0.0 < n.Edge.Start.X-10 {
			mx = 0.0
		} else {
			mx = n.Edge.Start.X - 10
		}
	}
	var end *Point
	if math.IsInf(float64(n.Edge.F), 1) {
		end = Pt(mx, v.Height)
	} else if math.IsInf(float64(n.Edge.F), -1) {
		end = Pt(mx, 0)
	} else {
		end = Pt(mx, mx*n.Edge.F+n.Edge.G)
	}
	n.Edge.End = end
	v.points = append(v.points, end)
	v.finishEdge(n.Left())
	v.finishEdge(n.Right())
}

func (v *Voronoi) getXOfEdge(par *Parabola, y float64) float64 {
	var (
		left  = par.GetLeftChild()
		right = par.GetRightChild()
		p     = left.Site
		r     = right.Site
		dp    = 2.0 * (p.Y - y)
		a1    = 1.0 / dp
		b1    = -2.0 * p.X / dp
		c1    = y + dp/4 + p.X*p.X/dp
	)
	dp = 2.0 * (r.Y - y)
	var (
		a2   = 1.0 / dp
		b2   = -2.0 * r.X / dp
		c2   = v.Y + dp/4 + r.X*r.X/dp
		a    = a1 - a2
		b    = b1 - b2
		c    = c1 - c2
		disc = b*b - 4*a*c
		x1   = (-b + float64(math.Sqrt(float64(disc)))) / (2 * a)
		x2   = (-b - float64(math.Sqrt(float64(disc)))) / (2 * a)
	)
	var ry float64
	if p.Y < r.Y {
		if x1 > x2 {
			ry = x1
		} else {
			ry = x2
		}
	} else {
		if x1 < x2 {
			ry = x1
		} else {
			ry = x2
		}
	}
	return ry
}
