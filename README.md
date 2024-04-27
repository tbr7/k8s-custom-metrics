# k8s-custom-metrics

This repository contains two primary components: `SimpleApp` and `AutoScalingDemo`, which demonstrate an end-to-end scenario for auto-scaling applications within a Kubernetes cluster based on custom metrics from RabbitMQ.

## Overview

The purpose of this project is to showcase how Kubernetes can dynamically scale applications based on the workload observed through custom metrics. This is particularly useful for applications that need to scale based on the number of messages in a queue, such as those processing tasks or transactions asynchronously.

- **SimpleApp**: A basic Go web server that responds to HTTP requests. It's designed to be simple and is used primarily to illustrate basic scaling and Kubernetes deployment.
- **AutoScalingDemo**: Integrates with RabbitMQ to simulate a scenario where the application needs to scale based on the number of messages in a queue. It includes a custom metrics exporter to facilitate this interaction with Kubernetes' Horizontal Pod Autoscaler (HPA).

## Architecture

The demo setup includes the following components:
- A **RabbitMQ** deployment that acts as the message queueing service.
- The **SimpleApp** which processes messages from RabbitMQ.
- A **Metrics Exporter** that reads the queue length from RabbitMQ and exposes it for the Kubernetes HPA.
- The **Kubernetes HPA** configured to scale the SimpleApp based on the queue length reported by the Metrics Exporter.

## Prerequisites

Before you start, ensure you have the following installed:
- Kubernetes cluster or minikube
- kubectl configured to interact with your cluster
- Helm, for deploying RabbitMQ
- Docker, for building and pushing the application containers
- Access to a container registry (optional, if pushing images)

## Installation

Follow these steps to set up the environment:

### Deploy RabbitMQ

Use Helm to deploy RabbitMQ into your Kubernetes cluster:

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-rabbitmq bitnami/rabbitmq
```

### Build and Deploy SimpleApp and Metrics Exporter

1. **Build the Docker images**:
    ```bash
    docker build -t your-registry/simple-app:v1 ./SimpleApp
    docker build -t your-registry/metrics-exporter:v1 ./AutoScalingDemo
    docker push your-registry/simple-app:v1
    docker push your-registry/metrics-exporter:v1
    ```

2. **Deploy the applications**:
    ```bash
    kubectl apply -f deployment.yaml
    ```

### Configure HPA

Apply the HPA configuration to enable auto-scaling:

```bash
kubectl apply -f hpa.yaml
```

## Usage

To simulate load, send messages to the RabbitMQ queue:

```bash
# Port-forward RabbitMQ management port
kubectl port-forward svc/my-rabbitmq 15672:15672

# Access RabbitMQ management at http://localhost:15672 and publish messages to the 'hello' queue.
```

## Configuration

- **RabbitMQ Credentials**: Set up the correct credentials in your application deployment configurations.
- **Scaling Thresholds**: Adjust the HPA thresholds in `hpa.yaml` based on observed needs.

## Contributing

Contributions to this project are welcome! Please fork the repository and submit a pull request with your changes or improvements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
