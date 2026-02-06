package disney

import "strings"

func (p *Page) String() string {
   var data strings.Builder
   if len(p.Containers[0].Seasons) >= 1 {
      var line bool
      for _, seasonItem := range p.Containers[0].Seasons {
         if line {
            data.WriteString("\n\n")
         } else {
            line = true
         }
         data.WriteString("name = ")
         data.WriteString(seasonItem.Visuals.Name)
         data.WriteString("\nid = ")
         data.WriteString(seasonItem.Id)
      }
   } else {
      data.WriteString("title = ")
      data.WriteString(p.Actions[0].InternalTitle)
   }
   return data.String()
}

type Page struct {
   Actions []struct {
      InternalTitle string // movie
   }
   Containers []struct {
      Seasons []struct { // series
         Visuals struct {
            Name string
         }
         Id string
      }
   }
}
