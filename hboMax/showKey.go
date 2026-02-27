package hboMax

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strings"
)

// https://hbomax.com/movies/one-battle-after-another/bebe611d-8178-481a-a4f2-de743b5b135a
// https://hbomax.com/at/en/movies/austin-powers-international-man-of-mystery/a979fb8b-f713-4de3-a625-d16ad4d37448
// https://hbomax.com/shows/white-lotus/14f9834d-bc23-41a8-ab61-5c8abdbea505
func (s *ShowKey) Parse(rawLink string) error {
   link, err := url.Parse(rawLink)
   if err != nil {
      return err
   }
   segments := strings.Split(strings.TrimPrefix(link.Path, "/"), "/")
   n := len(segments)
   // We expect structure: .../[type]/[slug]/[id]
   if n < 3 {
      return errors.New("invalid URL format: not enough path segments")
   }
   s.Id = segments[n-1]
   s.Type = segments[n-3]
   switch s.Type {
   case "movies", "shows":
      return nil
   default:
      return errors.New("unrecognized content type")
   }
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

type ShowKey struct {
   Type string
   Id   string
}
