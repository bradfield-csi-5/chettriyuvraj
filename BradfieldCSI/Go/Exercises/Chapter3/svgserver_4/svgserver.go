/* Following the approach of the Lissajous example in Section 1.7, construct a web server that computes surfaces and writes
SVG data to the client.  The server must set the Content-Type header like this:

w.Header().Set("Content-Type", "image/svg+xml")

(This step was not required in the Lissajous example because the server uses standard heuristics to
recognize common formats like PNG from the first 512 bytes of the response and generates the proper header.)
Allow the client to specify values like height, width, and color as HTTP request parameters.
*/

package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", baseHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	// set default values if not provided by parameters
	queryMap := r.URL.Query()
	heightStr, widthStr, color := queryMap.Get("height"), queryMap.Get("width"), queryMap.Get("color")
	if heightStr == "" {
		heightStr = "320"
	}
	if widthStr == "" {
		widthStr = "600"
	}
	if color == "" {
		color = "red"
	}
	height, _ := strconv.ParseFloat(heightStr, 64)
	width, _ := strconv.ParseFloat(widthStr, 64)

	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(200)
	servesvg(w, height, width, color)
}

const (
	cells = 100 // number of grid cells
)

func servesvg(out io.Writer, height float64, width float64, color string) {
	fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%f' height='%f'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j, width, height)
			bx, by := corner(i, j, width, height)
			cx, cy := corner(i, j+1, width, height)
			dx, dy := corner(i+1, j+1, width, height)
			fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintf(out, "</svg>")
}

func corner(i int, j int, width float64, height float64) (float64, float64) {
	const (
		xyrange = 30.0        // axis ranges (-xyrange..+xyrange)
		angle   = math.Pi / 6 // angle of x, y axes (=30°)
	)
	sin30, cos30 := math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)
	xyscale := width / 2 / xyrange                   // pixels per x or y unit
	zscale := height * 0.4                           // pixels per z unit

	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	val := math.Sin(r) / r
	return val
}
