# Local Development
## Docker compose

There is a small docker compose setup for the application. However, **development using docker compose should only be used when development using minikube is not possible**.

### Requirements
- [docker](https://docs.docker.com/engine/install/)

### Running the application
```sh
docker compose up -d
```

This should start all services. 
Relevant ports are:
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

<details>
<summary> Ingress usage </summary>

To use ingress instead of the service forwarding (for this the ingress has to be enabled in the values.yaml). Then the minikube plugin has to be used:
```sh
minikube addons enable ingress
# Validate the ingress works as intended
curl --resolve "smtp4dev.water-bottler.local:80:$( minikube ip )" -i http://smtp4dev.water-bottler.local:80/
# Upload an image using the ingress
curl -X POST -F image=@<image path here> -H "X-API-KEY: amVmZnMtd2F0ZXItYm90dGxlci1leGFtcGxlLWFwaS1rZXk=" --resolve "upload.water-bottler.local:80:$( minikube ip )" upload.water-bottler.local/upload -s -o /dev/null -w "%{http_code}"
```

For further info on ingresses in minikube see the [docs](https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/).
</details>

### Using the application
To upload an image an API request has to be send to the upload service.

After processing the image can be downloaded from the download service using the url received via email.

In a local development environment the steps are the following:

1. Ensure that port forwarding is set up. We will use port mapping equal to the docker compose stack:
```sh
kubectl port-forward service/upload-service 8081:8080 &
kubectl port-forward service/download-service 8083:8080 &
kubectl port-forward service/smtp4dev 80:80 &
```

2. Send an image of your choice to upload service (port 8081). There are some example images in the `object-recognition-service/example-data/example-images/`  directory:

```sh
curl -X POST -F image=@<image path here> -H "X-API-KEY: amVmZnMtd2F0ZXItYm90dGxlci1leGFtcGxlLWFwaS1rZXk=" localhost:8081/upload -s -o /dev/null -w "%{http_code}"
# Should return 202
```

3. Wait for email to be send to smtp4dev. Just check the [email inbox](http://localhost:80) in you browser. You should receive an 'You Image is ready' email shortly.

4. Go to the link within the email and receiver your processed image! (The link should look similar to this one `http://localhost:8083/download?file=`)
