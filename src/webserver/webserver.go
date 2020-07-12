package webserver

import (
	"bbq"
	"bytes"
	"fmt"
	"net/http"
	"outputs"
	"strconv"
)

type screenHandler struct {
	b *bbq.BBQ
}

func (sh *screenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")

	buf := sh.b.ScreenPNG()
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(buf)))
	w.Write(buf)
}

type rootHandler struct {
	b      *bbq.BBQ
	screen string
}

func (rh *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	v := r.FormValue("setpoint")
	if v != "" {
		t, e := strconv.ParseFloat(v, 32)
		if e == nil {
			rh.b.SetT(float32(t))
		}
	}

	v = r.FormValue("display")
	if v == "On" {
		rh.b.ScreenOn()
		rh.screen = "Off"
	} else if v == "Off" {
		rh.b.ScreenOff()
		rh.screen = "On"
	}

	p := new(bytes.Buffer)
	p.WriteString(`
<html>
<head>
<title>BBQ Temperature Control</title>
<meta http-equiv="refresh" content="10">
</head>

<body style="background-color:black ; color:#cccccc ;">
<center>
<img src="/png" style="width:50% ; height:50% ; image-rendering:crisp-edges ; object-fit:contain ;"></td></tr>
<p>
<hr>
<p>
<form action="/" method="post">
<label for="setpoint" style="font-size:200%;">Set point:</label> <input type="text" id="setpoint" name="setpoint" value="` +
		fmt.Sprintf("%.1f", rh.b.GetT()) +
		`" style="font-size:200%;"> <input type="submit" value="Set" style="font-size:200%;">
</form>

<form action="/" method="post">
<label for="display" style="font-size:200%;">Display:</label> <input type="submit" name="display" value="` +
		rh.screen +
		`" style="font-size:200%;">
</form>
</body>
</html>
`)

	w.Header().Set("Content-Length", fmt.Sprintf("%d", p.Len()))
	w.Write(p.Bytes())

}

func Serve(b *bbq.BBQ) {
	http.Handle("/png", &screenHandler{b})
	http.Handle("/", &rootHandler{b, "Off"})

	if e := http.ListenAndServe(":80", nil); e != nil {
		outputs.Err("Can not serve the screen:", e)
	}
}
