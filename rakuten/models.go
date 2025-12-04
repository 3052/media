package rakuten

import (
   "fmt"
   "net/url"
   "strings"
)

type Cache struct {
   Movie   *Movie
   Mpd     *url.URL
   MpdBody []byte
   TvShow  *TvShow
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

// String implements the fmt.Stringer interface.
// It returns the ID, Title, and a unique list of available audio languages.
func (v VideoItem) String() string {
   seen := make(map[string]bool)
   var langs []string

   for _, stream := range v.ViewOptions.Private.Streams {
      for _, lang := range stream.AudioLanguages {
         if !seen[lang.Id] {
            seen[lang.Id] = true
            langs = append(langs, lang.Id)
         }
      }
   }

   return fmt.Sprintf("%s - %s [%s]", v.Id, v.Title, strings.Join(langs, ", "))
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

// String implements the fmt.Stringer interface.
// It returns the TV Show Title followed by a list of Season IDs.
func (t TvShowData) String() string {
   if len(t.Seasons) == 0 {
      return "No seasons available"
   }

   // Assume the title is the same for all seasons in this response
   title := t.Seasons[0].TvShowTitle
   var seasonIDs []string

   for _, s := range t.Seasons {
      seasonIDs = append(seasonIDs, s.Id)
   }

   return fmt.Sprintf("%s - Seasons [%s]", title, strings.Join(seasonIDs, ", "))
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
