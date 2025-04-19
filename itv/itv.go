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

const query = `
query ProgrammePage( $brandLegacyId: BrandLegacyId ) {
   titles(
      filter: { brandLegacyId: $brandLegacyId }
      sortBy: SEQUENCE_ASC
   ) {
      ... on Episode {
         series {
            seriesNumber
         }
         episodeNumber
      }
      title
      latestAvailableVersion {
         playlistUrl
      }
   }
}
`

// this is better than strings.Replace and strings.ReplaceAll
func graphql_compact(data string) string {
   return strings.Join(strings.Fields(data), " ")
}

func programme_page(brand_legacy_id string) ([]title, error) {
   data, err := json.Marshal(map[string]string{
      "brandLegacyId": brand_legacy_id,
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
      "query":     {graphql_compact(query)},
      "variables": {string(data)},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         Titles []title
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return value.Data.Titles, nil
}

func (t *title) String() string {
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

type title struct {
   Series *struct {
      SeriesNumber int64
   }
   EpisodeNumber          int64
   Title                  string
   LatestAvailableVersion struct {
      PlaylistUrl string
   }
}
type legacy_id [1]string

func (v legacy_id) String() string {
   return v[0]
}

// itv.com/watch/joan/10a3918/10a3918a0001
// itv.com/watch/joan/10_3918/10a3918a0001
func (v *legacy_id) Set(data string) error {
   v[0] = strings.ReplaceAll(data, "a", "/")
   return nil
}

const programme_page = `
query ProgrammePage(
   $brandLegacyId: BrandLegacyId
) {
   titles(
      filter: { brandLegacyId: $brandLegacyId }
   ) {
      ... on Title {
         latestAvailableVersion { playlistUrl }
      }
   }
}
`

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

func (v legacy_id) programme_page() (*http.Response, error) {
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
      "query":     {programme_page},
      "variables": {string(data)},
   }.Encode()
   return http.DefaultClient.Do(req)
}

// hard geo block
func NewPlaylist() (Byte[Playlist], error) {
   value := map[string]any{
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
   }
   data, err := json.MarshalIndent(value, "", " ")
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://magni.itv.com", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/playlist/itvonline/ITV/SOMETHINNG"
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
