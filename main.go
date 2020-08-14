package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/mux"

	"time"

	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var handleSig bool

func handler(w http.ResponseWriter, r *http.Request) {
	var bgColor string
	var ok bool

	if bgColor, ok = os.LookupEnv("COLOR"); !ok {
		bgColor = "powderblue"
	}

	nodeName, clusterName, err := fetchDetails()
	if err != nil {
		log.Println(err)
		nodeName = "default"
		clusterName = "kubernetes"
	}

	fmt.Fprintf(w, "<html><body style=\"background-color:%s;\"><h1>Demo App</h1><p>Demo app running on node %s and cluster %s </p></body></html>", bgColor, nodeName, clusterName)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if !handleSig {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	} else {
		w.WriteHeader(http.StatusGone)
		json.NewEncoder(w).Encode(map[string]bool{"ok": false})
	}
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/health", healthCheck)
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	//Wait for sigterm
	<-c

	//Will cause healthcheck to unregister
	handleSig = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	log.Println("SIGTERM received.. waiting 10 seconds before shutdown")
	time.Sleep(10 * time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

func fetchDetails() (nodeName string, clusterName string, err error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nodeName, clusterName, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nodeName, clusterName, err
	}

	namespace, err := fetchNamespace()
	if err != nil {
		return nodeName, clusterName, err
	}

	podName, err := os.Hostname()
	if err != nil {
		return nodeName, clusterName, err
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), podName, v1.GetOptions{})
	if err != nil {
		return nodeName, clusterName, err
	}

	nodeName = pod.Spec.NodeName

	// Using an assumption the node name will be of format prefix-ClusterName-role-Count
	nodeArray := strings.Split(nodeName, "-")
	if len(nodeArray) >= 4 {
		clusterName = strings.Join(nodeArray[1:len(nodeArray)-2], "-")
	} else {
		clusterName = "kubernetes"
	}

	return nodeName, clusterName, err
}

func fetchNamespace() (namespace string, err error) {
	nsContent, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "", err
	}

	namespace = string(nsContent)
	return namespace, err
}
