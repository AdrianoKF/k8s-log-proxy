# k8s-log-proxy

Simply HTTP interface for fetching Kubernetes pod/container logs.

It exposes a single HTTP GET endpoint `/logs/<:namespace>/<:pod_name>`
under port 8080, which returns the most recent 8192 log lines
from the specified pod.

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
