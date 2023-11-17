package main

import (
	"fmt"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/msgs2"
	"os"
)

func main() {
	f, err := os.Open("./1-48050850-ad9e-4497-bf8b-d20a7d7cb232.dem")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cfg := demoinfocs.DefaultParserConfig
	playerAudioData := make(map[uint64][]byte)

	p := demoinfocs.NewParserWithConfig(f, cfg)
	defer p.Close()

	p.RegisterNetMessageHandler(func(m *msgs2.CSVCMsg_VoiceInit) {
		fmt.Println(m.Codec, m, m.String(), m.Version)
	})

	p.RegisterNetMessageHandler(func(m *msgs2.CSVCMsg_VoiceData) {
		audioBytes := m.Audio.GetVoiceData()
		if _, ok := playerAudioData[m.GetXuid()]; !ok {
			playerAudioData[m.GetXuid()] = audioBytes
		} else {
			playerAudioData[m.GetXuid()] = append(playerAudioData[m.GetXuid()], audioBytes...)
		}
	})

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		panic(err)
	}
}
