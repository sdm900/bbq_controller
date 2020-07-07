package outputs

import (
	"fmt"
	"os"
	"runtime"
)

var Debug bool

func msg(a ...interface{}) {
	switch len(a) {
	case 0:
		return
	case 1:
		fmt.Fprint(os.Stderr, a[0], " ")
	default:
		for _, v := range a {
			msg(v)
		}
	}
}

// TODO: use standard logging stuff so our server uses syslog

func Msg(a ...interface{}) {
	msg(a...)
	fmt.Fprint(os.Stderr, "\n")
}

func Err(a ...interface{}) {
	msg(a...)
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

func dbg(str string, a ...interface{}) {
	if os.Getenv("DEBUG") == "sdm900" {
		msg(str, " ")

		pc, file, line, ok := runtime.Caller(2)
		if ok {
			s := runtime.FuncForPC(pc).Name()
			msg(a...)
			Msg(" : ", s, "        ", file, " ", line)
		} else {
			Msg(a...)
		}
	}
}

func In(a ...interface{}) {
	dbg(" IN", a...)
}

func Out(a ...interface{}) {
	dbg("OUT", a...)
}

func Dbg(a ...interface{}) {
	dbg("DBG", a...)
}
