package alsa

// #cgo LDFLAGS: -lasound
//
// #include <alsa/asoundlib.h>
//
// static char *no_const(const char *s) { return (char *)s; }
import "C"
import (
	"sync/atomic"
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

type PCMHwParams struct {
	inner C.snd_pcm_hw_params_t
}
