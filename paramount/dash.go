package paramount

import (
   "encoding/json"
   "io"
   "net/http"
   "net/url"
)

func (v *Video) Dash() (*Dash, error) {
   resp, err := http.Get(v.ItemList[0].StreamingUrl)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Dash{Body: body, Url: resp.Request.URL}, nil
}

type Video struct {
   ItemList []struct {
      StreamingUrl string
   }
}

func FetchVideo(at, cid string, cbsCom *http.Cookie) (*Video, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "www.paramountplus.com",
      Path:     join("/apps-api/v2.0/androidphone/video/cid/", cid, ".json"),
      RawQuery: url.Values{"at": {at}}.Encode(),
   }
   if cbsCom != nil {
      req.AddCookie(cbsCom)
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Video{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

type Dash struct {
   Body []byte
   Url  *url.URL
}
