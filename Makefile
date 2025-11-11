K8S_NAMESPACE ?= app-stack
CHART_PATH := charts/parent-chart

kind-up:
	kind create cluster --config infra/kind/kind-cluster.yaml

kind-down:
	kind delete cluster --name app-stack

deploy:
	helm dependency update $(CHART_PATH)
	helm upgrade --install app-stack $(CHART_PATH) -n $(K8S_NAMESPACE) --create-namespace

test:
	kubectl get all -n $(K8S_NAMESPACE)

cleanup:
	helm uninstall app-stack -n $(K8S_NAMESPACE)
