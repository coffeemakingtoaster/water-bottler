#!/usr/bin/env bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
NC='\033[0m'

minikube status
if [ $? -eq 0 ]; then
	echo "${GREEN}Minikube is installed and running! Starting image builds...${NC}"
else
	echo "${RED}Minikube does not seem to be started (or installed). Make sure minikube is installed and your local cluster is running!${NC}"
	exit 1
fi

# Set the docker env to your local minikube docker daemon
eval $(minikube docker-env)

declare -a SERVICE_DIRECORIES=(authentication-service download-service notification-service upload-service object-recognition-service)

for dir in "${SERVICE_DIRECORIES[@]}"; do
	echo "Starting build for ${PURPLE}${dir}${NC}"
	if [ ! -f $dir/Dockerfile ]; then
    		echo "${YELLOW}No Dockerfile for ${dir}! Skipping service...${NC}"
	else
		docker build --quiet -t "github.com/coffeemakingtoaster/water-bottler/${dir}:local" ./${dir}
		echo "${GREEN}Build for ${PURPLE}${dir}${GREEN} done!${NC}"
	fi
done

echo "${GREEN}Services build. All images available in your cluster are displayed below:${NC}"
minikube image ls --format='table'
if [ "$1" != "deploy" ]; then
    echo "Skipping deployment"
    exit 0
fi 

CURRENT_CONTEXT=$(kubectl config current-context)
if [ "$CURRENT_CONTEXT" != "minikube" ]; then
    echo "${RED}The current kubectl context is NOT minikube! This script and application are not meant for production clusters as of now!${NC}"
    exit 1
fi

# Helm prepare
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# install rabbitmq cluster operator
kubectl apply -f "https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml"

# setup prometheus
helm install prometheus-stack prometheus-community/kube-prometheus-stack
helm install prometheus-adapter prometheus-community/prometheus-adapter --values ./overrides/prometheus-adapter.yaml

# start deployments
kubectl apply -f ./development-deployments/

#kubectl rollout restart deployment prometheus-adapter

for dir in "${SERVICE_DIRECORIES[@]}"; do
	echo "Starting build for ${PURPLE}${dir}${NC}"
	if [ ! -f $dir/deployment.yaml ]; then
    		echo "${YELLOW}No Deployment for ${dir}! Skipping service...${NC}"
	else
		(cd $dir && kubectl apply -f ./deployment.yaml)
		echo "${GREEN}Deployment for ${PURPLE}${dir}${GREEN} done!${NC}"
	fi
done

kubectl wait --for=condition=Available deployment/smtp4dev --timeout=300s

CMD="kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1 | jq '.resources.[].name' | grep services/rabbitmq_queue_messages_ready"

while true; do
  RESULT=$(eval "$CMD")
  if [[ -n "$RESULT" ]]; then
    echo "Autoscaler should now be ready!"
    break
  fi
  echo "Waiting custom metric to be available (this can take a few minutes)"
  sleep  10
done

echo ""
echo ""
echo "Done deploying"
echo "The api endpoint is accessible here:"
minikube service upload-service --url&
echo "The smtp dashboard is accessible here:"
minikube service smtp4dev --url
