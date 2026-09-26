package audio

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

var (
	speakerInitOnce sync.Once
	hardwareRate    beep.SampleRate
)

// InitGlobalAudio explicitly initializes the audio subsystem at startup
func InitGlobalAudio() {
	speakerInitOnce.Do(func() {
		hardwareRate = 44100 // Safe default baseline rate
		_ = speaker.Init(hardwareRate, hardwareRate.N(time.Second/10))
	})
}

// PlaySound locates and streams a wav file asynchronously from your assets folder
func PlaySound(filename string) {
	go func() {
		var audioPath string
		if execDir, err := os.Executable(); err == nil {
			audioPath = filepath.Join(filepath.Dir(execDir), ".ping", "assets", "audio", filename)
		}

		if audioPath == "" {
			audioPath = filepath.Join(".ping", "assets", "audio", filename)
		} else if _, err := os.Stat(audioPath); os.IsNotExist(err) {
			audioPath = filepath.Join(".ping", "assets", "audio", filename)
		}

		f, err := os.Open(audioPath)
		if err != nil {
			return
		}
		defer f.Close()

		streamer, format, err := wav.Decode(f)
		if err != nil {
			return
		}
		defer streamer.Close()

		// Ensure speaker is initialized even if InitGlobalAudio wasn't called
		InitGlobalAudio()

		var playable beep.Streamer = streamer
		if format.SampleRate != hardwareRate {
			playable = beep.Resample(4, format.SampleRate, hardwareRate, streamer)
		}

		done := make(chan bool)
		speaker.Play(beep.Seq(playable, beep.Callback(func() {
			close(done)
		})))
		<-done
	}()
}
