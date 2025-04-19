package itv

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "strconv"
   "strings"
)

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

type Playlist struct {
   Playlist struct {
      Video struct {
         MediaFiles []MediaFile
      }
   }
}

type Byte[T any] []byte

func (p *Playlist) Unmarshal(data Byte[Playlist]) error {
   return json.Unmarshal(data, p)
}

func (p *Playlist) FullHd() (*MediaFile, bool) {
   for _, file := range p.Playlist.Video.MediaFiles {
      if file.Resolution == "1080" {
         return &file, true
      }
   }
   return nil, false
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

// hard geo block
func (t *Title) Playlist() (Byte[Playlist], error) {
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

func (v LegacyId) String() string {
   return strings.ReplaceAll(v[0], "/", "a")
}

type LegacyId [1]string

// 18910
// 10a3918
// 10a3918a0001
func (v *LegacyId) Set(data string) error {
   split := strings.SplitN(data, "a", 3)
   v[0] = split[0]
   if len(split) >= 2 {
      v[0] += "/" + split[1]
   }
   return nil
}

func (v LegacyId) Titles() ([]Title, error) {
   data, err := json.Marshal(map[string]string{
      "brandLegacyId": v[0],
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
