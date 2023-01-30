# webhook

Inspired by [rancher webhook](https://github.com/rancher/webhook) and [harvester webhook](https://github.com/harvester/harvester/tree/master/pkg/webhook), it's a framework for developing a Kubernetes webhook easily.
Developers need to only implement interface and register it to the webhook server. It supports validator, mutator and CRD conversion webhook.

Go to [harvester/webhook-sample](https://github.com/harvester/webhook-sample) to look a simple example.
