package canal

import (
   "io"
   "net/http"
   "strings"
)

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
