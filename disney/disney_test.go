package disney

import "testing"

var entity_tests = []struct {
   disney    string
   format    string
   justWatch string
   country   []string
}{
   {
      format:    "HD",
      disney:    "https://disneyplus.com/browse/entity-b7751be5-e33f-4317-89d3-2a2f116ea36a",
      justWatch: "https://justwatch.com/us/movie/the-life-aquatic-with-steve-zissou",
      country: []string{
         "Argentina", "Australia", "Austria", "Belgium", "Bolivia", "Brazil",
         "Bulgaria", "Canada", "Chile", "Colombia", "Costa Rica", "Croatia",
         "Czech Republic", "Denmark", "Ecuador", "Finland", "France", "Germany",
         "Greece", "Guatemala", "Hungary", "Ireland", "Italy", "Japan", "Mexico",
         "Netherlands", "New Zealand", "Norway", "Peru", "Poland", "Portugal",
         "Romania", "Singapore", "Slovakia", "South Korea", "Spain", "Sweden",
         "Switzerland", "Taiwan", "Turkey", "United Kingdom", "Venezuela",
      },
   },
   {
      country: []string{
         "Argentina", "Australia", "Austria", "Belgium", "Bolivia", "Brazil",
         "Bulgaria", "Canada", "Chile", "Colombia", "Costa Rica", "Croatia",
         "Czech Republic", "Denmark", "Ecuador", "Finland", "France", "Germany",
         "Greece", "Guatemala", "Hungary", "Ireland", "Italy", "Japan", "Mexico",
         "Netherlands", "New Zealand", "Norway", "Peru", "Poland", "Portugal",
         "Romania", "Singapore", "Slovakia", "South Korea", "Spain", "Sweden",
         "Switzerland", "Taiwan", "Turkey", "United Kingdom", "United States",
         "Venezuela",
      },
      justWatch: "https://justwatch.com/us/tv-show/paradise-2025",
      disney:    "https://disneyplus.com/play/26fb7558-a0bb-4153-9dde-18800e1738a9",
      format:    "4K ULTRA HD",
   },
   {
      justWatch: "https://justwatch.com/us/movie/the-muppet-show",
      country: []string{
         "Argentina", "Australia", "Austria", "Belgium", "Bolivia", "Brazil",
         "Bulgaria", "Canada", "Chile", "Colombia", "Costa Rica", "Croatia",
         "Czech Republic", "Denmark", "Ecuador", "Finland", "France", "Germany",
         "Greece", "Guatemala", "Hungary", "Ireland", "Italy", "Japan", "Mexico",
         "Netherlands", "New Zealand", "Norway", "Peru", "Poland", "Portugal",
         "Romania", "Singapore", "Slovakia", "South Korea", "Spain", "Sweden",
         "Switzerland", "Taiwan", "Turkey", "United Kingdom", "United States",
         "Venezuela",
      },
      disney: "https://disneyplus.com/browse/entity-e7487044-477c-4a57-ae45-afc689f2a346",
      format: "4K ULTRA HD",
   },
}

func TestEntity(t *testing.T) {
   t.Log(entity_tests)
}
