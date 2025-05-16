package main

import (
	"41.neocities.org/media/nbc"
	"41.neocities.org/stream"
	"flag"
	"net/http"
	"os"
	"path/filepath"
)

func (f *flags) New() error {
	var err error
	f.media, err = os.UserHomeDir()
	if err != nil {
		return err
	}
	f.media = filepath.ToSlash(f.media) + "/media"
	f.e.ClientId = f.media + "/client_id.bin"
	f.e.PrivateKey = f.media + "/private_key.pem"
	return nil
}

type flags struct {
	dash  string
	e     stream.License
	media string
	nbc   int
}

func main() {
	var f flags
	err := f.New()
	if err != nil {
		panic(err)
	}
	flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
	flag.StringVar(&f.dash, "d", "", "dash ID")
	flag.IntVar(&f.nbc, "n", 0, "NBC ID")
	flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
	flag.IntVar(&stream.ThreadCount, "t", 1, "thread count")
	flag.Parse()
	switch {
	case f.nbc >= 1:
		err = f.do_nbc()
	case f.dash != "":
		err = f.do_dash()
	default:
		flag.Usage()
	}
	if err != nil {
		panic(err)
	}
}

func (f *flags) do_nbc() error {
	var metadata nbc.Metadata
	err := metadata.New(f.nbc)
	if err != nil {
		return err
	}
	vod, err := metadata.Vod()
	if err != nil {
		return err
	}
	resp, err := http.Get(vod.PlaybackUrl)
	if err != nil {
		return err
	}
	return stream.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_dash() error {
	f.e.Widevine = nbc.Widevine
	return f.e.Download(f.media+"/Mpd", f.dash)
}
