package pcm

//
// #include <alsa/asoundlib.h>
//
import "C"
import "github.com/ClarkGuan/go-alsa"

func (pcm *Dev) Dump(output alsa.Output) error {
	rc := C.snd_pcm_dump(pcm.inner, (*C.snd_output_t)(output.Ptr()))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *Dev) DumpHardwareSetup(output alsa.Output) error {
	rc := C.snd_pcm_dump_hw_setup(pcm.inner, (*C.snd_output_t)(output.Ptr()))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *Dev) DumpSoftwareSetup(output alsa.Output) error {
	rc := C.snd_pcm_dump_sw_setup(pcm.inner, (*C.snd_output_t)(output.Ptr()))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (pcm *Dev) DumpSetup(output alsa.Output) error {
	rc := C.snd_pcm_dump_setup(pcm.inner, (*C.snd_output_t)(output.Ptr()))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *HardwareParams) Dump(output alsa.Output) error {
	rc := C.snd_pcm_hw_params_dump(params.inner, (*C.snd_output_t)(output.Ptr()))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (params *SoftwareParams) Dump(output alsa.Output) error {
	rc := C.snd_pcm_sw_params_dump(params.inner, (*C.snd_output_t)(output.Ptr()))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}

func (s *Status) Dump(output alsa.Output) error {
	rc := C.snd_pcm_status_dump(s.inner, (*C.snd_output_t)(output.Ptr()))
	if rc < 0 {
		return alsa.NewError(int(rc))
	}
	return nil
}
