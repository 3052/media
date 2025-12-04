package rakuten

import (
   "fmt"
   "net/url"
   "strings"
)

// ParseURL attempts to parse a URL and populate the Movie struct.
// Usage:
//
//   var m Movie
//   err := m.ParseURL("https://rakuten.tv/...")
func (m *Movie) ParseURL(rawURL string) error {
   u, err := url.Parse(rawURL)
   if err != nil {
      return fmt.Errorf("invalid url: %w", err)
   }

   marketCode, err := extractMarketCode(u.Path)
   if err != nil {
      return err
   }

   // 1. Check Query Parameters
   // Pattern: ?content_type=movies&content_id=...
   q := u.Query()
   if q.Get("content_type") == "movies" {
      id := q.Get("content_id")
      if id == "" {
         return fmt.Errorf("url missing content_id param")
      }

      m.MarketCode = marketCode
      m.Id = id
      return nil
   }

   // 2. Check Path Segments
   // Pattern: /nl/movies/id OR /nl/player/movies/stream/id
   path := strings.Trim(u.Path, "/")
   segments := strings.Split(path, "/")

   for _, seg := range segments {
      if seg == "movies" {
         // Assuming the ID is the last segment
         id := segments[len(segments)-1]
         if id == "movies" {
            return fmt.Errorf("url does not contain a specific movie id")
         }

         m.MarketCode = marketCode
         m.Id = id
         return nil
      }
   }

   return fmt.Errorf("not a movie url")
}

// ParseURL attempts to parse a URL and populate the TvShow struct.
// Usage:
//
//   var t TvShow
//   err := t.ParseURL("https://rakuten.tv/...")
func (t *TvShow) ParseURL(rawURL string) error {
   u, err := url.Parse(rawURL)
   if err != nil {
      return fmt.Errorf("invalid url: %w", err)
   }

   marketCode, err := extractMarketCode(u.Path)
   if err != nil {
      return err
   }

   // 1. Check Query Parameters
   // Pattern: ?content_type=tv_shows&tv_show_id=...
   q := u.Query()
   if q.Get("content_type") == "tv_shows" {
      id := q.Get("tv_show_id")
      if id == "" {
         return fmt.Errorf("url missing tv_show_id param")
      }

      t.MarketCode = marketCode
      t.Id = id
      return nil
   }

   // 2. Check Path Segments
   // Pattern: /nl/tv_shows/id
   path := strings.Trim(u.Path, "/")
   segments := strings.Split(path, "/")

   for _, seg := range segments {
      if seg == "tv_shows" {
         // Assuming the ID is the last segment
         id := segments[len(segments)-1]
         if id == "tv_shows" {
            return fmt.Errorf("url does not contain a specific tv show id")
         }

         t.MarketCode = marketCode
         t.Id = id
         return nil
      }
   }

   return fmt.Errorf("not a tv show url")
}

// extractMarketCode extracts the first segment of the path (e.g., "nl", "uk").
func extractMarketCode(path string) (string, error) {
   trimmed := strings.Trim(path, "/")

   // Check if we have anything left after trimming
   if trimmed == "" {
      return "", fmt.Errorf("could not determine market code from path")
   }

   segments := strings.Split(trimmed, "/")
   return segments[0], nil
}
