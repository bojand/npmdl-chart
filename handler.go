package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	npmdl "github.com/bojand/go-npmdl"
	"github.com/valyala/fasthttp"
	chart "github.com/wcharczuk/go-chart"
)

const dateLayout = "2006-01-02"
const aspectRatio = 3.2

type TemplateData struct {
	Name string
}

// getChartDimensions gets the chart dimentions based on the w query param
func getChartDimensions(ctx *fasthttp.RequestCtx) (w int, h int) {
	width := 800
	height := 250

	if ctx.QueryArgs().Has("w") {
		wStr := string(ctx.QueryArgs().Peek("w"))
		w, err := strconv.Atoi(wStr)
		if err == nil && w > 0 {
			width = w
		}

		hasH := ctx.QueryArgs().Has("h")

		if hasH {
			hStr := string(ctx.QueryArgs().Peek("h"))
			h, err := strconv.Atoi(hStr)
			if err == nil && h > 0 {
				height = h
			}
		}

		if !hasH {
			height = int(float64(width) / aspectRatio)
		}
	}

	return width, height
}

// getPackageNameAndChartType gets package name and chart type based on name path param
func getPackageNameAndChartType(ctx *fasthttp.RequestCtx) (name string, imageType string) {
	nameParam := strings.ToLower(ctx.UserValue("name").(string))
	fmt.Println("nameParam: ", nameParam)
	baseName := filepath.Base(nameParam)
	fmt.Println("baseName: ", baseName)
	ext := filepath.Ext(baseName)
	fmt.Println("ext: ", ext)
	if ext == "" {
		return baseName, "svg"
	}

	imgType := strings.TrimPrefix(ext, ".")
	if imgType != "svg" && imgType != "png" {
		imgType = "svg"
	}

	pkg := strings.TrimSuffix(baseName, ext)
	return pkg, imgType
}

// DrawNPMChart is the handler for the request
func DrawNPMChart(ctx *fasthttp.RequestCtx) {
	name, imgType := getPackageNameAndChartType(ctx)
	width, height := getChartDimensions(ctx)

	rangeParam := "last-year"
	if ctx.QueryArgs().Has("range") {
		rangeParam = strings.ToLower(string(ctx.QueryArgs().Peek("range")))
	}

	fmt.Printf("name: %s range: %s imgType: %s\n", name, rangeParam, imgType)

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
		xValues[i], _ = time.Parse(dateLayout, dl.Day)
		yValues[i] = float64(dl.Downloads)
	}

	graph := CreateNPMChart(name, xValues, yValues, width, height)

	ctx.Response.Header.Add("cache-control", "no-cache, no-store, must-revalidate")
	ctx.Response.Header.Add("date", time.Now().Format(time.RFC1123))
	ctx.Response.Header.Add("expires", time.Now().Format(time.RFC1123))

	if imgType == "png" {
		ctx.Response.Header.Add("content-type", "image/png")
		graph.Render(chart.PNG, ctx)
	} else {
		ctx.Response.Header.Add("content-type", "image/svg+xml;charset=utf-8")
		graph.Render(chart.SVG, ctx)
	}
}

// Index serves index.html
func Index(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf-8")
	template.Must(template.ParseFiles("templates/index.html")).Execute(ctx, nil)
}

func GetNPMChart(ctx *fasthttp.RequestCtx) {
	name := strings.ToLower(ctx.UserValue("name").(string))
	data := TemplateData{Name: name}
	ctx.SetContentType("text/html; charset=utf-8")
	template.Must(template.ParseFiles("templates/index.html")).Execute(ctx, data)
}
