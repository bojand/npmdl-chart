package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	npmdl "github.com/bojand/go-npmdl"
	"github.com/go-chi/chi"
	chart "github.com/wcharczuk/go-chart"
)

const dateLayout = "2006-01-02"
const aspectRatio = 3.2

// TemplateData is used for filling the HTML template
type TemplateData struct {
	Name string
}

// getChartDimensions gets the chart dimentions based on the w query param
func getChartDimensions(req *http.Request) (w int, h int) {
	width := 800
	height := 250

	wStr := req.URL.Query().Get("w")

	if wStr != "" {
		w, err := strconv.Atoi(wStr)
		if err == nil && w > 0 {
			width = w
		}

		hStr := req.URL.Query().Get("h")

		if hStr != "" {
			h, err := strconv.Atoi(hStr)
			if err == nil && h > 0 {
				height = h
			}
		}

		if hStr == "" {
			height = int(float64(width) / aspectRatio)
		}
	}

	return width, height
}

// getPackageNameAndChartType gets package name and chart type based on name path param
func getPackageNameAndChartType(req *http.Request) (name string, imageType string) {
	nameParam := strings.ToLower(chi.URLParam(req, "*"))
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
func DrawNPMChart(res http.ResponseWriter, req *http.Request) {
	name, imgType := getPackageNameAndChartType(req)
	width, height := getChartDimensions(req)

	rangeParam := req.URL.Query().Get("range")
	if rangeParam == "" {
		rangeParam = "last-year"
	}

	fmt.Printf("name: %s range: %s imgType: %s\n", name, rangeParam, imgType)

	out, err := npmdl.GetRangeCounts(rangeParam, name)
	if err != nil {
		http.Error(res, http.StatusText(400), 400)
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

	res.Header().Set("cache-control", "no-cache, no-store, must-revalidate")
	res.Header().Set("date", time.Now().Format(time.RFC1123))
	res.Header().Set("expires", time.Now().Format(time.RFC1123))

	if imgType == "png" {
		res.Header().Set("content-type", "image/png")
		graph.Render(chart.PNG, res)
	} else {
		res.Header().Set("content-type", "image/svg+xml;charset=utf-8")
		graph.Render(chart.SVG, res)
	}
}

// Index serves index.html
func Index(res http.ResponseWriter, r *http.Request) {
	res.Header().Set("content-type", "text/html; charset=utf-8")
	template.Must(template.ParseFiles("templates/index.html")).Execute(res, nil)
}

// GetNPMChart gets the page with the chart
func GetNPMChart(res http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(chi.URLParam(r, "name"))
	data := TemplateData{Name: name}
	res.Header().Set("content-type", "text/html; charset=utf-8")
	template.Must(template.ParseFiles("templates/index.html")).Execute(res, data)
}

// FileServer for serving files
func FileServer(r chi.Router, path string, fsPath string) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, fsPath)

	fs := http.StripPrefix(path, http.FileServer(http.Dir(filesDir)))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
