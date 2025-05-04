package main

import (
   "41.neocities.org/media/cineMember"
   "41.neocities.org/media/internal"
   "41.neocities.org/platform/mullvad"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}
