// https://hulu.com/movie/05e76ad8-c3dd-4c3e-bab9-df3cf71c6871
// https://hulu.com/movie/alien-romulus-05e76ad8-c3dd-4c3e-bab9-df3cf71c6871
package hulu

import "path"

func Id(url string) string {
   part := path.Base(url)
   len_part := len(part)
   const len_uuid = 36
   if len_part > len_uuid {
      if part[len_part-len_uuid-1] == '-' {
         return part[len_part-len_uuid:]
      }
   }
   return part
}
