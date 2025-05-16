package main

import (
	"41.neocities.org/media/canal"
	"41.neocities.org/net"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func (f *flags) do_address() error {
	var fields canal.Fields
	err := fields.New(f.address)
	if err != nil {
		return err
	}
	fmt.Println(
		canal.AlgoliaConvertTracking, "=",
		fields.Get(canal.AlgoliaConvertTracking),
	)
	return nil
}

type flags struct {
	address  string
	asset    string
	dash     string
	e        net.License
	email    string
	media    string
	password string
	refresh  bool
	season   int64
}

func (f *flags) do_dash() error {
	data, err := os.ReadFile(f.media + "/canal/Play")
	if err != nil {
		return err
	}
	var play canal.Play
	err = play.Unmarshal(data)
	if err != nil {
		return err
	}
	f.e.Widevine = func(data []byte) ([]byte, error) {
		return play.Widevine(data)
	}
	return f.e.Download(f.media+"/Mpd", f.dash)
}

func (f *flags) do_asset() error {
	data, err := os.ReadFile(f.media + "/canal/Session")
	if err != nil {
		return err
	}
	var session canal.Session
	err = session.Unmarshal(data)
	if err != nil {
		return err
	}
	data, err = session.Play(f.asset)
	if err != nil {
		return err
	}
	var play canal.Play
	err = play.Unmarshal(data)
	if err != nil {
		return err
	}
	err = write_file(f.media+"/canal/Play", data)
	if err != nil {
		return err
	}
	resp, err := http.Get(play.Url)
	if err != nil {
		return err
	}
	return net.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_refresh() error {
	data, err := os.ReadFile(f.media + "/canal/Session")
	if err != nil {
		return err
	}
	var session canal.Session
	err = session.Unmarshal(data)
	if err != nil {
		return err
	}
	data, err = canal.NewSession(session.SsoToken)
	if err != nil {
		return err
	}
	return write_file(f.media+"/canal/Session", data)
}

func (f *flags) do_email() error {
	var ticket canal.Ticket
	err := ticket.New()
	if err != nil {
		return err
	}
	token, err := ticket.Token(f.email, f.password)
	if err != nil {
		return err
	}
	data, err := canal.NewSession(token.SsoToken)
	if err != nil {
		return err
	}
	return write_file(f.media+"/canal/Session", data)
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
	flag.StringVar(&f.address, "a", "", "address")
	flag.StringVar(&f.asset, "asset", "", "asset ID")
	flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
	flag.StringVar(&f.dash, "d", "", "dash ID")
	flag.StringVar(&f.email, "email", "", "email")
	flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
	flag.StringVar(&f.password, "password", "", "password")
	flag.BoolVar(&f.refresh, "r", false, "refresh")
	flag.Int64Var(&f.season, "s", 0, "season")
	flag.Parse()
	if f.email != "" {
		if f.password != "" {
			err = f.do_email()
		}
	} else if f.refresh {
		err = f.do_refresh()
	} else if f.address != "" {
		err = f.do_address()
	} else if f.asset != "" {
		if f.season >= 1 {
			err = f.do_season()
		} else {
			err = f.do_asset()
		}
	} else if f.dash != "" {
		err = f.do_dash()
	} else {
		flag.Usage()
	}
	if err != nil {
		panic(err)
	}
}
func (f *flags) do_season() error {
	data, err := os.ReadFile(f.media + "/canal/Session")
	if err != nil {
		return err
	}
	var session canal.Session
	err = session.Unmarshal(data)
	if err != nil {
		return err
	}
	assets, err := session.Assets(f.asset, f.season)
	if err != nil {
		return err
	}
	for i, asset := range assets {
		if i >= 1 {
			fmt.Println()
		}
		fmt.Println(&asset)
	}
	return nil
}
