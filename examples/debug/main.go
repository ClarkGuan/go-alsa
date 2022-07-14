package main

import (
	"fmt"
	"log"

	"github.com/ClarkGuan/go-alsa"
	"github.com/ClarkGuan/go-alsa/pcm"
)

func main() {
	pcmDev, err := pcm.Open("default", pcm.SndPcmStreamPlayback, 0)
	if err != nil {
		log.Fatalln(err)
	}
	stdout, err := alsa.AttachStdout()
	if err != nil {
		log.Fatalln(err)
	}
	if err := pcmDev.Dump(stdout); err != nil {
		log.Fatalln(err)
	}
	hardwareParams, err := pcm.AnyHardwareParamsFrom(pcmDev)
	if err != nil {
		log.Fatalln(err)
	}
	if err := hardwareParams.Dump(stdout); err != nil {
		log.Fatalln(err)
	}
	fmt.Println()
	if err := pcmDev.DumpSetup(stdout); err != nil {
		log.Fatalln(err)
	}
	//_, err = pcm.CurrentSoftwareParamsFrom(pcmDev)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//if err := softwareParams.Dump(stdout); err != nil {
	//	log.Fatalln(err)
	//}
}
