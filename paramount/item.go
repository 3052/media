package paramount

import (
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "strconv"
   "strings"
   "time"
)

const encoding = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func cms_account(id string) int64 {
   var (
      i = 0
      j = 1
   )
   for _, value := range id {
      i += strings.IndexRune(encoding, value) * j
      j *= len(encoding)
   }
   return int64(i)
}

// hard geo block
func (v *VideoItem) Mpd() string {
   b := []byte("https://link.theplatform.com/s/")
   b = append(b, v.CmsAccountId...)
   b = append(b, "/media/guid/"...)
   b = strconv.AppendInt(b, cms_account(v.CmsAccountId), 10)
   b = append(b, '/')
   b = append(b, v.ContentId...)
   b = append(b, "?assetTypes="...)
   b = append(b, v.asset_type()...)
   b = append(b, "&formats=MPEG-DASH"...)
   return string(b)
}

func (v *VideoItem) asset_type() string {
   if v.MediaType == "Movie" {
      return "DASH_CENC_PRECON"
   }
   return "DASH_CENC"
}

func (v *VideoItem) Unmarshal(data []byte) error {
   var value struct {
      Error string
      ItemList []VideoItem
   }
   err := json.Unmarshal(data, &value)
   if err != nil {
      return err
   }
   if value.Error != "" {
      return errors.New(value.Error)
   }
   if len(value.ItemList) == 0 {
      return errors.New(`"itemList":[]`)
   }
   *v = value.ItemList[0]
   return nil
}

// must use app token and IP address for correct location
func (*VideoItem) Marshal(token AppToken, cid string) ([]byte, error) {
   req, err := http.NewRequest("", "https://www.paramountplus.com", nil)
   if err != nil {
      return nil, err
   }
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/apps-api/v2.0/androidphone/video/cid/")
      b.WriteString(cid)
      b.WriteString(".json")
      return b.String()
   }()
   req.URL.RawQuery = token.Values.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var b strings.Builder
      resp.Write(&b)
      return nil, errors.New(b.String())
   }
   return io.ReadAll(resp.Body)
}

func (v *VideoItem) Title() string {
   return v.Label
}

func (v *VideoItem) Show() string {
   if v.MediaType == "Full Episode" {
      return v.SeriesTitle
   }
   return ""
}

func (v *VideoItem) Season() int {
   return v.SeasonNum.Data
}

func (v *VideoItem) Episode() int {
   return v.EpisodeNum.Data
}

func (v *VideoItem) Year() int {
   return v.AirDateIso.Year()
}

func (n *Number) UnmarshalText(data []byte) error {
   if len(data) >= 1 {
      var err error
      n.Data, err = strconv.Atoi(string(data))
      if err != nil {
         return err
      }
   }
   return nil
}

type VideoItem struct {
   AirDateIso time.Time `json:"_airDateISO"`
   CmsAccountId string
   ContentId string
   EpisodeNum Number
   Label string
   MediaType string
   SeasonNum Number
   SeriesTitle string
}

type Number struct {
   Data int
}

func (n Number) MarshalText() ([]byte, error) {
   return strconv.AppendInt(nil, int64(n.Data), 10), nil
}
