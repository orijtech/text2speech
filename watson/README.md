# watson

Convert text to speech using IBM's Watson engine.

## Requirements
Create a project and get API credentials for IBM Watson at:
https://www.ibm.com/watson/developercloud/

Then ensure to set in your environment:
* WATSON_TEXT_TO_SPEECH_USERNAME
* WATSON_TEXT_TO_SPEECH_PASSWORD

## Sample Usage
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
		Voice: watson.VoiceAmericanLisaFemale,

		OutputContentType: watson.WAV,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	io.Copy(f, rc)
}
```
