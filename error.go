package alsa

// #cgo LDFLAGS: -lasound
//
// #include <alsa/asoundlib.h>
//
// static char *_snd_strerror(int errnum) {
//     return (char *) snd_strerror(errnum);
// }
import "C"
import "fmt"

const (
	EAGAIN   = int(C.EAGAIN)
	ENOSYS   = int(C.ENOSYS)
	EBADFD   = int(C.EBADFD)
	EPIPE    = int(C.EPIPE)
	ESTRPIPE = int(C.ESTRPIPE)
)

type Error struct {
	Errno int
}

func NewError(errno int) *Error {
	return &Error{Errno: errno}
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%d]%s", err.Errno, C.GoString(C._snd_strerror(C.int(err.Errno))))
}
