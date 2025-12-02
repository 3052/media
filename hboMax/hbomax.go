package hboMax

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "path"
   "slices"
   "strconv"
   "strings"
)

func (p *Playback) Mpd(session *Cache) error {
   resp, err := http.Get(
      strings.Replace(p.Fallback.Manifest.Url, "_fallback", "", 1),
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   session.Mpd = resp.Request.URL
   session.MpdBody, err = io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   return nil
}

type Cache struct {
   Login    *Login
   Mpd      *url.URL
   MpdBody  []byte
   Playback *Playback
   St       *St
}

type Playback struct {
   Drm struct {
      Schemes struct {
         PlayReady *Scheme
         Widevine  *Scheme
      }
   }
   Errors   []Error
   Fallback struct {
      Manifest struct {
         Url string // _fallback.mpd:1080p, .mpd:4K
      }
   }
   Manifest struct {
      Url string // 1080p
   }
}

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func (s St) Login() (*Login, error) {
   req, _ := http.NewRequest("POST", prd_api, nil)
   req.URL.Path = "/authentication/linkDevice/login"
   req.AddCookie(s[0])
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   value := &Login{}
   err = json.NewDecoder(resp.Body).Decode(value)
   if err != nil {
      return nil, err
   }
   return value, nil
}

func (s *St) Fetch() error {
   req, _ := http.NewRequest("", prd_api+"/token?realm=bolt", nil)
   req.Header.Set("x-device-info", device_info)
   req.Header.Set("x-disco-client", disco_client)
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

func (l Login) Movie(show_id string) (*Videos, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.Path = "/cms/routes/movie/" + show_id
   req.URL.RawQuery = url.Values{
      "include":          {"default"},
      "page[items.size]": {"1"},
   }.Encode()
   req.Header.Set("authorization", "Bearer "+l.Data.Attributes.Token)
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
   if len(movie.Errors) >= 1 {
      return nil, &movie.Errors[0]
   }
   return &movie, nil
}

func (e *Error) Error() string {
   if e.Detail != "" {
      return e.Detail
   }
   return e.Message
}

func (v *Video) String() string {
   var data []byte
   if v.Attributes.SeasonNumber >= 1 {
      data = append(data, "season number = "...)
      data = strconv.AppendInt(data, int64(v.Attributes.SeasonNumber), 10)
   }
   if v.Attributes.EpisodeNumber >= 1 {
      data = append(data, "\nepisode number = "...)
      data = strconv.AppendInt(data, int64(v.Attributes.EpisodeNumber), 10)
   }
   if data != nil {
      data = append(data, '\n')
   }
   data = append(data, "name = "...)
   data = append(data, v.Attributes.Name...)
   data = append(data, "\nvideo type = "...)
   data = append(data, v.Attributes.VideoType...)
   data = append(data, "\nedit id = "...)
   data = append(data, v.Relationships.Edit.Data.Id...)
   return string(data)
}

func (i *Initiate) String() string {
   var data strings.Builder
   data.WriteString("target URL = ")
   data.WriteString(i.TargetUrl)
   data.WriteString("\nlinking code = ")
   data.WriteString(i.LinkingCode)
   return data.String()
}

type Scheme struct {
   LicenseUrl string
}

const (
   device_info  = "!/!(!/!;!/!;!/!)"
   disco_client = "!:!:beam:!"
   prd_api      = "https://default.prd.api.discomax.com"
)

type Initiate struct {
   LinkingCode string
   TargetUrl   string
}

type Videos struct {
   Errors   []Error
   Included []*Video
}

type Error struct {
   Detail  string // show was filtered by validator
   Message string // Token is missing or not valid
}

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

func (v *Videos) FilterAndSort() {
   v.Included = slices.DeleteFunc(v.Included, func(video *Video) bool {
      if video.Attributes == nil {
         return true // Remove videos with nil attributes.
      }
      videoType := video.Attributes.VideoType
      return videoType != "EPISODE" && videoType != "MOVIE"
   })
   slices.SortFunc(v.Included, func(a, b *Video) int {
      if a.Attributes == nil || b.Attributes == nil {
         return 0 // Consider them equal if attributes are missing.
      }
      return a.Attributes.EpisodeNumber - b.Attributes.EpisodeNumber
   })
}

// https://hbomax.com/movies/weapons/bcbb6e0d-ca89-43e4-a9b1-2fc728145beb
// https://play.hbomax.com/show/bcbb6e0d-ca89-43e4-a9b1-2fc728145beb
func ExtractId(rawUrl string) (string, error) {
   u, err := url.Parse(rawUrl)
   if err != nil {
      return "", err
   }
   if u.Scheme == "" {
      return "", errors.New("invalid URL: scheme is missing")
   }
   return path.Base(u.Path), nil
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
      Errors []Error
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if len(value.Errors) >= 1 {
      return nil, &value.Errors[0]
   }
   return &value.Data.Attributes, nil
}

func (l Login) Season(show_id string, number int) (*Videos, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.Path = "/cms/collections/generic-show-page-rail-episodes-tabbed-content"
   req.Header.Set("authorization", "Bearer "+l.Data.Attributes.Token)
   req.URL.RawQuery = url.Values{
      "include":          {"default"},
      "pf[seasonNumber]": {strconv.Itoa(number)},
      "pf[show.id]":      {show_id},
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

func (l *Login) PlayReady(edit_id string) (*Playback, error) {
   return l.playback(edit_id, "playready")
}

func (l *Login) Widevine(edit_id string) (*Playback, error) {
   return l.playback(edit_id, "widevine")
}

func (l *Login) playback(edit_id, drm string) (*Playback, error) {
   data, err := json.Marshal(map[string]any{
      "editId":               edit_id,
      "consumptionType":      "streaming",
      "appBundle":            "",         // required
      "applicationSessionId": "",         // required
      "firstPlay":            false,      // required
      "gdpr":                 false,      // required
      "playbackSessionId":    "",         // required
      "userPreferences":      struct{}{}, // required
      "capabilities": map[string]any{
         "contentProtection": map[string]any{
            "contentDecryptionModules": []any{
               map[string]string{
                  "drmKeySystem": drm,
               },
            },
         },
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
      var data bytes.Buffer
      data.WriteString("/playback-orchestrator/any/playback-orchestrator/v1")
      data.WriteString("/playbackInfo")
      return data.String()
   }()
   req.Header.Set("authorization", "Bearer "+l.Data.Attributes.Token)
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode == 504 {
      // bail since no response body
      return nil, errors.New(resp.Status)
   }
   var play Playback
   err = json.NewDecoder(resp.Body).Decode(&play)
   if err != nil {
      return nil, err
   }
   if len(play.Errors) >= 1 {
      return nil, &play.Errors[0]
   }
   return &play, nil
}

func (p *Playback) PlayReady(data []byte) ([]byte, error) {
   resp, err := http.Post(
      p.Drm.Schemes.PlayReady.LicenseUrl, "text/xml",
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

type Login struct {
   Data struct {
      Attributes struct {
         Token string
      }
   }
}

type St [1]*http.Cookie
