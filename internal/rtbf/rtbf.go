package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/rtbf"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(*http.Request) string {
      return "L"
   })
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
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/rtbf/user_cache.json"

   flag.StringVar(&c.config.ClientId, "C", c.config.ClientId, "client ID")
   flag.StringVar(&c.config.PrivateKey, "P", c.config.PrivateKey, "private key")
   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   flag.Parse()

   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   flag.Usage()
   return nil
}

type user_cache struct {
   Account rtbf.Account
   Mpd     *url.URL
   MpdBody []byte
}

func write(name string, cache *user_cache) error {
   data, err := json.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (c *command) do_email_password() error {
   var account rtbf.Account
   err = account.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return write(c.name, &user_cache{Account: account})
}

type command struct {
   config maya.Config
   name   string
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash string
}

///

func (c *command) do_address() error {
   data, err := os.ReadFile(c.name + "/rtbf/Login")
   if err != nil {
      return err
   }
   var login rtbf.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   jwt, err := login.Jwt()
   if err != nil {
      return err
   }
   gigya, err := jwt.Login()
   if err != nil {
      return err
   }
   path, err := rtbf.GetPath(c.address)
   if err != nil {
      return err
   }
   asset_id, err := rtbf.FetchAssetId(path)
   if err != nil {
      return err
   }
   data, err = gigya.Entitlement(asset_id)
   if err != nil {
      return err
   }
   var title rtbf.Entitlement
   err = title.Unmarshal(data)
   if err != nil {
      return err
   }
   format, ok := title.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(format.MediaLocator)
   if err != nil {
      return err
   }
}

func (c *command) do_dash() error {
   c.config.Send = func(data []byte) ([]byte, error) {
      return title.Widevine(data)
   }
   return c.filters.Filter(resp, &c.config)
}
