package pluto

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

// pluto.tv/on-demand/movies/64946365c5ae350013623630
// pluto.tv/on-demand/movies/disobedience-ca-2018-1-1
func (s *Series) Fetch(movieShow string) error {
   data := url.Values{}
   data.Set("appName", app_name)
   data.Set("appVersion", "9")
   data.Set("clientID", "9")
   data.Set("clientModelNumber", "9")
   data.Set("deviceMake", "9")
   data.Set("deviceModel", "9")
   data.Set("deviceVersion", "9")
   data.Set("drmCapabilities", drm_capabilities)
   if strings.Contains(movieShow, "-") {
      data.Set("episodeSlugs", movieShow)
   } else {
      data.Set("seriesIDs", movieShow)
   }
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "boot.pluto.tv",
      Path:     "/v4/start",
      RawQuery: data.Encode(),
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
   if strings.Contains(movieShow, "-") {
      if s.Vod[0].Slug != movieShow {
         return errors.New("slug mismatch")
      }
   } else if s.Vod[0].Id != movieShow {
      return errors.New("id mismatch")
   }
   return nil
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
   Slug     string
   Stitched *Stitched
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

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

type Stitched struct {
   Paths []struct {
      Path string
   }
}

type Series struct {
   SessionToken string
   Vod          []Vod
}
