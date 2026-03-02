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

func (l *Login) video_alias(alias string) (*video_alias, error) {
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
   result := &video_alias{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

type video_alias struct {
   Video struct {
      VideoId int
   }
}

const x_version  = "!/!/!/!"

type Video struct {
   Alias   string
   VideoID int
}

// https://kanopy.com/video/genius-party
// https://kanopy.com/video/6440418
func (v *Video) Parse(inputURL string) error {
   parsedURL, err := url.Parse(inputURL)
   if err != nil {
      return err
   }
   if !strings.Contains(parsedURL.Host, "kanopy.com") {
      return errors.New("invalid domain")
   }
   if !strings.HasPrefix(parsedURL.Path, "/video/") {
      return errors.New("invalid path structure")
   }
   identifier := path.Base(parsedURL.Path)
   numericID, err := strconv.Atoi(identifier)
   if err != nil {
      v.Alias = identifier
   } else {
      v.VideoID = numericID
   }
   return nil
}
