package rakuten

// DeviceID is the default identifier used for requests.
const DeviceID = "atvui40"

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

// VideoQuality defines the allowed video qualities for streaming.
type VideoQuality string

const (
   VideoQualityHD  VideoQuality = "HD"
   VideoQualityFHD VideoQuality = "FHD"
)

// PlayerType defines the allowed player types/DRM schemes.
type PlayerType string

const (
   PlayerPlayReady PlayerType = DeviceID + ":DASH-CENC:PR"
   PlayerWidevine  PlayerType = DeviceID + ":DASH-CENC:WVM"
)

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
   AudioQuality             string       `json:"audio_quality"`
   DeviceIdentifier         string       `json:"device_identifier"`
   DeviceSerial             string       `json:"device_serial"`
   SubtitleLanguage         string       `json:"subtitle_language"`
   VideoType                string       `json:"video_type"`
   Player                   PlayerType   `json:"player"`
   ClassificationId         int          `json:"classification_id"`
   ContentType              string       `json:"content_type"`
   DeviceStreamVideoQuality VideoQuality `json:"device_stream_video_quality"`
   AudioLanguage            string       `json:"audio_language"`
   ContentId                string       `json:"content_id"`
}
