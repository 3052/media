package main

import (
   "41.neocities.org/media/canal"
   "41.neocities.org/media/internal"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flags) do_address() error {
   data, err := os.ReadFile(f.media + "/canal/Session")
   if err != nil {
      return err
   }
   var session canal.Session
   err = session.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = canal.NewSession(session.SsoToken)
   if err != nil {
      return err
   }
   err = session.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.media+"/canal/Session", data)
   if err != nil {
      return err
   }
   var fields canal.Fields
   err = fields.New(f.address)
   if err != nil {
      return err
   }
   data, err = session.Play(fields.AlgoliaConvertTracking())
   if err != nil {
      return err
   }
   var play canal.Play
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.media+"/canal/Play", data)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.Url)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}
