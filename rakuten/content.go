package rakuten

import (
   "errors"
   "fmt"
   "net/url"
   "strings"
)

// extractMarketCode extracts the first segment of the path (e.g., "nl", "uk").
func extractMarketCode(path string) (string, error) {
   trimmed := strings.Trim(path, "/")
   // Check if we have anything left after trimming
   if trimmed == "" {
      return "", errors.New("could not determine market code from path")
   }
   segments := strings.Split(trimmed, "/")
   return segments[0], nil
}

// github.com/pandvan/rakuten-m3u-generator/blob/master/rakuten.py
var classificationMap = map[string]int{
   "cz": 272,
   "es": 5,
   "fr": 23,
   "ie": 41,
   "nl": 69,
   "pl": 277,
   "pt": 64,
   "se": 282,
   "uk": 18,
}

type alfa int

const (
   Movies alfa = iota
   TvShows
)

type Content struct {
   Id         string
   MarketCode string
   Type       alfa
}

func (c *Content) Parse(urlData string) error {
   urlParse, err := url.Parse(urlData)
   if err != nil {
      return err
   }
   marketCode, err := extractMarketCode(urlParse.Path)
   if err != nil {
      return err
   }
   c.MarketCode = marketCode
   // 1. Check Query Parameters
   query := urlParse.Query()
   // 'contentType' here is the URL parameter value (e.g. "movies", "tv_shows")
   contentType := query.Get("content_type")
   if contentType == "movies" || contentType == "tv_shows" {
      var id string
      if contentType == "movies" {
         id = query.Get("content_id")
         if id == "" {
            return errors.New("url missing content_id param")
         }
         c.Type = Movies
      } else {
         id = query.Get("tv_show_id")
         if id == "" {
            return errors.New("url missing tv_show_id param")
         }
         c.Type = TvShows
      }
      c.Id = id
      return nil
   }
   // 2. Check Path Segments
   path := strings.Trim(urlParse.Path, "/")
   segments := strings.Split(path, "/")
   for _, seg := range segments {
      if seg == "movies" || seg == "tv_shows" {
         id := segments[len(segments)-1]
         if id == seg {
            return fmt.Errorf("url does not contain a specific %s id", seg)
         }
         c.Id = id
         if seg == "movies" {
            c.Type = Movies
         } else {
            c.Type = TvShows
         }
         return nil
      }
   }
   return errors.New("not a movie or tv show url")
}
