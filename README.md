# webhook

Inspired by [rancher webhook](https://github.com/rancher/webhook) and [harvester webhook](https://github.com/harvester/harvester/tree/master/pkg/webhook), it's a framework for developing a Kubernetes webhook easily.
Developers need to only implement a validator or mutator interface and register it to the webhook server.

## Example
The example implements a simple pod validator and pod mutator.
- pod validator: log when a pod is deleted.
- pod mutator: add label `example:example` when a pod is created.

### How to deploy
```shell
export KUBECONFIG=<your kubeconfig>
make apply
```
