package rakuten

import (
   "bytes"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (s *StreamInfo) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      s.LicenseUrl, "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (c *Content) String() string {
   var b strings.Builder
   b.WriteString("title = ")
   b.WriteString(c.Title)
   b.WriteString("\ncontent id = ")
   b.WriteString(c.Id)
   id := map[string]struct{}{}
   for _, stream := range c.ViewOptions.Private.Streams {
      for _, language := range stream.AudioLanguages {
         _, ok := id[language.Id]
         if !ok {
            b.WriteString("\naudio language = ")
            b.WriteString(language.Id)
            id[language.Id] = struct{}{}
         }
      }
   }
   return b.String()
}

func (s *Season) String() string {
   var b strings.Builder
   b.WriteString("show title = ")
   b.WriteString(s.TvShowTitle)
   b.WriteString("\nseason id = ")
   b.WriteString(s.Id)
   return b.String()
}

const (
   Fhd Quality = "FHD"
   Hd  Quality = "HD"
)

func (a *Address) Wvm(
   content_id, audio_language string, video Quality,
) (*StreamInfo, error) {
   return a.streamInfo(content_id, audio_language, ":DASH-CENC:WVM", video)
}

func (a *Address) Pr(
   content_id, audio_language string, video Quality,
) (*StreamInfo, error) {
   return a.streamInfo(content_id, audio_language, ":DASH-CENC:PR", video)
}

func (a *Address) Parse(data string) error {
   parsed, err := url.Parse(data)
   if err != nil {
      return err
   }
   path := strings.Split(strings.Trim(parsed.Path, "/"), "/")
   query := parsed.Query()
   if len(path) > 0 {
      a.MarketCode = path[0]
   }
   if content_type := query.Get("content_type"); content_type != "" {
      a.ContentType = content_type
      a.ContentId = query.Get("content_id")
      a.TvShowId = query.Get("tv_show_id")
   } else {
      if len(path) > 1 {
         a.ContentType = path[1]
      }
      if len(path) > 2 {
         if a.ContentType == "tv_shows" {
            a.TvShowId = path[2]
         } else {
            a.ContentId = path[2]
         }
      }
   }
   return nil
}
