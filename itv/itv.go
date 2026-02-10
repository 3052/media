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

func (m *MediaFile) License(data []byte) ([]byte, error) {
   resp, err := http.Post(
      m.KeyServiceUrl, "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (p *Playlist) Resolution1080() (*MediaFile, bool) {
   for _, file := range p.Playlist.Video.MediaFiles {
      if file.Resolution == "1080" {
         return &file, true
      }
   }
   return nil, false
}

type MediaFile struct {
   Href          Href
   KeyServiceUrl string
   Resolution    string
}

type Playlist struct {
   Playlist struct {
      Video struct {
         MediaFiles []MediaFile
      }
   }
}

func (h Href) Mpd() (*http.Response, error) {
   var err error
   http.DefaultClient.Jar, err = cookiejar.New(nil)
   if err != nil {
      return nil, err
   }
   return http.Get(h[0])
}

type Href [1]string

func (h *Href) UnmarshalText(data []byte) error {
   (*h)[0] = strings.Replace(string(data), "itvpnpctv", "itvpnpdotcom", 1)
   return nil
}

// hard geo block
func (i LegacyId) Playlist() (*Playlist, error) {
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
   req.URL.Path = "/playlist/itvonline/ITV/" + i.String()
   req.Header.Set("accept", "application/vnd.itv.vod.playlist.v4+json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   play := &Playlist{}
   err = json.NewDecoder(resp.Body).Decode(play)
   if err != nil {
      return nil, err
   }
   return play, nil
}

type LegacyId [3]string

func (i LegacyId) String() string {
   return strings.Join(i[:], "_") + ".001"
}

// https://www.itv.com/watch/gone-girl/10a5503a0001B
func (i *LegacyId) Set(data string) error {
   data = strings.TrimSuffix(data, "B")
   var found bool
   (*i)[0], data, found = strings.Cut(data, "a")
   if !found {
      return errors.New(`"a" not found`)
   }
   (*i)[1], (*i)[2], found = strings.Cut(data, "a")
   if !found {
      (*i)[2] = "0001"
   }
   return nil
}
