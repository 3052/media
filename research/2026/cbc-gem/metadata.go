package cbc

import (
   "strconv"
   "strings"
)

type GemCatalog struct {
   Content []struct {
      Lineups []struct {
         Items []LineupItem
      }
   }
   SelectedUrl string
   StructuredMetadata Metadata
}

type Metadata struct {
   PartOfSeries *struct {
      Name string // The Fall
   }
   PartOfSeason *struct {
      SeasonNumber int
   }
   EpisodeNumber int
   Name string
   DateCreated string // 2014-01-01T00:00:00
}

func (Metadata) Owner() (string, bool) {
   return "", false
}

func (m Metadata) Title() (string, bool) {
   return m.Name, true
}

func (m Metadata) Episode() (string, bool) {
   if m.EpisodeNumber >= 1 {
      return strconv.Itoa(m.EpisodeNumber), true
   }
   return "", false
}

func (m Metadata) Season() (string, bool) {
   if p := m.PartOfSeason; p != nil {
      return strconv.Itoa(p.SeasonNumber), true
   }
   return "", false
}

func (m Metadata) Show() (string, bool) {
   if m.PartOfSeries != nil {
      return m.PartOfSeries.Name, true
   }
   return "", false
}

func (m Metadata) Year() (string, bool) {
   if m.PartOfSeries != nil {
      return "", false
   }
   year, _, _ := strings.Cut(m.DateCreated, "-")
   return year, true
}
