package itv

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "log"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "strconv"
   "strings"
)

func (t *Title) Playlist() (*Playlist, error) {
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
      "POST", t.LatestAvailableVersion.PlaylistUrl, bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("accept", "application/vnd.itv.vod.playlist.v4+json")
   req.Header.Set("user-agent", "!")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var play Playlist
   err = json.NewDecoder(resp.Body).Decode(&play)
   if err != nil {
      return nil, err
   }
   if play.Error != "" {
      return nil, errors.New(play.Error)
   }
   return &play, nil
}

type Playlist struct {
   Error string
   Playlist struct {
      Video struct {
         MediaFiles []MediaFile
      }
   }
}

// pass: https://www.itv.com/watch/joan/10a3918
// fail: https://www.itv.com/watch/joan/10a3918/10a3918a0001
func (l *LegacyId) Set(data string) error {
   data = strings.TrimPrefix(data, "https://")
   data = strings.TrimPrefix(data, "www.")
   data = strings.TrimPrefix(data, "itv.com")
   split := strings.SplitN(data, "/", 5)
   if len(split) != 4 {
      return errors.New("/watch/[programmeSlug]/[programmeId]")
   }
   l[0] = strings.ReplaceAll(split[3], "a", "/")
   return nil
}

type LegacyId [1]string
var Transport = http.Transport{
   Proxy: func(req *http.Request) (*url.URL, error) {
      log.Println(req.Method, req.URL)
      return http.ProxyFromEnvironment(req)
   },
}

func (l LegacyId) Titles() ([]Title, error) {
   data, err := json.Marshal(map[string]string{
      "brandLegacyId": l[0],
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "", "https://content-inventory.prd.oasvc.itv.com/discovery", nil,
   )
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = url.Values{
      "query":     {graphql_compact(programme_page)},
      "variables": {string(data)},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         Titles []Title
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return value.Data.Titles, nil
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

type MediaFile struct {
   Href          string
   KeyServiceUrl string
   Resolution    string
}

func (m *MediaFile) Mpd() (*http.Response, error) {
   var err error
   http.DefaultClient.Jar, err = cookiejar.New(nil)
   if err != nil {
      return nil, err
   }
   return http.Get(strings.Replace(m.Href, "itvpnpctv", "itvpnpdotcom", 1))
}

type Title struct {
   LatestAvailableVersion struct {
      PlaylistUrl string
   }
   Series *struct {
      SeriesNumber int64
   }
   EpisodeNumber          int64
   Title                  string
}

const programme_page = `
query ProgrammePage( $brandLegacyId: BrandLegacyId ) {
   titles(
      filter: { brandLegacyId: $brandLegacyId }
      sortBy: SEQUENCE_ASC
   ) {
      ... on Episode {
         series { seriesNumber }
         episodeNumber
      }
      title
      latestAvailableVersion { playlistUrl }
   }
}
`

// this is better than strings.Replace and strings.ReplaceAll
func graphql_compact(data string) string {
   return strings.Join(strings.Fields(data), " ")
}

func (t *Title) String() string {
   var b []byte
   if t.Series != nil {
      b = []byte("series = ")
      b = strconv.AppendInt(b, t.Series.SeriesNumber, 10)
      b = append(b, "\nepisode = "...)
      b = strconv.AppendInt(b, t.EpisodeNumber, 10)
   }
   if t.Title != "" {
      if b != nil {
         b = append(b, '\n')
      }
      b = append(b, "title = "...)
      b = append(b, t.Title...)
   }
   if b != nil {
      b = append(b, '\n')
   }
   b = append(b, "playlist = "...)
   b = append(b, t.LatestAvailableVersion.PlaylistUrl...)
   return string(b)
}

func (p *Playlist) FullHd() (*MediaFile, bool) {
   for _, file := range p.Playlist.Video.MediaFiles {
      if file.Resolution == "1080" {
         return &file, true
      }
   }
   return nil, false
}

func (p *Playlist) playReady(id string) error {
   data, err := json.Marshal(map[string]any{
      "client": map[string]string{
         "id": "browser",
      },
      "variantAvailability": map[string]any{
         "drm": map[string]string{
            "maxSupported": "SL3000",
            "system":       "playready",
         },
         "featureset": []string{
            "hd",
            "mpeg-dash",
            "single-track",
            "playready",
         },
         "platformTag": "ctv", // 1080p
      },
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://magni.itv.com/playlist/itvonline/ITV/" + id,
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Set("accept", "application/vnd.itv.vod.playlist.v4+json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(p)
}

