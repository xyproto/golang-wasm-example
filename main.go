// +build js,wasm

package main

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js .

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"syscall/js"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func minus256toplus256() int {
	return rand.Intn(512) - 256
}

func loadImage(path string) string {
	href := js.Global().Get("location").Get("href")
	u, err := url.Parse(href.String())
	if err != nil {
		log.Fatal(err)
	}
	u.Path = path
	u.RawQuery = fmt.Sprint(time.Now().UnixNano())

	log.Println("loading image file: " + u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func main() {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")
	canvas.Set("width", js.ValueOf(500))
	canvas.Set("height", js.ValueOf(500))

	images := make([]js.Value, 4)
	files := []string{
		"/data/out01.png",
		"/data/out02.png",
		"/data/out03.png",
		"/data/out02.png",
	}
	for i, file := range files {
		images[i] = js.Global().Call("eval", "new Image()")
		images[i].Set("src", "data:image/png;base64,"+loadImage(file))
	}

	//canvas.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	//	js.Global().Get("window").Call("alert", "Don't click me!")
	//	return nil
	//}))

	bgColors := []string{"red", "yellow", "blue", "green", "white"}
	bgColorIndex := 0

	n := 0
	js.Global().Call("setInterval", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ctx.Call("clearRect", 0, 0, 500, 500)
		ctx.Call("drawImage", images[n%4], 0, 0)
		n++

		style := canvas.Get("style")
		left := style.Get("left")
		if left == js.Undefined() {
			left = js.ValueOf("320px")
		} else {
			//n, _ := strconv.Atoi(strings.TrimRight(left.String(), "px"))
			left = js.ValueOf(fmt.Sprintf("%dpx", 320+minus256toplus256()))
		}
		style.Set("left", left)

		bgColor := bgColors[bgColorIndex]
		bgColorIndex++
		if bgColorIndex >= len(bgColors) {
			bgColorIndex = 0
		}
		style.Set("background-color", bgColor)
		body := document.Get("documentElement")
		body.Call("setAttribute", "style", "background-color: " + bgColor)

		return nil
	}), js.ValueOf(50))

	select {}
}
