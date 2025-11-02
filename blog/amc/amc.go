package main

import (
   "encoding/json"
   "fmt"
   "io/ioutil"
   "log"
)

// More accurate structs to match the JSON structure for the specific data needed.
type Response struct {
   Data struct {
      Children []Node `json:"children"`
   } `json:"data"`
}

type Node struct {
   Type     string      `json:"type"`
   Children []Node      `json:"children,omitempty"`
   Properties *Properties `json:"properties,omitempty"`
}

type Properties struct {
   Text     *Text     `json:"text,omitempty"`
   Metadata *Metadata `json:"metadata,omitempty"`
}

type Text struct {
   Title struct {
      Title string `json:"title"`
   } `json:"title"`
}

// Metadata struct now includes all fields we need to access at different levels.
type Metadata struct {
   Title string `json:"title"`
   NID   int    `json:"nid"`
}


func main() {
   // Read the JSON data from the file 'series-detail.json'
   jsonData, err := ioutil.ReadFile("series-detail.json")
   if err != nil {
      log.Fatalf("Error reading JSON file: %v", err)
   }

   var response Response

   // Unmarshal the JSON data into our Go structs
   if err := json.Unmarshal(jsonData, &response); err != nil {
      log.Fatalf("Error unmarshaling JSON: %v", err)
   }

   // Find the "Seasons" tab item by its title, which is more robust
   var seasonsTabNode Node
   foundSeasonsTab := false

   // The main tab bar is the second child in this specific JSON
   if len(response.Data.Children) > 1 && response.Data.Children[1].Type == "tab_bar" {
      mainTabBar := response.Data.Children[1]
      for _, tabItem := range mainTabBar.Children {
         // Check if the tab's title is "Seasons"
         if tabItem.Properties != nil && tabItem.Properties.Text != nil && tabItem.Properties.Text.Title.Title == "Seasons" {
            seasonsTabNode = tabItem
            foundSeasonsTab = true
            break
         }
      }
   }

   if !foundSeasonsTab {
      log.Fatal("Could not find the 'Seasons' tab in the JSON data.")
      return
   }

   // The actual list of seasons is nested within a child tab_bar
   if len(seasonsTabNode.Children) > 0 && seasonsTabNode.Children[0].Type == "tab_bar" {
      seasonsList := seasonsTabNode.Children[0].Children
      // Loop through the seasons and print the desired information.
      for _, season := range seasonsList {
         if season.Properties != nil && season.Properties.Metadata != nil {
            meta := season.Properties.Metadata
            fmt.Printf("title = %s\nnid = %d\n\n", meta.Title, meta.NID)
         }
      }
   } else {
      log.Fatal("Could not find the list of seasons inside the 'Seasons' tab.")
   }
}
