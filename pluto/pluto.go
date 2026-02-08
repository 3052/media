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

func (d *Dash) Fetch(link *url.URL) error {
   var req http.Request
   req.Method = "GET"
   req.URL = link
   req.Header = http.Header{}
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   d.Body, err = io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   d.Url = resp.Request.URL
   return nil
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

func (s *Series) Fetch(id string) error {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme: "https",
      Host: "boot.pluto.tv",
      Path: "/v4/start",
      RawQuery: url.Values{
         "appName":           {app_name},
         "appVersion":        {"9"},
         "clientID":          {"9"},
         "clientModelNumber": {"9"},
         "deviceMake":        {"9"},
         "deviceModel":       {"9"},
         "deviceVersion":     {"9"},
         "drmCapabilities":   {drm_capabilities},
         "seriesIDs":         {id},
      }.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
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

// Define constants for the hardcoded URL parts
const (
   stitcherScheme = "https"
   stitcherHost   = "cfd-v4-service-stitcher-dash-use1-1.prd.pluto.tv"
)

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

type Series struct {
   SessionToken string
   Vod          []Vod
}

// buildStitcherURL manually constructs the URL struct.
func (s *Series) buildStitcherURL(path string) *url.URL {
   stitcher := &url.URL{
      Host:   stitcherHost,
      Path:   "/v2" + path,
      Scheme: stitcherScheme,
   }
   values := url.Values{}
   values.Set("jwt", s.SessionToken)
   stitcher.RawQuery = values.Encode()
   return stitcher
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
