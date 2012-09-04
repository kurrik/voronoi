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

package voronoi

import (
	"container/heap"
	"fmt"
)

type Point struct {
	X float32
	Y float32
}

func Pt(x float32, y float32) *Point {
	return &Point{X: x, Y: y}
}

type Vertices []*Point

type Edge struct {
	Start     *Point
	End       *Point
	Direction *Point
	Left      *Point
	Right     *Point
	F         float32
	G         float32
	Neighbor  *Edge
}

type Edges []*Edge

type Event struct {
	Point   *Point
	IsPlace bool
	Y       float32
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
	// Sorted by Y ascending.
	return q[i].Y < q[j].Y
}

func (q EventQueue) Swap(i int, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *EventQueue) Push(x interface{}) {
	a := *q
	n := len(a)
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

type Voronoi struct {
	Edges    Edges
	Vertices Vertices
	Places   *Vertices
	Width    float32
	Height   float32
	Root     *Parabola
	Y        float32
	del      EventList
	points   Vertices
	queue    EventQueue
}

func (v *Voronoi) GetEdges(places *Vertices, w float32, h float32) Edges {
	v.Places = places
	v.Width = w
	v.Height = h
	v.Root = nil
	v.Edges = make(Edges, 0, 0)
	v.points = make(Vertices, 0, 0)

	v.queue = make(EventQueue, 0, len(*places))
	for _, p := range *places {
		heap.Push(&v.queue, NewEvent(p, true))
	}

	v.del = make(EventList, 0, 0)
	var e *Event
	for len(v.queue) > 0 {
		e = heap.Pop(&v.queue).(*Event)
		v.Y = e.Point.Y
		if i := v.del.Find(e); i != -1 && v.del[i] != v.del.Last() {
			v.del.Remove(e)
			continue
		}
		if e.IsPlace {
			v.insertParabola(e.Point)
		} else {
			v.removeParabola(e)
		}
		fmt.Println(v.Y)
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
}

func (v *Voronoi) removeParabola(e *Event) {
}

func (v *Voronoi) finishEdge(p *Parabola) {
}
