package rakuten

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strconv"
)

// classificationMap maps market codes to their internal classification IDs.
var classificationMap = map[string]int{
   "cz": 272,
   "dk": 283,
   "es": 5,
   "fr": 23,
   "nl": 69,
   "pl": 277,
   "pt": 64,
   "se": 282,
   "uk": 18,
}

// --- Shared Structs for Nested JSON ---

type AudioLanguage struct {
   Id string `json:"id"`
}

type Stream struct {
   AudioLanguages []AudioLanguage `json:"audio_languages"`
}

type ViewOptions struct {
   Private struct {
      Streams []Stream `json:"streams"`
   } `json:"private"`
}

// VideoItem represents the structure for both Movies and Episodes.
type VideoItem struct {
   Id          string      `json:"id"`
   Title       string      `json:"title"`
   ViewOptions ViewOptions `json:"view_options"`
}

// --- Response Data Structs ---

type SeasonData struct {
   Episodes []VideoItem `json:"episodes"`
}

type TvShowData struct {
   Seasons []struct {
      TvShowTitle string `json:"tv_show_title"`
      Id          string `json:"id"`
   } `json:"seasons"`
}

// StreamData represents the inner data object of a stream request.
type StreamData struct {
   StreamInfos []struct {
      LicenseUrl string `json:"license_url"`
      Url        string `json:"url"`
   } `json:"stream_infos"`
}

// --- Request Payload Struct ---

type StreamRequestPayload struct {
   AudioQuality             string `json:"audio_quality"`
   DeviceIdentifier         string `json:"device_identifier"`
   DeviceSerial             string `json:"device_serial"`
   SubtitleLanguage         string `json:"subtitle_language"`
   VideoType                string `json:"video_type"`
   Player                   string `json:"player"`
   ClassificationId         int    `json:"classification_id"`
   ContentType              string `json:"content_type"`
   DeviceStreamVideoQuality string `json:"device_stream_video_quality"`
   AudioLanguage            string `json:"audio_language"`
   ContentId                string `json:"content_id"`
}

// --- Helper Functions ---

// buildURL handles the common logic for constructing GET request URLs.
func buildURL(marketCode, endpoint, id string) (string, error) {
   classID, ok := classificationMap[marketCode]
   if !ok {
      return "", fmt.Errorf("unsupported market code: %s", marketCode)
   }

   baseURL := fmt.Sprintf("https://gizmo.rakuten.tv/v3/%s/%s", endpoint, id)

   params := url.Values{}
   params.Add("device_identifier", "atvui40")
   params.Add("market_code", marketCode)
   params.Add("classification_id", strconv.Itoa(classID))

   return fmt.Sprintf("%s?%s", baseURL, params.Encode()), nil
}

// makeStreamRequest handles the common logic for the POST stream request and parsing.
func makeStreamRequest(marketCode, contentType, contentId, player, videoQuality, audioLanguage string) (*StreamData, error) {
   classID, ok := classificationMap[marketCode]
   if !ok {
      return nil, fmt.Errorf("unsupported market code: %s", marketCode)
   }

   payload := StreamRequestPayload{
      AudioQuality:             "2.0",
      DeviceIdentifier:         "atvui40",
      DeviceSerial:             "not implemented",
      SubtitleLanguage:         "MIS",
      VideoType:                "stream",
      Player:                   player,
      ClassificationId:         classID,
      ContentType:              contentType,
      DeviceStreamVideoQuality: videoQuality,
      AudioLanguage:            audioLanguage,
      ContentId:                contentId,
   }

   jsonData, err := json.Marshal(payload)
   if err != nil {
      return nil, fmt.Errorf("failed to marshal request body: %w", err)
   }

   url := "https://gizmo.rakuten.tv/v3/avod/streamings"
   req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Content-Type", "application/json")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
   }

   var wrapper struct {
      Data StreamData `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, fmt.Errorf("failed to decode response: %w", err)
   }

   return &wrapper.Data, nil
}

// --- Movie Type and Methods ---

type Movie struct {
   Id         string // Matches "content_id" in URLs
   MarketCode string
}

// Request fetches movie details (GET).
func (m *Movie) Request() (*VideoItem, error) {
   fullURL, err := buildURL(m.MarketCode, "movies", m.Id)
   if err != nil {
      return nil, err
   }

   resp, err := http.Get(fullURL)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
   }

   var wrapper struct {
      Data VideoItem `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, fmt.Errorf("failed to decode response: %w", err)
   }

   return &wrapper.Data, nil
}

// RequestStream requests the stream for this movie (POST).
func (m *Movie) RequestStream(player, videoQuality, audioLanguage string) (*StreamData, error) {
   return makeStreamRequest(m.MarketCode, "movies", m.Id, player, videoQuality, audioLanguage)
}

// --- TvShow Type and Methods ---

type TvShow struct {
   Id         string // Matches "tv_show_id" in URLs
   MarketCode string
}

// Request fetches TV show details like seasons (GET).
func (t *TvShow) Request() (*TvShowData, error) {
   fullURL, err := buildURL(t.MarketCode, "tv_shows", t.Id)
   if err != nil {
      return nil, err
   }

   resp, err := http.Get(fullURL)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
   }

   var wrapper struct {
      Data TvShowData `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, fmt.Errorf("failed to decode response: %w", err)
   }

   return &wrapper.Data, nil
}

// RequestSeason fetches episodes for a specific season (GET).
func (t *TvShow) RequestSeason(seasonId string) (*SeasonData, error) {
   fullURL, err := buildURL(t.MarketCode, "seasons", seasonId)
   if err != nil {
      return nil, err
   }

   resp, err := http.Get(fullURL)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
   }

   var wrapper struct {
      Data SeasonData `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, fmt.Errorf("failed to decode response: %w", err)
   }

   return &wrapper.Data, nil
}

// RequestStream requests the stream for a specific Episode (POST).
// Note: You must provide the episodeId (retrieved from RequestSeason).
func (t *TvShow) RequestStream(episodeId, player, videoQuality, audioLanguage string) (*StreamData, error) {
   // For TV content, the standard Rakuten content_type for streaming an episode is "episodes"
   return makeStreamRequest(t.MarketCode, "episodes", episodeId, player, videoQuality, audioLanguage)
}
