package main

import (
	"41.neocities.org/media/hulu"
	"41.neocities.org/net"
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

func main() {
	var f flags
	err := f.New()
	if err != nil {
		panic(err)
	}
	flag.StringVar(&f.address, "a", "", "address")
	flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
	flag.StringVar(&f.dash, "d", "", "DASH ID")
	flag.StringVar(&f.email, "email", "", "email")
	flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
	flag.StringVar(&f.password, "password", "", "password")
	flag.Parse()
	if f.email != "" {
		if f.password != "" {
			err = f.do_email()
		}
	} else if f.address != "" {
		err = f.do_address()
	} else if f.dash != "" {
		err = f.do_dash()
	} else {
		flag.Usage()
	}
	if err != nil {
		panic(err)
	}
}

type flags struct {
	address  string
	dash     string
	e        net.License
	email    string
	media    string
	password string
}

func (f *flags) do_email() error {
	data, err := hulu.NewAuthenticate(f.email, f.password)
	if err != nil {
		return err
	}
	return write_file(f.media+"/hulu/Authenticate", data)
}

func (f *flags) do_address() error {
	data, err := os.ReadFile(f.media + "/hulu/Authenticate")
	if err != nil {
		return err
	}
	var auth hulu.Authenticate
	err = auth.Unmarshal(data)
	if err != nil {
		return err
	}
	err = auth.Refresh()
	if err != nil {
		return err
	}
	deep, err := auth.DeepLink(hulu.Id(f.address))
	if err != nil {
		return err
	}
	data, err = auth.Playlist(deep)
	if err != nil {
		return err
	}
	err = write_file(f.media+"/hulu/Playlist", data)
	if err != nil {
		return err
	}
	var play hulu.Playlist
	err = play.Unmarshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Get(play.StreamUrl)
	if err != nil {
		return err
	}
	return net.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_dash() error {
	data, err := os.ReadFile(f.media + "/hulu/Playlist")
	if err != nil {
		return err
	}
	var play hulu.Playlist
	err = play.Unmarshal(data)
	if err != nil {
		return err
	}
	f.e.Widevine = func(data []byte) ([]byte, error) {
		return play.Widevine(data)
	}
	return f.e.Download(f.media+"/Mpd", f.dash)
}
