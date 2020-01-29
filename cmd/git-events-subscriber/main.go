package main

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func webhookPayload() string {
	return "{\"repository\": {\"links\": {\"html\": {\"href\": \"http://bb.alutech-mc.com:8080/scm/as/infra.git\"}}},\"push\": {\"changes\": [{\"new\": {\"name\": \"master\"}}]}}"
}

func handlePush(w http.ResponseWriter, req *http.Request) {
	res, err := http.Post("https://argocd-server.argocd/api/webhook", "application/json", bytes.NewBufferString(webhookPayload()))

	if err != nil {
		log.Warnf("Failed to notify local ArgoCD. Error details: '%s'", err)
	} else if res.Status != "200" {
		log.Warnf("Local ArgoCD responded with non 200. Response code: %s", res.Status)
	} else {
		log.Infof("Local ArgoCD is notified successfully.")
	}
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
	res, err := http.Post(
		os.Getenv("PUBLISHER_URL")+"/subscribers",
		"application/json",
		bytes.NewBufferString(fmt.Sprintf("{\"WebhookUrl\":\"%s/push\"}", os.Getenv("MY_INGRESS_URL"))),
	)

	if err != nil {
		log.Warnf("Failed to register in publisher. Error details: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if res.StatusCode != 200 {
		log.Warnf("Failed to register in publisher. Non 200 response: %s", res.StatusCode)
		http.Error(w, err.Error(), res.StatusCode)
	} else {
		log.Debug("Registered successfully in publisher.")
	}
}

func main() {
	log.Infof("PUBLISHER_URL: %s", os.Getenv("PUBLISHER_URL"))
	log.Infof("MY_INGRESS_URL: %s", os.Getenv("MY_INGRESS_URL"))

	r := mux.NewRouter()
	r.HandleFunc("/push", handlePush).Methods("POST")
	r.HandleFunc("/health", healthCheck).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
