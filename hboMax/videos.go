package hboMax

type Videos struct {
   Errors   []Error
   Included []*Video
}

type Error struct {
   Detail  string // show was filtered by validator
   Message string // Token is missing or not valid
}

type Video struct {
   Attributes *struct {
      SeasonNumber  int
      EpisodeNumber int
      Name          string
      VideoType     string
   }
   Relationships *struct {
      Edit *struct {
         Data struct {
            Id string
         }
      }
   }
}

