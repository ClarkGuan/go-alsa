package pcm

import "C"

func fromBool(b bool) C.int {
	if b {
		return 1
	} else {
		return 0
	}
}
