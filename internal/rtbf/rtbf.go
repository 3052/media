package main

import (
	"41.neocities.org/media/rtbf"
	"41.neocities.org/stream"
	"errors"
	"flag"
	"log"
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

func write_file(name string, data []byte) error {
	log.Println("WriteFile", name)
	return os.WriteFile(name, data, os.ModePerm)
}

type flags struct {
	dash     string
	e        stream.License
	email    string
	media    string
	password string
	address  string
}

func (f *flags) do_password() error {
	data, err := rtbf.NewLogin(f.email, f.password)
	if err != nil {
		return err
	}
	return write_file(f.media+"/rtbf/Login", data)
}

func main() {
	var f flags
	err := f.New()
	if err != nil {
		panic(err)
	}
	flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
	flag.StringVar(&f.email, "e", "", "email")
	flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
	flag.StringVar(&f.password, "p", "", "password")
	flag.StringVar(&f.dash, "i", "", "DASH ID")
	flag.StringVar(&f.address, "a", "", "address")
	flag.Parse()
	switch {
	case f.password != "":
		err := f.do_password()
		if err != nil {
			panic(err)
		}
	case f.address != "":
		err := f.do_address()
		if err != nil {
			panic(err)
		}
	case f.dash != "":
		err := f.do_dash()
		if err != nil {
			panic(err)
		}
	default:
		flag.Usage()
	}
}

func (f *flags) do_address() error {
	data, err := os.ReadFile(f.media + "/rtbf/Login")
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
	var address rtbf.Address
	address.New(f.address)
	content, err := address.Content()
	if err != nil {
		return err
	}
	data, err = gigya.Entitlement(content)
	if err != nil {
		return err
	}
	var title rtbf.Entitlement
	err = title.Unmarshal(data)
	if err != nil {
		return err
	}
	err = write_file(f.media+"/rtbf/Entitlement", data)
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
	return stream.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_dash() error {
	data, err := os.ReadFile(f.media + "/rtbf/Entitlement")
	if err != nil {
		return err
	}
	var title rtbf.Entitlement
	err = title.Unmarshal(data)
	if err != nil {
		return err
	}
	f.e.Widevine = func(data []byte) ([]byte, error) {
		return title.Widevine(data)
	}
	return f.e.Download(f.media+"/Mpd", f.dash)
}
