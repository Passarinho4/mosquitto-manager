package internal

import (
	"bufio"
	"encoding/hex"
	"log"
	"os"
	"syscall"
)

func reloadConfig(config *Config) {
	process, err := os.FindProcess(config.mosquittoPid)
	if err != nil {
		log.Fatal(err)
	}
	err = process.Signal(syscall.SIGHUP)
	if err != nil {
		log.Fatal(err)
	}
}

func prepareConfigFile(client *ClientManager, config *Config) error {
	creds := client.getMosquittoCreds()
	pskfile, err := os.OpenFile(config.pskFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer pskfile.Close()
	w := bufio.NewWriter(pskfile)

	for _, cred := range creds {
		_, err = w.WriteString(cred.Login + ":" + hex.EncodeToString([]byte(cred.Password)) + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	return w.Flush()
}
