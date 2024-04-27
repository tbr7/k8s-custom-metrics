package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/streadway/amqp"
)

type APIResourceList struct {
	Kind         string        `json:"kind"`
	APIVersion   string        `json:"apiVersion"`
	GroupVersion string        `json:"groupVersion"`
	Resources    []APIResource `json:"resources"`
}

type APIResource struct {
	Name         string   `json:"name"`
	SingularName string   `json:"singularName"`
	Namespaced   bool     `json:"namespaced"`
	Kind         string   `json:"kind"`
	Verbs        []string `json:"verbs"`
}

type MetricValueList struct {
	Kind       string        `json:"kind"`
	APIVersion string        `json:"apiVersion"`
	Metadata   Metadata      `json:"metadata"`
	Items      []MetricValue `json:"items"`
}

type Metadata struct {
	SelfLink string `json:"selfLink"`
}

type MetricValue struct {
	DescribedObject DescribedObject `json:"describedObject"`
	MetricName      string          `json:"metricName"`
	Timestamp       time.Time       `json:"timestamp"`
	Value           string          `json:"value"` // This can be an integer or a string, for example "100" or "100m" if milli-units
}

type DescribedObject struct {
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	APIVersion string `json:"apiVersion"`
}

func main() {
	log.Println("Starting server on port 5001")
	http.HandleFunc("/", metricsHandler)
	// Generated via openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj '/CN=localhost'
	// apiservice requires https but doesn't actually care if insecure. see "insecureSkipTLSVerify: true" in apiservice.yaml
	err := http.ListenAndServeTLS(":5001", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request:", r.Method, r.URL.Path)

	// Respond with a list of available metrics at the root custom metrics API path
	if r.URL.Path == "/apis/custom.metrics.k8s.io/v1beta1" {
		listAvailableMetrics(w)
		return
	}

	// Specific metric value request
	if r.URL.Path == "/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/services/metrics-exporter/rabbitmq_queue_length" {
		provideMetricValue(w, r)
		return
	}

	// Respond with 404 if no valid endpoint is found
	http.NotFound(w, r)
}

func listAvailableMetrics(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResourceList{
		Kind:         "APIResourceList",
		APIVersion:   "v1",
		GroupVersion: "custom.metrics.k8s.io/v1beta1",
		Resources: []APIResource{
			{
				Name:         "namespaces/default/services/*/rabbitmq_queue_length",
				SingularName: "",
				Namespaced:   true,
				Kind:         "MetricValueList",
				Verbs:        []string{"get"},
			},
		},
	})
}

func getQueueLength() int {
	conn, err := amqp.Dial("amqp://user:pass@my-rabbitmq.default.svc.cluster.local:5672/")
	if err != nil {
		log.Printf("Error connecting to RabbitMQ: %s", err)
		return 0
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Error opening channel: %s", err)
		return 0
	}
	defer ch.Close()

	q, err := ch.QueueInspect("hello")
	if err != nil {
		log.Printf("Error inspecting queue: %s", err)
		return 0
	}

	return q.Messages
}

func provideMetricValue(w http.ResponseWriter, r *http.Request) {
	queueLength := getQueueLength()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MetricValueList{
		Kind:       "MetricValueList",
		APIVersion: "custom.metrics.k8s.io/v1beta1",
		Metadata: Metadata{
			SelfLink: r.URL.Path,
		},
		Items: []MetricValue{
			{
				DescribedObject: DescribedObject{
					Kind:       "Service",
					Namespace:  "default",
					Name:       "metrics-exporter",
					APIVersion: "v1",
				},
				MetricName: "rabbitmq_queue_length",
				Timestamp:  time.Now(),
				Value:      fmt.Sprintf("%d", queueLength),
			},
		},
	})
}
