package alsa

//
// #include <poll.h>
// #include <alsa/asoundlib.h>
//
// static char *no_const(const char *s) { return (char *)s; }
//
// static void _snd_pcm_hw_params_alloca(snd_pcm_hw_params_t **ptr) { snd_pcm_hw_params_alloca(ptr); }
//
// static void _snd_pcm_info_alloca(snd_pcm_info_t **ptr) { snd_pcm_info_alloca(ptr); }
import "C"
import (
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"
)

type StreamType int

func (stream StreamType) Name() string {
	return C.GoString(C.no_const(C.snd_pcm_stream_name(C.snd_pcm_stream_t(stream))))
}

func (stream StreamType) String() string {
	return stream.Name()
}

type PCMType int

func (ty PCMType) Name() string {
	return C.GoString(C.no_const(C.snd_pcm_type_name(C.snd_pcm_type_t(ty))))
}

func (ty PCMType) String() string {
	return ty.Name()
}

type PCMState int

func (state PCMState) Name() string {
	return C.GoString(C.no_const(C.snd_pcm_state_name(C.snd_pcm_state_t(state))))
}

type PCMAccess int

func (access PCMAccess) Name() string {
	return C.GoString(C.no_const(C.snd_pcm_access_name(C.snd_pcm_access_t(access))))
}

func (access PCMAccess) String() string {
	return access.Name()
}

type PCMFormat int

func BuildLinearFormat(width, physicalWidth int, unsigned, bigEndian bool) PCMFormat {
	return PCMFormat(C.snd_pcm_build_linear_format(C.int(width), C.int(physicalWidth),
		fromBool(unsigned), fromBool(bigEndian)))
}

func (format PCMFormat) Signed() (bool, error) {
	rc := C.snd_pcm_format_signed(C.snd_pcm_format_t(format))
	if rc < 0 {
		return false, NewError(int(rc))
	}
	return rc != 0, nil
}

func (format PCMFormat) Unsigned() (bool, error) {
	rc := C.snd_pcm_format_unsigned(C.snd_pcm_format_t(format))
	if rc < 0 {
		return false, NewError(int(rc))
	}
	return rc != 0, nil
}

func (format PCMFormat) Linear() bool {
	return C.snd_pcm_format_linear(C.snd_pcm_format_t(format)) != 0
}

func (format PCMFormat) Float() bool {
	return C.snd_pcm_format_float(C.snd_pcm_format_t(format)) != 0
}

func (format PCMFormat) LittleEndian() (bool, error) {
	rc := C.snd_pcm_format_little_endian(C.snd_pcm_format_t(format))
	if rc < 0 {
		return false, NewError(int(rc))
	}
	return rc != 0, nil
}

func (format PCMFormat) BigEndian() (bool, error) {
	rc := C.snd_pcm_format_big_endian(C.snd_pcm_format_t(format))
	if rc < 0 {
		return false, NewError(int(rc))
	}
	return rc != 0, nil
}

func (format PCMFormat) CPUEndian() (bool, error) {
	rc := C.snd_pcm_format_cpu_endian(C.snd_pcm_format_t(format))
	if rc < 0 {
		return false, NewError(int(rc))
	}
	return rc != 0, nil
}

func (format PCMFormat) Width() (int, error) {
	rc := C.snd_pcm_format_width(C.snd_pcm_format_t(format))
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (format PCMFormat) PhysicalWidth() (int, error) {
	rc := C.snd_pcm_format_physical_width(C.snd_pcm_format_t(format))
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (format PCMFormat) Size(samples int) (int, error) {
	rc := C.snd_pcm_format_size(C.snd_pcm_format_t(format), C.size_t(samples))
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (format PCMFormat) Silence() byte {
	return byte(C.snd_pcm_format_silence(C.snd_pcm_format_t(format)))
}

func (format PCMFormat) Silence16() uint16 {
	return uint16(C.snd_pcm_format_silence_16(C.snd_pcm_format_t(format)))
}

func (format PCMFormat) Silence32() uint32 {
	return uint32(C.snd_pcm_format_silence_32(C.snd_pcm_format_t(format)))
}

func (format PCMFormat) Silence64() uint64 {
	return uint64(C.snd_pcm_format_silence_64(C.snd_pcm_format_t(format)))
}

func (format PCMFormat) SetSilenceData(data unsafe.Pointer, samples int) error {
	rc := C.snd_pcm_format_set_silence(C.snd_pcm_format_t(format), data, C.uint(samples))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (format PCMFormat) Name() string {
	return C.GoString(C.no_const(C.snd_pcm_format_name(C.snd_pcm_format_t(format))))
}

func (format PCMFormat) Description() string {
	return C.GoString(C.no_const(C.snd_pcm_format_description(C.snd_pcm_format_t(format))))
}

func (format PCMFormat) String() string {
	return format.Description()
}

func PCMFormatName(name string) PCMFormat {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return PCMFormat(C.snd_pcm_format_value(cName))
}

type PCMClass int
type PCMSubClass int

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

	SndPcmAccessMmapInterleaved    = PCMAccess(C.SND_PCM_ACCESS_MMAP_INTERLEAVED)
	SndPcmAccessMmapNoninterleaved = PCMAccess(C.SND_PCM_ACCESS_MMAP_NONINTERLEAVED)
	SndPcmAccessMmapComplex        = PCMAccess(C.SND_PCM_ACCESS_MMAP_COMPLEX)
	SndPcmAccessRwInterleaved      = PCMAccess(C.SND_PCM_ACCESS_RW_INTERLEAVED)
	SndPcmAccessRwNoninterleaved   = PCMAccess(C.SND_PCM_ACCESS_RW_NONINTERLEAVED)

	SndPcmFormatUnknown          = PCMFormat(C.SND_PCM_FORMAT_UNKNOWN)
	SndPcmFormatS8               = PCMFormat(C.SND_PCM_FORMAT_S8)
	SndPcmFormatU8               = PCMFormat(C.SND_PCM_FORMAT_U8)
	SndPcmFormatS16Le            = PCMFormat(C.SND_PCM_FORMAT_S16_LE)
	SndPcmFormatS16Be            = PCMFormat(C.SND_PCM_FORMAT_S16_BE)
	SndPcmFormatU16Le            = PCMFormat(C.SND_PCM_FORMAT_U16_LE)
	SndPcmFormatU16Be            = PCMFormat(C.SND_PCM_FORMAT_U16_BE)
	SndPcmFormatS24Le            = PCMFormat(C.SND_PCM_FORMAT_S24_LE)
	SndPcmFormatS24Be            = PCMFormat(C.SND_PCM_FORMAT_S24_BE)
	SndPcmFormatU24Le            = PCMFormat(C.SND_PCM_FORMAT_U24_LE)
	SndPcmFormatU24Be            = PCMFormat(C.SND_PCM_FORMAT_U24_BE)
	SndPcmFormatS32Le            = PCMFormat(C.SND_PCM_FORMAT_S32_LE)
	SndPcmFormatS32Be            = PCMFormat(C.SND_PCM_FORMAT_S32_BE)
	SndPcmFormatU32Le            = PCMFormat(C.SND_PCM_FORMAT_U32_LE)
	SndPcmFormatU32Be            = PCMFormat(C.SND_PCM_FORMAT_U32_BE)
	SndPcmFormatFloatLe          = PCMFormat(C.SND_PCM_FORMAT_FLOAT_LE)
	SndPcmFormatFloatBe          = PCMFormat(C.SND_PCM_FORMAT_FLOAT_BE)
	SndPcmFormatFloat64Le        = PCMFormat(C.SND_PCM_FORMAT_FLOAT64_LE)
	SndPcmFormatFloat64Be        = PCMFormat(C.SND_PCM_FORMAT_FLOAT64_BE)
	SndPcmFormatIec958SubframeLe = PCMFormat(C.SND_PCM_FORMAT_IEC958_SUBFRAME_LE)
	SndPcmFormatIec958SubframeBe = PCMFormat(C.SND_PCM_FORMAT_IEC958_SUBFRAME_BE)
	SndPcmFormatMuLaw            = PCMFormat(C.SND_PCM_FORMAT_MU_LAW)
	SndPcmFormatALaw             = PCMFormat(C.SND_PCM_FORMAT_A_LAW)
	SndPcmFormatImaAdpcm         = PCMFormat(C.SND_PCM_FORMAT_IMA_ADPCM)
	SndPcmFormatMpeg             = PCMFormat(C.SND_PCM_FORMAT_MPEG)
	SndPcmFormatGsm              = PCMFormat(C.SND_PCM_FORMAT_GSM)
	SndPcmFormatS20Le            = PCMFormat(C.SND_PCM_FORMAT_S20_LE)
	SndPcmFormatS20Be            = PCMFormat(C.SND_PCM_FORMAT_S20_BE)
	SndPcmFormatU20Le            = PCMFormat(C.SND_PCM_FORMAT_U20_LE)
	SndPcmFormatU20Be            = PCMFormat(C.SND_PCM_FORMAT_U20_BE)
	SndPcmFormatSpecial          = PCMFormat(C.SND_PCM_FORMAT_SPECIAL)
	SndPcmFormatS243le           = PCMFormat(C.SND_PCM_FORMAT_S24_3LE)
	SndPcmFormatS243be           = PCMFormat(C.SND_PCM_FORMAT_S24_3BE)
	SndPcmFormatU243le           = PCMFormat(C.SND_PCM_FORMAT_U24_3LE)
	SndPcmFormatU243be           = PCMFormat(C.SND_PCM_FORMAT_U24_3BE)
	SndPcmFormatS203le           = PCMFormat(C.SND_PCM_FORMAT_S20_3LE)
	SndPcmFormatS203be           = PCMFormat(C.SND_PCM_FORMAT_S20_3BE)
	SndPcmFormatU203le           = PCMFormat(C.SND_PCM_FORMAT_U20_3LE)
	SndPcmFormatU203be           = PCMFormat(C.SND_PCM_FORMAT_U20_3BE)
	SndPcmFormatS183le           = PCMFormat(C.SND_PCM_FORMAT_S18_3LE)
	SndPcmFormatS183be           = PCMFormat(C.SND_PCM_FORMAT_S18_3BE)
	SndPcmFormatU183le           = PCMFormat(C.SND_PCM_FORMAT_U18_3LE)
	SndPcmFormatU183be           = PCMFormat(C.SND_PCM_FORMAT_U18_3BE)
	SndPcmFormatG72324           = PCMFormat(C.SND_PCM_FORMAT_G723_24)
	SndPcmFormatG723241b         = PCMFormat(C.SND_PCM_FORMAT_G723_24_1B)
	SndPcmFormatG72340           = PCMFormat(C.SND_PCM_FORMAT_G723_40)
	SndPcmFormatG723401b         = PCMFormat(C.SND_PCM_FORMAT_G723_40_1B)
	SndPcmFormatDsdU8            = PCMFormat(C.SND_PCM_FORMAT_DSD_U8)
	SndPcmFormatDsdU16Le         = PCMFormat(C.SND_PCM_FORMAT_DSD_U16_LE)
	SndPcmFormatDsdU32Le         = PCMFormat(C.SND_PCM_FORMAT_DSD_U32_LE)
	SndPcmFormatDsdU16Be         = PCMFormat(C.SND_PCM_FORMAT_DSD_U16_BE)
	SndPcmFormatDsdU32Be         = PCMFormat(C.SND_PCM_FORMAT_DSD_U32_BE)

	SndPcmClassGeneric   = PCMClass(C.SND_PCM_CLASS_GENERIC)
	SndPcmClassMulti     = PCMClass(C.SND_PCM_CLASS_MULTI)
	SndPcmClassModem     = PCMClass(C.SND_PCM_CLASS_MODEM)
	SndPcmClassDigitizer = PCMClass(C.SND_PCM_CLASS_DIGITIZER)

	SndPcmSubclassGenericMix = PCMSubClass(C.SND_PCM_SUBCLASS_GENERIC_MIX)
	SndPcmSubclassMultiMix   = PCMSubClass(C.SND_PCM_SUBCLASS_MULTI_MIX)
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
	if rc < 0 {
		return nil, NewError(int(rc))
	}
	// add gc flags
	runtime.SetFinalizer(pcm, (*PCM).Close)
	return pcm, nil
}

func (pcm *PCM) Close() error {
	for {
		p := unsafe.Pointer(pcm.inner)
		c := pcm.inner
		if p == nil {
			break
		}
		if atomic.CompareAndSwapPointer(&p, p, nil) {
			rc := C.snd_pcm_close(c)
			if rc < 0 {
				panic(NewError(int(rc)))
			}
			// clear gc flags
			runtime.SetFinalizer(pcm, nil)
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
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Info() (*PCMInfo, error) {
	info := new(PCMInfo)
	C._snd_pcm_info_alloca(&info.inner)
	rc := C.snd_pcm_info(pcm.inner, info.inner)
	if rc < 0 {
		return nil, NewError(int(rc))
	}
	runtime.SetFinalizer(info, (*PCMInfo).Close)
	return info, nil
}

func (pcm *PCM) Prepare() error {
	rc := C.snd_pcm_prepare(pcm.inner)
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Reset() error {
	rc := C.snd_pcm_reset(pcm.inner)
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Start() error {
	rc := C.snd_pcm_start(pcm.inner)
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Drop() error {
	rc := C.snd_pcm_drop(pcm.inner)
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Drain() error {
	rc := C.snd_pcm_drain(pcm.inner)
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Pause(enable bool) error {
	rc := C.snd_pcm_pause(pcm.inner, fromBool(enable))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) State() PCMState {
	return PCMState(C.snd_pcm_state(pcm.inner))
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
		return 0, NewError(int(rc))
	}
	return int(delay), nil
}

func (pcm *PCM) Resume() error {
	rc := C.snd_pcm_resume(pcm.inner)
	if rc < 0 {
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
	if rc < 0 {
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

func (pcm *PCM) Writei(data unsafe.Pointer, frames int) (int, error) {
	if data == nil {
		return 0, nil
	}
	rc := C.snd_pcm_writei(pcm.inner, data, C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Readi(data unsafe.Pointer, frames int) (int, error) {
	if data == nil {
		return 0, nil
	}
	rc := C.snd_pcm_readi(pcm.inner, data, C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Writen(data []unsafe.Pointer, frames int) (int, error) {
	if len(data) <= 0 {
		return 0, nil
	}
	rc := C.snd_pcm_writen(pcm.inner, &data[0], C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Readn(data []unsafe.Pointer, frames int) (int, error) {
	if len(data) <= 0 {
		return 0, nil
	}
	rc := C.snd_pcm_readn(pcm.inner, &data[0], C.snd_pcm_uframes_t(frames))
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (pcm *PCM) Link(other *PCM) error {
	rc := C.snd_pcm_link(pcm.inner, other.inner)
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) Unlink() error {
	rc := C.snd_pcm_unlink(pcm.inner)
	if rc < 0 {
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
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int16(revents), nil
}

func (pcm *PCM) InstallHwParams(params *PCMHwParams) error {
	rc := C.snd_pcm_hw_params(pcm.inner, params.inner)
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (pcm *PCM) UninstallHwParams() error {
	rc := C.snd_pcm_hw_free(pcm.inner)
	if rc < 0 {
		return NewError(int(rc))
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

type PCMHwParams struct {
	inner *C.snd_pcm_hw_params_t
	pcm   *PCM
}

func PCMHwParamsAny(pcm *PCM) (*PCMHwParams, error) {
	params := &PCMHwParams{}
	C._snd_pcm_hw_params_alloca(&params.inner)
	runtime.SetFinalizer(params, (*PCMHwParams).Close)

	rc := C.snd_pcm_hw_params_any(pcm.inner, params.inner)
	if rc < 0 {
		return nil, NewError(int(rc))
	}
	params.pcm = pcm

	return params, nil
}

func PCMHwParamsCurrent(pcm *PCM) (*PCMHwParams, error) {
	params := &PCMHwParams{}
	C._snd_pcm_hw_params_alloca(&params.inner)
	runtime.SetFinalizer(params, (*PCMHwParams).Close)

	rc := C.snd_pcm_hw_params_current(pcm.inner, params.inner)
	if rc < 0 {
		return nil, NewError(int(rc))
	}
	params.pcm = pcm

	return params, nil
}

func (params *PCMHwParams) Close() error {
	for {
		p := unsafe.Pointer(params.inner)
		i := params.inner
		if p == nil {
			break
		}
		if atomic.CompareAndSwapPointer(&p, p, nil) {
			C.snd_pcm_hw_params_free(i)
			runtime.SetFinalizer(params, nil)
		}
	}

	return nil
}

func (params *PCMHwParams) CanMmapSampleResolution() bool {
	return C.snd_pcm_hw_params_can_mmap_sample_resolution(params.inner) != 0
}

func (params *PCMHwParams) IsDouble() bool {
	return C.snd_pcm_hw_params_is_double(params.inner) != 0
}

func (params *PCMHwParams) IsBatch() bool {
	return C.snd_pcm_hw_params_is_batch(params.inner) != 0
}

func (params *PCMHwParams) IsBlockTransfer() bool {
	return C.snd_pcm_hw_params_is_block_transfer(params.inner) != 0
}

func (params *PCMHwParams) IsMonotonic() bool {
	return C.snd_pcm_hw_params_is_monotonic(params.inner) != 0
}

func (params *PCMHwParams) CanOverRange() bool {
	return C.snd_pcm_hw_params_can_overrange(params.inner) != 0
}

func (params *PCMHwParams) CanPause() bool {
	return C.snd_pcm_hw_params_can_pause(params.inner) != 0
}

func (params *PCMHwParams) CanResume() bool {
	return C.snd_pcm_hw_params_can_resume(params.inner) != 0
}

func (params *PCMHwParams) IsHalfDuplex() bool {
	return C.snd_pcm_hw_params_is_half_duplex(params.inner) != 0
}

func (params *PCMHwParams) IsJointDuplex() bool {
	return C.snd_pcm_hw_params_is_joint_duplex(params.inner) != 0
}

func (params *PCMHwParams) CanSyncStart() bool {
	return C.snd_pcm_hw_params_can_sync_start(params.inner) != 0
}

func (params *PCMHwParams) CanDisablePeriodWakeup() bool {
	return C.snd_pcm_hw_params_can_disable_period_wakeup(params.inner) != 0
}

func (params *PCMHwParams) SupportsAudioWallClockTimestamps() bool {
	return C.snd_pcm_hw_params_supports_audio_wallclock_ts(params.inner) != 0
}

func (params *PCMHwParams) SupportsAudioTimestampType(tsType int) bool {
	return C.snd_pcm_hw_params_supports_audio_ts_type(params.inner, C.int(tsType)) != 0
}

func (params *PCMHwParams) GetRateNumDen() (int, int, error) {
	var rateNum, rateDen C.uint
	rc := C.snd_pcm_hw_params_get_rate_numden(params.inner, &rateNum, &rateDen)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(rateNum), int(rateDen), nil
}

func (params *PCMHwParams) GetSBits() (int, error) {
	rc := C.snd_pcm_hw_params_get_sbits(params.inner)
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (params *PCMHwParams) GetFifoSize() (int, error) {
	rc := C.snd_pcm_hw_params_get_fifo_size(params.inner)
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(rc), nil
}

func (params *PCMHwParams) GetAccess() (PCMAccess, error) {
	var acs C.snd_pcm_access_t
	rc := C.snd_pcm_hw_params_get_access(params.inner, &acs)
	if rc < 0 {
		return -1, NewError(int(rc))
	}
	return PCMAccess(acs), nil
}

func (params *PCMHwParams) TestAccess(acs PCMAccess) error {
	rc := C.snd_pcm_hw_params_test_access(params.pcm.inner, params.inner, C.snd_pcm_access_t(acs))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetAccess(acs PCMAccess) error {
	rc := C.snd_pcm_hw_params_set_access(params.pcm.inner, params.inner, C.snd_pcm_access_t(acs))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetAccessFirst() (PCMAccess, error) {
	var access C.snd_pcm_access_t
	rc := C.snd_pcm_hw_params_set_access_first(params.pcm.inner, params.inner, &access)
	if rc < 0 {
		return -1, NewError(int(rc))
	}
	return PCMAccess(access), nil
}

func (params *PCMHwParams) SetAccessLast() (PCMAccess, error) {
	var access C.snd_pcm_access_t
	rc := C.snd_pcm_hw_params_set_access_last(params.pcm.inner, params.inner, &access)
	if rc < 0 {
		return -1, NewError(int(rc))
	}
	return PCMAccess(access), nil
}

func (params *PCMHwParams) GetFormat() (PCMFormat, error) {
	var format C.snd_pcm_format_t
	rc := C.snd_pcm_hw_params_get_format(params.inner, &format)
	if rc < 0 {
		return -1, NewError(int(rc))
	}
	return PCMFormat(format), nil
}

func (params *PCMHwParams) TestFormat(format PCMFormat) error {
	rc := C.snd_pcm_hw_params_test_format(params.pcm.inner, params.inner, C.snd_pcm_format_t(format))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetFormat(format PCMFormat) error {
	rc := C.snd_pcm_hw_params_set_format(params.pcm.inner, params.inner, C.snd_pcm_format_t(format))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetFormatFirst() (PCMFormat, error) {
	var format C.snd_pcm_format_t
	rc := C.snd_pcm_hw_params_set_format_first(params.pcm.inner, params.inner, &format)
	if rc < 0 {
		return SndPcmFormatUnknown, NewError(int(rc))
	}
	return PCMFormat(format), nil
}

func (params *PCMHwParams) SetFormatLast() (PCMFormat, error) {
	var format C.snd_pcm_format_t
	rc := C.snd_pcm_hw_params_set_format_last(params.pcm.inner, params.inner, &format)
	if rc < 0 {
		return SndPcmFormatUnknown, NewError(int(rc))
	}
	return PCMFormat(format), nil
}

func (params *PCMHwParams) GetChannels() (int, error) {
	var channels C.uint
	rc := C.snd_pcm_hw_params_get_channels(params.inner, &channels)
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(channels), nil
}

func (params *PCMHwParams) GetChannelsMin() (int, error) {
	var channels C.uint
	rc := C.snd_pcm_hw_params_get_channels_min(params.inner, &channels)
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(channels), nil
}

func (params *PCMHwParams) GetChannelsMax() (int, error) {
	var channels C.uint
	rc := C.snd_pcm_hw_params_get_channels_max(params.inner, &channels)
	if rc < 0 {
		return 0, NewError(int(rc))
	}
	return int(channels), nil
}

func (params *PCMHwParams) TestChannels(channels int) error {
	rc := C.snd_pcm_hw_params_test_channels(params.pcm.inner, params.inner, C.uint(channels))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetChannels(channels int) error {
	rc := C.snd_pcm_hw_params_set_channels(params.pcm.inner, params.inner, C.uint(channels))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetChannelsMin(min int) (int, error) {
	var channels = C.uint(min)
	rc := C.snd_pcm_hw_params_set_channels_min(params.pcm.inner, params.inner, &channels)
	if rc < 0 {
		return -1, NewError(int(rc))
	}
	return int(channels), nil
}

func (params *PCMHwParams) SetChannelsMax(max int) (int, error) {
	var channels = C.uint(max)
	rc := C.snd_pcm_hw_params_set_channels_max(params.pcm.inner, params.inner, &channels)
	if rc < 0 {
		return -1, NewError(int(rc))
	}
	return int(channels), nil
}

func (params *PCMHwParams) SetChannelsMinMax(min, max int) (int, int, error) {
	var cMin = C.uint(min)
	var cMax = C.uint(max)
	rc := C.snd_pcm_hw_params_set_channels_minmax(params.pcm.inner, params.inner, &cMin, &cMax)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cMin), int(cMax), nil
}

func (params *PCMHwParams) SetChannelsNear(channels int) (int, error) {
	cArg := C.uint(channels)
	rc := C.snd_pcm_hw_params_set_channels_near(params.pcm.inner, params.inner, &cArg)
	if rc < 0 {
		return -1, NewError(int(rc))
	}
	return int(cArg), nil
}

func (params *PCMHwParams) SetChannelsFirst() (int, error) {
	var channels C.uint
	rc := C.snd_pcm_hw_params_set_channels_first(params.pcm.inner, params.inner, &channels)
	if rc < 0 {
		return -1, NewError(int(rc))
	}
	return int(channels), nil
}

func (params *PCMHwParams) SetChannelsLast() (int, error) {
	var channels C.uint
	rc := C.snd_pcm_hw_params_set_channels_last(params.pcm.inner, params.inner, &channels)
	if rc < 0 {
		return -1, NewError(int(rc))
	}
	return int(channels), nil
}

func (params *PCMHwParams) GetRate() (int, int, error) {
	var rate C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_rate(params.inner, &rate, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(rate), int(dir), nil
}

func (params *PCMHwParams) GetRateMin() (int, int, error) {
	var rate C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_rate_min(params.inner, &rate, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(rate), int(dir), nil
}

func (params *PCMHwParams) GetRateMax() (int, int, error) {
	var rate C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_rate_max(params.inner, &rate, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(rate), int(dir), nil
}

func (params *PCMHwParams) TestRate(rate, dir int) error {
	rc := C.snd_pcm_hw_params_test_rate(params.pcm.inner, params.inner, C.uint(rate), C.int(dir))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetRate(rate, dir int) error {
	rc := C.snd_pcm_hw_params_set_rate(params.pcm.inner, params.inner, C.uint(rate), C.int(dir))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetRateMin(rate, dir int) (int, int, error) {
	var cRate = C.uint(rate)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_rate_min(params.pcm.inner, params.inner, &cRate, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cRate), int(cDir), nil
}

func (params *PCMHwParams) SetRateMax(rate, dir int) (int, int, error) {
	var cRate = C.uint(rate)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_rate_max(params.pcm.inner, params.inner, &cRate, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cRate), int(cDir), nil
}

func (params *PCMHwParams) SetRateMinMax(min, minDir, max, maxDir int) (int, int, int, int, error) {
	var cMin = C.uint(min)
	var cMinDir = C.int(minDir)
	var cMax = C.uint(max)
	var cMaxDir = C.int(maxDir)
	rc := C.snd_pcm_hw_params_set_rate_minmax(params.pcm.inner, params.inner, &cMin, &cMinDir, &cMax, &cMaxDir)
	if rc < 0 {
		return 0, 0, 0, 0, NewError(int(rc))
	}
	return int(cMin), int(cMinDir), int(cMax), int(cMaxDir), nil
}

func (params *PCMHwParams) SetRateNear(rate, dir int) (int, int, error) {
	var cRate = C.uint(rate)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_rate_near(params.pcm.inner, params.inner, &cRate, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cRate), int(cDir), nil
}

func (params *PCMHwParams) SetRateFirst() (int, int, error) {
	var rate C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_set_rate_first(params.pcm.inner, params.inner, &rate, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(rate), int(dir), nil
}

func (params *PCMHwParams) SetRateLast() (int, int, error) {
	var rate C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_set_rate_last(params.pcm.inner, params.inner, &rate, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(rate), int(dir), nil
}

func (params *PCMHwParams) SetRateResample(enable bool) error {
	rc := C.snd_pcm_hw_params_set_rate_resample(params.pcm.inner, params.inner, C.uint(fromBool(enable)))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) GetRateResample() (bool, error) {
	var enable C.uint
	rc := C.snd_pcm_hw_params_get_rate_resample(params.pcm.inner, params.inner, &enable)
	if rc < 0 {
		return false, NewError(int(rc))
	}
	return enable != 0, nil
}

func (params *PCMHwParams) SetExportBuffer(enable bool) error {
	rc := C.snd_pcm_hw_params_set_export_buffer(params.pcm.inner, params.inner, C.uint(fromBool(enable)))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) GetExportBuffer() (bool, error) {
	var enable C.uint
	rc := C.snd_pcm_hw_params_get_export_buffer(params.pcm.inner, params.inner, &enable)
	if rc < 0 {
		return false, NewError(int(rc))
	}
	return enable != 0, nil
}

func (params *PCMHwParams) SetPeriodWakeup(enable bool) error {
	rc := C.snd_pcm_hw_params_set_period_wakeup(params.pcm.inner, params.inner, C.uint(fromBool(enable)))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) GetPeriodWakeup() (bool, error) {
	var enable C.uint
	rc := C.snd_pcm_hw_params_get_period_wakeup(params.pcm.inner, params.inner, &enable)
	if rc < 0 {
		return false, NewError(int(rc))
	}
	return enable != 0, nil
}

func (params *PCMHwParams) GetPeriodTime() (int, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_period_time(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *PCMHwParams) GetPeriodTimeMin() (int, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_period_time_min(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *PCMHwParams) GetPeriodTimeMax() (int, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_period_time_max(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *PCMHwParams) TestPeriodTime(val, dir int) error {
	rc := C.snd_pcm_hw_params_test_period_time(params.pcm.inner, params.inner, C.uint(val), C.int(dir))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetPeriodTime(val, dir int) error {
	rc := C.snd_pcm_hw_params_set_period_time(params.pcm.inner, params.inner, C.uint(val), C.int(dir))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetPeriodTimeMin(val, dir int) (int, int, error) {
	var cVal = C.uint(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_period_time_min(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *PCMHwParams) SetPeriodTimeMax(val, dir int) (int, int, error) {
	var cVal = C.uint(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_period_time_max(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *PCMHwParams) SetPeriodTimeMinMax(min, minDir, max, maxDir int) (int, int, int, int, error) {
	var cMin = C.uint(min)
	var cMinDir = C.int(minDir)
	var cMax = C.uint(max)
	var cMaxDir = C.int(maxDir)
	rc := C.snd_pcm_hw_params_set_period_time_minmax(params.pcm.inner, params.inner, &cMin, &cMinDir, &cMax, &cMaxDir)
	if rc < 0 {
		return 0, 0, 0, 0, NewError(int(rc))
	}
	return int(cMin), int(cMinDir), int(cMax), int(cMaxDir), nil
}

func (params *PCMHwParams) SetPeriodTimeNear(val, dir int) (int, int, error) {
	var cVal = C.uint(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_period_time_near(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *PCMHwParams) SetPeriodTimeFirst() (int, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_set_period_time_first(params.pcm.inner, params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *PCMHwParams) SetPeriodTimeLast() (int, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_set_period_time_last(params.pcm.inner, params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *PCMHwParams) GetPeriodSize() (int, int, error) {
	var val C.snd_pcm_uframes_t
	var dir C.int
	rc := C.snd_pcm_hw_params_get_period_size(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *PCMHwParams) GetPeriodSizeMin() (int, int, error) {
	var val C.snd_pcm_uframes_t
	var dir C.int
	rc := C.snd_pcm_hw_params_get_period_size_min(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *PCMHwParams) GetPeriodSizeMax() (int, int, error) {
	var val C.snd_pcm_uframes_t
	var dir C.int
	rc := C.snd_pcm_hw_params_get_period_size_max(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *PCMHwParams) TestPeriodSize(val, dir int) error {
	rc := C.snd_pcm_hw_params_test_period_size(params.pcm.inner, params.inner, C.snd_pcm_uframes_t(val), C.int(dir))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetPeriodSize(val, dir int) error {
	rc := C.snd_pcm_hw_params_set_period_size(params.pcm.inner, params.inner, C.snd_pcm_uframes_t(val), C.int(dir))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetPeriodSizeMin(val, dir int) (int, int, error) {
	var cVal = C.snd_pcm_uframes_t(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_period_size_min(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *PCMHwParams) SetPeriodSizeMax(val, dir int) (int, int, error) {
	var cVal = C.snd_pcm_uframes_t(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_period_size_max(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *PCMHwParams) SetPeriodSizeMinMax(min, minDir, max, maxDir int) (int, int, int, int, error) {
	var cMin = C.snd_pcm_uframes_t(min)
	var cMinDir = C.int(minDir)
	var cMax = C.snd_pcm_uframes_t(max)
	var cMaxDir = C.int(maxDir)
	rc := C.snd_pcm_hw_params_set_period_size_minmax(params.pcm.inner, params.inner, &cMin, &cMinDir, &cMax, &cMaxDir)
	if rc < 0 {
		return 0, 0, 0, 0, NewError(int(rc))
	}
	return int(cMin), int(cMinDir), int(cMax), int(cMaxDir), nil
}

func (params *PCMHwParams) SetPeriodSizeNear(val, dir int) (int, int, error) {
	var cVal = C.snd_pcm_uframes_t(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_period_size_near(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *PCMHwParams) SetPeriodSizeFirst() (int, int, error) {
	var val C.snd_pcm_uframes_t
	var dir C.int
	rc := C.snd_pcm_hw_params_set_period_size_first(params.pcm.inner, params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *PCMHwParams) SetPeriodSizeLast() (int, int, error) {
	var val C.snd_pcm_uframes_t
	var dir C.int
	rc := C.snd_pcm_hw_params_set_period_size_last(params.pcm.inner, params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *PCMHwParams) SetPeriodSizeInteger() error {
	rc := C.snd_pcm_hw_params_set_period_size_integer(params.pcm.inner, params.inner)
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) GetPeriods() (int, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_periods(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *PCMHwParams) GetPeriodsMin() (int, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_periods_min(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *PCMHwParams) GetPeriodsMax() (int, int, error) {
	var val C.uint
	var dir C.int
	rc := C.snd_pcm_hw_params_get_periods_max(params.inner, &val, &dir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(val), int(dir), nil
}

func (params *PCMHwParams) TestPeriods(val, dir int) error {
	rc := C.snd_pcm_hw_params_test_periods(params.pcm.inner, params.inner, C.uint(val), C.int(dir))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetPeriods(val, dir int) error {
	rc := C.snd_pcm_hw_params_set_periods(params.pcm.inner, params.inner, C.uint(val), C.int(dir))
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) SetPeriodsMin(val, dir int) (int, int, error) {
	var cVal = C.uint(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_periods_min(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *PCMHwParams) SetPeriodsMax(val, dir int) (int, int, error) {
	var cVal = C.uint(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_periods_max(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *PCMHwParams) SetPeriodsMinMax(min, minDir, max, maxDir int) (int, int, int, int, error) {
	var cMin = C.uint(min)
	var cMinDir = C.int(minDir)
	var cMax = C.uint(max)
	var cMaxDir = C.int(maxDir)
	rc := C.snd_pcm_hw_params_set_periods_minmax(params.pcm.inner, params.inner, &cMin, &cMinDir, &cMax, &cMaxDir)
	if rc < 0 {
		return 0, 0, 0, 0, NewError(int(rc))
	}
	return int(cMin), int(cMinDir), int(cMax), int(cMaxDir), nil
}

func (params *PCMHwParams) SetPeriodsNear(val, dir int) (int, int, error) {
	var cVal = C.uint(val)
	var cDir = C.int(dir)
	rc := C.snd_pcm_hw_params_set_periods_near(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *PCMHwParams) SetPeriodsFirst() (int, int, error) {
	var cVal C.uint
	var cDir C.int
	rc := C.snd_pcm_hw_params_set_periods_first(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *PCMHwParams) SetPeriodsLast() (int, int, error) {
	var cVal C.uint
	var cDir C.int
	rc := C.snd_pcm_hw_params_set_periods_last(params.pcm.inner, params.inner, &cVal, &cDir)
	if rc < 0 {
		return 0, 0, NewError(int(rc))
	}
	return int(cVal), int(cDir), nil
}

func (params *PCMHwParams) SetPeriodsInteger() error {
	rc := C.snd_pcm_hw_params_set_periods_integer(params.pcm.inner, params.inner)
	if rc < 0 {
		return NewError(int(rc))
	}
	return nil
}

func (params *PCMHwParams) Install() error {
	return params.pcm.InstallHwParams(params)
}

func (params *PCMHwParams) Uninstall() error {
	return params.pcm.UninstallHwParams()
}

type PCMInfo struct {
	inner *C.snd_pcm_info_t
}

func (info *PCMInfo) Close() error {
	for {
		p := unsafe.Pointer(info.inner)
		i := info.inner
		if p == nil {
			break
		}
		if atomic.CompareAndSwapPointer(&p, p, nil) {
			C.snd_pcm_info_free(i)
			runtime.SetFinalizer(info, nil)
		}
	}
	return nil
}

func (info *PCMInfo) GetDevice() int {
	return int(C.snd_pcm_info_get_device(info.inner))
}

func (info *PCMInfo) GetSubDevice() int {
	return int(C.snd_pcm_info_get_subdevice(info.inner))
}

func (info *PCMInfo) GetStreamType() StreamType {
	return StreamType(C.snd_pcm_info_get_stream(info.inner))
}

func (info *PCMInfo) GetCard() (int, error) {
	rc := C.snd_pcm_info_get_card(info.inner)
	if rc < 0 {
		return -1, NewError(int(rc))
	}
	return int(rc), nil
}

func (info *PCMInfo) GetID() string {
	return C.GoString(C.no_const(C.snd_pcm_info_get_id(info.inner)))
}

func (info *PCMInfo) GetName() string {
	return C.GoString(C.no_const(C.snd_pcm_info_get_name(info.inner)))
}

func (info *PCMInfo) GetSubDeviceName() string {
	return C.GoString(C.no_const(C.snd_pcm_info_get_subdevice_name(info.inner)))
}

func (info *PCMInfo) GetClass() PCMClass {
	return PCMClass(C.snd_pcm_info_get_class(info.inner))
}

func (info *PCMInfo) GetSubClass() PCMSubClass {
	return PCMSubClass(C.snd_pcm_info_get_subclass(info.inner))
}

func (info *PCMInfo) GetSubDevicesCount() int {
	return int(C.snd_pcm_info_get_subdevices_count(info.inner))
}

func (info *PCMInfo) GetSubDevicesAvailableCount() int {
	return int(C.snd_pcm_info_get_subdevices_avail(info.inner))
}
