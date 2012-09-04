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
	"testing"
)

func TestEventQueue(t *testing.T) {
	queue := make(EventQueue, 0, 4)
	heap.Push(&queue, &Event{Y: 5})
	heap.Push(&queue, &Event{Y: 3})
	heap.Push(&queue, &Event{Y: 7})
	heap.Push(&queue, &Event{Y: 1})

	var e *Event
	e = heap.Pop(&queue).(*Event)
	if e.Y != 1 {
		t.Fatalf("Wanted priority 1, got %v", e.Y)
	}
	e = heap.Pop(&queue).(*Event)
	if e.Y != 3 {
		t.Fatalf("Wanted priority 3, got %v", e.Y)
	}
	e = heap.Pop(&queue).(*Event)
	if e.Y != 5 {
		t.Fatalf("Wanted priority 5, got %v", e.Y)
	}
	e = heap.Pop(&queue).(*Event)
	if e.Y != 7 {
		t.Fatalf("Wanted priority 7, got %v", e.Y)
	}
}

func TestEventList(t *testing.T) {
	list := make(EventList, 0, 0)
	var (
		e1 = &Event{Y: 1}
		e2 = &Event{Y: 2}
		e3 = &Event{Y: 3}
	)
	list = append(list, e1)
	list = append(list, e2)
	if ret := list.Find(e3); ret != -1 {
		t.Fatalf("Expected -1, got %v", ret)
	}
	if ret := list.Find(e2); ret != 1 {
		t.Fatalf("Expected 1, got %v", ret)
	}
	if list.Last() != e2 {
		t.Fatalf("Last element was not e2")
	}
	if list.Remove(e3) {
		t.Fatalf("Removing nonexistent item should return false")
	}
	if !list.Remove(e2) {
		t.Fatalf("Removing existent item should return true")
	}
	if list.Last() != e1 {
		t.Fatalf("Removing e2 should update last element")
	}
	list.Remove(e1)
	if list.Last() != nil {
		t.Fatalf("Empty list should not have a last")
	}
}

func TestGetEdges(t *testing.T) {
	v := Voronoi{}
	ver := &Vertices{
		Pt(1,1),
		Pt(2,2),
		Pt(5,1),
	}
	e := v.GetEdges(ver, 10, 10)
	t.Logf("Edges: %v", e)
	t.Fatal("Not implemented")
}
