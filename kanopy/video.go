package kanopy

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "path"
   "strconv"
   "strings"
)

type Video struct {
   Alias string
   VideoId int
}

// https://kanopy.com/video/genius-party
// https://kanopy.com/video/6440418
func (v *Video) Parse(inputUrl string) error {
   parsedUrl, err := url.Parse(inputUrl)
   if err != nil {
      return err
   }
   if !strings.Contains(parsedUrl.Host, "kanopy.com") {
      return errors.New("invalid domain")
   }
   if !strings.HasPrefix(parsedUrl.Path, "/video/") {
      return errors.New("invalid path structure")
   }
   identifier := path.Base(parsedUrl.Path)
   numericId, err := strconv.Atoi(identifier)
   if err != nil {
      v.Alias = identifier
   } else {
      v.VideoId = numericId
   }
   return nil
}

const x_version  = "!/!/!/!"

func (l *Login) Video(alias string) (*Video, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme: "https",
      Host: "www.kanopy.com",
      Path: "/kapi/videos/alias/" + alias,
   }
   req.Header.Set("x-version", x_version)
   req.Header.Set("authorization", "Bearer " + l.Jwt)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Video Video
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Video, nil
}

// good for 10 years
type Login struct {
   Jwt    string
   UserId int
}
