package watson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

const (
	endpointURL = "https://stream.watsonplatform.net/text-to-speech/api/v1/synthesize"
)

type Client struct {
	sync.RWMutex

	username string
	password string

	Transport http.RoundTripper
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	envUsernameKey = "WATSON_TEXT_TO_SPEECH_USERNAME"
	envPasswordKey = "WATSON_TEXT_TO_SPEECH_PASSWORD"
)

func envStr(envKey string, save *string) (errStr string) {
	if retr := os.Getenv(envKey); retr != "" {
		*save = retr
		return ""
	}

	return fmt.Sprintf("unset %q", envKey)
}

func NewClientFromEnv() (*Client, error) {
	var errsList []string
	var username, password string
	if errStr := envStr(envUsernameKey, &username); errStr != "" {
		errsList = append(errsList, errStr)
	}
	if errStr := envStr(envPasswordKey, &password); errStr != "" {
		errsList = append(errsList, errStr)
	}

	if len(errsList) > 0 {
		return nil, errors.New(strings.Join(errsList, "\n"))
	}

	client := new(Client)
	auth := &Auth{Username: username, Password: password}
	if err := client.SetAuth(auth); err != nil {
		return nil, err
	}

	return client, nil
}

var errNilCredentials = errors.New("nil credentials")

func (c *Client) SetAuth(auth *Auth) error {
	if auth == nil {
		return errNilCredentials
	}

	c.Lock()
	defer c.Unlock()

	c.username = auth.Username
	c.password = auth.Password
	return nil
}

type Voice string

const (
	VoiceGermanFemale          Voice = "de-DE_BirgitVoice"
	VoiceGermanMale            Voice = "de-DE_DieterVoice"
	VoiceBritishFemale         Voice = "en-GB_KateVoice"
	VoiceAmericanAllisonFemale Voice = "en-US_AllisonVoice"
	VoiceAmericanLisaFemale    Voice = "en-US_LisaVoice"
	VoiceAmericanMichaelMale   Voice = "en-US_MichaelVoice"

	VoiceSpanishCastilianMale       Voice = "es-ES_EnriqueVoice"
	VoiceSpanishCastilianFemale     Voice = "es-ES_LauraVoice"
	VoiceSpanishLatinAmericanFemale Voice = "es-LA_SofiaVoice"
	VoiceSpanishNorthAmericanFemale Voice = "es-US_SofiaVoice"
	VoiceFrenchFemale               Voice = "fr-FR_ReneeVoice"
	VoiceItalianFemale              Voice = "it-IT_FrancescaVoice"
	VoiceJapaneseFemale             Voice = "ja-JP_EmiVoice"
	VoiceBrazillianFemale           Voice = "pt-BR_IsabelaVoice"
)

type Request struct {
	Voice Voice `json:"voice"`

	Text string `json:"text"`

	OutputContentType ContentType `json:"requested_mimetype"`
}

type ContentType string

const (
	WAV        ContentType = "audio/wav"
	OGG        ContentType = "audio/ogg"
	OGGOpus    ContentType = "audio/ogg;codecs=opus"
	OGGVorbis  ContentType = "audio/ogg;codecs=vorbis"
	Mulaw      ContentType = "audio/mulaw;rate=rate"
	Basic      ContentType = "audio/basic"
	FLAC       ContentType = "audio/flac"
	WebmOpuS   ContentType = "audio/webm;codecs=opus"
	WebmVorbis ContentType = "audio/webm;codecs=vorbis"
	_116       ContentType = "audio/116;rate=rate"
)

func (c *Client) httpClient() *http.Client {
	c.RLock()
	defer c.RUnlock()
	if c.Transport == nil {
		return http.DefaultClient
	}
	return &http.Client{Transport: c.Transport}
}

func (c *Client) usernamePasswordCheck() (username, password string, err error) {
	c.RLock()
	username, password = c.username, c.password
	c.RUnlock()
	var errsList []string
	if username == "" {
		errsList = append(errsList, "expecting a username")
	}

	if password == "" {
		errsList = append(errsList, "expecting a password")
	}
	if len(errsList) > 0 {
		return "", "", errors.New(strings.Join(errsList, "\n"))
	}

	return username, password, nil
}

func (c *Client) SynthesizeAudio(req *Request) (io.ReadCloser, error) {
	username, password, err := c.usernamePasswordCheck()
	if err != nil {
		return nil, err
	}

	values := make(url.Values)
	if req.Voice != "" {
		values.Set("voice", string(req.Voice))
	}
	if ct := req.OutputContentType; ct != "" {
		values.Set("accept", string(ct))
	}

	theURL := endpointURL
	if len(values) > 0 {
		theURL += "?" + values.Encode()
	}

	bmap := map[string]interface{}{
		"text": req.Text,
	}

	blob, err := json.Marshal(bmap)
	if err != nil {
		return nil, err
	}

	hreq, _ := http.NewRequest("POST", theURL, bytes.NewReader(blob))
	hreq.SetBasicAuth(username, password)
	hreq.Header.Set("Content-Type", "application/json")

	httpClient := c.httpClient()
	res, err := httpClient.Do(hreq)
	if err != nil {
		return nil, err
	}

	warnings := res.Header.Get("Warnings")
	if len(warnings) > 0 {
		return nil, errors.New(warnings)
	}

	if !statusOK(res.StatusCode) {
		return nil, errors.New(res.Status)
	}

	return res.Body, nil
}

func statusOK(code int) bool { return code >= 200 && code <= 299 }
