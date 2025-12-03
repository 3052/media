package rakuten

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

// --- Helper Functions ---

// buildURL handles the common logic for constructing GET request URLs.
func buildURL(marketCode, endpoint, id string) (string, error) {
   classID, ok := classificationMap[marketCode]
   if !ok {
      return "", fmt.Errorf("unsupported market code: %s", marketCode)
   }

   baseURL := fmt.Sprintf("https://gizmo.rakuten.tv/v3/%s/%s", endpoint, id)

   params := url.Values{}
   params.Add("device_identifier", DeviceID)
   params.Add("market_code", marketCode)
   params.Add("classification_id", strconv.Itoa(classID))

   return fmt.Sprintf("%s?%s", baseURL, params.Encode()), nil
}

// makeStreamRequest handles the common logic for the POST stream request and parsing.
func makeStreamRequest(marketCode, contentType, contentId string, player PlayerType, quality VideoQuality, audioLanguage string) (*StreamData, error) {
   classID, ok := classificationMap[marketCode]
   if !ok {
      return nil, fmt.Errorf("unsupported market code: %s", marketCode)
   }

   payload := StreamRequestPayload{
      AudioQuality:             "2.0",
      DeviceIdentifier:         DeviceID,
      DeviceSerial:             "not implemented",
      SubtitleLanguage:         "MIS",
      VideoType:                "stream",
      Player:                   player,
      ClassificationId:         classID,
      ContentType:              contentType,
      DeviceStreamVideoQuality: quality,
      AudioLanguage:            audioLanguage,
      ContentId:                contentId,
   }

   jsonData, err := json.Marshal(payload)
   if err != nil {
      return nil, fmt.Errorf("failed to marshal request body: %w", err)
   }

   apiURL := "https://gizmo.rakuten.tv/v3/avod/streamings"
   req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
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
// Arguments: audioLanguage, player, quality.
func (m *Movie) RequestStream(audioLanguage string, player PlayerType, quality VideoQuality) (*StreamData, error) {
   return makeStreamRequest(m.MarketCode, "movies", m.Id, player, quality, audioLanguage)
}

func (s StreamData) Widevine(data []byte) ([]byte, error) {
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
// Arguments: episodeId, audioLanguage, player, quality.
func (t *TvShow) RequestStream(episodeId string, audioLanguage string, player PlayerType, quality VideoQuality) (*StreamData, error) {
   // For TV content, the standard Rakuten content_type for streaming an episode is "episodes"
   return makeStreamRequest(t.MarketCode, "episodes", episodeId, player, quality, audioLanguage)
}
