package main

import (
	"41.neocities.org/media/molotov"
	"41.neocities.org/stream"
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

func (f *flags) authenticate() error {
	var login molotov.Login
	err := login.New(f.email, f.password)
	if err != nil {
		return err
	}
	data, err := login.Auth.Refresh()
	if err != nil {
		return err
	}
	return write_file(f.media+"/molotov/Refresh", data)
}

func (f *flags) download() error {
	if f.dash != "" {
		data, err := os.ReadFile(f.media + "/molotov/Asset")
		if err != nil {
			return err
		}
		var asset molotov.Asset
		err = asset.Unmarshal(data)
		if err != nil {
			return err
		}
		f.e.Widevine = func(data []byte) ([]byte, error) {
			return asset.Widevine(data)
		}
		return f.e.Download(f.media+"/Mpd", f.dash)
	}
	data, err := os.ReadFile(f.media + "/molotov/Refresh")
	if err != nil {
		return err
	}
	var refresh molotov.Refresh
	err = refresh.Unmarshal(data)
	if err != nil {
		return err
	}
	data, err = refresh.Refresh()
	if err != nil {
		return err
	}
	err = write_file(f.media+"/molotov/Refresh", data)
	if err != nil {
		return err
	}
	view, err := refresh.View(&f.address)
	if err != nil {
		return err
	}
	data, err = refresh.Asset(view)
	if err != nil {
		return err
	}
	var asset molotov.Asset
	err = asset.Unmarshal(data)
	if err != nil {
		return err
	}
	err = write_file(f.media+"/molotov/Asset", data)
	if err != nil {
		return err
	}
	resp, err := http.Get(asset.FhdReady())
	if err != nil {
		return err
	}
	return stream.Mpd(f.media+"/Mpd", resp)
}

type flags struct {
	address  molotov.Address
	dash     string
	e        stream.License
	email    string
	media    string
	password string
}

func main() {
	var f flags
	err := f.New()
	if err != nil {
		panic(err)
	}
	flag.Var(&f.address, "a", "address")
	flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
	flag.StringVar(&f.email, "e", "", "email")
	flag.StringVar(&f.dash, "i", "", "dash ID")
	flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
	flag.StringVar(&f.password, "p", "", "password")
	flag.IntVar(&stream.ThreadCount, "t", 1, "thread count")
	flag.Parse()
	switch {
	case f.password != "":
		err := f.authenticate()
		if err != nil {
			panic(err)
		}
	case f.address.String() != "":
		err := f.download()
		if err != nil {
			panic(err)
		}
	default:
		flag.Usage()
	}
}
