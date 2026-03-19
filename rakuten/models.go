package rakuten

import "net/url"

// Content represents the parsed Rakuten URL data
type Content struct {
   Id               string
   MarketCode       string
   Type             string
   ClassificationId int
}

// Constants for device and player configuration
const DeviceId = "atvui40"

const (
   PlayReady Player = DeviceId + ":DASH-CENC:PR"
   Widevine  Player = DeviceId + ":DASH-CENC:WVM"
)

const (
   Fhd VideoQuality = "FHD"
   Hd  VideoQuality = "HD"
)

type VideoQuality string

type Player string

type Stream struct {
   StreamInfos []struct {
      LicenseUrl string `json:"license_url"`
      Url        string `json:"url"`
   } `json:"stream_infos"`
}

type Season struct {
   Episodes []MovieOrEpisode `json:"episodes"`
}

type TvShow struct {
   Seasons []struct {
      Id string `json:"id"`
   } `json:"seasons"`
}

type MovieOrEpisode struct {
   Title       string `json:"title"`
   Id          string `json:"id"`
   ViewOptions struct {
      Private struct {
         Streams []struct {
            AudioLanguages []struct {
               Id string `json:"id"`
            } `json:"audio_languages"`
         } `json:"streams"`
      } `json:"private"`
   } `json:"view_options"`
}

type Dash struct {
   Body []byte
   Url  *url.URL
}
