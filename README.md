# text2speech
Text to Speech packages

## Watson
Uses IBM's Watson to transcribe text to speech allowing for output
in different languages. See file [Watson README](./watson/README.md)
or a snippet below
```go
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/orijtech/text2speech/watson"
)

func main() {
	client, err := watson.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("output-audio.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	now := time.Now()
	hour, min, sec := now.Clock()
	text := fmt.Sprintf("The time in 24hr clock is %v:%v:%v", hour, min, sec)
	rc, err := client.SynthesizeAudio(&watson.Request{
		Text:  text,
		Voice: watson.VoiceAmericanAllisonFemale,

		OutputContentType: watson.WAV,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	io.Copy(f, rc)
}
```
