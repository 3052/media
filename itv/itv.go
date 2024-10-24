package itv

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

const query_discovery = `
{
   titles(filter: {
      legacyId: %q
   }) {
      ... on Episode {
         seriesNumber
         episodeNumber
      }
      ... on Film {
         productionYear
      }
      brand {
         title
      }
      latestAvailableVersion {
         playlistUrl
      }
      title
   }
}
`

// hard geo block
func (d *discovery_title) playlist() (*playlist, error) {
   var value struct {
      Client struct {
         Id string `json:"id"`
      } `json:"client"`
      VariantAvailability struct {
         Drm         struct {
            MaxSupported string `json:"maxSupported"`
            System       string `json:"system"`
         } `json:"drm"`
         FeatureSet  []string `json:"featureset"`
         PlatformTag string   `json:"platformTag"`
      } `json:"variantAvailability"`
   }
   value.Client.Id = "browser"
   value.VariantAvailability.Drm.MaxSupported = "L3"
   value.VariantAvailability.Drm.System = "widevine"
   // need all these to get 720:
   value.VariantAvailability.FeatureSet = []string{
      "hd",
      "mpeg-dash",
      "single-track",
      "widevine",
   }
   value.VariantAvailability.PlatformTag = "dotcom"
   data, err := json.MarshalIndent(value, "", " ")
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", d.LatestAvailableVersion.PlaylistUrl, bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("accept", "application/vnd.itv.vod.playlist.v4+json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   play := &playlist{}
   err = json.NewDecoder(resp.Body).Decode(play)
   if err != nil {
      return nil, err
   }
   return play, nil
}

type discovery_title struct {
   LatestAvailableVersion struct {
      PlaylistUrl string
   }
   Brand *struct {
      Title string
   }
   EpisodeNumber int
   ProductionYear int
   SeriesNumber int
   Title string
}

func (i legacy_id) String() string {
   return strings.Join(i[:], "/")
}

func (i legacy_id) discovery() (*discovery_title, error) {
   req, err := http.NewRequest(
      "", "https://content-inventory.prd.oasvc.itv.com/discovery", nil,
   )
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = url.Values{
      "query": {fmt.Sprintf(query_discovery, i)},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         Titles []discovery_title
      }
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if v := value.Errors; len(v) >= 1 {
      return nil, errors.New(v[0].Message)
   }
   return &value.Data.Titles[0], nil
}

func (i *legacy_id) Set(text string) error {
   var found bool
   (*i)[0], text, found = strings.Cut(text, "a")
   if !found {
      return errors.New(`"a" not found`)
   }
   (*i)[1], (*i)[2], found = strings.Cut(text, "a")
   if !found {
      (*i)[2] = "0001"
   }
   return nil
}

type legacy_id [3]string

type namer struct {
   title *discovery_title
}

func (n namer) Show() string {
   if n.title.Brand != nil {
      return n.title.Brand.Title
   }
   return ""
}

func (n namer) Season() int {
   return n.title.SeriesNumber
}

func (n namer) Episode() int {
   return n.title.EpisodeNumber
}

func (n namer) Title() string {
   return n.title.Title
}

func (n namer) Year() int {
   return n.title.ProductionYear
}

type playlist struct {
   Playlist struct {
      Video struct {
         MediaFiles []struct {
            Href string
            Resolution string
         }
      }
   }
}

func (p *playlist) resolution_720() (string, bool) {
   for _, file := range p.Playlist.Video.MediaFiles {
      if file.Resolution == "720" {
         return file.Href, true
      }
   }
   return "", false
}

func (poster) RequestUrl() (string, bool) {
   var u url.URL
   u.Host = "itvpnp.live.ott.irdeto.com"
   u.Path = "/Widevine/getlicense"
   u.RawQuery = "AccountId=itvpnp"
   u.Scheme = "https"
   return u.String(), true
}

type poster struct{}

func (poster) WrapRequest(b []byte) ([]byte, error) {
   return b, nil
}

func (poster) UnwrapResponse(b []byte) ([]byte, error) {
   return b, nil
}

func (poster) RequestHeader() (http.Header, error) {
   return http.Header{}, nil
}
