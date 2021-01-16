package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/AllenDang/giu/imgui"
	g "github.com/ianling/giu"
)

var (
	linedata     []float64
	linedata2    []float64
	lineTicks    []g.PlotTicker
	bardata      []float64
	bardata2     []float64
	bardata3     []float64
	timeDataMin  float64
	timeDataMax  float64
	timeDataX    []float64
	timeDataY    []float64
	timeScatterY []float64
	scatterdata  []float64
)

func loop() {
	var plotdata []float32
	const delta = 0.01
	for x := 0.0; len(plotdata) < 1000; x += delta {
		plotdata = append(plotdata, float32(math.Sin(x)))
	}
	g.SingleWindow("hello world").Layout(
		g.Label("Hello world from giu"),
		g.Label("Simple sin(x) plot:"),
		g.PlotLines("testplot", plotdata),
		g.Label("sin(x) plot with overlay text, and size:"),
		g.PlotLinesV("plot label", plotdata, 0, "overlay text", math.MaxFloat32, math.MaxFloat32, 500, 200),
	).Build()
}

func main() {
	delta := 0.1
	for x := 0.0; x < 10; x += delta {
		linedata = append(linedata, math.Sin(x))
		linedata2 = append(linedata2, math.Cos(x))
		scatterdata = append(scatterdata, math.Sin(x)+0.1)
	}

	for i := 0; i < 100; i += 5 {
		lineTicks = append(lineTicks, g.PlotTicker{Position: float64(i), Label: fmt.Sprintf("P%d", i)})
	}

	delta = 1
	for x := 0.0; x < 10; x += delta {
		bardata = append(bardata, math.Sin(x))
		bardata2 = append(bardata2, math.Sin(x)-0.2)
		bardata3 = append(bardata3, rand.Float64())
	}

	for i := 0; i < 100; i++ {
		timeDataX = append(timeDataX, float64(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Hour*time.Duration(24*i)).Unix()))
		timeDataY = append(timeDataY, rand.Float64())
		timeScatterY = append(timeScatterY, rand.Float64())
	}

	timeDataMin = timeDataX[0]
	timeDataMax = timeDataX[len(timeDataX)-1]

	wnd := g.NewMasterWindow("Plot Demo", 1000, 900, 0, nil)
	wnd.Run(loop)
}
