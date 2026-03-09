package disney

import (
   "net/url"
   "strings"
)

func FixUrl(inputUrl string) (string, error) {
   urlParse, err := url.Parse(inputUrl)
   if err != nil {
      return "", err
   }
   var segments []string
   for _, segment := range strings.Split(urlParse.EscapedPath(), "/") {
      if !strings.HasPrefix(segment, "dvt1=") {
         segments = append(segments, segment)
      }
   }
   urlParse.Path = "/int" + strings.Join(segments, "/")
   urlParse.RawQuery = ""
   return urlParse.String(), nil
}
