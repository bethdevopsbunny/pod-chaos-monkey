package main

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"math/rand"
	"time"
)

type config struct {
	KubernetesNamespace string `yaml:"kubernetesNamespace"`
	Cron                string `yaml:"cron"`
	LabelSelector       string `yaml:"labelSelector"`
	ConcurrentDeletes   int    `yaml:"concurrentDeletes"`
}

var yamlConfig config

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

func main() {

	yamlConfig.RetrieveConfig()

	clientset, err := ApiAuthentication()
	if err != nil {
		panic(err.Error())
	}

	StartCron(clientset)
}

// RetrieveConfig grab the pod chaos monkey config file from disk, loads it into the app.
func (config *config) RetrieveConfig() *config {

	yamlFile, err := ioutil.ReadFile("config/pod_chaos_monkey.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		panic(err.Error())
	}

	return config
}

// StartCron creates, configures and starts the cron.
func StartCron(clientset *kubernetes.Clientset) {

	pcmCron := cron.New()

	_, err := pcmCron.AddFunc(yamlConfig.Cron, func() {

		pods, err := SelectPods(clientset, yamlConfig.ConcurrentDeletes)
		if err != nil {
			panic(err.Error())
		}

		for i := 0; i < len(pods); i++ {
			err = DeletePod(clientset, pods[i])
			if err != nil {
				panic(err.Error())
			}
		}
	})
	if err != nil {
		panic(err.Error())
	}

	pcmCron.Start()
	log.Infof("Cron Info: %+v\n", pcmCron.Entries())
	// Indefinite timer
	time.Sleep(1<<63 - 1)

}

// SelectPods makes a request to the api server for a list of pods matching the provided label selector and namespace.
//				returns a slice of random pods the size requested in deleteCount.
func SelectPods(clientset *kubernetes.Clientset, deleteCount int) ([]string, error) {

	var podsForDeletion []string
	rand.Seed(time.Now().UnixNano())

	options := metav1.ListOptions{LabelSelector: yamlConfig.LabelSelector}
	podList, err := clientset.CoreV1().Pods(yamlConfig.KubernetesNamespace).List(context.Background(), options)
	if err != nil {
		return podsForDeletion, fmt.Errorf("failed to retrieve podlist")
	}
	if len(podList.Items) <= 0 {
		return podsForDeletion, fmt.Errorf("no pods match that selector")
	}
	// clamps delete count within available number of pods
	if deleteCount > len(podList.Items) {
		log.Errorf("concurrentDeletes of %d set too high for the number of available pods reduced to %d", deleteCount, len(podList.Items))
		deleteCount = len(podList.Items)
	}

	// slice generator
	for i := 0; i < deleteCount; i++ {
		deleteIndex := rand.Intn(len(podList.Items))
		podsForDeletion = append(podsForDeletion, podList.Items[deleteIndex].Name)
		podList.Items = append(podList.Items[:deleteIndex], podList.Items[deleteIndex+1:]...)
	}

	log.WithField("pods", podsForDeletion).
		WithField("namespace", yamlConfig.KubernetesNamespace).
		WithField("label", yamlConfig.LabelSelector).
		Infof("API Info: pod triggered for deletion")

	return podsForDeletion, nil
}

// DeletePod makes the request to the api server to delete the pod.
func DeletePod(clientset *kubernetes.Clientset, podForDeletion string) error {

	err := clientset.CoreV1().Pods(yamlConfig.KubernetesNamespace).Delete(context.TODO(), podForDeletion, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("%g pod failed to delete", podForDeletion)
	}
	log.WithField("pod", podForDeletion).
		WithField("namespace", yamlConfig.KubernetesNamespace).
		WithField("label", yamlConfig.LabelSelector).
		Infof("API Info: pod successfully deleted")

	return nil
}

// ApiAuthentication establishes authentication with the kubeapi server.
func ApiAuthentication() (*kubernetes.Clientset, error) {

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to gather cluster config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset with cluster config")
	}

	return clientset, nil
}
