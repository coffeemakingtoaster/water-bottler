# Local dev


## Docker compose

There is a small docker compose setup for the application. However, development using docker compose should only be used when development using minikube is not possible.

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
- [helm]()

### Starting

Start your local cluster using minikube:

```sh
minikube start 
```

    Optionally you can specify the resource amount the cluster should use. Use `minikube start --help` for more information.

Build all images of local services into the cluster using the script:

```sh
sh ./build-all-local.sh
```

This publishes the version of your local services to the cluster. In some cases however this may not be needed i.e. if you want to use the image versions published to ghcr.

### Running the application

Use helm (run this within the `helm-chart` directory): 

```sh
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm dependency build
helm install release-name .
```

For details on the helm chart see the respective [readme](../helm-chart/README.md).

To get the adress of a service use:
```sh
# for upload service
minikube service upload-service --url 
```
To use ingress instead of the service forwarding (for this the ingress has to be enabled in the values.yaml). Then the minikube plugin has to be used:
```sh
minikube addons enable ingress
# Validate the ingress works as intended
curl --resolve "smtp4dev.water-bottler.local:80:$( minikube ip )" -i http://smtp4dev.water-bottler.local:80/
# Upload an image using the ingress
curl -X POST -F image=@<image path here> -H "X-API-KEY: amVmZnMtd2F0ZXItYm90dGxlci1leGFtcGxlLWFwaS1rZXk=" --resolve "upload.water-bottler.local:80:$( minikube ip )" upload.water-bottler.local/upload -s -o /dev/null -w "%{http_code}"
```

For further info on ingresses in minikube see the [docs](https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/).
