package view

import (
	"fmt"
	"github.com/gophergala/golab/model"
	"html/template"
	"image"
	"image/jpeg"
	"net/http"
	"strconv"
	"time"
)

var params = struct {
	Title         string
	Width, Height int
	RunId         int64
}{AppTitle, ViewWidth, ViewHeight, time.Now().Unix()}

var playTempl = template.Must(template.New("t").Parse(play_html))

// The client's (browser's) view position inside the labyrinth image.
var Pos image.Point

// init registers the http handlers.
func init() {
	http.HandleFunc("/", playHtmlHandle)
	http.HandleFunc("/runid", runIdHandle)
	http.HandleFunc("/img", imgHandle)
	http.HandleFunc("/clicked", clickedHandle)
	http.HandleFunc("/cheat", cheatHandle)
	http.HandleFunc("/new", newGameHandle)
}

// InitNew initializes a new game.
func InitNew() {
	Pos = image.Point{}
}

// playHtmlHandle serves the html page where the user can play.
func playHtmlHandle(w http.ResponseWriter, r *http.Request) {
	playTempl.Execute(w, params)
}

// runidHandle serves the running app id which changes if app is restarted
// (so browser clients can detect if app was restarted).
func runIdHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", params.RunId)
}

// imgHandle serves images of the player's view.
func imgHandle(w http.ResponseWriter, r *http.Request) {
	quality, err := strconv.Atoi(r.FormValue("quality"))
	if err != nil || quality < 0 || quality > 100 {
		quality = 70
	}
	
	rect := image.Rect(0, 0, params.Width, params.Height)
	jpeg.Encode(w, model.LabImg.SubImage(rect), &jpeg.Options{quality})
}

// clickedHandle receives mouse click (mouse button pressed) events with mouse coordinates.
func clickedHandle(w http.ResponseWriter, r *http.Request) {
	x, err := strconv.Atoi(r.FormValue("x"))
	if err != nil {
		return
	}
	y, err := strconv.Atoi(r.FormValue("y"))
	if err != nil {
		return
	}

	fmt.Println("Clicked:", x, y)
	model.TargetPos.X, model.TargetPos.Y = x, y
}

// cheatHandle serves the whole image of the labyrinth
func cheatHandle(w http.ResponseWriter, r *http.Request) {
	jpeg.Encode(w, model.LabImg, &jpeg.Options{70})
}

// newGameHandle signals to start a newgame
func newGameHandle(w http.ResponseWriter, r *http.Request) {
	// Use non-blocking send
	select {
	case model.NewGameCh <- 1:
	default:
	}
}
