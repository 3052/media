package rakuten

import "strings"

// String implements the fmt.Stringer interface.
// It returns the ID, Title, and a unique list of available audio languages.
func (v *VideoItem) String() string {
   seen := make(map[string]bool)
   var data strings.Builder
   data.WriteString("title = ")
   data.WriteString(v.Title)
   data.WriteString("\nid = ")
   data.WriteString(v.Id)
   for _, stream := range v.ViewOptions.Private.Streams {
      for _, lang := range stream.AudioLanguages {
         if !seen[lang.Id] {
            seen[lang.Id] = true
            data.WriteString("\naudio language = ")
            data.WriteString(lang.Id)
         }
      }
   }
   return data.String()
}

type VideoItem struct {
   Title       string      `json:"title"`
   Id          string      `json:"id"`
   ViewOptions ViewOptions `json:"view_options"`
}

type ViewOptions struct {
   Private struct {
      Streams []Stream `json:"streams"`
   } `json:"private"`
}

type AudioLanguage struct {
   Id string `json:"id"`
}

type Stream struct {
   AudioLanguages []AudioLanguage `json:"audio_languages"`
}

func (t TvShowData) String() string {
   var data strings.Builder
   for i, season := range t.Seasons {
      if i >= 1 {
         data.WriteByte('\n')
      }
      data.WriteString("id = ")
      data.WriteString(season.Id)
   }
   return data.String()
}

type TvShowData struct {
   Seasons []struct {
      Id          string `json:"id"`
   } `json:"seasons"`
}

type StreamData struct {
   StreamInfos []struct {
      LicenseUrl string `json:"license_url"`
      Url        string `json:"url"`
   } `json:"stream_infos"`
}

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

var Quality = struct {
   FHD VideoQuality
   HD  VideoQuality
}{
   FHD: "FHD",
   HD:  "HD",
}

// PlayerType defines the allowed player types/DRM schemes.
type PlayerType string

var Player = struct {
   PlayReady PlayerType
   Widevine  PlayerType
}{
   PlayReady: DeviceID + ":DASH-CENC:PR",
   Widevine:  DeviceID + ":DASH-CENC:WVM",
}

// --- Shared Structs for Nested JSON ---

// --- Response Data Structs ---

type SeasonData struct {
   Episodes []VideoItem `json:"episodes"`
}

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
