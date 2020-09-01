package internal

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/mitchellh/go-ps"
	"io"
	"log"
	"mosquitto-manager/internal/api/types/v1alpha1"
	"net/http"
	"strconv"
	"strings"
)

type CredsWithId struct {
	Id       string `json:"Id" bson:"_id"`
	Login    string
	Password string
	Acls     []v1alpha1.Acl
}

type Creds struct {
	Login    string
	Password string
	Acls     []v1alpha1.Acl
}

type Id struct {
	Id string
}

type ManagerService struct {
	manager Manager
	config  Config
}

func NewManagerService() ManagerService {
	var service = ManagerService{}
	var kubeconfig, mongoUri, mongoDatabase, mongoCollection,
		pskFilePath, basicAuthLogin, basicAuthPass, port, crt, key, aclFile *string
	var mosquittoPid *int
	var acl *bool
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	mongoUri = flag.String("mongoUri", "", "MongoDB Uri if empty Kubernetes CRDs are used")
	mongoDatabase = flag.String("mongoDatabase", "mosquittoManager", "Mongo database used to store data, default mosquittoManager")
	mongoCollection = flag.String("mongoCollection", "data", "Mongo collection used to store data, default data")
	mosquittoPid = flag.Int("mosquittoPid", 0, "pid of mosquitto process")
	pskFilePath = flag.String("pskFilePath", "", "path to pskfile")
	basicAuthLogin = flag.String("basicAuthLogin", "", "basic auth login")
	basicAuthPass = flag.String("basicAuthPass", "", "basic auth password")
	port = flag.String("port", "8080", "port for mosquitto manager api")
	crt = flag.String("crt", "", "TLS crt path if empty http")
	key = flag.String("key", "", "TLS key path if empty http")
	acl = flag.Bool("acl", false, "If true the acls are created and managed")
	aclFile = flag.String("aclFile", "", "Path to mosquitto acl file if empty and acl=true, the default path is used")
	flag.Parse()

	log.Printf("ACL is set to " + strconv.FormatBool(*acl))

	if *mosquittoPid == 0 {
		mosquittoPid = tryFindMosquittoPidByName()
	}
	if *pskFilePath == "" {
		*pskFilePath = "/proc/" + strconv.Itoa(*mosquittoPid) + "/root/etc/mosquitto/pskfile"
	}
	if *acl && *aclFile == "" {
		*aclFile = "/proc/" + strconv.Itoa(*mosquittoPid) + "/root/etc/mosquitto/acl.conf"
	}

	log.Printf("Mosquitto pid - " + strconv.Itoa(*mosquittoPid))
	log.Printf("Pskfile path - " + *pskFilePath)
	var manager = createManager(kubeconfig, mongoUri, mongoDatabase, mongoCollection)
	var config = NewConfig(*mosquittoPid, *pskFilePath, *basicAuthLogin, *basicAuthPass, *port, *crt, *key, *aclFile)
	service.config = config
	service.manager = manager
	return service
}

func createManager(kubeconfig *string, mongoUri *string, mongoDatabase *string, mongoCollection *string) Manager {
	if *mongoUri == "" {
		return NewK8sManager(kubeconfig)
	} else {
		return NewMongoManager(MongoConfig{
			Uri:        *mongoUri,
			Database:   *mongoDatabase,
			Collection: *mongoCollection,
		})
	}
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
	err := service.checkAuth(w, r)
	if err != nil {
		return
	}
	var lp Creds
	err = json.NewDecoder(r.Body).Decode(&lp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Adding Login=" + lp.Login + " Password=" + lp.Password)
	id, err := service.manager.Create(lp)
	if err != nil {
		log.Printf("Error during adding")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Added Login=" + lp.Login + " Password=" + lp.Password + " ID=" + *id)
	if !service.manager.IsObserveSupported() {
		service.reloadAfterChange()
	}
	_, err = io.WriteString(w, *id)
	if err != nil {
		log.Print("Error during writing the response.", err)
	}

}

func (service *ManagerService) remove(w http.ResponseWriter, r *http.Request) {
	err := service.checkAuth(w, r)
	if err != nil {
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/creds/")
	log.Printf("Removing Id=" + id)
	err = service.manager.Remove(Id{id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Removed Id=" + id)
	if !service.manager.IsObserveSupported() {
		service.reloadAfterChange()
	}
}

func (service *ManagerService) list(w http.ResponseWriter, r *http.Request) {
	err := service.checkAuth(w, r)
	if err != nil {
		return
	}
	crds := service.manager.GetAll()
	js, err := json.Marshal(crds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

func (service *ManagerService) getById(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := service.checkAuth(w, r)
		if err != nil {
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/creds/")
		log.Println("Trying to get creds by Id=" + id + " and path is " + r.URL.Path)
		creds, err := service.manager.Get(Id{Id: id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		js, err := json.Marshal(*creds)
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

func (service *ManagerService) update(w http.ResponseWriter, r *http.Request) {
	err := service.checkAuth(w, r)
	if err != nil {
		return
	}
	var lp Creds
	err = json.NewDecoder(r.Body).Decode(&lp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/creds/")
	log.Print("Updating " + id + " Login=" + lp.Login + " Password=" + lp.Password)
	err = service.manager.Update(Id{Id: id}, lp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Print("Updated " + id + " Login=" + lp.Login + " Password=" + lp.Password)
	if !service.manager.IsObserveSupported() {
		service.reloadAfterChange()
	}
}

func (service *ManagerService) reloadAfterChange() {
	log.Printf("Trying to prepare PSK file")
	err := preparePskFile(service.manager, &service.config)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("PSK file prepared")
	if service.config.aclFile != "" {
		log.Printf("Trying to prepare ACL file")
		err = prepareAclFile(service.manager, &service.config)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("ACL file prepared")
	} else {
		log.Printf("Ignoring ACLs because ACL file path not set")
	}
	reloadConfig(&service.config)
	log.Printf("Mosquitto config reloaded.")
}

func (service *ManagerService) credsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		service.list(w, r)
	case "POST":
		service.add(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (service *ManagerService) credsIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		service.update(w, r)
	case "DELETE":
		service.remove(w, r)
	case "GET":
		service.getById(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func StartServer() {
	mux := http.NewServeMux()
	service := NewManagerService()
	mux.HandleFunc("/creds", service.credsHandler)
	mux.HandleFunc("/creds/", service.credsIdHandler)

	go watchAndSyncCredsWithPskFile(service)

	if service.config.isTLS() {
		log.Printf("Starting mosquitto-manager TLS server on port " + service.config.port)
		err := http.ListenAndServeTLS(service.config.port, service.config.crt, service.config.key, mux)
		log.Fatal(err)
	} else {
		log.Printf("Starting mosquitto-manager server on port " + service.config.port)
		err := http.ListenAndServe(service.config.port, mux)
		log.Fatal(err)
	}

}

func watchAndSyncCredsWithPskFile(service ManagerService) {
	if service.manager.IsObserveSupported() {
		service.manager.ObserveIfSupported(service)
	} else {
		service.reloadAfterChange()
	}
}
