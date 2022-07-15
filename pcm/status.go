package pcm

// #include <alsa/asoundlib.h>
//
// static int _snd_pcm_status_calloc(snd_pcm_status_t **ptr) {
//     int rc = snd_pcm_status_malloc(ptr);
//     if (rc < 0) return rc;
//     memset(*ptr, 0, snd_pcm_status_sizeof());
//     return 0;
// }
//
// typedef struct {
//     unsigned int valid;
//     unsigned int actual_type;
//     unsigned int accuracy_report;
//     unsigned int accuracy;
// } slow_snd_pcm_audio_tstamp_report_t;
//
// static void _snd_pcm_status_get_audio_htstamp_report(const snd_pcm_status_t *obj,
//         slow_snd_pcm_audio_tstamp_report_t *report) {
//     snd_pcm_audio_tstamp_report_t report2;
//     snd_pcm_status_get_audio_htstamp_report(obj, &report2);
//     report->valid = report2.valid;
//     report->actual_type = report2.actual_type;
//     report->accuracy_report = report2.accuracy_report;
//     report->accuracy = report2.accuracy;
// }
//
// static void _snd_pcm_status_set_audio_htstamp_config(snd_pcm_status_t *obj,
//         unsigned int type_requested, unsigned int report_delay) {
//     snd_pcm_audio_tstamp_config_t config2 = { .type_requested = type_requested, .report_delay = report_delay };
//     snd_pcm_status_set_audio_htstamp_config(obj, &config2);
// }
import "C"
import (
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ClarkGuan/go-alsa"
)

func (pcm *Dev) Status() (*Status, error) {
	s := new(Status)
	C._snd_pcm_status_calloc(&s.inner)
	rc := C.snd_pcm_status(pcm.inner, s.inner)
	if rc < 0 {
		return nil, alsa.NewError(int(rc))
	}
	runtime.SetFinalizer(s, (*Status).Close)
	return s, nil
}

type Status struct {
	inner *C.snd_pcm_status_t
}

func (s *Status) Close() error {
	for {
		p := unsafe.Pointer(s.inner)
		if p == nil {
			break
		}
		tmp := s.inner
		if atomic.CompareAndSwapPointer(&p, p, nil) {
			C.snd_pcm_status_free(tmp)
			runtime.SetFinalizer(s, nil)
		}
	}
	return nil
}

func (s *Status) GetState() State {
	return State(C.snd_pcm_status_get_state(s.inner))
}

func (s *Status) GetTriggerTimestamp() time.Time {
	var tm C.snd_timestamp_t
	C.snd_pcm_status_get_trigger_tstamp(s.inner, &tm)
	// 1us = 1000ns
	return time.Unix(int64(tm.tv_sec), int64(tm.tv_usec*1000))
}

func (s *Status) GetTriggerHighTimestamp() time.Time {
	var tm C.snd_htimestamp_t
	C.snd_pcm_status_get_trigger_htstamp(s.inner, &tm)
	return time.Unix(int64(tm.tv_sec), int64(tm.tv_nsec))
}

func (s *Status) GetAudioHighTimestamp() time.Time {
	var tm C.snd_htimestamp_t
	C.snd_pcm_status_get_audio_htstamp(s.inner, &tm)
	return time.Unix(int64(tm.tv_sec), int64(tm.tv_nsec))
}

type HighTimestampReport struct {
	Valid          bool
	ActualType     byte
	AccuracyReport bool
	Accuracy       int
}

func (s *Status) GetAudioHighTimestampReport() HighTimestampReport {
	var r C.slow_snd_pcm_audio_tstamp_report_t
	C._snd_pcm_status_get_audio_htstamp_report(s.inner, &r)
	return HighTimestampReport{
		Valid:          r.valid != 0,
		ActualType:     byte(r.actual_type),
		AccuracyReport: r.accuracy_report != 0,
		Accuracy:       int(r.accuracy),
	}
}

func (s *Status) GetDriverHighTimestamp() time.Time {
	var tm C.snd_htimestamp_t
	C.snd_pcm_status_get_driver_htstamp(s.inner, &tm)
	return time.Unix(int64(tm.tv_sec), int64(tm.tv_nsec))
}

func (s *Status) SetAudioHighTimestampConfig(typeRequested, reportDelay int8) {
	C._snd_pcm_status_set_audio_htstamp_config(s.inner, C.uint(typeRequested), C.uint(reportDelay))
}

func (s *Status) GetDelayFrames() int {
	return int(C.snd_pcm_status_get_delay(s.inner))
}

func (s *Status) GetAvailFrames() int {
	return int(C.snd_pcm_status_get_avail(s.inner))
}

func (s *Status) GetMaxAvailFrames() int {
	return int(C.snd_pcm_status_get_avail_max(s.inner))
}

func (s *Status) GetOverRangeCount() int {
	return int(C.snd_pcm_status_get_overrange(s.inner))
}
