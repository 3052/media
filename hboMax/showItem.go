// https://hbomax.com/at/en/movies/austin-powers-international-man-of-mystery/a979fb8b-f713-4de3-a625-d16ad4d37448
// https://hbomax.com/movies/one-battle-after-another/bebe611d-8178-481a-a4f2-de743b5b135a
// https://hbomax.com/shows/white-lotus/14f9834d-bc23-41a8-ab61-5c8abdbea505
// https://play.hbomax.com/movie/b7b66574-c6e3-4ed3-a266-6bc44180252e
// https://play.hbomax.com/show/31cb4b84-951a-4daf-8925-746fcdcddcb8
package hboMax

import (
   "errors"
   "net/url"
   "strings"
)

func ParseUrl(inputUrl string) (*ShowItem, error) {
   parsedUrl, err := url.Parse(inputUrl)
   if err != nil {
      return nil, err
   }
   path := strings.TrimPrefix(parsedUrl.Path, "/")
   segments := strings.Split(path, "/")
   count := len(segments)
   if count < 2 {
      return nil, errors.New("invalid url path")
   }
   // Create the instance
   show := ShowItem{
      Id: segments[count-1],
   }
   // Check immediate parent (e.g., /movie/id)
   if count >= 2 && isCategory(segments[count-2]) {
      show.Category = segments[count-2]
      return &show, nil
   }
   // Check grandparent (e.g., /movies/slug/id)
   if count >= 3 && isCategory(segments[count-3]) {
      show.Category = segments[count-3]
      return &show, nil
   }
   return nil, errors.New("category not found")
}

type ShowItem struct {
   Category string
   Id       string
}

func isCategory(segment string) bool {
   switch segment {
   case "movies", "shows", "movie", "show":
      return true
   default:
      return false
   }
}
