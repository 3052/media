package max

import (
   "bytes"
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strings"
   "time"
)

func (w *WatchUrl) MarshalText() ([]byte, error) {
   var b bytes.Buffer
   if w.VideoId != "" {
      b.WriteString("/video/watch/")
      b.WriteString(w.VideoId)
   }
   if w.EditId != "" {
      b.WriteByte('/')
      b.WriteString(w.EditId)
   }
   return b.Bytes(), nil
}

type DefaultRoutes struct {
   Data struct {
      Attributes struct {
         Url WatchUrl
      }
   }
   Included []RouteInclude
}

type RouteInclude struct {
   Attributes struct {
      AirDate       time.Time
      Name          string
      EpisodeNumber int
      SeasonNumber  int
   }
   Id            string
   Relationships *struct {
      Show *struct {
         Data struct {
            Id string
         }
      }
   }
}

type LinkLogin struct {
   Data struct {
      Attributes struct {
         Token string
      }
   }
}

func (v *LinkLogin) Routes(watch *WatchUrl) (*DefaultRoutes, error) {
   req, err := http.NewRequest("", prd_api, nil)
   if err != nil {
      return nil, err
   }
   req.URL.Path = func() string {
      data, _ := watch.MarshalText()
      var b strings.Builder
      b.WriteString("/cms/routes")
      b.Write(data)
      return b.String()
   }()
   req.URL.RawQuery = url.Values{
      "include": {"default"},
      // this is not required, but results in a smaller response
      "page[items.size]": {"1"},
   }.Encode()
   req.Header.Set("authorization", "Bearer " + v.Data.Attributes.Token)
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
   route := &DefaultRoutes{}
   err = json.NewDecoder(resp.Body).Decode(route)
   if err != nil {
      return nil, err
   }
   return route, nil
}
