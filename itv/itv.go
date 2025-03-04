package itv

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/http/cookiejar"
   "strings"
)

func (m *MediaFile) Mpd() (*http.Response, error) {
   var err error
   http.DefaultClient.Jar, err = cookiejar.New(nil)
   if err != nil {
      return nil, err
   }
   return http.Get(strings.Replace(m.Href, "itvpnpctv", "itvpnpdotcom", 1))
}

// hard geo block
func (e EpisodeId) Playlist() (Byte[Playlist], error) {
   data, err := json.Marshal(map[string]any{
      "client": map[string]string{
         "id": "browser",
      },
      "variantAvailability": map[string]any{
         "drm": map[string]string{
            "maxSupported": "L3",
            "system":       "widevine",
         },
         "featureset": []string{ // need all these to get 720p
            "hd",
            "mpeg-dash",
            "single-track",
            "widevine",
         },
         "platformTag": "ctv", // 1080p
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://magni.itv.com", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/playlist/itvonline/ITV/")
      b.WriteString(strings.Join(e[:], "_"))
      b.WriteString(".001")
      return b.String()
   }()
   req.Header.Set("accept", "application/vnd.itv.vod.playlist.v4+json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

func (m *MediaFile) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      m.KeyServiceUrl, "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (p *Playlist) FullHd() (*MediaFile, bool) {
   for _, file := range p.Playlist.Video.MediaFiles {
      if file.Resolution == "1080" {
         return &file, true
      }
   }
   return nil, false
}

type Playlist struct {
   Playlist struct {
      Video struct {
         MediaFiles []MediaFile
      }
   }
}

type EpisodeId [3]string

// https://www.itv.com/watch/gone-girl/10a5503a0001B
func (e *EpisodeId) Set(data string) error {
   data = strings.TrimSuffix(data, "B")
   var found bool
   (*e)[0], data, found = strings.Cut(data, "a")
   if !found {
      return errors.New(`"a" not found`)
   }
   (*e)[1], (*e)[2], found = strings.Cut(data, "a")
   if !found {
      (*e)[2] = "0001"
   }
   return nil
}

func (e EpisodeId) String() string {
   return strings.Join(e[:], "a")
}

type MediaFile struct {
   Href          string
   KeyServiceUrl string
   Resolution    string
}

type Byte[T any] []byte

func (p *Playlist) Unmarshal(data Byte[Playlist]) error {
   return json.Unmarshal(data, p)
}
