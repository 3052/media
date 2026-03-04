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

// ParseURL attempts to parse a URL and populate the Media struct.
func (m *Media) ParseURL(rawLink string) error {
   link, err := url.Parse(rawLink)
   if err != nil {
      return err
   }
   marketCode, err := extractMarketCode(link.Path)
   if err != nil {
      return err
   }
   m.MarketCode = marketCode

   // 1. Check Query Parameters
   query := link.Query()
   contentType := query.Get("content_type")
   if contentType == "movies" || contentType == "tv_shows" {
      var id string
      if contentType == "movies" {
         id = query.Get("content_id")
         if id == "" {
            return errors.New("url missing content_id param")
         }
      } else {
         id = query.Get("tv_show_id")
         if id == "" {
            return errors.New("url missing tv_show_id param")
         }
      }
      m.Id = id
      m.Type = ContentType(contentType)
      return nil
   }

   // 2. Check Path Segments
   path := strings.Trim(link.Path, "/")
   segments := strings.Split(path, "/")

   for _, seg := range segments {
      if seg == "movies" || seg == "tv_shows" {
         id := segments[len(segments)-1]
         if id == seg {
            return fmt.Errorf("url does not contain a specific %s id", seg)
         }
         m.Id = id
         m.Type = ContentType(seg)
         return nil
      }
   }
   return errors.New("not a movie or tv show url")
}
