# k8s-log-proxy

Simply HTTP interface for fetching Kubernetes pod/container logs.

It exposes a single HTTP GET endpoint `/logs/<:namespace>/<:pod_name>`
under port 8080, which returns the most recent 8192 log lines
from the specified pod.

## Usage

[Helm](https://helm.sh) must be installed to use the chart.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

  helm repo add k8s-log-proxy https://adrianokf.github.io/k8s-log-proxy

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages.  You can then run `helm search repo
k8s-log-proxy` to see the charts.

To install the `k8s-log-proxy` chart:

    helm install logproxy k8s-log-proxy/k8s-log-proxy

To uninstall the chart:

    helm delete logproxy

## Development

### Skaffold

If you have access to a running Kubernetes cluster, you can
use [Skaffold](https://skaffold.dev) for rapid development.

Skaffold handles the workflow for building, pushing and deploying
the application, by simply invoking the following command
in the root folder of the repository:

```
$ skaffold dev
```

The `deploy/` folder contains all necessary Kubernetes resources
for deployment, including a service account with appropriate
role bindings to read all container logs.

## Notes

This tool is a very quick solution, with no consideration of
security or privacy concerns. Any user with access to the HTTP
interface can access the logs of **any** pod in the cluster.

Proceed with caution!
