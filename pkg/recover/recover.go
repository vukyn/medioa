package recover

import (
	"fmt"

	"github.com/vukyn/kuery/log"
)

func RecoverPanic() {
	log := log.New("routine", "recoverPanic")
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		log.Error("", err)
	}
}
