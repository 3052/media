package main

import (
   "encoding/base64"
   "net/http"
   "net/url"
   "os"
)

func cache_hash() string {
   return base64.StdEncoding.EncodeToString([]byte("ff="))
}

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "gw.cds.amcn.com"
   req.URL.Path = "/content-compiler-cr/api/v1/content/amcn/amcplus/type/season-episodes/id/1010638"
   req.URL.Scheme = "https"
   
   req.Header["Authorization"] = []string{"Bearer eyJraWQiOiJwcm9kLTEiLCJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJlbnRpdGxlbWVudHMiOiJ1bmF1dGgiLCJhdWQiOiJyZXNvdXJjZV9zZXJ2ZXIiLCJhdXRoX3R5cGUiOiJiZWFyZXIiLCJmZWF0dXJlX2ZsYWdzIjoiZXlKaGJXTndiSFZ6TFcxMmNHUWlPbnNpWTJoaGNuUmxjaTF0ZG5BdFpXNWhZbXhsWkNJNlptRnNjMlY5TENKaGJXTndiSFZ6TFhScGVtVnVMV04zSWpwN0ltVnVZV0pzWldRaU9uUnlkV1Y5TENKallXSmZZVzFqY0MxamIyNTBaVzUwTFdkaGRHVmtMVzlqZEMweU1ESTBJam9pUWlJc0ltTmhZbDloYldOd0xXTnZiblJsYm5RdFpYaDBjbUZ6TFcxaGNpMHlNREkxSWpvaVFTSXNJbUZ0WTNCc2RYTXRZV1F0ZEdsbGNpSTZleUpoWkMxMGFXVnlMWEIxY21Ob1lYTmxMVzl1SWpwbVlXeHpaWDBzSW1GdFkzQnNkWE10YzJ0cGNDMXdjbTl0YnkxaFpITWlPbnNpYzJ0cGNDMXdjbTl0YnkxaFpITXRaVzVoWW14bFpDSTZabUZzYzJVc0luWmhjbWxoZEdsdmJpSTZJa0ZOUXlzZ1YybDBhQ0JCWkhNaWZTd2lZMjl0WTJGemRDMWhaQzFpYkc5amEyVnlMWE5qY21WbGJpSTZleUp6ZFdKMGFYUnNaU0k2SWtadmNpQm9aV3h3TENCbGJXRnBiQ0JqZFhOMGIyMWxjbk5sY25acFkyVkFZVzFqY0d4MWN5NWpiMjB1SWl3aWRHbDBiR1VpT2lKVWFHVWdUVzl1ZEdoc2VTQjNhWFJvSUVGa2N5QndiR0Z1SUdseklHNXZkQ0JqZFhKeVpXNTBiSGtnYzNWd2NHOXlkR1ZrSUc5dUlGaG1hVzVwZEhrZ1pHVjJhV05sY3k0aUxDSmxibUZpYkdWa0lqcDBjblZsZlN3aVlXMWpjR3gxY3kxMmFYcHBieTF3Y205dGIzUnBiMjRpT25zaVkyOTFjRzl1WDJOdlpHVWlPaUp5WlhSaGFXd3hJaXdpWTI5MWNHOXVYMk52WkdWZmJXOXVkR2hzZVNJNkluSmxkR0ZwYkRNaUxDSmxibUZpYkdWa0lqcDBjblZsZlgwPSIsInJvbGVzIjpbInVuYXV0aCJdLCJpc3MiOiJpcC0xMC0yLTQ3LTM4LmVjMi5pbnRlcm5hbCIsInRva2VuX3R5cGUiOiJhdXRoIiwiZXhwIjoxNzQ1NzY0NzM3LCJkZXZpY2UtaWQiOiJiNjhlNzE0MS1lZGU1LTQxMzctYWNkMS0xOThkNGRjMWJjODkiLCJieXBhc3NfcGVyc2lzdCI6dHJ1ZSwiaWF0IjoxNzQ1NzYxMTM3LCJqdGkiOiIwN2FmY2M3Mi02ODg0LTRmNmYtYWFmOC0zMzMyYThlNGJlNWUifQ.QmSZLBpVUj47DdD7kjpe1I9-OFzDwRi4fSoGVzFHO2vHvILzsG_wkNFg8eE-hVJ_hSMBVCf37_w_U8XHBViQcIIm7U_oV8at0BUtvXBl29aNsjT7kKX72JAGzNbA_UC2RxWPX31rw5Bl84496qowid7c-mlji-9xLZT4YqgD_UrnMADh9AmsnHSkoeGDyvZr5vSkDm4DQHH0ca-2kwwdBbe1W4GX4vJJjeme1Z02qL1oR-uWp9NPvenyc37cQPpRUvfQnJej1M66dn_icTtLY0yQnPDjchrogyj7VVbmZUplHa9W87d8PkQTjs9bsHyjqXW1_Czs8_cStuWbn0ACXw"}
   req.Header["X-Amcn-Cache-Hash"] = []string{cache_hash()}
   req.Header["X-Amcn-Network"] = []string{"amcplus"}
   req.Header["X-Amcn-Platform"] = []string{"web"}
   req.Header["X-Amcn-Tenant"] = []string{"amcn"}
   req.Header["X-Amcn-User-Cache-Hash"] = []string{cache_hash()}
   
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}
