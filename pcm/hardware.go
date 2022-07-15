package pcm

//
// #include <alsa/asoundlib.h>
//
// static int _snd_pcm_hw_params_calloc(snd_pcm_hw_params_t **ptr) {
//     int rc = snd_pcm_hw_params_malloc(ptr);
//     if (rc < 0) return rc;
//     memset(*ptr, 0, snd_pcm_hw_params_sizeof());
//     return 0;
// }
//
import "C"
import (
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ClarkGuan/go-alsa"
)

type HardwareParams struct {
	inner *C.snd_pcm_hw_params_t
	pcm   *Dev
}

func (pcm *Dev) AnyHardwareParams() (*HardwareParams, error) {
	params := new(HardwareParams)
	rc := C._snd_pcm_hw_params_calloc(&params.inner)
	if rc < 0 {
		return nil, alsa.NewError(int(rc))
	}
	rc = C.snd_pcm_hw_params_any(pcm.inner, params.inner)
	if rc < 0 {
		return nil, alsa.NewError(int(rc))
	}
	params.pcm = pcm
	runtime.SetFinalizer(params, (*HardwareParams).Close)
	return params, nil
}

func (pcm *Dev) CurrentHardwareParams() (*HardwareParams, error) {
	params := new(HardwareParams)
	rc := C._snd_pcm_hw_params_calloc(&params.inner)
	if rc < 0 {
		return nil, alsa.NewError(int(rc))
	}
	rc = C.snd_pcm_hw_params_current(pcm.inner, params.inner)
	if rc < 0 {
		return nil, alsa.NewError(int(rc))
	}
	params.pcm = pcm
	runtime.SetFinalizer(params, (*HardwareParams).Close)
	return params, nil
}

func (params *HardwareParams) Close() error {
	for {
		p := unsafe.Pointer(params.inner)
		if p == nil {
			break
		}
		i := params.inner
		if atomic.CompareAndSwapPointer(&p, p, nil) {
			C.snd_pcm_hw_params_free(i)
			runtime.SetFinalizer(params, nil)
		}
	}

	return nil
}

func (params *HardwareParams) CanMmapSampleResolution() bool {
	return C.snd_pcm_hw_params_can_mmap_sample_resolution(params.inner) != 0
}

func (params *HardwareParams) IsDouble() bool {
	return C.snd_pcm_hw_params_is_double(params.inner) != 0
}

func (params *HardwareParams) IsBatch() bool {
	return C.snd_pcm_hw_params_is_batch(params.inner) != 0
}

func (params *HardwareParams) IsBlockTransfer() bool {
	return C.snd_pcm_hw_params_is_block_transfer(params.inner) != 0
}

func (params *HardwareParams) IsMonotonic() bool {
	return C.snd_pcm_hw_params_is_monotonic(params.inner) != 0
}

func (params *HardwareParams) CanOverRange() bool {
	return C.snd_pcm_hw_params_can_overrange(params.inner) != 0
}

func (params *HardwareParams) CanPause() bool {
	return C.snd_pcm_hw_params_can_pause(params.inner) != 0
}

func (params *HardwareParams) CanResume() bool {
	return C.snd_pcm_hw_params_can_resume(params.inner) != 0
}

func (params *HardwareParams) IsHalfDuplex() bool {
	return C.snd_pcm_hw_params_is_half_duplex(params.inner) != 0
}

func (params *HardwareParams) IsJointDuplex() bool {
	return C.snd_pcm_hw_params_is_joint_duplex(params.inner) != 0
}

func (params *HardwareParams) CanSyncStart() bool {
	return C.snd_pcm_hw_params_can_sync_start(params.inner) != 0
}

func (params *HardwareParams) CanDisablePeriodWakeup() bool {
	return C.snd_pcm_hw_params_can_disable_period_wakeup(params.inner) != 0
}

func (params *HardwareParams) SupportsAudioWallClockTimestamps() bool {
	return C.snd_pcm_hw_params_supports_audio_wallclock_ts(params.inner) != 0
}

func (params *HardwareParams) SupportsAudioTimestampType(tsType int) bool {
	return C.snd_pcm_hw_params_supports_audio_ts_type(params.inner, C.int(tsType)) != 0
}

func (params *HardwareParams) GetRateNumDen() (int, int, error) {
	var rateNum, rateDen C.uint
	rc := C.snd_pcm_hw_params_get_rate_numden(params.inner, &rateNum, &rateDen)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(rateNum), int(rateDen), nil
}

func (params *HardwareParams) GetSBits() (int, error) {
	rc := C.snd_pcm_hw_params_get_sbits(params.inner)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(rc), nil
}

func (params *HardwareParams) GetFifoSize() (int, error) {
	rc := C.snd_pcm_hw_params_get_fifo_size(params.inner)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(rc), nil
}

func (params *HardwareParams) GetAccess() (Access, error) {
	var acs C.snd_pcm_access_t
	rc := C.snd_pcm_hw_params_get_access(params.inner, &acs)
	if rc < 0 {
		return -1, alsa.NewError(int(rc))
	}
	return Access(acs), nil
}

func (params *HardwareParams) TestAccess(acs Access) error {
	rc := C.snd_pcm_hw_params_test_access(params.pcm.inner, params.inner, C.snd_pcm_access_t(acs))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetAccess(acs Access) error {
	rc := C.snd_pcm_hw_params_set_access(params.pcm.inner, params.inner, C.snd_pcm_access_t(acs))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetAccessFirst() (Access, error) {
	var access C.snd_pcm_access_t
	rc := C.snd_pcm_hw_params_set_access_first(params.pcm.inner, params.inner, &access)
	if rc < 0 {
		return -1, alsa.NewError(int(rc))
	}
	return Access(access), nil
}

func (params *HardwareParams) SetAccessLast() (Access, error) {
	var access C.snd_pcm_access_t
	rc := C.snd_pcm_hw_params_set_access_last(params.pcm.inner, params.inner, &access)
	if rc < 0 {
		return -1, alsa.NewError(int(rc))
	}
	return Access(access), nil
}

func (params *HardwareParams) GetFormat() (Format, error) {
	var format C.snd_pcm_format_t
	rc := C.snd_pcm_hw_params_get_format(params.inner, &format)
	if rc < 0 {
		return -1, alsa.NewError(int(rc))
	}
	return Format(format), nil
}

func (params *HardwareParams) TestFormat(format Format) error {
	rc := C.snd_pcm_hw_params_test_format(params.pcm.inner, params.inner, C.snd_pcm_format_t(format))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetFormat(format Format) error {
	rc := C.snd_pcm_hw_params_set_format(params.pcm.inner, params.inner, C.snd_pcm_format_t(format))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetFormatFirst() (Format, error) {
	var format C.snd_pcm_format_t
	rc := C.snd_pcm_hw_params_set_format_first(params.pcm.inner, params.inner, &format)
	if rc < 0 {
		return SndPcmFormatUnknown, alsa.NewError(int(rc))
	}
	return Format(format), nil
}

func (params *HardwareParams) SetFormatLast() (Format, error) {
	var format C.snd_pcm_format_t
	rc := C.snd_pcm_hw_params_set_format_last(params.pcm.inner, params.inner, &format)
	if rc < 0 {
		return SndPcmFormatUnknown, alsa.NewError(int(rc))
	}
	return Format(format), nil
}

func (params *HardwareParams) GetChannels() (int, error) {
	var channels C.uint
	rc := C.snd_pcm_hw_params_get_channels(params.inner, &channels)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(channels), nil
}

func (params *HardwareParams) GetChannelsMin() (int, error) {
	var channels C.uint
	rc := C.snd_pcm_hw_params_get_channels_min(params.inner, &channels)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(channels), nil
}

func (params *HardwareParams) GetChannelsMax() (int, error) {
	var channels C.uint
	rc := C.snd_pcm_hw_params_get_channels_max(params.inner, &channels)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(channels), nil
}

func (params *HardwareParams) TestChannels(channels int) error {
	rc := C.snd_pcm_hw_params_test_channels(params.pcm.inner, params.inner, C.uint(channels))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetChannels(channels int) error {
	rc := C.snd_pcm_hw_params_set_channels(params.pcm.inner, params.inner, C.uint(channels))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetChannelsMin(min int) (int, error) {
	var channels = C.uint(min)
	rc := C.snd_pcm_hw_params_set_channels_min(params.pcm.inner, params.inner, &channels)
	if rc < 0 {
		return -1, alsa.NewError(int(rc))
	}
	return int(channels), nil
}

func (params *HardwareParams) SetChannelsMax(max int) (int, error) {
	var channels = C.uint(max)
	rc := C.snd_pcm_hw_params_set_channels_max(params.pcm.inner, params.inner, &channels)
	if rc < 0 {
		return -1, alsa.NewError(int(rc))
	}
	return int(channels), nil
}

func (params *HardwareParams) SetChannelsMinMax(min, max int) (int, int, error) {
	var cMin = C.uint(min)
	var cMax = C.uint(max)
	rc := C.snd_pcm_hw_params_set_channels_minmax(params.pcm.inner, params.inner, &cMin, &cMax)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(cMin), int(cMax), nil
}

func (params *HardwareParams) SetChannelsNear(channels int) (int, error) {
	cArg := C.uint(channels)
	rc := C.snd_pcm_hw_params_set_channels_near(params.pcm.inner, params.inner, &cArg)
	if rc < 0 {
		return -1, alsa.NewError(int(rc))
	}
	return int(cArg), nil
}

func (params *HardwareParams) SetChannelsFirst() (int, error) {
	var channels C.uint
	rc := C.snd_pcm_hw_params_set_channels_first(params.pcm.inner, params.inner, &channels)
	if rc < 0 {
		return -1, alsa.NewError(int(rc))
	}
	return int(channels), nil
}

func (params *HardwareParams) SetChannelsLast() (int, error) {
	var channels C.uint
	rc := C.snd_pcm_hw_params_set_channels_last(params.pcm.inner, params.inner, &channels)
	if rc < 0 {
		return -1, alsa.NewError(int(rc))
	}
	return int(channels), nil
}

func (params *HardwareParams) GetRate() (int, int, error) {
	var rate C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_rate(params.inner, &rate, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(rate), int(dir), nil
}

func (params *HardwareParams) GetRateMin() (int, int, error) {
	var rate C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_rate_min(params.inner, &rate, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(rate), int(dir), nil
}

func (params *HardwareParams) GetRateMax() (int, int, error) {
	var rate C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_rate_max(params.inner, &rate, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(rate), int(dir), nil
}

func (params *HardwareParams) TestRate(rate, dir int) error {
	rc := C.snd_pcm_hw_params_test_rate(params.pcm.inner, params.inner, C.uint(rate), C.int(dir))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetRate(rate, dir int) error {
	rc := C.snd_pcm_hw_params_set_rate(params.pcm.inner, params.inner, C.uint(rate), C.int(dir))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetRateMin(rate, dir int) (int, int, error) {
	var cRate = C.uint(rate)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_rate_min(params.pcm.inner, params.inner, &cRate, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(cRate), int(cDir), nil
}

func (params *HardwareParams) SetRateMax(rate, dir int) (int, int, error) {
	var cRate = C.uint(rate)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_rate_max(params.pcm.inner, params.inner, &cRate, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(cRate), int(cDir), nil
}

func (params *HardwareParams) SetRateMinMax(min, minDir, max, maxDir int) (int, int, int, int, error) {
	var cMin = C.uint(min)
	var cMinDir = C.int(minDir)
	var cMax = C.uint(max)
	var cMaxDir = C.int(maxDir)
	rc := C.snd_pcm_hw_params_set_rate_minmax(params.pcm.inner, params.inner, &cMin, &cMinDir, &cMax, &cMaxDir)
	if rc < 0 {
		return 0, 0, 0, 0, alsa.NewError(int(rc))
	}
	return int(cMin), int(cMinDir), int(cMax), int(cMaxDir), nil
}

func (params *HardwareParams) SetRateNear(rate, dir int) (int, int, error) {
	var cRate = C.uint(rate)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_rate_near(params.pcm.inner, params.inner, &cRate, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(cRate), int(cDir), nil
}

func (params *HardwareParams) SetRateFirst() (int, int, error) {
	var rate C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_set_rate_first(params.pcm.inner, params.inner, &rate, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(rate), int(dir), nil
}

func (params *HardwareParams) SetRateLast() (int, int, error) {
	var rate C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_set_rate_last(params.pcm.inner, params.inner, &rate, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(rate), int(dir), nil
}

func (params *HardwareParams) SetRateResample(enable bool) error {
	rc := C.snd_pcm_hw_params_set_rate_resample(params.pcm.inner, params.inner, C.uint(fromBool(enable)))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) GetRateResample() (bool, error) {
	var enable C.uint
	rc := C.snd_pcm_hw_params_get_rate_resample(params.pcm.inner, params.inner, &enable)
	if rc < 0 {
		return false, alsa.NewError(int(rc))
	}
	return enable != 0, nil
}

func (params *HardwareParams) SetExportBuffer(enable bool) error {
	rc := C.snd_pcm_hw_params_set_export_buffer(params.pcm.inner, params.inner, C.uint(fromBool(enable)))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) GetExportBuffer() (bool, error) {
	var enable C.uint
	rc := C.snd_pcm_hw_params_get_export_buffer(params.pcm.inner, params.inner, &enable)
	if rc < 0 {
		return false, alsa.NewError(int(rc))
	}
	return enable != 0, nil
}

func (params *HardwareParams) SetPeriodWakeup(enable bool) error {
	rc := C.snd_pcm_hw_params_set_period_wakeup(params.pcm.inner, params.inner, C.uint(fromBool(enable)))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) GetPeriodWakeup() (bool, error) {
	var enable C.uint
	rc := C.snd_pcm_hw_params_get_period_wakeup(params.pcm.inner, params.inner, &enable)
	if rc < 0 {
		return false, alsa.NewError(int(rc))
	}
	return enable != 0, nil
}

func (params *HardwareParams) GetPeriodTime() (time.Duration, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_period_time(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(val) * time.Microsecond, int(dir), nil
}

func (params *HardwareParams) GetPeriodTimeMin() (time.Duration, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_period_time_min(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(val) * time.Microsecond, int(dir), nil
}

func (params *HardwareParams) GetPeriodTimeMax() (time.Duration, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_period_time_max(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(val) * time.Microsecond, int(dir), nil
}

func (params *HardwareParams) TestPeriodTime(val time.Duration, dir int) error {
	rc := C.snd_pcm_hw_params_test_period_time(params.pcm.inner, params.inner, C.uint(val.Microseconds()), C.int(dir))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetPeriodTime(val time.Duration, dir int) error {
	rc := C.snd_pcm_hw_params_set_period_time(params.pcm.inner, params.inner, C.uint(val.Microseconds()), C.int(dir))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetPeriodTimeMin(val time.Duration, dir int) (time.Duration, int, error) {
	var cVal = C.uint(val.Microseconds())
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_period_time_min(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(cVal) * time.Microsecond, int(cDir), nil
}

func (params *HardwareParams) SetPeriodTimeMax(val time.Duration, dir int) (time.Duration, int, error) {
	var cVal = C.uint(val.Microseconds())
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_period_time_max(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(cVal) * time.Microsecond, int(cDir), nil
}

func (params *HardwareParams) SetPeriodTimeMinMax(
	min time.Duration, minDir int, max time.Duration, maxDir int) (time.Duration, int, time.Duration, int, error) {
	var cMin = C.uint(min.Microseconds())
	var cMinDir = C.int(minDir)
	var cMax = C.uint(max.Microseconds())
	var cMaxDir = C.int(maxDir)
	rc := C.snd_pcm_hw_params_set_period_time_minmax(params.pcm.inner, params.inner, &cMin, &cMinDir, &cMax, &cMaxDir)
	if rc < 0 {
		return 0, 0, 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(cMin) * time.Microsecond, int(cMinDir), time.Duration(cMax) * time.Microsecond, int(cMaxDir), nil
}

func (params *HardwareParams) SetPeriodTimeNear(val time.Duration, dir int) (time.Duration, int, error) {
	var cVal = C.uint(val.Microseconds())
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_period_time_near(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(cVal) * time.Microsecond, int(cDir), nil
}

func (params *HardwareParams) SetPeriodTimeFirst() (time.Duration, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_set_period_time_first(params.pcm.inner, params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(val) * time.Microsecond, int(dir), nil
}

func (params *HardwareParams) SetPeriodTimeLast() (time.Duration, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_set_period_time_last(params.pcm.inner, params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(val) * time.Microsecond, int(dir), nil
}

func (params *HardwareParams) GetPeriodSize() (int, int, error) {
	var val C.snd_pcm_uframes_t
	var dir C.int
	rc := C.snd_pcm_hw_params_get_period_size(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *HardwareParams) GetPeriodSizeMin() (int, int, error) {
	var val C.snd_pcm_uframes_t
	var dir C.int
	rc := C.snd_pcm_hw_params_get_period_size_min(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *HardwareParams) GetPeriodSizeMax() (int, int, error) {
	var val C.snd_pcm_uframes_t
	var dir C.int
	rc := C.snd_pcm_hw_params_get_period_size_max(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *HardwareParams) TestPeriodSize(val, dir int) error {
	rc := C.snd_pcm_hw_params_test_period_size(params.pcm.inner, params.inner, C.snd_pcm_uframes_t(val), C.int(dir))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetPeriodSize(val, dir int) error {
	rc := C.snd_pcm_hw_params_set_period_size(params.pcm.inner, params.inner, C.snd_pcm_uframes_t(val), C.int(dir))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetPeriodSizeMin(val, dir int) (int, int, error) {
	var cVal = C.snd_pcm_uframes_t(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_period_size_min(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *HardwareParams) SetPeriodSizeMax(val, dir int) (int, int, error) {
	var cVal = C.snd_pcm_uframes_t(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_period_size_max(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *HardwareParams) SetPeriodSizeMinMax(min, minDir, max, maxDir int) (int, int, int, int, error) {
	var cMin = C.snd_pcm_uframes_t(min)
	var cMinDir = C.int(minDir)
	var cMax = C.snd_pcm_uframes_t(max)
	var cMaxDir = C.int(maxDir)
	rc := C.snd_pcm_hw_params_set_period_size_minmax(params.pcm.inner, params.inner, &cMin, &cMinDir, &cMax, &cMaxDir)
	if rc < 0 {
		return 0, 0, 0, 0, alsa.NewError(int(rc))
	}
	return int(cMin), int(cMinDir), int(cMax), int(cMaxDir), nil
}

func (params *HardwareParams) SetPeriodSizeNear(val, dir int) (int, int, error) {
	var cVal = C.snd_pcm_uframes_t(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_period_size_near(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *HardwareParams) SetPeriodSizeFirst() (int, int, error) {
	var val C.snd_pcm_uframes_t
	var dir C.int
	rc := C.snd_pcm_hw_params_set_period_size_first(params.pcm.inner, params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *HardwareParams) SetPeriodSizeLast() (int, int, error) {
	var val C.snd_pcm_uframes_t
	var dir C.int
	rc := C.snd_pcm_hw_params_set_period_size_last(params.pcm.inner, params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *HardwareParams) SetPeriodSizeInteger() error {
	rc := C.snd_pcm_hw_params_set_period_size_integer(params.pcm.inner, params.inner)
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) GetPeriodsPerBuffer() (int, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_periods(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *HardwareParams) GetPeriodsPerBufferMin() (int, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_periods_min(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *HardwareParams) GetPeriodsPerBufferMax() (int, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_periods_max(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *HardwareParams) TestPeriodsPerBuffer(val, dir int) error {
	rc := C.snd_pcm_hw_params_test_periods(params.pcm.inner, params.inner, C.uint(val), C.int(dir))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetPeriodsPerBuffer(val, dir int) error {
	rc := C.snd_pcm_hw_params_set_periods(params.pcm.inner, params.inner, C.uint(val), C.int(dir))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetPeriodsPerBufferMin(val, dir int) (int, int, error) {
	var cVal = C.uint(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_periods_min(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *HardwareParams) SetPeriodsPerBufferMax(val, dir int) (int, int, error) {
	var cVal = C.uint(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_periods_max(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *HardwareParams) SetPeriodsPerBufferMinMax(min, minDir, max, maxDir int) (int, int, int, int, error) {
	var cMin = C.uint(min)
	var cMinDir = C.int(minDir)
	var cMax = C.uint(max)
	var cMaxDir = C.int(maxDir)
	rc := C.snd_pcm_hw_params_set_periods_minmax(params.pcm.inner, params.inner, &cMin, &cMinDir, &cMax, &cMaxDir)
	if rc < 0 {
		return 0, 0, 0, 0, alsa.NewError(int(rc))
	}
	return int(cMin), int(cMinDir), int(cMax), int(cMaxDir), nil
}

func (params *HardwareParams) SetPeriodsPerBufferNear(val, dir int) (int, int, error) {
	var cVal = C.uint(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_periods_near(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *HardwareParams) SetPeriodsPerBufferFirst() (int, int, error) {
	var cVal C.uint
	var cDir C.int
	rc := C.snd_pcm_hw_params_set_periods_first(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *HardwareParams) SetPeriodsPerBufferLast() (int, int, error) {
	var cVal C.uint
	var cDir C.int
	rc := C.snd_pcm_hw_params_set_periods_last(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *HardwareParams) SetPeriodsPerBufferInteger() error {
	rc := C.snd_pcm_hw_params_set_periods_integer(params.pcm.inner, params.inner)
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) GetBufferTime() (time.Duration, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_buffer_time(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(val) * time.Microsecond, int(dir), nil
}

func (params *HardwareParams) GetBufferTimeMin() (time.Duration, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_buffer_time_min(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(val) * time.Microsecond, int(dir), nil
}

func (params *HardwareParams) GetBufferTimeMax() (time.Duration, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_buffer_time_max(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(val) * time.Microsecond, int(dir), nil
}

func (params *HardwareParams) TestBufferTime(val time.Duration, dir int) error {
	rc := C.snd_pcm_hw_params_test_buffer_time(params.pcm.inner, params.inner, C.uint(val.Microseconds()), C.int(dir))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetBufferTime(val time.Duration, dir int) error {
	rc := C.snd_pcm_hw_params_set_buffer_time(params.pcm.inner, params.inner, C.uint(val.Microseconds()), C.int(dir))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetBufferTimeMin(val time.Duration, dir int) (time.Duration, int, error) {
	var cVal = C.uint(val.Microseconds())
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_buffer_time_min(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(cVal) * time.Microsecond, int(cDir), nil
}

func (params *HardwareParams) SetBufferTimeMax(val time.Duration, dir int) (time.Duration, int, error) {
	var cVal = C.uint(val.Microseconds())
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_buffer_time_max(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(cVal) * time.Microsecond, int(cDir), nil
}

func (params *HardwareParams) SetBufferTimeMinMax(
	min time.Duration, minDir int, max time.Duration, maxDir int) (time.Duration, int, time.Duration, int, error) {
	var cMin = C.uint(min.Microseconds())
	var cMinDir = C.int(minDir)
	var cMax = C.uint(max.Microseconds())
	var cMaxDir = C.int(maxDir)
	rc := C.snd_pcm_hw_params_set_buffer_time_minmax(params.pcm.inner, params.inner, &cMin, &cMinDir, &cMax, &cMaxDir)
	if rc < 0 {
		return 0, 0, 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(cMin) * time.Microsecond, int(cMinDir), time.Duration(cMax) * time.Microsecond, int(cMaxDir), nil
}

func (params *HardwareParams) SetBufferTimeNear(val time.Duration, dir int) (time.Duration, int, error) {
	var cVal = C.uint(val.Microseconds())
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_buffer_time_near(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(cVal) * time.Microsecond, int(cDir), nil
}

func (params *HardwareParams) SetBufferTimeFirst() (time.Duration, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_set_buffer_time_first(params.pcm.inner, params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(val) * time.Microsecond, int(dir), nil
}

func (params *HardwareParams) SetBufferTimeLast() (time.Duration, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_set_buffer_time_last(params.pcm.inner, params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return time.Duration(val) * time.Microsecond, int(dir), nil
}

func (params *HardwareParams) GetBufferSize() (int, error) {
	var val C.snd_pcm_uframes_t
	rc := C.snd_pcm_hw_params_get_buffer_size(params.inner, &val)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(val), nil
}

func (params *HardwareParams) GetBufferSizeMin() (int, error) {
	var val C.snd_pcm_uframes_t
	rc := C.snd_pcm_hw_params_get_buffer_size_min(params.inner, &val)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(val), nil
}

func (params *HardwareParams) GetBufferSizeMax() (int, error) {
	var val C.snd_pcm_uframes_t
	rc := C.snd_pcm_hw_params_get_buffer_size_max(params.inner, &val)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(val), nil
}

func (params *HardwareParams) TestBufferSize(val int) error {
	rc := C.snd_pcm_hw_params_test_buffer_size(params.pcm.inner, params.inner, C.snd_pcm_uframes_t(val))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetBufferSize(val int) error {
	rc := C.snd_pcm_hw_params_set_buffer_size(params.pcm.inner, params.inner, C.snd_pcm_uframes_t(val))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) SetBufferSizeMin(val int) (int, error) {
	var cVal = C.snd_pcm_uframes_t(val)
	rc := C.snd_pcm_hw_params_set_buffer_size_min(params.pcm.inner, params.inner, &cVal)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(cVal), nil
}

func (params *HardwareParams) SetBufferSizeMax(val int) (int, error) {
	var cVal = C.snd_pcm_uframes_t(val)
	rc := C.snd_pcm_hw_params_set_buffer_size_max(params.pcm.inner, params.inner, &cVal)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(cVal), nil
}

func (params *HardwareParams) SetBufferSizeMinMax(min, max int) (int, int, error) {
	var cMin = C.snd_pcm_uframes_t(min)
	var cMax = C.snd_pcm_uframes_t(max)
	rc := C.snd_pcm_hw_params_set_buffer_size_minmax(params.pcm.inner, params.inner, &cMin, &cMax)
	if rc < 0 {
		return 0, 0, alsa.NewError(int(rc))
	}
	return int(cMin), int(cMax), nil
}

func (params *HardwareParams) SetBufferSizeNear(val int) (int, error) {
	var cVal = C.snd_pcm_uframes_t(val)
	rc := C.snd_pcm_hw_params_set_buffer_size_near(params.pcm.inner, params.inner, &cVal)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(cVal), nil
}

func (params *HardwareParams) SetBufferSizeFirst() (int, error) {
	var val C.snd_pcm_uframes_t
	rc := C.snd_pcm_hw_params_set_buffer_size_first(params.pcm.inner, params.inner, &val)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(val), nil
}

func (params *HardwareParams) SetBufferSizeLast() (int, error) {
	var val C.snd_pcm_uframes_t
	rc := C.snd_pcm_hw_params_set_buffer_size_last(params.pcm.inner, params.inner, &val)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(val), nil
}

func (params *HardwareParams) GetMinAlign() (int, error) {
	var val C.snd_pcm_uframes_t
	rc := C.snd_pcm_hw_params_get_min_align(params.inner, &val)
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(val), nil
}

func (params *HardwareParams) Install() error {
	return params.pcm.InstallHardwareParams(params)
}

func (params *HardwareParams) Uninstall() error {
	return params.pcm.UninstallHardwareParams()
}
