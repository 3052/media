package canal

import (
   "io"
   "net/http"
   "strings"
)

func (f Fields) AssetId() string {
   var key, value string
   for _, field := range f {
      switch key {
      case "data-algolia-convert-tracking":
         value = field
      case "/web/signup/":
         return value
      }
      key = field
   }
   return ""
}

type Fields []string

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
