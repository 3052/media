package main

import (
   "encoding/xml"
   "fmt"
   "os"
)

// =========================================================================
// STRUCT DEFINITIONS
// =========================================================================

type Client struct {
   Data struct {
      AccessToken  string `json:"access_token"`
      RefreshToken string `json:"refresh_token"`
   }
}

type Source struct {
   KeySystems struct {
      ComWidevineAlpha struct {
         LicenseUrl string `json:"license_url"`
      } `json:"com.widevine.alpha"`
   } `json:"key_systems"`
   Src  string
   Type string
}

type user_cache struct {
   Source []Source
   Client Client
   Mpd    struct {
      Body string
      Url  string
   }
   AmcnBcJwt string
}

// =========================================================================
// MAIN
// =========================================================================

func main() {
   // 1. Create Dummy Data
   payload := user_cache{
      Source: []Source{
         {
            Src:  "https://example.com/video.mpd",
            Type: "application/dash+xml",
            KeySystems: struct {
               ComWidevineAlpha struct {
                  LicenseUrl string `json:"license_url"`
               } `json:"com.widevine.alpha"`
            }{
               ComWidevineAlpha: struct {
                  LicenseUrl string `json:"license_url"`
               }{
                  LicenseUrl: "https://license.example.com/widevine",
               },
            },
         },
      },
      Client: Client{
         Data: struct {
            AccessToken  string `json:"access_token"`
            RefreshToken string `json:"refresh_token"`
         }{
            AccessToken:  "abcdef123456",
            RefreshToken: "ghijk7891011",
         },
      },
      Mpd: struct {
         Body string
         Url  string
      }{
         Body: "<MPD>string_content</MPD>",
         Url:  "https://example.com/manifest.mpd",
      },
      AmcnBcJwt: "header.payload.signature",
   }

   // ---------------------------------------------------------
   // XML Marshal
   // ---------------------------------------------------------
   xmlBytes, err := xml.MarshalIndent(payload, "", "  ")
   if err != nil {
      panic(err)
   }

   if err := os.WriteFile("user_cache.xml", xmlBytes, 0644); err != nil {
      panic(err)
   }
   fmt.Println("Written user_cache.xml")
}
