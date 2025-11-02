package amc

import (
   "encoding/json"
   "os"
   "testing"
)

func TestNode_ExtractSeasons(t *testing.T) {
   jsonData, err := os.ReadFile("series-detail.json")
   if err != nil {
      t.Fatalf("Failed to read test data file: %v", err)
   }

   var topLevel struct {
      Data json.RawMessage `json:"data"`
   }
   if err := json.Unmarshal(jsonData, &topLevel); err != nil {
      t.Fatalf("Failed to unmarshal top level of test data: %v", err)
   }

   var rootNode Node
   if err := json.Unmarshal(topLevel.Data, &rootNode); err != nil {
      t.Fatalf("Failed to unmarshal data into Node: %v", err)
   }

   expectedMetadata := []*Metadata{
      {Title: "Season 1", NID: 1010634, EpisodeNumber: -1},
      {Title: "Season 2", NID: 1010638, EpisodeNumber: -1},
      {Title: "Season 3", NID: 1010635, EpisodeNumber: -1},
      {Title: "Season 4", NID: 1010643, EpisodeNumber: -1},
      {Title: "Season 5", NID: 1010637, EpisodeNumber: -1},
   }

   // Act: Call the specific method.
   actualMetadata, err := rootNode.ExtractSeasons()
   if err != nil {
      t.Fatalf("ExtractSeasons() returned an unexpected error: %v", err)
   }

   if len(actualMetadata) != len(expectedMetadata) {
      t.Fatalf("ExtractSeasons() returned wrong number of items: got %d, want %d", len(actualMetadata), len(expectedMetadata))
   }

   for i := range actualMetadata {
      if *actualMetadata[i] != *expectedMetadata[i] {
         t.Errorf("Mismatch at index %d:\ngot:  %+v\nwant: %+v", i, *actualMetadata[i], *expectedMetadata[i])
      }
   }
}
