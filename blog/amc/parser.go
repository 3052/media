// Package amc provides tools to parse specific AMC data structures.
package amc

import (
   "encoding/json"
   "errors"
)

// Season represents the clean, extracted data for a single season.
// This is part of the public API of our package.
type Season struct {
   Title string
   NID   int
}

// response holds the top-level structure of the JSON data.
// It is unexported as it's an internal implementation detail.
type response struct {
   Data struct {
      Children []*node `json:"children"`
   } `json:"data"`
}

// node represents a single element in the JSON tree.
// It is also unexported.
type node struct {
   Type       string  `json:"type"`
   Children   []*node `json:"children,omitempty"`
   Properties *struct {
      Text *struct {
         Title struct {
            Title string `json:"title"`
         } `json:"title"`
      } `json:"text,omitempty"`
      Metadata *struct {
         Title string `json:"title"`
         NID   int    `json:"nid"`
      } `json:"metadata,omitempty"`
   } `json:"properties,omitempty"`
}

// findSeasonsTabNode is an unexported helper method to find the "Seasons" tab.
func (r *response) findSeasonsTabNode() (*node, bool) {
   for _, topLevelChild := range r.Data.Children {
      if topLevelChild.Type == "tab_bar" {
         for _, tabItem := range topLevelChild.Children {
            if tabItem.Properties != nil {
               if tabItem.Properties.Text != nil {
                  if tabItem.Properties.Text.Title.Title == "Seasons" {
                     return tabItem, true
                  }
               }
            }
         }
      }
   }
   return nil, false
}

// ExtractSeasons parses the raw JSON data and returns a slice of Seasons.
// This is the main exported function for our package.
func ExtractSeasons(jsonData []byte) ([]Season, error) {
   var res response
   if err := json.Unmarshal(jsonData, &res); err != nil {
      return nil, err // Return error on bad JSON
   }

   seasonsTabNode, found := res.findSeasonsTabNode()
   if !found {
      return nil, errors.New("could not find the 'Seasons' tab in the JSON data")
   }

   if len(seasonsTabNode.Children) > 0 {
      firstSeasonChild := seasonsTabNode.Children[0]
      if firstSeasonChild.Type == "tab_bar" {
         seasonsList := firstSeasonChild.Children
         // Pre-allocate slice with the right capacity for efficiency
         extractedSeasons := make([]Season, 0, len(seasonsList))

         for _, seasonNode := range seasonsList {
            if seasonNode.Properties != nil && seasonNode.Properties.Metadata != nil {
               meta := seasonNode.Properties.Metadata
               extractedSeasons = append(extractedSeasons, Season{
                  Title: meta.Title,
                  NID:   meta.NID,
               })
            }
         }
         return extractedSeasons, nil
      }
   }

   return nil, errors.New("could not find the list of seasons inside the 'Seasons' tab")
}
