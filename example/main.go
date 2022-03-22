package main

import (
	"context"
	"os"

	"github.com/rancher/wrangler/pkg/kubeconfig"
	"github.com/rancher/wrangler/pkg/leader"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/yaocw2020/webhook/example/admitter"
	"github.com/yaocw2020/webhook/pkg/config"
	"github.com/yaocw2020/webhook/pkg/server"
	"github.com/yaocw2020/webhook/pkg/types"
)

func main() {
	var options config.Options
	var logLevel string

	flags := []cli.Flag{
		cli.StringFlag{
			Name:        "loglevel",
			Usage:       "Specify log level",
			EnvVar:      "LOGLEVEL",
			Value:       "info",
			Destination: &logLevel,
		},
		cli.IntFlag{
			Name:        "threadiness",
			EnvVar:      "THREADINESS",
			Usage:       "Specify controller threads",
			Value:       5,
			Destination: &options.Threadiness,
		},
		cli.IntFlag{
			Name:        "https-port",
			EnvVar:      "WEBHOOK_SERVER_HTTPS_PORT",
			Usage:       "HTTPS listen port",
			Value:       8443,
			Destination: &options.HTTPSListenPort,
		},
		cli.StringFlag{
			Name:        "namespace",
			EnvVar:      "NAMESPACE",
			Destination: &options.Namespace,
			Usage:       "The harvester namespace",
			Value:       "harvester-system",
			Required:    true,
		},
		cli.StringFlag{
			Name:        "controller-user",
			EnvVar:      "CONTROLLER_USER_NAME",
			Destination: &options.ControllerUsername,
			Value:       "harvester-load-balancer",
			Usage:       "The harvester controller username",
		},
		cli.StringFlag{
			Name:        "gc-user",
			EnvVar:      "GARBAGE_COLLECTION_USER_NAME",
			Destination: &options.GarbageCollectionUsername,
			Usage:       "The system username that performs garbage collection",
			Value:       "system:serviceaccount:kube-system:generic-garbage-collector",
		},
	}

	cfg, err := kubeconfig.GetNonInteractiveClientConfig(os.Getenv("KUBECONFIG")).ClientConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	client := kubernetes.NewForConfigOrDie(cfg)
	ctx := signals.SetupSignalContext()

	app := cli.NewApp()
	app.Flags = flags
	app.Action = func(c *cli.Context) {
		setLogLevel(logLevel)
		if err := run(ctx, cfg, &options); err != nil {
			logrus.Fatalf("run webhook server failed: %v", err)
		}
	}

	leader.RunOrDie(ctx, "default", "webhook-example", client, func(ctx context.Context) {
		if err := app.Run(os.Args); err != nil {
			logrus.Fatalf("run webhook server failed: %v", err)
		}
	})
}

func run(ctx context.Context, cfg *rest.Config, options *config.Options) error {
	webhookServer := server.New(ctx, cfg, "webhook-example", options)
	admitters := []types.Admitter{
		admitter.NewPoolValidator().ToAdmitter(),
		admitter.NewPoolMutator(),
	}
	webhookServer.Register(admitters)
	if err := webhookServer.Start(); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func setLogLevel(level string) {
	ll, err := logrus.ParseLevel(level)
	if err != nil {
		ll = logrus.DebugLevel
	}
	// set global log level
	logrus.SetLevel(ll)
}
