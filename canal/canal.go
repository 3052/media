package canal

import (
   "bytes"
   "crypto/hmac"
   "crypto/sha256"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
   "time"
)

func (s *Session) Episodes(tracking string, season int) ([]Episode, error) {
   req, _ := http.NewRequest("", "https://tvapi-hlm2.solocoo.tv/v1/assets", nil)
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.URL.RawQuery = fmt.Sprint(
      "limit=99&query=episodes,", tracking,
      ",season,", season,
   )
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Assets  []Episode
      Message string
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return result.Assets, nil
}

func (l *Login) Error() string {
   var data strings.Builder
   data.WriteString("label = ")
   data.WriteString(l.Label)
   data.WriteString("\nmessage = ")
   data.WriteString(l.Message)
   return data.String()
}

type Player struct {
   Drm struct {
      LicenseUrl string
   }
   Message   string
   Subtitles []struct {
      Url string
   }
   Url string // MPD
}

func (s *Session) Player(tracking string) (*Player, error) {
   data, err := json.Marshal(map[string]any{
      "player": map[string]any{
         "capabilities": map[string]any{
            "drmSystems": []string{"Widevine"},
            "mediaTypes": []string{"DASH"},
         },
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://tvapi-hlm2.solocoo.tv", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = func() string {
      var data strings.Builder
      data.WriteString("/v1/assets/")
      data.WriteString(tracking)
      data.WriteString("/play")
      return data.String()
   }()
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Player
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}

func Tracking(address string) (string, error) {
   resp, err := http.Get(address)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return "", errors.New(resp.Status)
   }
   var data strings.Builder
   _, err = io.Copy(&data, resp.Body)
   if err != nil {
      return "", err
   }
   const startKey = `data-algolia-convert-tracking="`
   _, after, found := strings.Cut(data.String(), startKey)
   if !found {
      return "", fmt.Errorf("attribute key '%s' not found", startKey)
   }
   before, _, found := strings.Cut(after, `"`)
   if !found {
      return "", fmt.Errorf("could not find closing quote for the attribute")
   }
   return before, nil
}

type Session struct {
   Message  string
   SsoToken string
   Token    string // this last one hour
}

func (p *Player) Mpd() (*Mpd, error) {
   resp, err := http.Get(p.Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Mpd{data, resp.Request.URL}, nil
}

type Mpd struct {
   Body []byte
   Url  *url.URL
}

type Login struct {
   Label    string
   Message  string
   SsoToken string // this last one day
}

const device_serial = "!!!!"

func get_client(requestURL *url.URL, requestBody []byte) (string, error) {
   encoding := base64.RawURLEncoding
   // 1. base64 raw URL decode secret key
   decodedSecretKey, err := encoding.DecodeString(SecretKeyBase64)
   if err != nil {
      return "", fmt.Errorf("failed to decode secret key: %v", err)
   }
   // Prepare timestamp and hash the body
   currentTimestamp := time.Now().Unix()
   bodyChecksum := sha256.Sum256(requestBody)
   encodedBodyHash := encoding.EncodeToString(bodyChecksum[:])
   // 2. hmac.New(sha256.New, secret key)
   hmacHasher := hmac.New(sha256.New, decodedSecretKey)
   // 3, 4, 5. Write components to the hasher
   // fmt.Fprint handles string conversion for requestURL and currentTimestamp
   fmt.Fprint(hmacHasher, requestURL, encodedBodyHash, currentTimestamp)
   // 6. base64 raw URL encode the hmac sum
   finalSignature := encoding.EncodeToString(hmacHasher.Sum(nil))
   // Construct final result string
   result := fmt.Sprintf(
      "Client key=%s,time=%d,sig=%s",
      ClientKey,
      currentTimestamp,
      finalSignature,
   )
   return result, nil
}

func (t *Ticket) Login(username, password string) (*Login, error) {
   data, err := json.Marshal(map[string]any{
      "ticket": t.Ticket,
      "userInput": map[string]string{
         "username": username,
         "password": password,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://m7cp.login.solocoo.tv/login", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   client, err := get_client(req.URL, data)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", client)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Login
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, &result
   }
   return &result, nil
}

func (t *Ticket) Fetch() error {
   data, err := json.Marshal(map[string]any{
      "deviceInfo": map[string]string{
         "brand":        "m7cp", // sg.ui.sso.fatal.internal_error
         "deviceModel":  "Firefox",
         "deviceOem":    "Firefox",
         "deviceSerial": device_serial,
         "deviceType":   "PC",
         "osVersion":    "Windows 10",
      },
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://m7cplogin.solocoo.tv/login", bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   client, err := get_client(req.URL, data)
   if err != nil {
      return err
   }
   req.Header.Set("authorization", client)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   err = json.NewDecoder(resp.Body).Decode(t)
   if err != nil {
      return err
   }
   if t.Message != "" {
      return errors.New(t.Message)
   }
   return nil
}
func (e *Episode) String() string {
   data := &strings.Builder{}
   data.WriteString("episode = ")
   fmt.Fprint(data, e.Params.SeriesEpisode)
   data.WriteString("\ntitle = ")
   data.WriteString(e.Title)
   data.WriteString("\ndesc = ")
   data.WriteString(e.Desc)
   data.WriteString("\ntracking = ")
   data.WriteString(e.Id)
   return data.String()
}

type Episode struct {
   Desc   string
   Id     string
   Params struct {
      SeriesEpisode int
   }
   Title string
}

// Global variables for authentication
var (
   ClientKey       = "web.NhFyz4KsZ54"
   SecretKeyBase64 = "OXh0-pIwu3gEXz1UiJtqLPscZQot3a0q"
)

type Ticket struct {
   Message string
   Ticket  string
}

func (p *Player) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(p.Drm.LicenseUrl, "", bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (s *Session) Fetch(ssoToken string) error {
   data, err := json.Marshal(map[string]string{
      "brand":        "m7cp",
      "deviceSerial": device_serial,
      "deviceType":   "PC",
      "ssoToken":     ssoToken,
   })
   if err != nil {
      return err
   }
   resp, err := http.Post(
      "https://tvapi-hlm2.solocoo.tv/v1/session", "", bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   err = json.NewDecoder(resp.Body).Decode(s)
   if err != nil {
      return err
   }
   if s.Message != "" {
      return errors.New(s.Message)
   }
   return nil
}
