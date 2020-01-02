# Helm Chart for bk-cmdb

**Notes:** The master branch is in heavy development, please use the codes on other branch instead. A high available solution for bk-cmdb based on chart can be find [here](docs/High%20Availability.md). And refer to the [guide](docs/Upgrade.md) to upgrade the existing deployment.

## Introduction

This [Helm](https://github.com/kubernetes/helm) chart installs [bk-cmdb](https://github.com/Tencent/bk-cmdb) in a Kubernetes cluster. Welcome to [contribute](CONTRIBUTING.md) to Helm Chart for bk-cmdb.

## Prerequisites

- Kubernetes cluster 1.10+
- Helm 2.8.0+

## Installation

### Download the chart

Download bk-cmdb helm chart code.

```bash
git clone https://github.com/Tencent/bk-cmdb-helm
```

Checkout the branch.

```bash
cd bk-cmdb-helm
git checkout branch_name
```

### Configure the chart

The following items can be configured in `values.yaml` or set via `--set` flag during installation.

#### Configure the way how to expose bk-cmdb service:

- **Ingress**: The ingress controller must be installed in the Kubernetes cluster.
- **ClusterIP**: Exposes the service on a cluster-internal IP. Choosing this value makes the service only reachable from within the cluster.
- **NodePort**: Exposes the service on each Node’s IP at a static port (the NodePort). You’ll be able to contact the NodePort service, from outside the cluster, by requesting `NodeIP:NodePort`.
- **LoadBalancer**: Exposes the service externally using a cloud provider’s load balancer.

#### Configure the external URL

The external URL for bk-cmdb core service is used to:

1. populate the docker/helm commands showed on portal
2. populate the token service URL returned to docker/notary client

Format: `protocol://domain[:port]`. Usually:

- if expose the service via `Ingress`, the `domain` should be the value of `expose.ingress.hosts.core`
- if expose the service via `ClusterIP`, the `domain` should be the value of `expose.clusterIP.name`
- if expose the service via `NodePort`, the `domain` should be the IP address of one Kubernetes node
- if expose the service via `LoadBalancer`, set the `domain` as your own domain name and add a CNAME record to map the domain name to the one you got from the cloud provider  

If bk-cmdb is deployed behind the proxy, set it as the URL of proxy.

#### Configure the way how to persistent data:

- **Disable**: The data does not survive the termination of a pod.
- **Persistent Volume Claim(default)**: A default `StorageClass` is needed in the Kubernetes cluster to dynamic provision the volumes. Specify another StorageClass in the `storageClass` or set `existingClaim` if you have already existing persistent volumes to use.
- **External Storage(only for images and charts)**: For images and charts, the external storages are supported: `azure`, `gcs`, `s3` `swift` and `oss`.

#### Configure the other items listed in [configuration](#configuration) section.

### Install the chart

Install the bk-cmdb helm chart with a release name `my-release`:

```bash
helm install --name my-release .
```

## Uninstallation

To uninstall/delete the `my-release` deployment:

```bash
helm delete --purge my-release
```

## Configuration

The following table lists the configurable parameters of the bk-cmdb chart and the default values.

| Parameter                                                                   | Description                                                                                                                                                                                                                                                                                                                                     | Default                         |
| --------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------- |
| **Expose**                                                                  |
| `registry.registry.resources`                                               | The [resources] to allocate for container                                                                                                                                                                                                                                                                                                       | undefined                       |  | `dev` |
| `registry.controller.image.repository`                                      | Repository for registry controller image                                                                                                                                                                                                                                                                                                        | `goharbor/harbor-registryctl`   |
| `redis.internal.affinity`                                                   | Node/Pod affinities                                                                                                                                                                                                                                                                                                                             | `{}`                            |
| `redis.external.host`                                                       | The hostname of external Redis                                                                                                                                                                                                                                                                                                                  | `192.168.0.2`                   |
| `redis.external.port`                                                       | The port of external Redis                                                                                                                                                                                                                                                                                                                      | `6379`                          |
| `redis.external.coreDatabaseIndex`                                          | The database index for core                                                                                                                                                                                                                                                                                                                     | `0`                             |
| `redis.external.jobserviceDatabaseIndex`                                    | The database index for jobservice                                                                                                                                                                                                                                                                                                               | `1`                             |
| `redis.external.registryDatabaseIndex`                                      | The database index for registry                                                                                                                                                                                                                                                                                                                 | `2`                             |
| `redis.external.chartmuseumDatabaseIndex`                                   | The database index for chartmuseum                                                                                                                                                                                                                                                                                                              | `3`                             |
| `redis.external.password`                                                   | The password of external Redis                                                                                                                                                                                                                                                                                                                  |                                 |
| `redis.podAnnotations`                                                      | Annotations to add to the redis pod                                                                                                                                                                                                                                                                                                             | `{}`                            |

[resources]: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/
