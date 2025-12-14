package main

import (
   "41.neocities.org/media/cineMember"
   "41.neocities.org/net"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4s" {
         return ""
      }
      return "L"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) New() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   
   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   flag.BoolVar(&c.vtt, "v", false, "VTT")
   flag.Parse()
   
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.vtt {
      return c.do_vtt()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
}

type command struct {
   name    string
   config   net.Config
   
   // 1
   email    string
   password string
   // 2
   address  string
   // 3
   vtt      bool
   // 4
   dash string
}

///

func (c *command) do_email_password() error {
   var session cineMember.Session
   err := session.Fetch()
   if err != nil {
      return err
   }
   err = session.Login(c.email, c.password)
   if err != nil {
      return err
   }
   return write_file(
      c.name+"/cineMember/Session", []byte(session.String()),
   )
}

func (c *command) do_address() error {
   data, err := os.ReadFile(c.name + "/cineMember/Session")
   if err != nil {
      return err
   }
   var session cineMember.Session
   err = session.Set(string(data))
   if err != nil {
      return err
   }
   id, err := cineMember.Id(c.address)
   if err != nil {
      return err
   }
   stream, err := session.Stream(id)
   if err != nil {
      return err
   }
   if c.vtt {
      return vtt(stream)
   }
   address, ok := stream.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(address)
   if err != nil {
      return err
   }
   return c.filters.Filter(resp, &c.config)
}

func vtt(stream *cineMember.Stream) error {
   address, ok := stream.Vtt()
   if !ok {
      return errors.New(".Vtt()")
   }
   resp, err := http.Get(address)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   file, err := os.Create(filepath.Base(address))
   if err != nil {
      return err
   }
   defer file.Close()
   _, err = file.ReadFrom(resp.Body)
   if err != nil {
      return err
   }
   return nil
}
