package watson_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/orijtech/text2speech/watson"
)

func Example_client_SynthesizeAudio_AmericanEnglish() {
	client, err := watson.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("english-time-audio.ogg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	now := time.Now()
	hour, min, sec := now.Clock()
	text := fmt.Sprintf("The time in 24hr clock is %v:%v:%v", hour, min, sec)
	rc, err := client.SynthesizeAudio(&watson.Request{
		Text:  text,
		Voice: watson.VoiceAmericanMichaelMale,

		OutputContentType: watson.OGG,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	io.Copy(f, rc)
}

func Example_client_SynthesizeAudio_BrazillianPortuguese() {
	client, err := watson.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("brazillian-output-audio.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	text := fmt.Sprintf("This is not a drill. Good evening!")
	rc, err := client.SynthesizeAudio(&watson.Request{
		Text:  text,
		Voice: watson.VoiceBrazillianFemale,

		OutputContentType: watson.WAV,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	io.Copy(f, rc)
}
