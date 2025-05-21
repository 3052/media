package canal

import (
   "io"
   "net/http"
   "strings"
)

const AlgoliaConvertTracking = "data-algolia-convert-tracking"

func (f *Fields) New(address string) error {
   resp, err := http.Get(address)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   *f = strings.FieldsFunc(string(data), func(r rune) bool {
      return strings.ContainsRune(` "=`, r)
   })
   return nil
}

type Fields []string

func (f Fields) Get(key string) string {
   var found bool
   for _, field := range f {
      switch {
      case field == key:
         found = true
      case found:
         return field
      }
   }
   return ""
}
