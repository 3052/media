package amc

import (
   "encoding/json"
   "os"
   "testing"
)

func TestNode_ExtractEpisodes(t *testing.T) {
   jsonData, err := os.ReadFile("season-episodes.json")
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
      {Title: "Nature Under Constraint and Vexed", NID: 1011160, EpisodeNumber: 1},
      {Title: "Governed by Sound Reason and True Religion", NID: 1011163, EpisodeNumber: 2},
      {Title: "Mingling Its Own Nature With It", NID: 1011167, EpisodeNumber: 3},
      {Title: "Governed As It Were by Chance", NID: 1011165, EpisodeNumber: 4},
      {Title: "Ipsa Scientia Potestas Est", NID: 1011164, EpisodeNumber: 5},
      {Title: "To Hound Nature in Her Wanderings", NID: 1011176, EpisodeNumber: 6},
      {Title: "Knowledge of Causes, and Secret Motion of Things", NID: 1011166, EpisodeNumber: 7},
      {Title: "Variable and Full of Perturbation", NID: 1011170, EpisodeNumber: 8},
      {Title: "Things Which Have Never Yet Been Done", NID: 1011171, EpisodeNumber: 9},
      {Title: "By Means Which Have Never Yet Been Tried", NID: 1011178, EpisodeNumber: 10},
   }

   // Act: Call the specific method.
   actualMetadata, err := rootNode.ExtractEpisodes()
   if err != nil {
      t.Fatalf("ExtractEpisodes() returned an unexpected error: %v", err)
   }

   if len(actualMetadata) != len(expectedMetadata) {
      t.Fatalf("ExtractEpisodes() returned wrong number of items: got %d, want %d", len(actualMetadata), len(expectedMetadata))
   }

   for i := range actualMetadata {
      if *actualMetadata[i] != *expectedMetadata[i] {
         t.Errorf("Mismatch at index %d:\ngot:  %+v\nwant: %+v", i, *actualMetadata[i], *expectedMetadata[i])
      }
   }
}
