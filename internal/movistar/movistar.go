package main

import (
	"41.neocities.org/media/movistar"
	"41.neocities.org/net"
	"flag"
	"log"
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

type flags struct {
	dash     string
	e        net.License
	email    string
	media    string
	movistar int64
	password string
}

func main() {
	var f flags
	err := f.New()
	if err != nil {
		panic(err)
	}
	flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
	flag.StringVar(&f.dash, "d", "", "dash ID")
	flag.StringVar(&f.email, "email", "", "email")
	flag.Int64Var(&f.movistar, "m", 0, "movistar ID")
	flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
	flag.StringVar(&f.password, "password", "", "password")
	flag.IntVar(&net.ThreadCount, "t", 1, "thread count")
	flag.Parse()
	if f.email != "" {
		if f.password != "" {
			err := f.do_email()
			if err != nil {
				panic(err)
			}
		}
	} else if f.movistar >= 1 {
		err := f.do_movistar()
		if err != nil {
			panic(err)
		}
	} else if f.dash != "" {
		err := f.do_dash()
		if err != nil {
			panic(err)
		}
	} else {
		flag.Usage()
	}
}

func (f *flags) do_email() error {
	data, err := movistar.NewToken(f.email, f.password)
	if err != nil {
		return err
	}
	var token movistar.Token
	err = token.Unmarshal(data)
	if err != nil {
		return err
	}
	err = write_file(f.media+"/movistar/Token", data)
	if err != nil {
		return err
	}
	oferta, err := token.Oferta()
	if err != nil {
		return err
	}
	data1, err := token.Device(oferta)
	if err != nil {
		return err
	}
	return write_file(f.media+"/movistar/Device", data1)
}

func (f *flags) do_movistar() error {
	data, err := movistar.NewDetails(f.movistar)
	if err != nil {
		return err
	}
	var details movistar.Details
	err = details.Unmarshal(data)
	if err != nil {
		return err
	}
	err = write_file(f.media+"/movistar/Details", data)
	if err != nil {
		return err
	}
	resp, err := details.Mpd()
	if err != nil {
		return err
	}
	return net.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_dash() error {
	data, err := os.ReadFile(f.media + "/movistar/Token")
	if err != nil {
		return err
	}
	var token movistar.Token
	err = token.Unmarshal(data)
	if err != nil {
		return err
	}
	data, err = os.ReadFile(f.media + "/movistar/Device")
	if err != nil {
		return err
	}
	var device movistar.Device
	err = device.Unmarshal(data)
	if err != nil {
		return err
	}
	oferta, err := token.Oferta()
	if err != nil {
		return err
	}
	init1, err := oferta.InitData(device)
	if err != nil {
		return err
	}
	data, err = os.ReadFile(f.media + "/movistar/Details")
	if err != nil {
		return err
	}
	var details movistar.Details
	err = details.Unmarshal(data)
	if err != nil {
		return err
	}
	session, err := device.Session(init1, &details)
	if err != nil {
		return err
	}
	f.e.Widevine = func(data []byte) ([]byte, error) {
		return session.Widevine(data)
	}
	return f.e.Download(f.media+"/Mpd", f.dash)
}
