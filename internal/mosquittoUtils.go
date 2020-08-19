package internal

import (
	"bufio"
	"encoding/hex"
	"log"
	"os"
)

func reloadConfig(config *Config) {
	//process, err := os.FindProcess(config.mosquittoPid)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = process.Signal(syscall.SIGHUP)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func preparePskFile(client Manager, config *Config) error {
	creds := client.GetAll()
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

	log.Printf("Write %d lines to pskfile", len(creds))

	return w.Flush()
}

func prepareAclFile(client Manager, config *Config) error {
	creds := client.GetAll()
	aclfile, err := os.OpenFile(config.aclFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer aclfile.Close()
	w := bufio.NewWriter(aclfile)

	for _, cred := range creds {
		if cred.Acls == nil {
			break
		}
		_, err = w.WriteString("user " + cred.Login + "\n")
		if err != nil {
			log.Fatal(err)
		}
		for _, acl := range cred.Acls {
			_, err = w.WriteString(acl.AclType + " " + acl.AccessType + " " + acl.Topic + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
		_, err = w.WriteString("\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	return w.Flush()
}
