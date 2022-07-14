package pcm

//
// #include <alsa/asoundlib.h>
//
// static char *no_const(const char *s) { return (char *)s; }
// static void _snd_pcm_info_alloca(snd_pcm_info_t **ptr) { snd_pcm_info_alloca(ptr); }
//
import "C"
import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/ClarkGuan/go-alsa"
)

func (pcm *PCM) Info() (*Info, error) {
	info := new(Info)
	C._snd_pcm_info_alloca(&info.inner)
	rc := C.snd_pcm_info(pcm.inner, info.inner)
	if rc < 0 {
		return nil, alsa.NewError(int(rc))
	}
	runtime.SetFinalizer(info, (*Info).Close)
	return info, nil
}

type Info struct {
	inner *C.snd_pcm_info_t
}

func (info *Info) Close() error {
	for {
		p := unsafe.Pointer(info.inner)
		if p == nil {
			break
		}
		i := info.inner
		if atomic.CompareAndSwapPointer(&p, p, nil) {
			C.snd_pcm_info_free(i)
			runtime.SetFinalizer(info, nil)
		}
	}
	return nil
}

func (info *Info) GetDevice() int {
	return int(C.snd_pcm_info_get_device(info.inner))
}

func (info *Info) GetSubDevice() int {
	return int(C.snd_pcm_info_get_subdevice(info.inner))
}

func (info *Info) GetStreamType() StreamType {
	return StreamType(C.snd_pcm_info_get_stream(info.inner))
}

func (info *Info) GetCard() (int, error) {
	rc := C.snd_pcm_info_get_card(info.inner)
	if rc < 0 {
		return -1, alsa.NewError(int(rc))
	}
	return int(rc), nil
}

func (info *Info) GetID() string {
	return C.GoString(C.no_const(C.snd_pcm_info_get_id(info.inner)))
}

func (info *Info) GetName() string {
	return C.GoString(C.no_const(C.snd_pcm_info_get_name(info.inner)))
}

func (info *Info) GetSubDeviceName() string {
	return C.GoString(C.no_const(C.snd_pcm_info_get_subdevice_name(info.inner)))
}

func (info *Info) GetClass() Class {
	return Class(C.snd_pcm_info_get_class(info.inner))
}

func (info *Info) GetSubClass() SubClass {
	return SubClass(C.snd_pcm_info_get_subclass(info.inner))
}

func (info *Info) GetSubDevicesCount() int {
	return int(C.snd_pcm_info_get_subdevices_count(info.inner))
}

func (info *Info) GetSubDevicesAvailableCount() int {
	return int(C.snd_pcm_info_get_subdevices_avail(info.inner))
}
