package main

import (
	_ "embed"
	"log"
	"time"
	"unsafe"

	"github.com/ClarkGuan/go-alsa"
	"github.com/ClarkGuan/go-alsa/pcm"
	"github.com/pkg/errors"
)

//go:embed pcms/zh_long_driving.pcm
var ufcwData []byte

var stdout, _ = alsa.AttachStdout()

func main() {
	dev, err := pcm.Open("default", pcm.SndPcmStreamPlayback, 0)
	if err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}

	hardwareParams, err := dev.AnyHardwareParams()
	if err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}
	if err := hardwareParams.SetAccess(pcm.SndPcmAccessRwInterleaved); err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}
	if err := hardwareParams.SetFormat(pcm.SndPcmFormatS16Le); err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}
	if trueRate, _, err := hardwareParams.SetRateNear(16000, 0); err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	} else {
		log.Printf("true rate: %d\n", trueRate)
	}
	if _, err := hardwareParams.SetChannelsNear(1); err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}
	if truePeriodTime, _, err := hardwareParams.SetPeriodTimeNear(100*time.Millisecond, 0); err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	} else {
		log.Printf("true period time: %v\n", truePeriodTime)
	}
	if truePeriodsPerBuffer, _, err := hardwareParams.SetPeriodsPerBufferNear(2, 0); err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	} else {
		log.Printf("true periodsPerBuffer: %v\n", truePeriodsPerBuffer)
	}
	if err := hardwareParams.SetRateResample(true); err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}

	log.Println(dev.State())

	if err := hardwareParams.Install(); err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}

	if bufferTime, _, err := hardwareParams.GetBufferTime(); err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	} else {
		log.Printf("true buffer time: %v\n", bufferTime)
	}
	bufferSize, err := hardwareParams.GetBufferSize()
	if err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	} else {
		log.Printf("buffer size: %v\n", bufferSize)
	}

	log.Println("pcm dump:")
	if err := dev.Dump(stdout); err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}
	log.Println(dev.State())

	if err := dev.Start(); err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}

	log.Println(dev.State())

	// play
	tmp := ufcwData
	transfer := bufferSize

	for {
		if len(tmp) < dev.FramesToBytes(bufferSize) {
			transfer = dev.BytesToFrames(len(tmp))
			if transfer == 0 {
				break
			}
		} else {
			transfer = bufferSize
		}
		n, err := dev.Writei(unsafe.Pointer(&tmp[0]), transfer)
		if err != nil {
			log.Fatalf("%+v\n", errors.WithStack(err))
		}
		if n > 0 {
			tmp = tmp[dev.FramesToBytes(n):]
		}
	}

	if err := dev.Drain(); err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}

	log.Println(dev.State())

	log.Println("over")
}
