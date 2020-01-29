package main

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func handlePush(w http.ResponseWriter, req *http.Request) {
	log.Info("ArgoCD notified")

	//

	// res, err := http.Post("https://argocd-server.argocd/api/webhook", "application/json", bytes.NewBufferString(""))

	// https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html

	// if err != nil {
	// 	log.Printf("Failed to notify a subscriber '%s'", webhook)
	// } else if res.Status != "200" {
	// 	log.Printf("Subscriber '%s' responded with non 200. Response code: %s", webhook, res.Status)
	// } else {
	// 	log.Printf("Subscriber '%s' notified successfully.", webhook)
	// }
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
	http.ListenAndServe(":8081", nil)
}
