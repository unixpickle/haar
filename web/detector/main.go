package main

import (
	"encoding/json"

	"github.com/gopherjs/gopherjs/js"
	"github.com/unixpickle/haar"
)

const overlapThreshold = 0.3

var cascade *haar.Cascade

func main() {
	js.Global.Set("onmessage", js.MakeFunc(messageHandler))
}

func messageHandler(this *js.Object, dataArg []*js.Object) interface{} {
	if len(dataArg) != 1 {
		panic("expected one argument")
	}
	data := dataArg[0].Get("data")

	if cascade == nil {
		cascade = new(haar.Cascade)
		jsonData := js.Global.Get("JSON").Call("stringify", data)
		json.Unmarshal([]byte(jsonData.String()), cascade)
	} else {
		width := data.Index(0).Int()
		height := data.Index(1).Int()
		buffer := data.Index(2)
		recognized := recognize(width, height, buffer)
		js.Global.Call("postMessage", recognized)
	}

	return nil
}

func recognize(width, height int, buffer *js.Object) interface{} {
	bitmap := make([]float64, buffer.Length())
	for i := range bitmap {
		bitmap[i] = buffer.Index(i).Float()
	}
	img := haar.BitmapIntegralImage(bitmap, width, height)
	dualImg := haar.NewDualImage(img)

	matches := cascade.Scan(dualImg, 0, 1.5).JoinOverlaps(overlapThreshold)
	data, _ := json.Marshal(matches)
	return js.Global.Get("JSON").Call("parse", string(data))
}
