package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

var configDir string
var k2pdfopt string

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	configDir = filepath.Join(usr.HomeDir, ".config", "kindle-paper")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		panic(err)
	}
	k2pdfopt = "k2pdfopt"
}

func convert(in, out string) string {
	out = filepath.Join(filepath.Dir(in), out)
	log.Printf("Running k2pdfopt...")
	cmd := exec.Command(k2pdfopt, "-ui- -odpi 300 -dev kv -x", in, "-o", out)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println("Error while converting", "-", err)
	}
	return out + filepath.Ext(in)
}

func main() {
	var in, out string
	flag.StringVar(&in, "in", "", "input")
	flag.StringVar(&out, "out", "", "output")

	flag.Parse()

	mail, err := restoreMailSettings()
	if err != nil {
		panic(err)
	}

	converted := convert(in, out)
	fmt.Println("done")
	if err := sendToKindle(mail, converted); err != nil {
		log.Println(err)
	}
	if err := os.Remove(converted); err != nil {
		log.Println(err)
	}
}
