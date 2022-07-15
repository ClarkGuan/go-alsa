package pcm

//
// #cgo LDFLAGS: -lasound
//
// #include <alsa/asoundlib.h>
//
// static char *no_const(const char *s) { return (char *)s; }
//
import "C"
import (
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ClarkGuan/go-alsa"
)

type PCM struct {
	inner *C.snd_pcm_t
}

func Open(name string, stream StreamType, mode int) (*PCM, error) {
	pcm := new(PCM)
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	rc := C.snd_pcm_open(&pcm.inner, cName, C.snd_pcm_stream_t(stream), C.int(mode))
	if rc < 0 {
		return nil, alsa.NewError(int(rc))
	}
	runtime.SetFinalizer(pcm, (*PCM).Close)
	return pcm, nil
}

func (pcm *PCM) Close() error {
	for {
		p := unsafe.Pointer(pcm.inner)
		if p == nil {
			break
		}
		c := pcm.inner
		if atomic.CompareAndSwapPointer(&p, p, nil) {
			rc := C.snd_pcm_close(c)
			if rc < 0 {
				panic(alsa.NewError(int(rc)))
			}
			runtime.SetFinalizer(pcm, nil)
		}
	}

	return nil
}

func (pcm *PCM) Name() string {
	return C.GoString(C.no_const(C.snd_pcm_name(pcm.inner)))
}

func (pcm *PCM) Type() Type {
	return Type(C.snd_pcm_type(pcm.inner))
}

func (pcm *PCM) StreamType() StreamType {
	return StreamType(C.snd_pcm_stream(pcm.inner))
}

func (pcm *PCM) NonBlock(enable bool) error {
	rc := C.snd_pcm_nonblock(pcm.inner, fromBool(enable))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

// TODO Async API
// snd_async_add_pcm_handler
// snd_async_handler_get_pcm

func (pcm *PCM) Prepare() error {
	rc := C.snd_pcm_prepare(pcm.inner)
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Reset() error {
	rc := C.snd_pcm_reset(pcm.inner)
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Start() error {
	rc := C.snd_pcm_start(pcm.inner)
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Drop() error {
	rc := C.snd_pcm_drop(pcm.inner)
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Drain() error {
	rc := C.snd_pcm_drain(pcm.inner)
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Pause(enable bool) error {
	rc := C.snd_pcm_pause(pcm.inner, fromBool(enable))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) State() State {
	return State(C.snd_pcm_state(pcm.inner))
}

//func (pcm *PCM) HwSync() error {
//	rc := C.snd_pcm_hwsync(pcm.inner)
//	if rc < 0 {
//		return NewError(int(rc))
//	}
//	return nil
//}

func (pcm *PCM) Delay() (int, error) {
	var delay C.snd_pcm_sframes_t
	rc := C.snd_pcm_delay(pcm.inner, &delay)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(delay), nil
}

func (pcm *PCM) Resume() error {
	rc := C.snd_pcm_resume(pcm.inner)
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Avail() (int, error) {
	rc := C.snd_pcm_avail(pcm.inner)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) AvailUpdate() (int, error) {
	rc := C.snd_pcm_avail_update(pcm.inner)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) AvailDelay() (int, int, error) {
	var availp, delayp C.snd_pcm_sframes_t
	rc := C.snd_pcm_avail_delay(pcm.inner, &availp, &delayp)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(availp), int(delayp), nil
}

func (pcm *PCM) Rewindable() (int, error) {
	rc := C.snd_pcm_rewindable(pcm.inner)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Rewind(frames int) (int, error) {
	rc := C.snd_pcm_rewind(pcm.inner, C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Forwardable() (int, error) {
	rc := C.snd_pcm_forwardable(pcm.inner)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Forward(frames int) (int, error) {
	rc := C.snd_pcm_forward(pcm.inner, C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Writei(data unsafe.Pointer, frames int) (int, error) {
	if data == nil {
		return 0, nil
	}
	rc := C.snd_pcm_writei(pcm.inner, data, C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		if code := pcm.recover(int(rc), true); code < 0 {
			return 0, alsa.NewError(code)
		}
		return 0, nil
	}
	return int(rc), nil
}

func (pcm *PCM) Readi(data unsafe.Pointer, frames int) (int, error) {
	if data == nil {
		return 0, nil
	}
	rc := C.snd_pcm_readi(pcm.inner, data, C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		if code := pcm.recover(int(rc), true); code < 0 {
			return 0, alsa.NewError(code)
		}
		return 0, nil
	}
	return int(rc), nil
}

func (pcm *PCM) Writen(data []unsafe.Pointer, frames int) (int, error) {
	if len(data) <= 0 {
		return 0, nil
	}
	rc := C.snd_pcm_writen(pcm.inner, &data[0], C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		if code := pcm.recover(int(rc), true); code < 0 {
			return 0, alsa.NewError(code)
		}
		return 0, nil
	}
	return int(rc), nil
}

func (pcm *PCM) Readn(data []unsafe.Pointer, frames int) (int, error) {
	if len(data) <= 0 {
		return 0, nil
	}
	rc := C.snd_pcm_readn(pcm.inner, &data[0], C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		if code := pcm.recover(int(rc), true); code < 0 {
			return 0, alsa.NewError(code)
		}
		return 0, nil
	}
	return int(rc), nil
}

func (pcm *PCM) Link(other *PCM) error {
	rc := C.snd_pcm_link(pcm.inner, other.inner)
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Unlink() error {
	rc := C.snd_pcm_unlink(pcm.inner)
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) HighTimestamp() (int, *time.Time, error) {
	var avail C.snd_pcm_uframes_t
	var tstamp C.snd_htimestamp_t
	rc := C.snd_pcm_htimestamp(pcm.inner, &avail, &tstamp)
	if rc < 0 {
		return 0, nil, alsa.NewError(int(rc))
	}
	tm := time.Unix(int64(tstamp.tv_sec), int64(tstamp.tv_nsec))
	return int(avail), &tm, nil
}

func (pcm *PCM) Wait(timeout int) (int, error) {
	rc := C.snd_pcm_wait(pcm.inner, C.int(timeout))
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) PollDescriptorsCount() int {
	return int(C.snd_pcm_poll_descriptors_count(pcm.inner))
}

func (pcm *PCM) PollDescriptors(fds []PollFd) int {
	if len(fds) <= 0 {
		return 0
	}
	return int(C.snd_pcm_poll_descriptors(pcm.inner, (*C.struct_pollfd)(unsafe.Pointer(&fds[0])), C.uint(len(fds))))
}

func (pcm *PCM) PollDescriptorsREvents(fds []PollFd) (int16, error) {
	if len(fds) <= 0 {
		return 0, nil
	}
	var revents C.ushort
	rc := C.snd_pcm_poll_descriptors_revents(pcm.inner, (*C.struct_pollfd)(unsafe.Pointer(&fds[0])), C.uint(len(fds)), &revents)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int16(revents), nil
}

func (pcm *PCM) InstallHardwareParams(params *HardwareParams) error {
	rc := C.snd_pcm_hw_params(pcm.inner, params.inner)
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) UninstallHardwareParams() error {
	rc := C.snd_pcm_hw_free(pcm.inner)
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) InstallSoftwareParams(params *SoftwareParams) error {
	rc := C.snd_pcm_sw_params(pcm.inner, params.inner)
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) BytesToFrames(byteCount int) int {
	return int(C.snd_pcm_bytes_to_frames(pcm.inner, C.ssize_t(byteCount)))
}

func (pcm *PCM) FramesToBytes(frames int) int {
	return int(C.snd_pcm_frames_to_bytes(pcm.inner, C.snd_pcm_sframes_t(frames)))
}

func (pcm *PCM) BytesToSamples(byteCount int) int {
	return int(C.snd_pcm_bytes_to_samples(pcm.inner, C.ssize_t(byteCount)))
}

func (pcm *PCM) SamplesToBytes(samples int) int {
	return int(C.snd_pcm_samples_to_bytes(pcm.inner, C.long(samples)))
}

func (pcm *PCM) Recover(err error, silent bool) error {
	if alsaErr, b := err.(*alsa.Error); b {
		rc := pcm.recover(alsaErr.Errno, silent)
		if rc < 0 {
			return alsa.NewError(rc)
		}
		return nil
	} else {
		return err
	}
}

func (pcm *PCM) recover(errno int, silent bool) int {
	return int(C.snd_pcm_recover(pcm.inner, C.int(errno), fromBool(silent)))
}

func (pcm *PCM) SetParams(format Format, access Access, channels, rate int, resample bool, latency time.Duration) error {
	rc := C.snd_pcm_set_params(pcm.inner,
		C.snd_pcm_format_t(format),
		C.snd_pcm_access_t(access),
		C.uint(channels),
		C.uint(rate),
		fromBool(resample),
		C.uint(latency.Microseconds()))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) GetParams() (int, int, error) {
	var bufferFrames C.snd_pcm_uframes_t
	var periodFrames C.snd_pcm_uframes_t
	rc := C.snd_pcm_get_params(pcm.inner, &bufferFrames, &periodFrames)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(bufferFrames), int(periodFrames), nil
}
