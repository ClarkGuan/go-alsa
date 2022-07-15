package pcm

//
// #include <alsa/asoundlib.h>
//
// static int _snd_pcm_sw_params_calloc(snd_pcm_sw_params_t **ptr) {
//     int rc = snd_pcm_sw_params_malloc(ptr);
//     if (rc < 0) return rc;
//     memset(*ptr, 0, snd_pcm_sw_params_sizeof());
//     return 0;
// }
//
import "C"
import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/ClarkGuan/go-alsa"
)

type SoftwareParams struct {
	inner *C.snd_pcm_sw_params_t
	pcm   *Dev
}

func (pcm *Dev) CurrentSoftwareParams() (*SoftwareParams, error) {
	params := new(SoftwareParams)
	C._snd_pcm_sw_params_calloc(&params.inner)
	rc := C.snd_pcm_sw_params_current(pcm.inner, params.inner)
	if rc < 0 {
		return nil, alsa.NewError(int(rc))
	}
	params.pcm = pcm
	runtime.SetFinalizer(params, (*SoftwareParams).Close)
	return params, nil
}

func (params *SoftwareParams) Close() error {
	for {
		p := unsafe.Pointer(params.inner)
		if p == nil {
			break
		}
		tmp := params.inner
		if atomic.CompareAndSwapPointer(&p, p, nil) {
			C.snd_pcm_sw_params_free(tmp)
			runtime.SetFinalizer(params, nil)
			params.pcm = nil
		}
	}
	return nil
}

func (params *SoftwareParams) GetBoundary() (int, error) {
	var val C.snd_pcm_uframes_t
	rc := C.snd_pcm_sw_params_get_boundary(params.inner, &val)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(val), nil
}

func (params *SoftwareParams) SetTimestampMode(mode TimestampMode) error {
	rc := C.snd_pcm_sw_params_set_tstamp_mode(params.pcm.inner, params.inner, C.snd_pcm_tstamp_t(mode))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *SoftwareParams) GetTimestampMode() (TimestampMode, error) {
	var val C.snd_pcm_tstamp_t
	rc := C.snd_pcm_sw_params_get_tstamp_mode(params.inner, &val)
	if rc < 0 {
		return -1, alsa.NewError(int(rc))
	}
	return TimestampMode(val), nil
}

func (params *SoftwareParams) SetTimestampType(t TimestampType) error {
	rc := C.snd_pcm_sw_params_set_tstamp_type(params.pcm.inner, params.inner, C.snd_pcm_tstamp_type_t(t))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *SoftwareParams) GetTimestampType() (TimestampType, error) {
	var val C.snd_pcm_tstamp_type_t
	rc := C.snd_pcm_sw_params_get_tstamp_type(params.inner, &val)
	if rc < 0 {
		return -1, alsa.NewError(int(rc))
	}
	return TimestampType(val), nil
}

func (params *SoftwareParams) SetMinAvailFrames(frames int) error {
	rc := C.snd_pcm_sw_params_set_avail_min(params.pcm.inner, params.inner, C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *SoftwareParams) GetMinAvailFrames() (int, error) {
	var val C.snd_pcm_uframes_t
	rc := C.snd_pcm_sw_params_get_avail_min(params.inner, &val)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(val), nil
}

func (params *SoftwareParams) SetPeriodEvent(enable bool) error {
	rc := C.snd_pcm_sw_params_set_period_event(params.pcm.inner, params.inner, fromBool(enable))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *SoftwareParams) GetPeriodEvent() (bool, error) {
	var enable C.int
	rc := C.snd_pcm_sw_params_get_period_event(params.inner, &enable)
	if rc < 0 {
		return false, alsa.NewError(int(rc))
	}
	return enable != 0, nil
}

func (params *SoftwareParams) SetStartThresholdFrames(frames int) error {
	rc := C.snd_pcm_sw_params_set_start_threshold(params.pcm.inner, params.inner, C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *SoftwareParams) GetStartThresholdFrames() (int, error) {
	var val C.snd_pcm_uframes_t
	rc := C.snd_pcm_sw_params_get_start_threshold(params.inner, &val)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(val), nil
}

func (params *SoftwareParams) SetStopThresholdFrames(frames int) error {
	rc := C.snd_pcm_sw_params_set_stop_threshold(params.pcm.inner, params.inner, C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *SoftwareParams) GetStopThresholdFrames() (int, error) {
	var val C.snd_pcm_uframes_t
	rc := C.snd_pcm_sw_params_get_stop_threshold(params.inner, &val)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(val), nil
}

func (params *SoftwareParams) SetSilenceThresholdFrames(frames int) error {
	rc := C.snd_pcm_sw_params_set_silence_threshold(params.pcm.inner, params.inner, C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *SoftwareParams) GetSilenceThresholdFrames() (int, error) {
	var val C.snd_pcm_uframes_t
	rc := C.snd_pcm_sw_params_get_silence_threshold(params.inner, &val)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(val), nil
}

func (params *SoftwareParams) SetSilenceSizeInFrames(frames int) error {
	rc := C.snd_pcm_sw_params_set_silence_size(params.pcm.inner, params.inner, C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *SoftwareParams) GetSilenceSizeInFrames() (int, error) {
	var val C.snd_pcm_uframes_t
	rc := C.snd_pcm_sw_params_get_silence_size(params.inner, &val)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(val), nil
}
