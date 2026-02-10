package peacock

import (
   "bytes"
   "crypto/hmac"
   "crypto/md5"
   "crypto/sha1"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "slices"
   "strings"
   "time"
)

func (v VideoPlayout) Akamai() (string, bool) {
   for _, endpoint := range v.Asset.Endpoints {
      if endpoint.CDN == "AKAMAI" {
         return endpoint.URL, true
      }
   }
   return "", false
}

type VideoPlayout struct {
   Asset struct {
      Endpoints []struct {
         CDN string
         URL string
      }
   }
   Protection struct {
      LicenceAcquisitionUrl string // wikipedia.org/wiki/License
   }
}

func (VideoPlayout) RequestHeader() (http.Header, error) {
   return http.Header{}, nil
}

func (VideoPlayout) RequestBody(b []byte) ([]byte, error) {
   return b, nil
}

func (VideoPlayout) ResponseBody(b []byte) ([]byte, error) {
   return b, nil
}

func (v VideoPlayout) RequestUrl() (string, bool) {
   return v.Protection.LicenceAcquisitionUrl, true
}

func (a AuthToken) Video(content_id string) (*VideoPlayout, error) {
   body, err := func() ([]byte, error) {
      type capability struct {
         Acodec string `json:"acodec"`
         Container string `json:"container"`
         Protection string `json:"protection"`
         Transport string `json:"transport"`
         Vcodec string `json:"vcodec"`
      }
      var s struct {
         ContentId string `json:"contentId"`
         Device struct {
            Capabilities []capability `json:"capabilities"`
         } `json:"device"`
      }
      s.ContentId = content_id
      s.Device.Capabilities = []capability{
         {
            Acodec: "AAC",
            Container: "ISOBMFF",
            Protection: "WIDEVINE",
            Transport: "DASH",
            Vcodec: "H264",
         },
      }
      return json.Marshal(s)
   }()
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://ovp.peacocktv.com/video/playouts/vod",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   // `application/json` fails
   req.Header.Set("content-type", "application/vnd.playvod.v1+json")
   req.Header.Set("x-skyott-usertoken", a.UserToken)
   req.Header.Set(
      "x-sky-signature", sign(req.Method, req.URL.Path, req.Header, body),
   )
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      var b bytes.Buffer
      res.Write(&b)
      return nil, errors.New(b.String())
   }
   video := new(VideoPlayout)
   if err := json.NewDecoder(res.Body).Decode(video); err != nil {
      return nil, err
   }
   return video, nil
}

// userToken is good for one day
type AuthToken struct {
   UserToken string
}

func (q *QueryNode) New(content_id string) error {
   req, err := http.NewRequest("GET", "https://atom.peacocktv.com", nil)
   if err != nil {
      return err
   }
   req.URL.Path = "/adapter-calypso/v3/query/node/content_id/" + content_id
   req.Header = http.Header{
      "X-Skyott-Proposition": {"NBCUOTT"},
      "X-Skyott-Territory": {Territory},
   }
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   return json.NewDecoder(res.Body).Decode(q)
}

func sign(method, path string, head http.Header, body []byte) string {
   timestamp := time.Now().Unix()
   text_headers := func() string {
      var s []string
      for k := range head {
         k = strings.ToLower(k)
         if strings.HasPrefix(k, "x-skyott-") {
            s = append(s, k + ": " + head.Get(k) + "\n")
         }
      }
      slices.Sort(s)
      return strings.Join(s, "")
   }()
   headers_md5 := md5.Sum([]byte(text_headers))
   payload_md5 := md5.Sum(body)
   signature := func() string {
      h := hmac.New(sha1.New, []byte(sky_key))
      fmt.Fprintln(h, method)
      fmt.Fprintln(h, path)
      fmt.Fprintln(h)
      fmt.Fprintln(h, sky_client)
      fmt.Fprintln(h, sky_version)
      fmt.Fprintf(h, "%x\n", headers_md5)
      fmt.Fprintln(h, timestamp)
      fmt.Fprintf(h, "%x\n", payload_md5)
      hashed := h.Sum(nil)
      return base64.StdEncoding.EncodeToString(hashed[:])
   }
   sky_ott := func() string {
      b := []byte("SkyOTT")
      // must be quoted
      b = fmt.Appendf(b, " client=%q", sky_client)
      // must be quoted
      b = fmt.Appendf(b, ",signature=%q", signature())
      // must be quoted
      b = fmt.Appendf(b, `,timestamp="%v"`, timestamp)
      // must be quoted
      b = fmt.Appendf(b, ",version=%q", sky_version)
      return string(b)
   }
   return sky_ott()
}

const (
   sky_client = "NBCU-ANDROID-v3"
   sky_key = "JuLQgyFz9n89D9pxcN6ZWZXKWfgj2PNBUb32zybj"
   sky_version = "1.0"
)

var Territory = "US"

func (q QueryNode) Show() string {
   return q.Attributes.SeriesName
}

func (q QueryNode) Season() int {
   return q.Attributes.SeasonNumber
}

func (q QueryNode) Episode() int {
   return q.Attributes.EpisodeNumber
}

func (q QueryNode) Title() string {
   return q.Attributes.Title
}

type QueryNode struct {
   Attributes struct {
      EpisodeNumber int
      SeasonNumber int
      SeriesName string
      Title string
      Year int
   }
}

func (q QueryNode) Year() int {
   return q.Attributes.Year
}
