package rakuten

import (
   "errors"
   "fmt"
   "net/url"
   "strings"
)

// Media represents a piece of content, which can be a Movie or a TV Show.
type Media struct {
   Id         string // Matches "content_id" or "tv_show_id" in URLs
   MarketCode string
   Type       ContentType
}

func ParseMedia(rawLink string) (*Media, error) {
   link, err := url.Parse(rawLink)
   if err != nil {
      return nil, err
   }
   // Assuming extractMarketCode is defined elsewhere in your package
   marketCode, err := extractMarketCode(link.Path)
   if err != nil {
      return nil, err
   }
   // Initialize the struct here
   m := Media{MarketCode: marketCode}
   // 1. Check Query Parameters
   query := link.Query()
   contentType := query.Get("content_type")
   if contentType == "movies" || contentType == "tv_shows" {
      var id string
      if contentType == "movies" {
         id = query.Get("content_id")
         if id == "" {
            return nil, errors.New("url missing content_id param")
         }
      } else {
         id = query.Get("tv_show_id")
         if id == "" {
            return nil, errors.New("url missing tv_show_id param")
         }
      }
      m.Id = id
      m.Type = ContentType(contentType)
      return &m, nil
   }
   // 2. Check Path Segments
   path := strings.Trim(link.Path, "/")
   segments := strings.Split(path, "/")
   for _, seg := range segments {
      if seg == "movies" || seg == "tv_shows" {
         id := segments[len(segments)-1]
         if id == seg {
            return nil, fmt.Errorf("url does not contain a specific %s id", seg)
         }
         m.Id = id
         m.Type = ContentType(seg)
         return &m, nil
      }
   }
   return nil, errors.New("not a movie or tv show url")
}
