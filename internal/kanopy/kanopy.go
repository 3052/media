package main

import (
	"41.neocities.org/media/kanopy"
	"41.neocities.org/net"
	"errors"
	"flag"
	"log"
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

func main() {
	var f flags
	err := f.New()
	if err != nil {
		panic(err)
	}
	flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
	flag.StringVar(&f.dash, "d", "", "DASH ID")
	flag.StringVar(&f.email, "email", "", "email")
	flag.IntVar(&f.kanopy, "k", 0, "Kanopy ID")
	flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
	flag.StringVar(&f.password, "password", "", "password")
	flag.Parse()
	if f.email != "" {
		if f.password != "" {
			err = f.do_email()
		}
	} else if f.kanopy >= 1 {
		err = f.do_kanopy()
	} else if f.dash != "" {
		err = f.do_dash()
	} else {
		flag.Usage()
	}
	if err != nil {
		panic(err)
	}
}

func (f *flags) do_email() error {
	data, err := kanopy.NewLogin(f.email, f.password)
	if err != nil {
		return err
	}
	return write_file(f.media+"/kanopy/Login", data)
}

type flags struct {
	dash     string
	e        net.License
	email    string
	kanopy   int
	media    string
	password string
}

func (f *flags) do_kanopy() error {
	data, err := os.ReadFile(f.media + "/kanopy/Login")
	if err != nil {
		return err
	}
	var login kanopy.Login
	err = login.Unmarshal(data)
	if err != nil {
		return err
	}
	member, err := login.Membership()
	if err != nil {
		return err
	}
	data, err = login.Plays(member, f.kanopy)
	if err != nil {
		return err
	}
	var plays kanopy.Plays
	err = plays.Unmarshal(data)
	if err != nil {
		return err
	}
	err = write_file(f.media+"/kanopy/Plays", data)
	if err != nil {
		return err
	}
	manifest, ok := plays.Dash()
	if !ok {
		return errors.New(".Dash()")
	}
	resp, err := manifest.Mpd()
	if err != nil {
		return err
	}
	return net.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_dash() error {
	data, err := os.ReadFile(f.media + "/kanopy/Login")
	if err != nil {
		return err
	}
	var login kanopy.Login
	err = login.Unmarshal(data)
	if err != nil {
		return err
	}
	data, err = os.ReadFile(f.media + "/kanopy/Plays")
	if err != nil {
		return err
	}
	var plays kanopy.Plays
	err = plays.Unmarshal(data)
	if err != nil {
		return err
	}
	manifest, _ := plays.Dash()
	f.e.Widevine = func(data []byte) ([]byte, error) {
		return login.Widevine(manifest, data)
	}
	return f.e.Download(f.media+"/Mpd", f.dash)
}
