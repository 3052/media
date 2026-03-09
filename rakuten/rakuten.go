package rakuten

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

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

// Parse extracts metadata from a Rakuten URL and populates the Content struct
func (c *Content) Parse(urlData string) error {
   urlParse, err := url.Parse(urlData)
   if err != nil {
      return err
   }

   // Trim prefix once and extract the market code
   path := strings.TrimPrefix(urlParse.Path, "/")
   c.MarketCode, _, _ = strings.Cut(path, "/")

   // Check if the market code exists in the map and set ClassificationId
   var ok bool
   c.ClassificationId, ok = classificationMap[c.MarketCode]
   if !ok {
      return errors.New("unknown market code")
   }

   // 1. Check Query Parameters
   query := urlParse.Query()
   contentType := query.Get("content_type")
   switch contentType {
   case "movies":
      c.Id = query.Get("content_id")
      c.Type = contentType
      return nil
   case "tv_shows":
      c.Id = query.Get("tv_show_id")
      c.Type = contentType
      return nil
   }

   // 2. Check Path Segments
   segments := strings.Split(path, "/")
   for _, segment := range segments {
      switch segment {
      case "movies", "tv_shows":
         c.Id = segments[len(segments)-1]
         c.Type = segment
         return nil
      }
   }
   return errors.New("not a movie or tv show url")
}

// buildUrl constructs the API URL using the stored ClassificationId
func (c *Content) buildUrl(endpoint, id string) string {
   urlData := url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   join("/v3/", endpoint, "/", id),
      RawQuery: url.Values{
         "classification_id": {strconv.Itoa(c.ClassificationId)},
         "device_identifier": {DeviceId},
         "market_code":       {c.MarketCode},
      }.Encode(),
   }
   return urlData.String()
}

// RequestSeason fetches episodes for a specific season (GET).
// This method is only applicable to TV Shows.
func (c *Content) RequestSeason(seasonId string) (*Season, error) {
   if c.Type != "tv_shows" {
      return nil, errors.New("cannot request season for a non-tv show content type")
   }

   fullUrl := c.buildUrl("seasons", seasonId)

   resp, err := http.Get(fullUrl)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }

   var result struct {
      Data Season
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Data, nil
}

func (c *Content) RequestTvShow() (*TvShow, error) {
   if c.Type != "tv_shows" {
      return nil, errors.New("cannot request tv show details for a non-tv show content type")
   }

   fullUrl := c.buildUrl("tv_shows", c.Id)

   resp, err := http.Get(fullUrl)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }

   var result struct {
      Data TvShow
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Data, nil
}

// RequestMovie fetches movie details (GET).
func (c *Content) RequestMovie() (*MovieOrEpisode, error) {
   if c.Type != "movies" {
      return nil, errors.New("cannot request movie details for a non-movie content type")
   }

   fullUrl := c.buildUrl("movies", c.Id)

   resp, err := http.Get(fullUrl)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }

   var result struct {
      Data MovieOrEpisode
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Data, nil
}

func (c *Content) makeStreamRequest(contentType, contentId string, playerData Player, quality VideoQuality, audioLanguage string) (*Stream, error) {
   data, err := json.Marshal(map[string]string{
      "audio_language":              audioLanguage,
      "audio_quality":               "2.0",
      "classification_id":           strconv.Itoa(c.ClassificationId),
      "content_id":                  contentId,
      "content_type":                contentType,
      "device_identifier":           DeviceId,
      "device_serial":               "not implemented",
      "device_stream_video_quality": string(quality),
      "player":                      string(playerData),
      "subtitle_language":           "MIS",
      "video_type":                  "stream",
   })
   if err != nil {
      return nil, err
   }

   resp, err := http.Post(
      "https://gizmo.rakuten.tv/v3/avod/streamings", "application/json",
      bytes.NewBuffer(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      Data   Stream
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, errors.New(result.Errors[0].Message)
   }
   return &result.Data, nil
}

func (c *Content) EpisodeStream(episodeId, audioLanguage string, playerData Player, quality VideoQuality) (*Stream, error) {
   if c.Type != "tv_shows" {
      return nil, errors.New("cannot request an episode stream for non-tv-show content")
   }
   return c.makeStreamRequest("episodes", episodeId, playerData, quality, audioLanguage)
}

func (c *Content) MovieStream(audioLanguage string, playerData Player, quality VideoQuality) (*Stream, error) {
   if c.Type != "movies" {
      return nil, errors.New("cannot request a movie stream for non-movie content")
   }
   return c.makeStreamRequest("movies", c.Id, playerData, quality, audioLanguage)
}

// join takes a variable number of strings and returns them combined into one string
func join(strs ...string) string {
   return strings.Join(strs, "")
}

// String implementation for MovieOrEpisode to pretty print details
func (m *MovieOrEpisode) String() string {
   seen := make(map[string]bool)
   var data strings.Builder
   data.WriteString("title = ")
   data.WriteString(m.Title)
   data.WriteString("\nid = ")
   data.WriteString(m.Id)
   for _, streamData := range m.ViewOptions.Private.Streams {
      for _, language := range streamData.AudioLanguages {
         if !seen[language.Id] {
            seen[language.Id] = true
            data.WriteString("\naudio language = ")
            data.WriteString(language.Id)
         }
      }
   }
   return data.String()
}

func (t TvShow) String() string {
   var data strings.Builder
   for i, season := range t.Seasons {
      if i >= 1 {
         data.WriteByte('\n')
      }
      data.WriteString("season id = ")
      data.WriteString(season.Id)
   }
   return data.String()
}

func (s Stream) Dash() (*Dash, error) {
   resp, err := http.Get(s.StreamInfos[0].Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Dash{Body: body, Url: resp.Request.URL}, nil
}

func (s Stream) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      s.StreamInfos[0].LicenseUrl, "application/x-protobuf",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
