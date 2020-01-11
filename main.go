package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type entry struct {
	Group      string
	Host       string
	KeepPwd    bool
	Login      string
	Note       string
	Password   string
	SecureMode bool
	Type       string
}

var (
	/*
		Program version, set at the build time
	*/
	version string
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("MikroTik Winbox 3 export file parser, ver. %s\n\n", version)
		fmt.Printf("Usage: %s [ Winbox 3 export file ]\n", os.Args[0])
		return
	}

	f, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	b := bytes.Split(f[4:], []byte{0x00, 0x00})

	var db []entry

	for _, v := range b[:len(b)-1] {
		x, err := parseEntry(v[2:])
		if err != nil {
			fmt.Println(err)
			continue
		}
		db = append(db, x)
	}

	j, err := json.Marshal(db)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}

func parseEntry(in []byte) (ret entry, err error) {
	var p int

	for p < len(in) {
		l := int(in[p])
		p++
		switch string(in[p : p+int(l)]) {
		case "group":
			p += int(l)
			e := bytes.IndexByte(in[p:], 0x00)
			if e == -1 {
				err = errors.New("Malformed input")
				return
			}
			ret.Group = string(in[p : p+e-1])
			p += e + 1
		case "host":
			p += int(l)
			e := bytes.IndexByte(in[p:], 0x00)
			if e == -1 {
				err = errors.New("Malformed input")
				return
			}
			ret.Host = string(in[p : p+e-1])
			p += e + 1
		case "keep-pwd":
			p += int(l)
			if in[p : p+1][0] == 0x01 {
				ret.KeepPwd = true
			}
			p += 3
		case "login":
			p += int(l)
			e := bytes.IndexByte(in[p:], 0x00)
			if e == -1 {
				err = errors.New("Malformed input")
				return
			}
			ret.Login = string(in[p : p+e-1])
			p += e + 1
		case "note":
			p += int(l)
			e := bytes.IndexByte(in[p:], 0x00)
			if e == -1 {
				err = errors.New("Malformed input")
				return
			}
			ret.Note = string(in[p : p+e-1])
			p += e + 1
		case "pwd":
			p += int(l)
			e := bytes.IndexByte(in[p:], 0x00)
			if e == -1 {
				err = errors.New("Malformed input")
				return
			}
			ret.Password = string(in[p : p+e-1])
			p += e + 1
		case "secure-mode":
			p += int(l)
			if in[p : p+1][0] == 0x01 {
				ret.SecureMode = true
			}
			p += 3
		case "type":
			p += int(l)
			ret.Type = string(in[p:])
			p = len(in)
		default:
			p = len(in) + 1
		}
	}

	return
}
