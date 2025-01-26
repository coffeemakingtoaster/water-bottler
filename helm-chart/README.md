# Water bottler helm chart

## Prerequisites

If you are using the `local` image tag it is expected that the images are present in your cluster.
If you are using minikube this can be done with the `build-all-local.sh` script...if you are using something else you will have to figure this out yourself.

Alternatively there are `latest` packages on github.

    NOTE: These are not supported for arm because the github arm runners take about 10x the time to build arm docker images and we have limited runner minutes :)

Several deployments expect the cluster to support persisitent volume (claims). We use the default storage class of the cluster. If there is no storage class installed on the cluster/no default storage class we recommend [OpenEBS Local PV](https://openebs.io/docs/quickstart-guide/installation). This is NOT part of the helm chart because installing random storage classes without understanding what they do/how they work could cause harm to you k8s setup!

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

### Configuration

Use [overrides](https://helm.sh/docs/chart_template_guide/values_files/) to set values for the install.
For details on what can be configured, check the `values.yaml`.

### Known issues

Sometimes the rabbitmq does not seem to terminate properly after uninstalling.
In this case: Reinstall the helm chart -> uninstall it again.
