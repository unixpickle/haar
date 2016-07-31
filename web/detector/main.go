package main

import (
	"encoding/json"

	"github.com/gopherjs/gopherjs/js"
	"github.com/unixpickle/haar"
)

const overlapThreshold = 0.3

var cascade haar.Cascade

func main() {
	cascadeObj := js.Global.Get("app").Get("cascade")
	jsonData := js.Global.Get("JSON").Call("stringify", cascadeObj)
	json.Unmarshal([]byte(jsonData.String()), &cascade)
	js.Global.Get("app").Set("detect", js.MakeFunc(recognize))
}

func recognize(this *js.Object, dataArg []*js.Object) interface{} {
	width := dataArg[0].Int()
	height := dataArg[1].Int()
	buffer := dataArg[2]

	bitmap := make([]float64, buffer.Length())
	for i := range bitmap {
		bitmap[i] = buffer.Index(i).Float()
	}
	img := haar.BitmapIntegralImage(bitmap, width, height)
	dualImg := haar.NewDualImage(img)

	matches := cascade.Scan(dualImg, 0, 1.5).JoinOverlaps(2)
	data, _ := json.Marshal(matches)
	return js.Global.Get("JSON").Call("parse", string(data))
}
