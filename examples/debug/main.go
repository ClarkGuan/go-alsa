package main

import (
	"log"

	"github.com/ClarkGuan/go-alsa/pcm"
)

func main() {
	pcmDev, err := pcm.Open("default", pcm.SndPcmStreamPlayback, 0)
	if err != nil {
		log.Fatalln(err)
	}
	//stdout, err := alsa.AttachStdout()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//if err := pcmDev.Dump(stdout); err != nil {
	//	log.Fatalln(err)
	//}
	hardwareParams, err := pcm.AnyHardwareParamsFrom(pcmDev)
	if err != nil {
		log.Fatalln(err)
	}
	if err := hardwareParams.SetAccess(pcm.SndPcmAccessRwInterleaved); err != nil {
		log.Fatalln(err)
	}
	//if err := hardwareParams.Dump(stdout); err != nil {
	//	log.Fatalln(err)
	//}
}
