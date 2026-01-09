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

func (v *Videos) FilterAndSort() {
   v.Included = slices.DeleteFunc(v.Included, func(vid *Video) bool {
      if vid.Attributes == nil {
         return true // Remove videos with nil attributes.
      }
      videoType := vid.Attributes.VideoType
      // Keep EPISODE, MOVIE, and STANDALONE_EVENT
      return videoType != "EPISODE" && videoType != "MOVIE" && videoType != "STANDALONE_EVENT"
   })
   slices.SortFunc(v.Included, func(a, b *Video) int {
      if a.Attributes == nil || b.Attributes == nil {
         return 0 // Consider them equal if attributes are missing.
      }
      return a.Attributes.EpisodeNumber - b.Attributes.EpisodeNumber
   })
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

func (login Login) Movie(showId string) (*Videos, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+login.Data.Attributes.Token)
   req.URL = &url.URL{
      Scheme: "https",
      Host:   api_host,
      Path:   "/cms/routes/movie/" + showId,
      RawQuery: url.Values{
         "include":          {"default"},
         "page[items.size]": {"1"},
      }.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Videos
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result, nil
}

type Videos struct {
   Errors   []Error
   Included []*Video
}

func (st St) Initiate(market string) (*Initiate, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("x-device-info", device_info)
   req.AddCookie(st.Cookie)
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   join("default.beam-", market, ".prd.api.discomax.com"),
      Path:   "/authentication/linkDevice/initiate",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data struct {
         Attributes Initiate
      }
      Errors []Error
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Data.Attributes, nil
}

func (login *Login) playback(editID, drm string) (*Playback, error) {
   data, err := json.Marshal(map[string]any{
      "editId":               editID,
      "consumptionType":      "streaming",
      "appBundle":            "",
      "applicationSessionId": "",
      "firstPlay":            false,
      "gdpr":                 false,
      "playbackSessionId":    "",
      "userPreferences":      struct{}{},
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
               "dash": struct{}{},
            },
         },
      },
      "deviceInfo": map[string]any{
         "player": map[string]any{
            "mediaEngine": map[string]string{
               "name":    "",
               "version": "",
            },
            "playerView": map[string]int{
               "height": 0,
               "width":  0,
            },
            "sdk": map[string]string{
               "name":    "",
               "version": "",
            },
         },
      },
   })
   if err != nil {
      return nil, err
   }
   var req http.Request
   req.Body = io.NopCloser(bytes.NewReader(data))
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+login.Data.Attributes.Token)
   req.Header.Set("content-type", "application/json")
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   api_host,
      Path:   "/playback-orchestrator/any/playback-orchestrator/v1/playbackInfo",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode == 504 {
      return nil, errors.New(resp.Status)
   }
   var result Playback
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result, nil
}

type Login struct {
   Data struct {
      Attributes struct {
         Token string
      }
   }
}

func (err *Error) Error() string {
   if err.Detail != "" {
      return err.Detail
   }
   return err.Message
}

type Error struct {
   Detail  string
   Message string
}

type Scheme struct {
   LicenseUrl string
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
         Url string
      }
   }
   Manifest struct {
      Url string
   }
}

type St struct {
   Cookie *http.Cookie
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

func (login *Login) Widevine(editID string) (*Playback, error) {
   return login.playback(editID, "widevine")
}

func (i *Initiate) String() string {
   var data strings.Builder
   data.WriteString("target URL = ")
   data.WriteString(i.TargetUrl)
   data.WriteString("\nlinking code = ")
   data.WriteString(i.LinkingCode)
   return data.String()
}

type Initiate struct {
   LinkingCode string
   TargetUrl   string
}

// Example URLs:
// https://hbomax.com/movies/weapons/bcbb6e0d-ca89-43e4-a9b1-2fc728145beb
// https://play.hbomax.com/show/bcbb6e0d-ca89-43e4-a9b1-2fc728145beb
func ExtractId(rawUrl string) (string, error) {
   parsedUrl, err := url.Parse(rawUrl)
   if err != nil {
      return "", err
   }
   if parsedUrl.Scheme == "" {
      return "", errors.New("invalid URL: scheme is missing")
   }
   return path.Base(parsedUrl.Path), nil
}

func join(data ...string) string {
   return strings.Join(data, "")
}

func (vid *Video) String() string {
   var data strings.Builder
   if vid.Attributes.SeasonNumber >= 1 {
      data.WriteString("season number = ")
      data.WriteString(strconv.Itoa(vid.Attributes.SeasonNumber))
   }
   if vid.Attributes.EpisodeNumber >= 1 {
      data.WriteString("\nepisode number = ")
      data.WriteString(strconv.Itoa(vid.Attributes.EpisodeNumber))
   }
   if data.Len() >= 1 {
      data.WriteByte('\n')
   }
   data.WriteString("name = ")
   data.WriteString(vid.Attributes.Name)
   data.WriteString("\nvideo type = ")
   data.WriteString(vid.Attributes.VideoType)
   data.WriteString("\nedit id = ")
   data.WriteString(vid.Relationships.Edit.Data.Id)
   return data.String()
}

type Mpd struct {
   Body []byte
   Url  *url.URL
}

func (p *Playback) Mpd() (*Mpd, error) {
   resp, err := http.Get(
      strings.Replace(p.Fallback.Manifest.Url, "_fallback", "", 1),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Mpd{data, resp.Request.URL}, nil
}

func (login *Login) PlayReady(editID string) (*Playback, error) {
   return login.playback(editID, "playready")
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

var Markets = []string{
   "amer",
   "apac",
   "emea",
   "latam",
}

const (
   api_host     = "default.prd.api.hbomax.com"
   disco_client = "!:!:beam:!"
   device_info  = "!/!(!/!;!/!;!/!)"
)

// You must call /authentication/linkDevice/initiate
// first or this will always fail.
func (st St) Login() (*Login, error) {
   var req http.Request
   req.Header = http.Header{}
   req.AddCookie(st.Cookie)
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   api_host,
      Path:   "/authentication/linkDevice/login",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Login{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func (login Login) Season(showId string, number int) (*Videos, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+login.Data.Attributes.Token)
   req.URL = &url.URL{
      Scheme: "https",
      Host:   api_host,
      Path:   "/cms/collections/generic-show-page-rail-episodes-tabbed-content",
      RawQuery: url.Values{
         "include":          {"default"},
         "pf[seasonNumber]": {strconv.Itoa(number)},
         "pf[show.id]":      {showId},
      }.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Videos{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func (st *St) Fetch() error {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("x-device-info", device_info)
   req.Header.Set("x-disco-client", disco_client)
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     api_host,
      Path:     "/token",
      RawQuery: "realm=bolt",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "st" {
         st.Cookie = cookie
         return nil
      }
   }
   return http.ErrNoCookie
}
