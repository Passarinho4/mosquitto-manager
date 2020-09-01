Mosquitto Manager is a simple addon to mosquitto which allows you to manage mosquitto pskfile via json over HTTP API.
All credentials are stored as Kubernetes Custom Resources or in MongoDB.

The code is in pre-alpha version please do not use it in production.

Command line options:

| Parameter | Default value | Comment |
| --------- | ------------- | ------- |
| kubeconfig | InClusterConfig (https://godoc.org/k8s.io/client-go/rest#InClusterConfig) | absolute path to the kubeconfig file |
| mongoUri |  | MongoDB Uri if empty Kubernetes CRDs are used (details - https://docs.mongodb.com/manual/reference/connection-string/) |
| mongoDatabase | mosquittoManager | Mongo database used to store data|
| mongoCollection | data | Mongo collection used to store data |
| mosquittoPid | 0 | pid of mosquitto process (just for development) |
| pskFilePath | "/proc/" + mosquittoPid + "/root/etc/mosquitto/pskfile" | path to pskfile (just for development) |
| basicAuthLogin |  | basic auth login if empty auth is disabled |
| basicAuthPass |  | basic auth password if empty auth is disabled  |
| port | 8080 | port for mosquitto manager api |
| crt |  | TLS crt path if empty http |
| key |  | TLS key path if empty http |
| acl | false | If true the acls are created and managed |
| aclFile | "/proc/" + mosquittoPid + "/root/etc/mosquitto/acl.conf" | Path to mosquitto acl file if empty and acl=true (just for development) |

if basic auth is not set endpoints are not secured.

Parameters marked as "just for development" have their default values configured for K8s deployment which is recommended 
in `/yamls/` examples. You can override them during development process to run mosquitto-manager locally.  

In `/yamls/pod.yaml` you can find two not obvious options: 
`shareProcessNamespace: true` and `SYS_PTRACE`. These options are required to allow mosquitto manager application 
(launched in different container than mosquitto process) edit the pskfile and send the reload config signal. 
The exact details about how it works - `managerServer.go`. 

Helpful commands:
Subscribe on topic news using mosquitto console client:

`mosquitto_sub -h 127.0.0.1 -p 1883 -t news -u mosquitto`

Publish on topic news using mosquitto console client with auth. 

`mosquitto_pub -h 127.0.0.1 -p 8883 -t news -m Hello I am alive xd --psk-identity l --psk 70 --insecure --debug -u mosquitto`


Http request to create new mosquitto user. 

`curl --location --request POST 'http://localhost:8080/creds' \
 --header 'Content-Type: application/json' \
 --data-raw '{
     "Login": "login-test",
     "Password": "password-test"
 }'`

Http request to list mosquitto users.

`curl --location --request GET 'localhost:8080/creds'`

Http request to get mosquitto user by ID (returned by add enpoint).

`curl --location --request GET 'http://localhost:8080/creds/5f4e48c7e184b1045c373c30'`

Http request to remove mosquitto user.

`curl --location --request DELETE 'localhost:8080/creds/5f4e48c7e184b1045c373c30' \
 --data-raw ''`

To enable ACL you need to build mosquitto image (from mosquitto directory) 
with acl_file `/etc/mosquitto/acl.conf` option in `mosquitto.conf` and
pass `acl=true` to mosquitto manager. 

API documentation is available here:
https://documenter.getpostman.com/view/109540/TVCe29My

Postman tests are in tests directory. 

