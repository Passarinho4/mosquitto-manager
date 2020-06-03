package internal

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/mitchellh/go-ps"
	"log"
	"net/http"
	"strconv"
)

type LoginPassword struct {
	Login    string
	Password string
}

type ManagerService struct {
	client *ClientManager
	config Config
}

func NewManagerService() ManagerService {
	var service = ManagerService{}
	var kubeconfig *string
	var mosquittoPid *int
	var pskFilePath *string
	var basicAuthLogin *string
	var basicAuthPass *string
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	mosquittoPid = flag.Int("mosquittoPid", 0, "pid of mosquitto process")
	pskFilePath = flag.String("pskFilePath", "", "path to pskfile")
	basicAuthLogin = flag.String("basicAuthLogin", "", "basic auth login")
	basicAuthPass = flag.String("basicAuthPass", "", "basic auth password")
	if *mosquittoPid == 0 {
		mosquittoPid = tryFindMosquittoPidByName()
	}
	if *pskFilePath == "" {
		*pskFilePath = "/proc/" + strconv.Itoa(*mosquittoPid) + "/root/etc/mosquitto/pskfile"
	}
	log.Printf("Mosquitto pid - " + strconv.Itoa(*mosquittoPid))
	log.Printf("Pskfile path - " + *pskFilePath)
	flag.Parse()
	var client = NewClientManager(kubeconfig)
	var config = NewConfig(*mosquittoPid, *pskFilePath, *basicAuthLogin, *basicAuthPass)
	service.config = config
	service.client = client
	return service
}

func tryFindMosquittoPidByName() *int {
	p, err := ps.Processes()
	if err != nil {
		log.Fatal(err)
	}
	pid := 0
	for _, process := range p {
		if process.Executable() == "mosquitto" {
			pid = process.Pid()
		}
	}
	return &pid
}

func (service *ManagerService) checkAuth(w http.ResponseWriter, r *http.Request) error {
	if service.config.basicAuthHeader == "" || r.Header.Get("Authorization") == service.config.basicAuthHeader {
		return nil
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return errors.New("unauthorized")
	}
}

func (service *ManagerService) add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := service.checkAuth(w, r)
		if err != nil {
			return
		}
		var lp LoginPassword
		err = json.NewDecoder(r.Body).Decode(&lp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Adding Login=" + lp.Login + " Password=" + lp.Password)
		err = service.client.createMosquittoCred(lp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Added Login=" + lp.Login + " Password=" + lp.Password)
		service.reloadAfterChange()
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (service *ManagerService) remove(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := service.checkAuth(w, r)
		if err != nil {
			return
		}
		var lp LoginPassword
		err = json.NewDecoder(r.Body).Decode(&lp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Removing Login=" + lp.Login + " Password=" + lp.Password)
		removeCRD(lp)
		service.reloadAfterChange()
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (service *ManagerService) list(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := service.checkAuth(w, r)
		if err != nil {
			return
		}
		crds := service.client.getMosquittoCreds()
		js, err := json.Marshal(crds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (service *ManagerService) reloadAfterChange() {
	err := prepareConfigFile(service.client, &service.config)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("PSK file prepared. Trying to reload mosquitto.")
	reloadConfig(&service.config)
	log.Printf("Mosquitto config reloaded.")
}

func StartServer() {
	mux := http.NewServeMux()
	service := NewManagerService()
	mux.HandleFunc("/add", service.add)
	mux.HandleFunc("/remove", service.remove)
	mux.HandleFunc("/list", service.list)

	log.Printf("Starting mosquitto-manager server...")

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
