package main

import (
	"fmt"
	"time"

	"github.com/bojand/go-npmdl"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

const layout = "2006-01-02"

func drawNPMChart(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	rangeParam := "last-year"
	if ctx.QueryArgs().Has("range") {
		rangeParam = string(ctx.QueryArgs().Peek("range"))
	}

	fmt.Printf("name: %s range: %s\n", name, rangeParam)

	out, err := npmdl.GetRangeCounts(rangeParam, name)
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
		return
	}

	n := len(out.Downloads)
	fmt.Printf("COUNT: %d\n", n)

	xValues := make([]time.Time, n)
	yValues := make([]float64, n)

	for i, dl := range out.Downloads {
		xValues[i], _ = time.Parse(layout, dl.Day)
		yValues[i] = float64(dl.Downloads)
	}

	graph := chart.Chart{
		Width:  800,
		Height: 250,
		XAxis: chart.XAxis{
			Name:      "Time",
			NameStyle: chart.StyleShow(),
			Style: chart.Style{
				Show:        true,
				StrokeWidth: 1,
				StrokeColor: drawing.Color{
					R: 85,
					G: 85,
					B: 85,
					A: 180,
				},
			},
		},
		YAxis: chart.YAxis{
			Name:      "Downloads",
			NameStyle: chart.StyleShow(),
			Style: chart.Style{
				Show:        true,
				StrokeWidth: 1,
				StrokeColor: drawing.Color{
					R: 85,
					G: 85,
					B: 85,
					A: 180,
				},
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name: name,
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.ColorRed,
					FillColor:   chart.ColorRed.WithAlpha(16),
				},
				XValues: xValues,
				YValues: yValues,
			},
		},
	}

	//note we have to do this as a separate step because we need a reference to graph
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	ctx.SetContentType("image/svg+xml;charset=utf-8")
	graph.Render(chart.SVG, ctx)

}

func main() {
	router := fasthttprouter.New()
	router.GET("/:name", drawNPMChart)

	fasthttp.ListenAndServe(":8080", router.Handler)
}
