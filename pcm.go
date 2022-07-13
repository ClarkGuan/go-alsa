package alsa

// #cgo LDFLAGS: -lasound
//
// #include <poll.h>
// #include <alsa/asoundlib.h>
//
// static char *no_const(const char *s) { return (char *)s; }
import "C"
import (
	"sync/atomic"
	"time"
	"unsafe"
)

type StreamType int
type PCMType int
type PCMState int

const (
	SndPcmStreamPlayback = StreamType(C.SND_PCM_STREAM_PLAYBACK)
	SndPcmStreamCapture  = StreamType(C.SND_PCM_STREAM_CAPTURE)

	SndPcmNonblock = int(C.SND_PCM_NONBLOCK)
	SndPcmAsync    = int(C.SND_PCM_ASYNC)

	SndPcmTypeHw          = PCMType(C.SND_PCM_TYPE_HW)
	SndPcmTypeHooks       = PCMType(C.SND_PCM_TYPE_HOOKS)
	SndPcmTypeMulti       = PCMType(C.SND_PCM_TYPE_MULTI)
	SndPcmTypeFile        = PCMType(C.SND_PCM_TYPE_FILE)
	SndPcmTypeNull        = PCMType(C.SND_PCM_TYPE_NULL)
	SndPcmTypeShm         = PCMType(C.SND_PCM_TYPE_SHM)
	SndPcmTypeInet        = PCMType(C.SND_PCM_TYPE_INET)
	SndPcmTypeCopy        = PCMType(C.SND_PCM_TYPE_COPY)
	SndPcmTypeLinear      = PCMType(C.SND_PCM_TYPE_LINEAR)
	SndPcmTypeAlaw        = PCMType(C.SND_PCM_TYPE_ALAW)
	SndPcmTypeMulaw       = PCMType(C.SND_PCM_TYPE_MULAW)
	SndPcmTypeAdpcm       = PCMType(C.SND_PCM_TYPE_ADPCM)
	SndPcmTypeRate        = PCMType(C.SND_PCM_TYPE_RATE)
	SndPcmTypeRoute       = PCMType(C.SND_PCM_TYPE_ROUTE)
	SndPcmTypePlug        = PCMType(C.SND_PCM_TYPE_PLUG)
	SndPcmTypeShare       = PCMType(C.SND_PCM_TYPE_SHARE)
	SndPcmTypeMeter       = PCMType(C.SND_PCM_TYPE_METER)
	SndPcmTypeMix         = PCMType(C.SND_PCM_TYPE_MIX)
	SndPcmTypeDroute      = PCMType(C.SND_PCM_TYPE_DROUTE)
	SndPcmTypeLbserver    = PCMType(C.SND_PCM_TYPE_LBSERVER)
	SndPcmTypeLinearFloat = PCMType(C.SND_PCM_TYPE_LINEAR_FLOAT)
	SndPcmTypeLadspa      = PCMType(C.SND_PCM_TYPE_LADSPA)
	SndPcmTypeDmix        = PCMType(C.SND_PCM_TYPE_DMIX)
	SndPcmTypeJack        = PCMType(C.SND_PCM_TYPE_JACK)
	SndPcmTypeDsnoop      = PCMType(C.SND_PCM_TYPE_DSNOOP)
	SndPcmTypeDshare      = PCMType(C.SND_PCM_TYPE_DSHARE)
	SndPcmTypeIec958      = PCMType(C.SND_PCM_TYPE_IEC958)
	SndPcmTypeSoftvol     = PCMType(C.SND_PCM_TYPE_SOFTVOL)
	SndPcmTypeIoplug      = PCMType(C.SND_PCM_TYPE_IOPLUG)
	SndPcmTypeExtplug     = PCMType(C.SND_PCM_TYPE_EXTPLUG)
	SndPcmTypeMmapEmul    = PCMType(C.SND_PCM_TYPE_MMAP_EMUL)

	SndPcmStateOpen         = PCMState(C.SND_PCM_STATE_OPEN)
	SndPcmStateSetup        = PCMState(C.SND_PCM_STATE_SETUP)
	SndPcmStatePrepared     = PCMState(C.SND_PCM_STATE_PREPARED)
	SndPcmStateRunning      = PCMState(C.SND_PCM_STATE_RUNNING)
	SndPcmStateXrun         = PCMState(C.SND_PCM_STATE_XRUN)
	SndPcmStateDraining     = PCMState(C.SND_PCM_STATE_DRAINING)
	SndPcmStatePaused       = PCMState(C.SND_PCM_STATE_PAUSED)
	SndPcmStateSuspended    = PCMState(C.SND_PCM_STATE_SUSPENDED)
	SndPcmStateDisconnected = PCMState(C.SND_PCM_STATE_DISCONNECTED)
)

type PollFd struct {
	Fd      int32
	Events  int16
	REvents int16
}

const (
	PollIn   = int16(C.POLLIN)
	PollOut  = int16(C.POLLOUT)
	PollErr  = int16(C.POLLERR)
	PollPri  = int16(C.POLLPRI)
	PollHup  = int16(C.POLLHUP)
	PollNval = int16(C.POLLNVAL)
	//PollMsg   = int16(C.POLLMSG)
	//PollRdHup = int16(C.POLLRDHUP)
)

type PCM struct {
	inner *C.snd_pcm_t
}

func OpenPCM(name string, stream StreamType, mode int) (*PCM, error) {
	pcm := &PCM{}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	rc := C.snd_pcm_open(&pcm.inner, cName, C.snd_pcm_stream_t(stream), C.int(mode))
	if rc != 0 {
		return nil, NewError(int(rc))
	}
	return pcm, nil
}

func (pcm *PCM) Close() error {
	for {
		p := unsafe.Pointer(pcm.inner)
		if p == nil {
			break
		}
		if atomic.CompareAndSwapPointer(&p, p, nil) {
			rc := C.snd_pcm_close((*C.snd_pcm_t)(p))
			if rc != 0 {
				return NewError(int(rc))
			}
		}
	}

	return nil
}

func (pcm *PCM) Name() string {
	return C.GoString(C.no_const(C.snd_pcm_name(pcm.inner)))
}

func (pcm *PCM) Type() PCMType {
	return PCMType(C.snd_pcm_type(pcm.inner))
}

func (pcm *PCM) StreamType() StreamType {
	return StreamType(C.snd_pcm_stream(pcm.inner))
}

func (pcm *PCM) NonBlock(enable bool) error {
	rc := C.snd_pcm_nonblock(pcm.inner, fromBool(enable))
	if rc != 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Prepare() error {
	rc := C.snd_pcm_prepare(pcm.inner)
	if rc != 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Reset() error {
	rc := C.snd_pcm_reset(pcm.inner)
	if rc != 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Start() error {
	rc := C.snd_pcm_start(pcm.inner)
	if rc != 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Drop() error {
	rc := C.snd_pcm_drop(pcm.inner)
	if rc != 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Drain() error {
	rc := C.snd_pcm_drain(pcm.inner)
	if rc != 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Pause(enable bool) error {
	rc := C.snd_pcm_pause(pcm.inner, fromBool(enable))
	if rc != 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) State() PCMState {
	return PCMState(C.snd_pcm_state(pcm.inner))
}

//func (pcm *PCM) HwSync() error {
//	rc := C.snd_pcm_hwsync(pcm.inner)
//	if rc != 0 {
//		return NewError(int(rc))
//	}
//	return nil
//}

func (pcm *PCM) Delay() (int, error) {
	var delay C.snd_pcm_sframes_t
	rc := C.snd_pcm_delay(pcm.inner, &delay)
	if rc != 0 {
		return 0, NewError(int(rc))
	}
	return int(delay), nil
}

func (pcm *PCM) Resume() error {
	rc := C.snd_pcm_resume(pcm.inner)
	if rc != 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Avail() (int, error) {
	rc := C.snd_pcm_avail(pcm.inner)
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) AvailUpdate() (int, error) {
	rc := C.snd_pcm_avail_update(pcm.inner)
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) AvailDelay() (int, int, error) {
	var availp, delayp C.snd_pcm_sframes_t
	rc := C.snd_pcm_avail_delay(pcm.inner, &availp, &delayp)
	if rc != 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(availp), int(delayp), nil
}

func (pcm *PCM) Rewindable() (int, error) {
	rc := C.snd_pcm_rewindable(pcm.inner)
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Rewind(frames int) (int, error) {
	rc := C.snd_pcm_rewind(pcm.inner, C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Forwardable() (int, error) {
	rc := C.snd_pcm_forwardable(pcm.inner)
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Forward(frames int) (int, error) {
	rc := C.snd_pcm_forward(pcm.inner, C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Writei(data []byte, frames int) (int, error) {
	if len(data) <= 0 {
		return 0, nil
	}
	rc := C.snd_pcm_writei(pcm.inner, unsafe.Pointer(&data[0]), C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Readi(data []byte, frames int) (int, error) {
	if len(data) <= 0 {
		return 0, nil
	}
	rc := C.snd_pcm_readi(pcm.inner, unsafe.Pointer(&data[0]), C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Writen(data [][]byte, frames int) (int, error) {
	if len(data) <= 0 {
		return 0, nil
	}
	ps := make([]unsafe.Pointer, len(data))
	for i, s := range data {
		ps[i] = unsafe.Pointer(&s[0])
	}
	rc := C.snd_pcm_writen(pcm.inner, &ps[0], C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Readn(data [][]byte, frames int) (int, error) {
	if len(data) <= 0 {
		return 0, nil
	}
	ps := make([]unsafe.Pointer, len(data))
	for i, s := range data {
		ps[i] = unsafe.Pointer(&s[0])
	}
	rc := C.snd_pcm_readn(pcm.inner, &ps[0], C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Link(other *PCM) error {
	rc := C.snd_pcm_link(pcm.inner, other.inner)
	if rc != 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Unlink() error {
	rc := C.snd_pcm_unlink(pcm.inner)
	if rc != 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) HTimestamp() (int, *time.Time, error) {
	var avail C.snd_pcm_uframes_t
	var tstamp C.snd_htimestamp_t
	rc := C.snd_pcm_htimestamp(pcm.inner, &avail, &tstamp)
	if rc < 0 {
		return 0, nil, NewError(int(rc))
	}
	tm := time.Unix(int64(tstamp.tv_sec), int64(tstamp.tv_nsec))
	return int(avail), &tm, nil
}

func (pcm *PCM) Wait(timeout int) (int, error) {
	rc := C.snd_pcm_wait(pcm.inner, C.int(timeout))
	if rc < 0 {
		return 0, NewError(int(rc))
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
	if rc != 0 {
		return 0, NewError(int(rc))
	}
	return int16(revents), nil
}

type PCMHwParams struct {
	inner C.snd_pcm_hw_params_t
}
