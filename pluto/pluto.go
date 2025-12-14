package pluto

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

func Mpd(address *url.URL) (*url.URL, []byte, error) {
   var req http.Request
   req.URL = address
   // The Request's URL and Header fields must be initialized
   req.Header = http.Header{}
   resp, err := http.DefaultClient.Do(&req)
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

// Define constants for the hardcoded URL parts
const (
   stitcherScheme = "https"
   stitcherHost   = "cfd-v4-service-stitcher-dash-use1-1.prd.pluto.tv"
)

// buildStitcherURL manually constructs the URL struct.
func (s *Series) buildStitcherURL(path string) *url.URL {
   u := &url.URL{
      Scheme: stitcherScheme,
      Host:   stitcherHost,
      Path:   path,
   }
   values := url.Values{}
   values.Set("jwt", s.SessionToken)
   u.RawQuery = values.Encode()
   return u
}

// GetMovieURL generates the Stitcher URL object for a movie.
// It assumes Vod and Stitched.Paths always have at least one entry.
func (s *Series) GetMovieURL() *url.URL {
   // Directly access the required path based on the data guarantees
   path := s.Vod[0].Stitched.Paths[0].Path
   return s.buildStitcherURL(path)
}
// GetEpisodeURL generates the Stitcher URL object for a specific episode by its ID.
func (s *Series) GetEpisodeURL(episodeID string) (*url.URL, error) {
   // Iterate through all seasons and episodes to find the matching ID
   for _, season := range s.Vod[0].Seasons {
      for _, episode := range season.Episodes {
         if episode.Id == episodeID {
            // Directly access the path based on the data guarantees
            path := episode.Stitched.Paths[0].Path
            return s.buildStitcherURL(path), nil
         }
      }
   }
   return nil, errors.New("episode not found")
}

type Stitched struct {
   Paths []struct {
      Path string
   }
}

type Vod struct {
   Id      string
   Seasons []struct {
      Episodes []struct {
         Id       string `json:"_id"`
         Name     string
         Number   int64
         Stitched Stitched
      }
      Number int64
   }
   Stitched *Stitched
}

type Series struct {
   Servers struct {
      StitcherDash string
   }
   SessionToken string
   Vod          []Vod
}

var (
   app_name         = "androidtv"
   drm_capabilities = "widevine:L1"
)

func Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      "https://service-concierge.clusters.pluto.tv/v1/wv/alt",
      "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (v *Vod) String() string {
   var (
      data  []byte
      lines bool
   )
   for _, season := range v.Seasons {
      for _, episode := range season.Episodes {
         if lines {
            data = append(data, "\n\n"...)
         } else {
            lines = true
         }
         data = append(data, "season = "...)
         data = strconv.AppendInt(data, season.Number, 10)
         data = append(data, "\nepisode = "...)
         data = strconv.AppendInt(data, episode.Number, 10)
         data = append(data, "\nname = "...)
         data = append(data, episode.Name...)
         data = append(data, "\nid = "...)
         data = append(data, episode.Id...)
      }
   }
   return string(data)
}

func (s *Series) Fetch(id string) error {
   req, _ := http.NewRequest("", "https://boot.pluto.tv/v4/start", nil)
   req.URL.RawQuery = url.Values{
      "appName":           {app_name},
      "appVersion":        {"9"},
      "clientID":          {"9"},
      "clientModelNumber": {"9"},
      "deviceMake":        {"9"},
      "deviceModel":       {"9"},
      "deviceVersion":     {"9"},
      "drmCapabilities":   {drm_capabilities},
      "seriesIDs":         {id},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   err = json.NewDecoder(resp.Body).Decode(s)
   if err != nil {
      return err
   }
   if s.Vod[0].Id != id {
      return errors.New("id mismatch")
   }
   return nil
}
