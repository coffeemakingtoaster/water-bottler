# Water bottler helm chart

## Prerequisites

If you are using the `local` image tag it is expected that the images are present in your cluster.
If you are using minikube this can be done with the `build-all-local.sh` script...if you are using something else you will have to figure this out yourself.

Alternatively there are `latest` packages on github.

    NOTE: These are not supported for arm because the github arm runners take about 10x the time to build arm docker images and we have limited runner minutes :)

```sh
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm dependency build
helm install release-name .
```

When installed it can take a bit until the system is entirely ready (especially the metric propagation is rather heavy). 
Wait for this command to return something if in doubt:

```sh
kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1 | jq '.resources.[].name' | grep services/rabbitmq_queue_messages_ready
```

### Uninstall

Remove from cluster

```sh
helm uninstall <release name>
```
