# Local dev

To enable local development, that already takes the k8s deployment into consideration we use minikube.

## Requirements

Therefore the following must be installed:

- [Minikube](https://github.com/kubernetes/minikube)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [docker](https://docs.docker.com/engine/install/)

## Starting

Start your local cluster using minikube:

```sh
minikube start
```

    Optionally you can specify the resource amount the cluster should use. Use `minikube start --help` for more information.

Publish all images of local services into the cluster using the script:

```sh
sh ./build-all-local.sh
```

This publishes the version of your local services to the cluster. In some cases however this may not be needed i.e. if you want to use the image versions published to ghcr.

## Running the application

TODO
