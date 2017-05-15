package watson_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/orijtech/text2speech/watson"
)

type customTransport struct{}

type input struct {
	Text string `json:"text"`
}

func (ct *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	defer req.Body.Close()

	blob, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	in := new(input)
	if err := json.Unmarshal(blob, in); err != nil {
		return nil, err
	}

	f, err := os.Open("./testdata/12_15PM.ogg")
	if err != nil {
		return nil, err
	}

	res := &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       f,
		Header:     make(http.Header),
	}

	return res, nil
}

func TestTranscription(t *testing.T) {
	testAuth := &watson.Auth{
		Username: "username",
		Password: "password",
	}

	tests := [...]struct {
		text    string
		wantErr bool
		auth    *watson.Auth

		wantContentType string

		requestedContentType watson.ContentType
	}{
		0: {
			text: "What time is it right now?",
			auth: testAuth,

			wantContentType: "application/ogg",
		},
	}

	c := &watson.Client{
		Transport: new(customTransport),
	}

	for i, tt := range tests {
		c.SetAuth(tt.auth)
		req := &watson.Request{
			Text: tt.text,

			OutputContentType: tt.requestedContentType,
		}

		rc, err := c.SynthesizeAudio(req)
		if tt.wantErr {
			if err == nil {
				t.Errorf("#%d: wantErr", i)
			}
			continue
		}

		if err != nil {
			t.Errorf("#%d: err: %v", i, err)
			continue
		}

		sniffBuf := make([]byte, 512)
		_, err = io.ReadAtLeast(rc, sniffBuf, 10)
		if err != nil {
			t.Errorf("#%d: err: %v", i, err)
		}

		gotCt := http.DetectContentType(sniffBuf)
		if gotCt != tt.wantContentType {
			t.Errorf("#%d: gotContentType=%q wantContentType: %q", i, gotCt, tt.wantContentType)
		}
		_ = rc.Close()
	}
}
