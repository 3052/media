package amc

import "errors"

// Metadata represents the clean, extracted data for a single item.
type Metadata struct {
   Title         string `json:"title"`
   NID           int    `json:"nid"`
   EpisodeNumber int    `json-:"episodeNumber,omitempty"`
}

// Node represents any element in the AMC JSON manifest tree, including the root.
type Node struct {
   Type     string `json:"type"`
   Children   []*Node `json:"children,omitempty"`
   Properties struct {
      ManifestType string `json:"manifestType,omitempty"`
      Text *struct {
         Title struct {
            Title string `json:"title"`
         } `json:"title"`
      } `json:"text,omitempty"`
      Metadata *Metadata `json:"metadata,omitempty"`
   } `json:"properties"`
}

// --- Public Methods ---

// ExtractSeasons contains the logic for parsing a "series-detail" manifest.
func (n *Node) ExtractSeasons() ([]*Metadata, error) {
   for _, child := range n.Children {
      // Guard: Skip any root child that is not a tab_bar.
      if child.Type != "tab_bar" {
         continue
      }

      for _, tabItem := range child.Children {
         // Guard: Skip any tab that isn't the "Seasons" tab.
         if tabItem.Type != "tab_bar_item" {
            continue
         }
         if tabItem.Properties.Text == nil {
            continue
         }
         if tabItem.Properties.Text.Title.Title != "Seasons" {
            continue
         }

         // We've found the "Seasons" tab item. Now find the list inside it.
         for _, seasonListContainer := range tabItem.Children {
            // Guard: Skip any child that is not the tab_bar list container.
            if seasonListContainer.Type != "tab_bar" {
               continue
            }

            // Success: We found the list. Extract and return.
            seasonList := seasonListContainer.Children
            extractedMetadata := make([]*Metadata, 0, len(seasonList))
            for _, seasonNode := range seasonList {
               if seasonNode.Properties.Metadata != nil {
                  extractedMetadata = append(extractedMetadata, seasonNode.Properties.Metadata)
               }
            }
            return extractedMetadata, nil
         }
      }
   }

   // If all loops complete without returning, the target was not found.
   return nil, errors.New("could not find the seasons list within the manifest")
}

// ExtractEpisodes contains the logic for parsing a "season-episodes" manifest.
func (n *Node) ExtractEpisodes() ([]*Metadata, error) {
   for _, listNode := range n.Children {
      // Guard: Skip any child that is not the main list container.
      if listNode.Type != "list" {
         continue
      }

      // Success: We found the list. Extract and return.
      list := listNode.Children
      extractedMetadata := make([]*Metadata, 0, len(list))
      for _, cardNode := range list {
         if cardNode.Type == "card" && cardNode.Properties.Metadata != nil {
            extractedMetadata = append(extractedMetadata, cardNode.Properties.Metadata)
         }
      }
      return extractedMetadata, nil
   }
   return nil, errors.New("could not find episode list in the manifest")
}
