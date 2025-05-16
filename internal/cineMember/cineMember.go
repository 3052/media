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

func (f *flags) do_address() error {
	data, err := os.ReadFile(f.media + "/cineMember/User")
	if err != nil {
		return err
	}
	var user cineMember.User
	err = user.Unmarshal(data)
	if err != nil {
		return err
	}
	var address cineMember.Address
	err = address.Parse(f.address)
	if err != nil {
		return err
	}
	article, err := address.Article()
	if err != nil {
		return err
	}
	asset, ok := article.Film()
	if !ok {
		return errors.New(".Film()")
	}
	data, err = user.Play(article, asset)
	if err != nil {
		return err
	}
	var play cineMember.Play
	err = play.Unmarshal(data)
	if err != nil {
		return err
	}
	err = write_file(f.media+"/cineMember/Play", data)
	if err != nil {
		return err
	}
	title, ok := play.Dash()
	if !ok {
		return errors.New(".Dash()")
	}
	resp, err := http.Get(title.Manifest)
	if err != nil {
		return err
	}
	return net.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_email() error {
	data, err := cineMember.NewUser(f.email, f.password)
	if err != nil {
		return err
	}
	return write_file(f.media+"/cineMember/User", data)
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
	flag.IntVar(&net.ThreadCount, "t", 1, "thread count")
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
func (f *flags) do_dash() error {
	data, err := os.ReadFile(f.media + "/cineMember/Play")
	if err != nil {
		return err
	}
	var play cineMember.Play
	err = play.Unmarshal(data)
	if err != nil {
		return err
	}
	title, _ := play.Dash()
	f.e.Widevine = func(data []byte) ([]byte, error) {
		return title.Widevine(data)
	}
	return f.e.Download(f.media+"/Mpd", f.dash)
}

type flags struct {
	address  string
	dash     string
	e        net.License
	email    string
	media    string
	password string
}
