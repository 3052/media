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

///

type command struct {
   // 1
   email string
   password string
   // 2
   representation string
   s internal.Stream
   home string
   peacock string
   v log.Level
}

func (c *command) New() error {
   flag.StringVar(&c.peacock, "b", "", "Peacock ID")
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.representation, "i", "", "representation")
   flag.StringVar(&c.password, "p", "", "password")
   flag.TextVar(&c.v.Level, "v", c.v.Level, "level")
   flag.StringVar(&c.s.ClientId, "c", c.s.ClientId, "client ID")
   flag.StringVar(&c.s.PrivateKey, "k", c.s.PrivateKey, "private key")
   flag.Parse()
   c.v.Set()
   log.Transport{}.Set()
   switch {
   case c.password != "":
      err := c.authenticate()
      if err != nil {
         panic(err)
      }
   case c.peacock != "":
      err := c.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
   var err error
   c.home, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   c.home = filepath.ToSlash(c.home)
   c.s.ClientId = c.home + "/widevine/client_id.bin"
   c.s.PrivateKey = c.home + "/widevine/private_key.pem"
   return nil
}

func (c command) download() error {
   text, err := os.ReadFile(c.home + "/peacock.json")
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
   media, err := c.s.DASH(req)
   if err != nil {
      return err
   }
   for _, medium := range media {
      if medium.ID == c.representation {
         var node peacock.QueryNode
         err := node.New(c.peacock)
         if err != nil {
            return err
         }
         c.s.Name = node
         c.s.Poster = video
         return c.s.Download(medium)
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

func (c command) authenticate() error {
   var sign peacock.SignIn
   err := sign.New(c.email, c.password)
   if err != nil {
      return err
   }
   text, err := sign.Marshal()
   if err != nil {
      return err
   }
   return os.WriteFile(c.home + "/peacock.json", text, 0666)
}
