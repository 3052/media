package main

import (
   "errors"
   "flag"
   "fmt"
   "net/http"
   "os"
   "path/filepath"
)

func main() {
   maya.Transport(func(*http.Request) string {
      return "L"
   })
   log.SetFlags(log.Ltime)
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/peacock/userCache.xml"
   // 1
   flag.StringVar(&c.email, "email", "", "email")
   flag.StringVar(&c.password, "password", "", "password")
   // 2
   flag.StringVar(&c.peacock, "p", "", "Peacock ID")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.peacock != "" {
      return c.do_peacock()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"email", "password"},
      {"p"},
      {"d", "C", "P"},
   })
}

type command struct {
   name string
   
   // 1
   email string
   password string
   // 2
   peacock string
   // 3
   dash string
   job maya.WidevineJob
}

///

func (c command) do_email_password() error {
   var sign peacock.SignIn
   err := sign.New(c.email, c.password)
   if err != nil {
      return err
   }
   text, err := sign.Marshal()
   if err != nil {
      return err
   }
   return os.WriteFile(c.name + "/peacock.json", text, 0666)
}

func (c command) do_peacock() error {
   text, err := os.ReadFile(c.name + "/peacock.json")
   if err != nil {
      return err
   }
   var sign peacock.SignIn
   sign.Unmarshal(text)
   auth, err := sign.Auth()
   if err != nil {
      return err
   }
   video, err := auth.Video(c.peacock)
   if err != nil {
      return err
   }
   akamai, ok := video.Akamai()
   if !ok {
      return errors.New("peacock.VideoPlayout.Akamai")
   }
   req, err := http.NewRequest("", akamai, nil)
   if err != nil {
      return err
   }
   media, err := c.job.DASH(req)
   if err != nil {
      return err
   }
   for _, medium := range media {
      if medium.ID == c.dash {
         var node peacock.QueryNode
         err := node.New(c.peacock)
         if err != nil {
            return err
         }
         c.job.Name = node
         c.job.Poster = video
         return c.job.Download(medium)
      }
   }
   // 2 MPD all
   for i, medium := range media {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(medium)
   }
   return nil
}
