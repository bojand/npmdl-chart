package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	npmdl "github.com/bojand/go-npmdl"
	"github.com/valyala/fasthttp"
	chart "github.com/wcharczuk/go-chart"
)

const layout = "2006-01-02"
const ratio float64 = 3.2

// GetChartDimensions gets the chart dimentions based on the w query param
func GetChartDimensions(ctx *fasthttp.RequestCtx) (w int, h int, e error) {
	width := 800
	height := 250

	if ctx.QueryArgs().Has("w") {
		wStr := string(ctx.QueryArgs().Peek("w"))
		w, err := strconv.Atoi(wStr)
		if err != nil {
			return 0, 0, err
		}

		if w > 0 {
			width = w
			height = int(float64(width) / ratio)
		}
	}

	return width, height, nil
}

// DrawNPMChart is the handler for the request
func DrawNPMChart(ctx *fasthttp.RequestCtx) {
	name := strings.ToLower(ctx.UserValue("name").(string))

	rangeParam := "last-year"
	if ctx.QueryArgs().Has("range") {
		rangeParam = strings.ToLower(string(ctx.QueryArgs().Peek("range")))
	}

	width, height, err := GetChartDimensions(ctx)
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
		return
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

	graph := CreateNPMChart(name, xValues, yValues, width, height)

	ctx.Response.Header.Add("content-type", "image/svg+xml;charset=utf-8")
	ctx.Response.Header.Add("cache-control", "no-cache, no-store, must-revalidate")
	ctx.Response.Header.Add("date", time.Now().Format(time.RFC1123))
	ctx.Response.Header.Add("expires", time.Now().Format(time.RFC1123))

	graph.Render(chart.SVG, ctx)
}
