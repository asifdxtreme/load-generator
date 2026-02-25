# Load Generator

A Kubernetes demo that shows how horizontal pod scaling improves query performance.

## How It Works

The application exposes a `/query` endpoint that performs CPU-intensive work (500K SHA-256 hash iterations per request). When traffic is spread across more replicas, individual requests complete faster because the load is distributed.

## Prerequisites

- Sign in : https://killercoda.com/playgrounds/course/kubernetes-playgrounds/two-node --> 1.34 Two Nodes setup

## Quick Start

### 1. Clone this repo

```bash
git clone https://github.com/asifdxtreme/load-generator.git
cd load-generator
```

### 2. Deploy with 2 replicas

```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl get pods -l app=load-generator
```

Wait for all pods to be `Running`.

### 3. Phase 1 — Load test with 2 replicas

```bash
kubectl apply -f k8s/loadtest-job.yaml
kubectl wait --for=condition=complete job/load-test --timeout=120s
kubectl logs job/load-test
```

Note the response times.

### 4. Scale to 5 replicas

```bash
kubectl delete job load-test
kubectl scale deployment load-generator --replicas=5
kubectl get pods -l app=load-generator
```

Wait for all 5 pods to be `Running`.

### 5. Phase 2 — Load test with 5 replicas

```bash
kubectl apply -f k8s/loadtest-job.yaml
kubectl wait --for=condition=complete job/load-test --timeout=120s
kubectl logs job/load-test
```

Compare the results — response times should be noticeably lower with 5 replicas since the CPU-intensive work is distributed across more pods.

## Load Tester Options

You can customize the load test by editing `k8s/loadtest-job.yaml`:

```yaml
command: ["load-generator", "loadtest", "-c", "10", "-url", "http://load-generator/query"]
```

| Flag   | Default                           | Description                    |
|--------|-----------------------------------|--------------------------------|
| `-c`   | `10`                              | Number of concurrent requests  |
| `-url` | `http://load-generator/query`     | Target URL                     |

## Cleanup

```bash
kubectl delete -f k8s/
```
