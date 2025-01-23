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
