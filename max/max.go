package max

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "iter"
   "net/http"
   "net/url"
   "path"
   "strconv"
   "strings"
)

func (n *Login) Playback(edit_id string) (Byte[Playback], error) {
   data, err := json.Marshal(map[string]any{
      "editId": edit_id,
      "consumptionType":      "streaming",
      "appBundle":            "",         // required
      "applicationSessionId": "",         // required
      "firstPlay":            false,      // required
      "gdpr":                 false,      // required
      "playbackSessionId":    "",         // required
      "userPreferences":      struct{}{}, // required
      "capabilities": map[string]any{
         "manifests": map[string]any{
            "formats": map[string]any{
               "dash": struct{}{}, // required
            }, // required
         }, // required
      }, // required
      "deviceInfo": map[string]any{
         "player": map[string]any{
            "mediaEngine": map[string]string{
               "name":    "", // required
               "version": "", // required
            }, // required
            "playerView": map[string]int{
               "height": 0, // required
               "width":  0, // required
            }, // required
            "sdk": map[string]string{
               "name":    "", // required
               "version": "", // required
            }, // required
         }, // required
      }, // required
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("POST", prd_api, bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   req.URL.Path = func() string {
      var b bytes.Buffer
      b.WriteString("/playback-orchestrator/any/playback-orchestrator/v1")
      b.WriteString("/playbackInfo")
      return b.String()
   }()
   // .Set to match .Get
   req.Header.Set("content-type", "application/json")
   req.Header.Set("authorization", "Bearer "+n.Data.Attributes.Token)
   req.Header.Set("proxy", "true")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode == 504 {
      // bail since no response body
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

func (p *Playback) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      p.Drm.Schemes.Widevine.LicenseUrl, "application/x-protobuf",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func (s St) Login() (Byte[Login], error) {
   req, _ := http.NewRequest("POST", prd_api, nil)
   req.URL.Path = "/authentication/linkDevice/login"
   req.AddCookie(s[0])
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (s St) Initiate() (*Initiate, error) {
   req, _ := http.NewRequest("POST", prd_api, nil)
   req.URL.Path = "/authentication/linkDevice/initiate"
   req.AddCookie(s[0])
   req.Header.Set("x-device-info", device_info)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         Attributes Initiate
      }
      Errors []struct {
         Detail string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if len(value.Errors) >= 1 {
      return nil, errors.New(value.Errors[0].Detail)
   }
   return &value.Data.Attributes, nil
}

func (s *St) New() error {
   req, _ := http.NewRequest("", prd_api+"/token?realm=bolt", nil)
   req.Header = http.Header{
      "x-device-info":  {device_info},
      "x-disco-client": {disco_client},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "st" {
         s[0] = cookie
         return nil
      }
   }
   return http.ErrNoCookie
}

type St [1]*http.Cookie

func (s *St) Set(data string) error {
   var err error
   s[0], err = http.ParseSetCookie(data)
   if err != nil {
      return err
   }
   return nil
}

func (s St) String() string {
   return s[0].String()
}

const (
   device_info  = "!/!(!/!;!/!;!/!)"
   disco_client = "!:!:beam:!"
   prd_api      = "https://default.prd.api.discomax.com"
)

type Byte[T any] []byte

func (i *Initiate) String() string {
   var b strings.Builder
   b.WriteString("target URL = ")
   b.WriteString(i.TargetUrl)
   b.WriteString("\nlinking code = ")
   b.WriteString(i.LinkingCode)
   return b.String()
}

type Initiate struct {
   LinkingCode string
   TargetUrl   string
}

type Login struct {
   Data struct {
      Attributes struct {
         Token string
      }
   }
}

func (n *Login) Unmarshal(data Byte[Login]) error {
   return json.Unmarshal(data, n)
}

func (p *Playback) Unmarshal(data Byte[Playback]) error {
   err := json.Unmarshal(data, p)
   if err != nil {
      return err
   }
   if len(p.Errors) >= 1 {
      return errors.New(p.Errors[0].Detail)
   }
   return nil
}

type Videos struct {
   Errors []struct {
      Detail string // show was filtered by validator
      Message string // Token is missing or not valid
   }
   Included []Video
}

func (n Login) Season(id ShowId, number int) (*Videos, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.Path = "/cms/collections/generic-show-page-rail-episodes-tabbed-content"
   req.Header.Set("authorization", "Bearer "+n.Data.Attributes.Token)
   req.URL.RawQuery = url.Values{
      "include":          {"default"},
      "pf[seasonNumber]": {strconv.Itoa(number)},
      "pf[show.id]":      {string(id)},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   season := &Videos{}
   err = json.NewDecoder(resp.Body).Decode(season)
   if err != nil {
      return nil, err
   }
   return season, nil
}

func (s ShowId) String() string {
   return string(s)
}

// max.com/movies/12199308-9afb-460b-9d79-9d54b5d2514c
// max.com/movies/heretic/12199308-9afb-460b-9d79-9d54b5d2514c
// max.com/shows/14f9834d-bc23-41a8-ab61-5c8abdbea505
// max.com/shows/white-lotus/14f9834d-bc23-41a8-ab61-5c8abdbea505
func (s *ShowId) Set(data string) error {
   switch {
   case strings.Contains(data, "/movies/"):
   case strings.Contains(data, "/shows/"):
   default:
      return errors.New("/movies/ or /shows/ not found")
   }
   *s = ShowId(path.Base(data))
   return nil
}

type ShowId string

type Video struct {
   Attributes *struct {
      SeasonNumber  int
      EpisodeNumber int
      Name          string
      VideoType     string
   }
   Relationships *struct {
      Edit *struct {
         Data struct {
            Id string
         }
      }
   }
}

func (v *Video) String() string {
   var b []byte
   if v.Attributes.SeasonNumber >= 1 {
      b = append(b, "season number = "...)
      b = strconv.AppendInt(b, int64(v.Attributes.SeasonNumber), 10)
   }
   if v.Attributes.EpisodeNumber >= 1 {
      b = append(b, "\nepisode number = "...)
      b = strconv.AppendInt(b, int64(v.Attributes.EpisodeNumber), 10)
   }
   if b != nil {
      b = append(b, '\n')
   }
   b = append(b, "name = "...)
   b = append(b, v.Attributes.Name...)
   b = append(b, "\nvideo type = "...)
   b = append(b, v.Attributes.VideoType...)
   b = append(b, "\nedit id = "...)
   b = append(b, v.Relationships.Edit.Data.Id...)
   return string(b)
}

func (v *Videos) Seq() iter.Seq[Video] {
   return func(yield func(Video) bool) {
      for _, video1 := range v.Included {
         if video1.Attributes != nil {
            switch video1.Attributes.VideoType {
            case "EPISODE", "MOVIE":
               if !yield(video1) {
                  return
               }
            }
         }
      }
   }
}

func (n Login) Movie(id ShowId) (*Videos, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.Path = "/cms/routes/movie/" + string(id)
   req.URL.RawQuery = url.Values{
      "include":          {"default"},
      "page[items.size]": {"1"},
   }.Encode()
   req.Header.Set("authorization", "Bearer "+n.Data.Attributes.Token)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var movie Videos
   err = json.NewDecoder(resp.Body).Decode(&movie)
   if err != nil {
      return nil, err
   }
   if movie.Error() != "" {
      return nil, &movie
   }
   return &movie, nil
}

func (v *Videos) Error() string {
   if len(v.Errors) >= 1 {
      err := v.Errors[0]
      if err.Detail != "" {
         return err.Detail
      }
      return err.Message
   }
   return ""
}

type Playback struct {
   Drm struct {
      Schemes struct {
         Widevine struct {
            LicenseUrl string
         }
      }
   }
   Errors []struct {
      Detail string
   }
   Fallback struct {
      Manifest struct {
         Url string
      }
   }
   // MPD shows higher bandwidth but its exact same, and extremely throttled
   Manifest struct {
      Url string
   }
}
