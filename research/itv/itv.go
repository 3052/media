package itv

import (
   "bytes"
   "encoding/json"
   "errors"
)

type next_data struct {
   Page string `json:"page"`
}

// ExtractFromHTML searches the provided byte slice for the __NEXT_DATA__ script
// and unmarshals the JSON content directly into the receiver.
func (n *next_data) ExtractFromHTML(htmlContent []byte) error {
   var (
      startToken = []byte(`<script id="__NEXT_DATA__" type="application/json">`)
      endToken   = []byte(`</script>`)
   )
   // Find start of the tag
   _, after, found := bytes.Cut(htmlContent, startToken)
   if !found {
      return errors.New("__NEXT_DATA__ tag not found")
   }
   // Find end of the tag
   jsonData, _, found := bytes.Cut(after, endToken)
   if !found {
      return errors.New("closing script tag not found")
   }
   // Unmarshal directly into the receiver 'n'
   return json.Unmarshal(jsonData, n)
}
