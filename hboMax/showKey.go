package hboMax

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strings"
)

type ShowKey struct {
   Type string
   Id   string
}

/*
https://hbomax.com/at/en/movies/austin-powers-international-man-of-mystery/a979fb8b-f713-4de3-a625-d16ad4d37448
https://hbomax.com/movies/one-battle-after-another/bebe611d-8178-481a-a4f2-de743b5b135a
https://hbomax.com/shows/white-lotus/14f9834d-bc23-41a8-ab61-5c8abdbea505
https://play.hbomax.com/show/31cb4b84-951a-4daf-8925-746fcdcddcb8
*/
func (s *ShowKey) Parse(rawLink string) error {
   link, err := url.Parse(rawLink)
   if err != nil {
      return err
   }
   segments := strings.Split(strings.TrimPrefix(link.Path, "/"), "/")
   count := len(segments)
   if count < 2 {
      return errors.New("invalid URL format: not enough path segments")
   }
   s.Id = segments[count-1]
   // 1. Check for standard catalog types (plural) at position -3
   // URL: .../movies/slug/id
   if count >= 3 {
      segmentType := segments[count-3]
      if segmentType == "movies" || segmentType == "shows" {
         s.Type = segmentType
         return nil
      }
   }
   // 2. Check for player types (singular) at position -2
   // URL: .../show/id
   segmentType := segments[count-2]
   if segmentType == "show" {
      s.Type = segmentType
      return nil
   }
   return errors.New("unrecognized content type")
}

func (l Login) Movie(show *ShowKey) (*Videos, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+l.Data.Attributes.Token)
   req.URL = &url.URL{
      Scheme: "https",
      Host:   api_host,
      RawQuery: url.Values{
         "include":          {"default"},
         "page[items.size]": {"1"},
      }.Encode(),
      Path: join(
         "/cms/routes/", strings.TrimSuffix(show.Type, "s"), "/", show.Id,
      ),
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
