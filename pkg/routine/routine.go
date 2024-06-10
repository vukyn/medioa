package routine

import (
	"fmt"
	"medioa/pkg/log"
)

func Run(fn func()) {
	go func() {
		defer recoverPanic()
		fn()
	}()
}

func recoverPanic() {
	log := log.New("routine", "recoverPanic")
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		log.Error("", err)
	}
}
