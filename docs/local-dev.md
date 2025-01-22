# Local dev


## Docker compose

There is a small docker compose setup for the application. However development using docker compose should only be used when development using minikube is not possible.

### Requirements

Its docker compose...so you need:

- [docker](https://docs.docker.com/engine/install/)

### Running the application

```sh
docker compose up -d
```

This should start all services. Relevant ports are:

- 8081 (upload service)
- 3111 (rabbitmq web interface with login values water:bottler)
- 80 (smtp4dev web interface)

## k8s

To enable local development, that already takes the k8s deployment into consideration we use minikube.

### Requirements

Therefor the following must be installed:

- [Minikube](https://github.com/kubernetes/minikube)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [docker](https://docs.docker.com/engine/install/)

### Starting

Start your local cluster using minikube:

```sh
minikube start --extra-config kubelet.EnableCustomMetrics=true
```

    Optionally you can specify the resource amount the cluster should use. Use `minikube start --help` for more information.

Publish all images of local services into the cluster using the script:

```sh
sh ./build-all-local.sh
```

This publishes the version of your local services to the cluster. In some cases however this may not be needed i.e. if you want to use the image versions published to ghcr.

### Running the application

Run the application by passing the `deploy` keyword to the build-all-local.sh script.

```sh
sh ./build-all-local.sh
```
