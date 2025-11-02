package amc

import (
   "os" // CHANGED: Replaced "io/ioutil" with "os"
   "reflect"
   "testing"
)

func TestExtractSeasons(t *testing.T) {
   // Arrange: Use the modern os.ReadFile function.
   jsonData, err := os.ReadFile("series-detail.json") // CHANGED: ioutil.ReadFile -> os.ReadFile
   if err != nil {
      t.Fatalf("Failed to read test data file: %v", err)
   }

   // Arrange: Define the exact output we expect to get.
   expectedSeasons := []Season{
      {Title: "Season 1", NID: 1010634},
      {Title: "Season 2", NID: 1010638},
      {Title: "Season 3", NID: 1010635},
      {Title: "Season 4", NID: 1010643},
      {Title: "Season 5", NID: 1010637},
   }

   // Act: Run the function we are testing.
   actualSeasons, err := ExtractSeasons(jsonData)

   // Assert: Check the results.
   if err != nil {
      t.Fatalf("ExtractSeasons() returned an unexpected error: %v", err)
   }

   if !reflect.DeepEqual(actualSeasons, expectedSeasons) {
      t.Errorf("ExtractSeasons() got = %v, want %v", actualSeasons, expectedSeasons)
   }
}
