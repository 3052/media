package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "flag"
   "fmt"
   "log"
   "os"
   "path/filepath"
)

func (p *program) run() error {
   cache_dir, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache_dir = filepath.ToSlash(cache_dir)
   p.cache_file = cache_dir + "/rosso/hboMax.xml"
   p.job.CertificateChain = cache_dir + "/SL3000/CertificateChain"
   p.job.EncryptSignKey = cache_dir + "/SL3000/EncryptSignKey"
   // 1
   flag.StringVar(&p.proxy, "x", "", "proxy")
   // 2
   flag.BoolVar(&p.initiate, "i", false, "device initiate")
   flag.StringVar(
      &p.market, "m", hboMax.Markets[0], fmt.Sprint(hboMax.Markets),
   )
   // 3
   flag.BoolVar(&p.login, "l", false, "device login")
   // 4
   flag.StringVar(&p.address, "a", "", "address")
   flag.IntVar(&p.season, "s", 0, "season")
   // 5
   flag.StringVar(&p.edit, "e", "", "edit ID")
   // 6
   flag.StringVar(&p.dash, "d", "", "DASH ID")
   flag.StringVar(&p.job.CertificateChain, "C", p.job.CertificateChain, "certificate chain")
   flag.StringVar(&p.job.EncryptSignKey, "E", p.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   err = p.run_proxy()
   if err != nil {
      return err
   }
   if p.initiate {
      return p.run_initiate()
   }
   if p.login {
      return p.run_login()
   }
   if p.address != "" {
      return p.run_address()
   }
   if p.edit != "" {
      return p.run_edit()
   }
   if p.dash != "" {
      return p.run_dash()
   }
   return maya.Usage([][]string{
      {"x"},
      {"i", "m"},
      {"l"},
      {"a", "s"},
      {"e"},
      {"d", "C", "E"},
   })
}
func main() {
   err := new(program).run()
   if err != nil {
      log.Fatal(err)
   }
}

type program struct {
   cache_file string
   // 1
   proxy string
   // 2
   initiate bool
   market   string
   // 3
   login bool
   // 4
   address string
   season  int
   // 5
   edit string
   // 6
   dash string
   job  maya.PlayReadyJob
}
