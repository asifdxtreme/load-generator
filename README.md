# Load Generator

A Kubernetes demo that shows how horizontal pod scaling improves query performance.

## How It Works

The application exposes a `/query` endpoint that performs CPU-intensive work (500K SHA-256 hash iterations per request). When traffic is spread across more replicas, individual requests complete faster because the load is distributed.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [kind](https://kind.sigs.k8s.io/) or [minikube](https://minikube.sigs.k8s.io/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Go 1.22+](https://go.dev/dl/)

## Quick Start

### 1. Create a cluster (if you don't have one)

```bash
kind create cluster
```

### 2. Build and load the image

```bash
docker build -t load-generator:latest .
kind load docker-image load-generator:latest
```

### 3. Deploy with 2 replicas

```bash
kubectl apply -f k8s/
kubectl get pods -l app=load-generator
```

Wait for all pods to be `Running`.

### 4. Port-forward the service

```bash
kubectl port-forward svc/load-generator 8080:80
```

Keep this running in a separate terminal.

### 5. Phase 1 — Load test with 2 replicas

```bash
go run loadtest/main.go
```

Note the response times.

### 6. Scale to 5 replicas

```bash
kubectl scale deployment load-generator --replicas=5
kubectl get pods -l app=load-generator
```

Wait for all 5 pods to be `Running`.

### 7. Phase 2 — Load test with 5 replicas

```bash
go run loadtest/main.go
```

Compare the results — response times should be noticeably lower with 5 replicas since the CPU-intensive work is distributed across more pods.

## Load Tester Options

```
Usage:
  -url string    Target URL (default "http://localhost:8080/query")
  -c int         Number of concurrent requests (default 100)
```

## Cleanup

```bash
kubectl delete -f k8s/
kind delete cluster
```
