// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"math"

	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// Align returns a two-dimensional row-major array of Canvases which will
// produce tiled plots with DataCanvases that are evenly sized and spaced.
// The arguments to the function are a two-dimensional row-major array
// of plots, a tile configuration, and the canvas to which the tiled
// plots are to be drawn.
func Align(plots [][]*Plot, t draw.Tiles, dc draw.Canvas) [][]draw.Canvas {
	o := make([][]draw.Canvas, len(plots))

	// Create the initial tiles.
	for j := 0; j < t.Rows; j++ {
		o[j] = make([]draw.Canvas, len(plots[j]))
		for i := 0; i < t.Cols; i++ {
			o[j][i] = t.At(dc, i, j)
		}
	}

	type posNeg struct {
		p, n float64
	}
	xSpacing := make([]posNeg, t.Cols)
	ySpacing := make([]posNeg, t.Rows)

	// Calculate the maximum spacing between data canvases
	// for each row and column.
	for j, row := range plots {
		for i, p := range row {
			c := o[j][i]
			dataC := p.DataCanvas(o[j][i])
			xSpacing[i].n = math.Max(float64(dataC.Min.X-c.Min.X), xSpacing[i].n)
			xSpacing[i].p = math.Max(float64(c.Max.X-dataC.Max.X), xSpacing[i].p)
			ySpacing[j].n = math.Max(float64(dataC.Min.Y-c.Min.Y), ySpacing[j].n)
			ySpacing[j].p = math.Max(float64(c.Max.Y-dataC.Max.Y), ySpacing[j].p)
		}
	}

	// Adjust the horizontal and vertical spacing between
	// canvases to match the maximum for each column and row,
	// respectively.
	for j, row := range plots {
		for i, p := range row {
			c := o[j][i]
			dataC := p.DataCanvas(o[j][i])
			o[j][i] = draw.Crop(c,
				vg.Length(xSpacing[i].n)-(dataC.Min.X-c.Min.X),
				c.Max.X-dataC.Max.X-vg.Length(xSpacing[i].p),
				vg.Length(ySpacing[j].n)-(dataC.Min.Y-c.Min.Y),
				c.Max.Y-dataC.Max.Y-vg.Length(ySpacing[j].p),
			)
		}
	}

	// Calculate the total row and column spacing.
	var xTotalSpace float64
	for _, s := range xSpacing {
		xTotalSpace += s.n + s.p
	}
	var yTotalSpace float64
	for _, s := range ySpacing {
		yTotalSpace += s.n + s.p
	}

	avgWidth := vg.Length((float64(dc.Max.X-dc.Min.X) - xTotalSpace) / float64(t.Cols))
	avgHeight := vg.Length((float64(dc.Max.Y-dc.Min.Y) - yTotalSpace) / float64(t.Rows))

	// Adjust canvases so the width of each DataCanvas is the
	// same.
	for j := 0; j < t.Rows; j++ {
		var moveHorizontal vg.Length
		for i := 0; i < t.Cols; i++ {
			c := o[j][i]
			dataC := plots[j][i].DataCanvas(o[j][i])
			width := dataC.Max.X - dataC.Min.X
			o[j][i] = draw.Crop(c,
				moveHorizontal,
				moveHorizontal+avgWidth-width,
				0,
				0,
			)
			moveHorizontal += avgWidth - width
		}
	}

	// Adjust canvases so the height of each DataCanvas
	// is the same.
	for i := 0; i < t.Cols; i++ {
		var moveVertical vg.Length
		for j := t.Rows - 1; j >= 0; j-- {
			c := o[j][i]
			dataC := plots[j][i].DataCanvas(o[j][i])
			height := dataC.Max.Y - dataC.Min.Y
			o[j][i] = draw.Crop(c,
				0,
				0,
				moveVertical,
				moveVertical+avgHeight-height,
			)
			moveVertical += avgHeight - height
		}
	}

	return o
}
