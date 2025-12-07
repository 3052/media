package rtbf

import (
   "encoding/json"
   "errors"
   "net/http"
   "strings"
)

func (g *GigyaLogin) Entitlement(assetId string) (*Entitlement, error) {
   req, _ := http.NewRequest("", "https://exposure.api.redbee.live", nil)
   req.URL.Path = func() string {
      var data strings.Builder
      data.WriteString("/v2/customer/RTBF/businessunit/Auvio/entitlement/")
      data.WriteString(assetId)
      data.WriteString("/play")
      return data.String()
   }()
   req.Header.Set("x-forwarded-for", "91.90.123.17")
   req.Header.Set("authorization", "Bearer "+g.SessionToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   title := &Entitlement{}
   err = json.NewDecoder(resp.Body).Decode(title)
   if err != nil {
      return nil, err
   }
   if title.Message != "" {
      return nil, errors.New(title.Message)
   }
   return title, nil
}
