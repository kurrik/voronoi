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

// Outputs an SVG image of random points and the corresponding Voronoi edges.
// Run with:
//   go run examples/basic.go > ~/Desktop/test.svg

package main

import (
	".."
	"fmt"
	"math/rand"
	"os"
	"text/template"
)

const TEMPLATE = `<?xml version="1.0" ?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN"
  "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="{{.Width}}px" height="{{.Height}}px" viewBox="0 0 {{.Width}} {{.Height}}"
     xmlns="http://www.w3.org/2000/svg" version="1.1">
  <title>{{.Title}}</title>
  <desc>{{.Description}}</desc>
  <!-- Edges -->
  <g stroke="red" stroke-width="{{.StrokeWidth}}" fill="none">
    {{range .Edges}}<path d="M{{.Start.X}},{{.Start.Y}} L{{.End.X}},{{.End.Y}}" />
    {{end}}</g>
  <!-- Vertices -->
  <g fill="black" >
    {{range .Vertices}}<circle cx="{{.X}}" cy="{{.Y}}" r="{{$.PointRadius}}" />
    {{end}}</g>
</svg>`

type SVG struct {
	Width       float64
	Height      float64
	Edges       voronoi.Edges
	Vertices    voronoi.Vertices
	Title       string
	Description string
	StrokeWidth float64
	PointRadius float64
}

func main() {
	pts := 600
	vor := voronoi.Voronoi{}
	svg := SVG{
		Title:       "Voronoi diagram",
		Description: "Edges and points",
		Width:       500,
		Height:      500,
		StrokeWidth: 1,
		PointRadius: 1,
		Vertices:    make([]*voronoi.Point, pts),
	}
	rnd := rand.New(rand.NewSource(7584))
	for i := 0; i < pts; i++ {
		var (
			x = rnd.Float64() * svg.Width
			y = rnd.Float64() * svg.Height
		)
		str := fmt.Sprintf("Point at %v,%v\n", x, y)
		os.Stderr.Write([]byte(str))
		svg.Vertices[i] = voronoi.Pt(x, y)
	}
	svg.Edges = vor.GetEdges(&svg.Vertices, svg.Width, svg.Height)
	tmpl := template.Must(template.New("svg").Parse(TEMPLATE))
	if err := tmpl.Execute(os.Stdout, svg); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
