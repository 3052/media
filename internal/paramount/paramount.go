package main

import (
   "41.neocities.org/media/paramount"
   "41.neocities.org/net"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
   "strings"
)

func main() {
   log.SetFlags(log.Ltime)
   http.DefaultTransport = &http.Transport{
      Protocols: &http.Protocols{},
      Proxy: func(req *http.Request) (*url.URL, error) {
         log.Println(req.Method, req.URL)
         if strings.HasSuffix(req.URL.Path, "/anonymous-session-token.json") {
            return nil, nil
         }
         if strings.HasSuffix(req.URL.Path, "/getlicense") {
            return nil, nil
         }
         return http.ProxyFromEnvironment(req)
      },
   }
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.paramount != "" {
      err = set.do_paramount()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

func (f *flag_set) New() error {
   media, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   media = filepath.ToSlash(media) + "/media"
   f.cdm.ClientId = media + "/client_id.bin"
   f.cdm.PrivateKey = media + "/private_key.pem"
   flag.StringVar(&f.cdm.ClientId, "C", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.cdm.PrivateKey, "P", f.cdm.PrivateKey, "private key")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.BoolVar(&f.intl, "i", false, "intl")
   flag.StringVar(&f.paramount, "p", "", "paramount ID")
   flag.Parse()
   return nil
}

func (f *flag_set) do_paramount() error {
   // INTL does NOT allow anonymous key request, so if you are INTL you
   // will need to use US VPN until someone codes the INTL login
   secret := paramount.ComCbsApp
   at, err := secret.At()
   if err != nil {
      return err
   }
   session, err := at.Session(f.paramount)
   if err != nil {
      return err
   }
   f.cdm.License = func(data []byte) ([]byte, error) {
      return session.License(data)
   }
   if f.intl {
      secret = paramount.ComCbsCa
   }
   at, err = secret.At()
   if err != nil {
      return err
   }
   item, err := at.Item(f.paramount)
   if err != nil {
      return err
   }
   resp, err := item.Mpd()
   if err != nil {
      return err
   }
   return f.filters.Filter(resp, &f.cdm)
}
type flag_set struct {
   cdm       net.Cdm
   filters   net.Filters
   paramount string
   intl bool
}

