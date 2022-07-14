package pcm

//
// #include <poll.h>
// #include <alsa/asoundlib.h>
//
// static char *no_const(const char *s) { return (char *)s; }
//
import "C"
import (
	"unsafe"

	"github.com/ClarkGuan/go-alsa"
)

type StreamType int

func (stream StreamType) Name() string {
	return C.GoString(C.no_const(C.snd_pcm_stream_name(C.snd_pcm_stream_t(stream))))
}

func (stream StreamType) String() string {
	return stream.Name()
}

type Type int

func (ty Type) Name() string {
	return C.GoString(C.no_const(C.snd_pcm_type_name(C.snd_pcm_type_t(ty))))
}

func (ty Type) String() string {
	return ty.Name()
}

type State int

func (state State) Name() string {
	return C.GoString(C.no_const(C.snd_pcm_state_name(C.snd_pcm_state_t(state))))
}

func (state State) String() string {
	return state.Name()
}

type Access int

func (access Access) Name() string {
	return C.GoString(C.no_const(C.snd_pcm_access_name(C.snd_pcm_access_t(access))))
}

func (access Access) String() string {
	return access.Name()
}

type Format int

func BuildLinearFormat(width, physicalWidth int, unsigned, bigEndian bool) Format {
	return Format(C.snd_pcm_build_linear_format(C.int(width), C.int(physicalWidth),
		fromBool(unsigned), fromBool(bigEndian)))
}

func (format Format) Signed() (bool, error) {
	rc := C.snd_pcm_format_signed(C.snd_pcm_format_t(format))
	if rc < 0 {
		return false, alsa.NewError(int(rc))
	}
	return rc != 0, nil
}

func (format Format) Unsigned() (bool, error) {
	rc := C.snd_pcm_format_unsigned(C.snd_pcm_format_t(format))
	if rc < 0 {
		return false, alsa.NewError(int(rc))
	}
	return rc != 0, nil
}

func (format Format) Linear() bool {
	return C.snd_pcm_format_linear(C.snd_pcm_format_t(format)) != 0
}

func (format Format) Float() bool {
	return C.snd_pcm_format_float(C.snd_pcm_format_t(format)) != 0
}

func (format Format) LittleEndian() (bool, error) {
	rc := C.snd_pcm_format_little_endian(C.snd_pcm_format_t(format))
	if rc < 0 {
		return false, alsa.NewError(int(rc))
	}
	return rc != 0, nil
}

func (format Format) BigEndian() (bool, error) {
	rc := C.snd_pcm_format_big_endian(C.snd_pcm_format_t(format))
	if rc < 0 {
		return false, alsa.NewError(int(rc))
	}
	return rc != 0, nil
}

func (format Format) CPUEndian() (bool, error) {
	rc := C.snd_pcm_format_cpu_endian(C.snd_pcm_format_t(format))
	if rc < 0 {
		return false, alsa.NewError(int(rc))
	}
	return rc != 0, nil
}

func (format Format) Width() (int, error) {
	rc := C.snd_pcm_format_width(C.snd_pcm_format_t(format))
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(rc), nil
}

func (format Format) PhysicalWidth() (int, error) {
	rc := C.snd_pcm_format_physical_width(C.snd_pcm_format_t(format))
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(rc), nil
}

func (format Format) Size(samples int) (int, error) {
	rc := C.snd_pcm_format_size(C.snd_pcm_format_t(format), C.size_t(samples))
	if rc < 0 {
		return 0, alsa.NewError(int(rc))
	}
	return int(rc), nil
}

func (format Format) Silence() byte {
	return byte(C.snd_pcm_format_silence(C.snd_pcm_format_t(format)))
}

func (format Format) Silence16() uint16 {
	return uint16(C.snd_pcm_format_silence_16(C.snd_pcm_format_t(format)))
}

func (format Format) Silence32() uint32 {
	return uint32(C.snd_pcm_format_silence_32(C.snd_pcm_format_t(format)))
}

func (format Format) Silence64() uint64 {
	return uint64(C.snd_pcm_format_silence_64(C.snd_pcm_format_t(format)))
}

func (format Format) SetSilenceData(data unsafe.Pointer, samples int) error {
	rc := C.snd_pcm_format_set_silence(C.snd_pcm_format_t(format), data, C.uint(samples))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (format Format) Name() string {
	return C.GoString(C.no_const(C.snd_pcm_format_name(C.snd_pcm_format_t(format))))
}

func (format Format) Description() string {
	return C.GoString(C.no_const(C.snd_pcm_format_description(C.snd_pcm_format_t(format))))
}

func (format Format) String() string {
	return format.Description()
}

func FromFormatName(name string) Format {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return Format(C.snd_pcm_format_value(cName))
}

type Class int
type SubClass int
type TimestampMode int

func (mode TimestampMode) Name() string {
	return C.GoString(C.no_const(C.snd_pcm_tstamp_mode_name(C.snd_pcm_tstamp_t(mode))))
}

func (mode TimestampMode) String() string {
	return mode.Name()
}

type TimestampType int

const (
	SndPcmStreamPlayback = StreamType(C.SND_PCM_STREAM_PLAYBACK)
	SndPcmStreamCapture  = StreamType(C.SND_PCM_STREAM_CAPTURE)

	SndPcmNonblock = int(C.SND_PCM_NONBLOCK)
	SndPcmAsync    = int(C.SND_PCM_ASYNC)

	SndPcmTypeHw          = Type(C.SND_PCM_TYPE_HW)
	SndPcmTypeHooks       = Type(C.SND_PCM_TYPE_HOOKS)
	SndPcmTypeMulti       = Type(C.SND_PCM_TYPE_MULTI)
	SndPcmTypeFile        = Type(C.SND_PCM_TYPE_FILE)
	SndPcmTypeNull        = Type(C.SND_PCM_TYPE_NULL)
	SndPcmTypeShm         = Type(C.SND_PCM_TYPE_SHM)
	SndPcmTypeInet        = Type(C.SND_PCM_TYPE_INET)
	SndPcmTypeCopy        = Type(C.SND_PCM_TYPE_COPY)
	SndPcmTypeLinear      = Type(C.SND_PCM_TYPE_LINEAR)
	SndPcmTypeAlaw        = Type(C.SND_PCM_TYPE_ALAW)
	SndPcmTypeMulaw       = Type(C.SND_PCM_TYPE_MULAW)
	SndPcmTypeAdpcm       = Type(C.SND_PCM_TYPE_ADPCM)
	SndPcmTypeRate        = Type(C.SND_PCM_TYPE_RATE)
	SndPcmTypeRoute       = Type(C.SND_PCM_TYPE_ROUTE)
	SndPcmTypePlug        = Type(C.SND_PCM_TYPE_PLUG)
	SndPcmTypeShare       = Type(C.SND_PCM_TYPE_SHARE)
	SndPcmTypeMeter       = Type(C.SND_PCM_TYPE_METER)
	SndPcmTypeMix         = Type(C.SND_PCM_TYPE_MIX)
	SndPcmTypeDroute      = Type(C.SND_PCM_TYPE_DROUTE)
	SndPcmTypeLbserver    = Type(C.SND_PCM_TYPE_LBSERVER)
	SndPcmTypeLinearFloat = Type(C.SND_PCM_TYPE_LINEAR_FLOAT)
	SndPcmTypeLadspa      = Type(C.SND_PCM_TYPE_LADSPA)
	SndPcmTypeDmix        = Type(C.SND_PCM_TYPE_DMIX)
	SndPcmTypeJack        = Type(C.SND_PCM_TYPE_JACK)
	SndPcmTypeDsnoop      = Type(C.SND_PCM_TYPE_DSNOOP)
	SndPcmTypeDshare      = Type(C.SND_PCM_TYPE_DSHARE)
	SndPcmTypeIec958      = Type(C.SND_PCM_TYPE_IEC958)
	SndPcmTypeSoftvol     = Type(C.SND_PCM_TYPE_SOFTVOL)
	SndPcmTypeIoplug      = Type(C.SND_PCM_TYPE_IOPLUG)
	SndPcmTypeExtplug     = Type(C.SND_PCM_TYPE_EXTPLUG)
	SndPcmTypeMmapEmul    = Type(C.SND_PCM_TYPE_MMAP_EMUL)

	SndPcmStateOpen         = State(C.SND_PCM_STATE_OPEN)
	SndPcmStateSetup        = State(C.SND_PCM_STATE_SETUP)
	SndPcmStatePrepared     = State(C.SND_PCM_STATE_PREPARED)
	SndPcmStateRunning      = State(C.SND_PCM_STATE_RUNNING)
	SndPcmStateXrun         = State(C.SND_PCM_STATE_XRUN)
	SndPcmStateDraining     = State(C.SND_PCM_STATE_DRAINING)
	SndPcmStatePaused       = State(C.SND_PCM_STATE_PAUSED)
	SndPcmStateSuspended    = State(C.SND_PCM_STATE_SUSPENDED)
	SndPcmStateDisconnected = State(C.SND_PCM_STATE_DISCONNECTED)

	SndPcmAccessMmapInterleaved    = Access(C.SND_PCM_ACCESS_MMAP_INTERLEAVED)
	SndPcmAccessMmapNoninterleaved = Access(C.SND_PCM_ACCESS_MMAP_NONINTERLEAVED)
	SndPcmAccessMmapComplex        = Access(C.SND_PCM_ACCESS_MMAP_COMPLEX)
	SndPcmAccessRwInterleaved      = Access(C.SND_PCM_ACCESS_RW_INTERLEAVED)
	SndPcmAccessRwNoninterleaved   = Access(C.SND_PCM_ACCESS_RW_NONINTERLEAVED)

	SndPcmFormatUnknown          = Format(C.SND_PCM_FORMAT_UNKNOWN)
	SndPcmFormatS8               = Format(C.SND_PCM_FORMAT_S8)
	SndPcmFormatU8               = Format(C.SND_PCM_FORMAT_U8)
	SndPcmFormatS16Le            = Format(C.SND_PCM_FORMAT_S16_LE)
	SndPcmFormatS16Be            = Format(C.SND_PCM_FORMAT_S16_BE)
	SndPcmFormatU16Le            = Format(C.SND_PCM_FORMAT_U16_LE)
	SndPcmFormatU16Be            = Format(C.SND_PCM_FORMAT_U16_BE)
	SndPcmFormatS24Le            = Format(C.SND_PCM_FORMAT_S24_LE)
	SndPcmFormatS24Be            = Format(C.SND_PCM_FORMAT_S24_BE)
	SndPcmFormatU24Le            = Format(C.SND_PCM_FORMAT_U24_LE)
	SndPcmFormatU24Be            = Format(C.SND_PCM_FORMAT_U24_BE)
	SndPcmFormatS32Le            = Format(C.SND_PCM_FORMAT_S32_LE)
	SndPcmFormatS32Be            = Format(C.SND_PCM_FORMAT_S32_BE)
	SndPcmFormatU32Le            = Format(C.SND_PCM_FORMAT_U32_LE)
	SndPcmFormatU32Be            = Format(C.SND_PCM_FORMAT_U32_BE)
	SndPcmFormatFloatLe          = Format(C.SND_PCM_FORMAT_FLOAT_LE)
	SndPcmFormatFloatBe          = Format(C.SND_PCM_FORMAT_FLOAT_BE)
	SndPcmFormatFloat64Le        = Format(C.SND_PCM_FORMAT_FLOAT64_LE)
	SndPcmFormatFloat64Be        = Format(C.SND_PCM_FORMAT_FLOAT64_BE)
	SndPcmFormatIec958SubframeLe = Format(C.SND_PCM_FORMAT_IEC958_SUBFRAME_LE)
	SndPcmFormatIec958SubframeBe = Format(C.SND_PCM_FORMAT_IEC958_SUBFRAME_BE)
	SndPcmFormatMuLaw            = Format(C.SND_PCM_FORMAT_MU_LAW)
	SndPcmFormatALaw             = Format(C.SND_PCM_FORMAT_A_LAW)
	SndPcmFormatImaAdpcm         = Format(C.SND_PCM_FORMAT_IMA_ADPCM)
	SndPcmFormatMpeg             = Format(C.SND_PCM_FORMAT_MPEG)
	SndPcmFormatGsm              = Format(C.SND_PCM_FORMAT_GSM)
	SndPcmFormatS20Le            = Format(C.SND_PCM_FORMAT_S20_LE)
	SndPcmFormatS20Be            = Format(C.SND_PCM_FORMAT_S20_BE)
	SndPcmFormatU20Le            = Format(C.SND_PCM_FORMAT_U20_LE)
	SndPcmFormatU20Be            = Format(C.SND_PCM_FORMAT_U20_BE)
	SndPcmFormatSpecial          = Format(C.SND_PCM_FORMAT_SPECIAL)
	SndPcmFormatS243le           = Format(C.SND_PCM_FORMAT_S24_3LE)
	SndPcmFormatS243be           = Format(C.SND_PCM_FORMAT_S24_3BE)
	SndPcmFormatU243le           = Format(C.SND_PCM_FORMAT_U24_3LE)
	SndPcmFormatU243be           = Format(C.SND_PCM_FORMAT_U24_3BE)
	SndPcmFormatS203le           = Format(C.SND_PCM_FORMAT_S20_3LE)
	SndPcmFormatS203be           = Format(C.SND_PCM_FORMAT_S20_3BE)
	SndPcmFormatU203le           = Format(C.SND_PCM_FORMAT_U20_3LE)
	SndPcmFormatU203be           = Format(C.SND_PCM_FORMAT_U20_3BE)
	SndPcmFormatS183le           = Format(C.SND_PCM_FORMAT_S18_3LE)
	SndPcmFormatS183be           = Format(C.SND_PCM_FORMAT_S18_3BE)
	SndPcmFormatU183le           = Format(C.SND_PCM_FORMAT_U18_3LE)
	SndPcmFormatU183be           = Format(C.SND_PCM_FORMAT_U18_3BE)
	SndPcmFormatG72324           = Format(C.SND_PCM_FORMAT_G723_24)
	SndPcmFormatG723241b         = Format(C.SND_PCM_FORMAT_G723_24_1B)
	SndPcmFormatG72340           = Format(C.SND_PCM_FORMAT_G723_40)
	SndPcmFormatG723401b         = Format(C.SND_PCM_FORMAT_G723_40_1B)
	SndPcmFormatDsdU8            = Format(C.SND_PCM_FORMAT_DSD_U8)
	SndPcmFormatDsdU16Le         = Format(C.SND_PCM_FORMAT_DSD_U16_LE)
	SndPcmFormatDsdU32Le         = Format(C.SND_PCM_FORMAT_DSD_U32_LE)
	SndPcmFormatDsdU16Be         = Format(C.SND_PCM_FORMAT_DSD_U16_BE)
	SndPcmFormatDsdU32Be         = Format(C.SND_PCM_FORMAT_DSD_U32_BE)

	SndPcmClassGeneric   = Class(C.SND_PCM_CLASS_GENERIC)
	SndPcmClassMulti     = Class(C.SND_PCM_CLASS_MULTI)
	SndPcmClassModem     = Class(C.SND_PCM_CLASS_MODEM)
	SndPcmClassDigitizer = Class(C.SND_PCM_CLASS_DIGITIZER)

	SndPcmSubclassGenericMix = SubClass(C.SND_PCM_SUBCLASS_GENERIC_MIX)
	SndPcmSubclassMultiMix   = SubClass(C.SND_PCM_SUBCLASS_MULTI_MIX)

	SndPcmTstampNone   = TimestampMode(C.SND_PCM_TSTAMP_NONE)
	SndPcmTstampEnable = TimestampMode(C.SND_PCM_TSTAMP_ENABLE)
	SndPcmTstampMmap   = TimestampMode(C.SND_PCM_TSTAMP_MMAP)

	SndPcmTstampTypeGettimeofday = TimestampType(C.SND_PCM_TSTAMP_TYPE_GETTIMEOFDAY)
	SndPcmTstampTypeMonotonic    = TimestampType(C.SND_PCM_TSTAMP_TYPE_MONOTONIC)
	SndPcmTstampTypeMonotonicRaw = TimestampType(C.SND_PCM_TSTAMP_TYPE_MONOTONIC_RAW)
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
