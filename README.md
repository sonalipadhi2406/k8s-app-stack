# Kubernetes App Stack

This project deploys a complete internal Kubernetes application stack using Helm, featuring a Golang monitoring app, MySQL database, and Nginx web server with secure networking, pod scheduling, and disaster recovery configurations.

## As per the tasks, I have created the following items:
Kubernetes cluster (Kind multi-node)
DB cluster with persistent data (MySQL StatefulSet with PVCs, 2 replicas)
Web Server (Nginx) with conditions
1 Multiple replicas 
2 Accessible from browser 
3 Custom config mounted (ConfigMap → nginx.conf)
4 Page shows Pod IP + serving-host from init container
Only allows web pods to connect to DB on 3306; deny others
Added DB disaster recovery
Flexible way to connect a Pod to a new network (Calico secondary IP pool + “external-net-test” pod routing; service still ClusterIP/NodePort)
Schedule specific DB replicas on specific nodes(node labels + StatefulSet affinity/nodeSelectorTerms)

## Golang application
Golang app that logs Pod create/delete/update(client-go informer; InClusterConfig; RBAC)
Used Helm to deploy all components (one parent chart + subcharts)
 
---

## Requirements
- Docker (latest)
- kind ≥ 0.22
- kubectl ≥ 1.29
- Helm ≥ 3.14
- make 

---

## Setup & Deployment Steps
```bash
NOTE: Use Helm to deploy all components

1. Create Kind Cluster

make kind-up
kubectl get nodes

2. Install Helm Dependencies & Deploy Application Stack
make deploy

3. make test

4. Verify Pods and Services

kubectl get pods -n app-stack
kubectl get svc -n app-stack

5. Access the Web App

kubectl port-forward svc/app-stack-webserver 8080:80 -n app-stack

- Then open in browser:
http://localhost:8080

Output EX:
Pod IP: 10.244.2.2
serving-host=Host-4v4tv

6. Validate NetworkPolicy
- Only webserver pods can connect to MySQL.
- Run this to test:

kubectl run test-client --rm -it --image=busybox -n app-stack -- /bin/sh
/ # nc -zv app-stack-mysql 3306

7. Test Backup Job

kubectl create job --from=cronjob/mysql-backup test-backup -n app-stack
kubectl logs -n app-stack -l job-name=test-backup


8. Golang application
Golang app that logs Pod create/delete/update

kubectl get pods -n app-stack -l app=pod-watcher
kubectl logs -n app-stack -l app=pod-watcher

---

# Trigger events
kubectl delete pod -n app-stack -l app=webserver
kubectl logs -n app-stack -l app=pod-watcher

## Where everything lives (quick map)
Cluster (Kind): infra/kind/kind-cluster.yaml
Parent Helm chart: charts/parent-chart/
MySQL: charts/parent-chart/charts/mysql/
Webserver: charts/parent-chart/charts/webserver/
Pod Watcher: charts/parent-chart/charts/pod-watcher/
NetworkPolicy: charts/parent-chart/charts/mysql/templates/networkpolicy.yaml
DB Backup CronJob: charts/parent-chart/charts/mysql/templates/cronjob-backup.yaml
Go watcher source/Dockerfile: charts/parent-chart/charts/pod-watcher/app/
