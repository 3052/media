package itv

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "path"
   "strconv"
   "strings"
)

func (m *MediaFile) Mpd() (*url.URL, []byte, error) {
   var err error
   http.DefaultClient.Jar, err = cookiejar.New(nil)
   if err != nil {
      return nil, nil, err
   }
   resp, err := http.Get(strings.Replace(m.Href, "itvpnpctv", "itvpnpdotcom", 1))
   if err != nil {
      return nil, nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, nil, err
   }
   return resp.Request.URL, data, nil
}

func (p *Playlist) FullHd() (*MediaFile, bool) {
   for _, file := range p.Playlist.Video.MediaFiles {
      if file.Resolution == "1080" {
         return &file, true
      }
   }
   return nil, false
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
func LegacyId(rawUrl string) (string, error) {
   parsed, err := url.Parse(rawUrl)
   if err != nil {
      return "", err
   }
   if parsed.Scheme == "" {
      return "", errors.New("invalid URL: scheme is missing")
   }
   return strings.ReplaceAll(path.Base(parsed.Path), "a", "/"), nil
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
      "POST", "https://magni.itv.com/playlist/itvonline/ITV/"+id,
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Set("accept", "application/vnd.itv.vod.playlist.v4+json")
   req.Header.Set("user-agent", "!")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(p)
}

func Titles(legacyId string) ([]Title, error) {
   data, err := json.Marshal(map[string]string{
      "brandLegacyId": legacyId,
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
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var payload struct {
      Data struct {
         Titles []Title
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&payload)
   if err != nil {
      return nil, err
   }
   return payload.Data.Titles, nil
}

type Playlist struct {
   Error    string
   Playlist struct {
      Video struct {
         MediaFiles []MediaFile
      }
   }
}

type MediaFile struct {
   Href          string
   KeyServiceUrl string
   Resolution    string
}

type Title struct {
   LatestAvailableVersion struct {
      PlaylistUrl string
   }
   Series *struct {
      SeriesNumber int64
   }
   EpisodeNumber int64
   Title         string
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

func (t *Title) String() string {
   var data []byte
   if t.Series != nil {
      data = []byte("series = ")
      data = strconv.AppendInt(data, t.Series.SeriesNumber, 10)
      data = append(data, "\nepisode = "...)
      data = strconv.AppendInt(data, t.EpisodeNumber, 10)
   }
   if t.Title != "" {
      if data != nil {
         data = append(data, '\n')
      }
      data = append(data, "title = "...)
      data = append(data, t.Title...)
   }
   if data != nil {
      data = append(data, '\n')
   }
   data = append(data, "playlist = "...)
   data = append(data, t.LatestAvailableVersion.PlaylistUrl...)
   return string(data)
}
