# glue-jobs-operator

[glue-jobs-operator](https://github.com/90poe/glue-jobs-operator) is AWS Glue Job controller for Kubernetes

![Version: 1.0.0](https://img.shields.io/badge/Version-1.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.3.2](https://img.shields.io/badge/AppVersion-0.3.2-informational?style=flat-square)

To use, create role, which will allow operator on K8S to access AWS Glue and create jobs.

This chart bootstraps an glue-jobs-operator deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Chart version 1.x.x: Kubernetes v1.27+

## Get Repo Info

```
git clone git@github.com:90poe/glue-jobs-operator.git
```

## Install Chart

**Important:** only helm3 is supported

```console
cd helm/glue-jobs-operator
helm install [RELEASE_NAME] .
```

The command deploys glue-jobs-operator on the Kubernetes cluster in the default configuration.

_See [configuration](#configuration) below._

_See [helm install](https://helm.sh/docs/helm/helm_install/) for command documentation._

## Uninstall Chart

```console
helm uninstall [RELEASE_NAME]
```

This removes all the Kubernetes components associated with the chart and deletes the release.

_See [helm uninstall](https://helm.sh/docs/helm/helm_uninstall/) for command documentation._

## Upgrading Chart

```console
helm upgrade [RELEASE_NAME] [CHART] --install
```

_See [helm upgrade](https://helm.sh/docs/helm/helm_upgrade/) for command documentation._

## Configuration

See [Customizing the Chart Before Installing](https://helm.sh/docs/intro/using_helm/#customizing-the-chart-before-installing). To see all configurable options with detailed comments, visit the chart's [values.yaml](./values.yaml), or run these configuration commands:

```console
helm show values ingress-nginx/ingress-nginx
```

### PodDisruptionBudget

Note that the PodDisruptionBudget resource will only be defined if the replicaCount is greater than one,
else it would make it impossible to evacuate a node. See [gh issue #7127](https://github.com/helm/charts/issues/7127) for more info.

### Prometheus Metrics

The Nginx ingress controller can export Prometheus metrics, by setting `controller.metrics.enabled` to `true`.

You can add Prometheus annotations to the metrics service using `controller.metrics.service.annotations`.
Alternatively, if you use the Prometheus Operator, you can enable ServiceMonitor creation using `controller.metrics.serviceMonitor.enabled`. And set `controller.metrics.serviceMonitor.additionalLabels.release="prometheus"`. "release=prometheus" should match the label configured in the prometheus servicemonitor ( see `kubectl get servicemonitor prometheus-kube-prom-prometheus -oyaml -n prometheus`)
