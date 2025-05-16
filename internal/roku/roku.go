package main

import (
	"41.neocities.org/media/roku"
	"41.neocities.org/stream"
	"flag"
	"fmt"
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

func main() {
	var f flags
	err := f.New()
	if err != nil {
		panic(err)
	}
	flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
	flag.StringVar(&f.dash, "d", "", "dash ID")
	flag.BoolVar(&f.code_write, "code", false, "write code")
	flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
	flag.StringVar(&f.roku, "r", "", "Roku ID")
	flag.BoolVar(&f.token_read, "t", false, "read token")
	flag.BoolVar(&f.token_write, "token", false, "write token")
	flag.Parse()
	switch {
	case f.code_write:
		err := f.do_code()
		if err != nil {
			panic(err)
		}
	case f.token_write:
		err := f.do_token()
		if err != nil {
			panic(err)
		}
	case f.roku != "":
		err := f.do_roku()
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

func (f *flags) do_code() error {
	data, err := (*roku.Code).AccountToken(nil)
	if err != nil {
		return err
	}
	var token roku.AccountToken
	err = token.Unmarshal(data)
	if err != nil {
		return err
	}
	err = write_file(f.media+"/roku/AccountToken", data)
	if err != nil {
		return err
	}
	data1, err := token.Activation()
	if err != nil {
		return err
	}
	var activation roku.Activation
	err = activation.Unmarshal(data1)
	if err != nil {
		return err
	}
	fmt.Println(&activation)
	return write_file(f.media+"/roku/Activation", data1)
}

func (f *flags) do_token() error {
	data, err := os.ReadFile(f.media + "/roku/AccountToken")
	if err != nil {
		return err
	}
	var token roku.AccountToken
	err = token.Unmarshal(data)
	if err != nil {
		return err
	}
	data, err = os.ReadFile(f.media + "/roku/Activation")
	if err != nil {
		return err
	}
	var activation roku.Activation
	err = activation.Unmarshal(data)
	if err != nil {
		return err
	}
	data, err = token.Code(&activation)
	if err != nil {
		return err
	}
	return write_file(f.media+"/roku/Code", data)
}

type flags struct {
	e          stream.License
	media      string
	token_read bool

	code_write  bool
	token_write bool
	roku        string
	dash        string
}

func (f *flags) do_roku() error {
	var code *roku.Code
	if f.token_read {
		data, err := os.ReadFile(f.media + "/roku/Code")
		if err != nil {
			return err
		}
		code = &roku.Code{}
		err = code.Unmarshal(data)
		if err != nil {
			return err
		}
	}
	data, err := code.AccountToken()
	if err != nil {
		return err
	}
	var token roku.AccountToken
	err = token.Unmarshal(data)
	if err != nil {
		return err
	}
	data1, err := token.Playback(f.roku)
	if err != nil {
		return err
	}
	err = write_file(f.media+"/roku/Playback", data1)
	if err != nil {
		return err
	}
	var play roku.Playback
	err = play.Unmarshal(data1)
	if err != nil {
		return err
	}
	resp, err := http.Get(play.Url)
	if err != nil {
		return err
	}
	return stream.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_dash() error {
	data, err := os.ReadFile(f.media + "/roku/Playback")
	if err != nil {
		return err
	}
	var play roku.Playback
	err = play.Unmarshal(data)
	if err != nil {
		return err
	}
	f.e.Widevine = func(data []byte) ([]byte, error) {
		return play.Widevine(data)
	}
	return f.e.Download(f.media+"/Mpd", f.dash)
}
