package main

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/tliron/commonlog"
	cobrautil "github.com/tliron/kutil/cobra"
	"github.com/tliron/kutil/util"
	"k8s.io/klog/v2"
)

var logTo string
var verbose int
var colorize string

var masterUrl string
var kubeconfigPath string
var context string

var version bool
var cluster bool
var namespace string
var concurrency uint
var resyncPeriod time.Duration
var cachePath string
var healthPort uint

func init() {
	command.Flags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	command.Flags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	command.Flags().StringVarP(&colorize, "colorize", "z", "true", "colorize output (boolean or \"force\")")

	// Conventional flags for Kubernetes controllers
	command.Flags().StringVar(&masterUrl, "master", "", "address of Kubernetes API server")
	command.Flags().StringVar(&kubeconfigPath, "kubeconfig", "", "path to Kubernetes configuration")
	command.Flags().StringVarP(&context, "context", "x", "", "name of context in Kubernetes configuration")

	// Our additional flags
	command.Flags().BoolVar(&version, "version", false, "print version")
	command.Flags().BoolVar(&cluster, "cluster", false, "enable cluster mode")
	command.Flags().StringVar(&namespace, "namespace", "", "namespace (overrides context namespace in Kubernetes configuration)")
	command.Flags().UintVar(&concurrency, "concurrency", 1, "number of concurrent workers per processor")
	command.Flags().DurationVar(&resyncPeriod, "resync", time.Second*30, "informer resync period")
	command.Flags().StringVar(&cachePath, "cache", "", "cache path")
	command.Flags().UintVar(&healthPort, "health-port", 8086, "HTTP port for health check (for liveness and readiness probes)")

	cobrautil.SetFlagsFromEnvironment("CNCK_OPERATOR_", command)
}

var command = &cobra.Command{
	Use:   toolName,
	Short: "Start the CNCK operator",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		util.InitializeColorization(colorize)
		commonlog.Initialize(verbose, logTo)
		if writer := commonlog.GetWriter(); writer != nil {
			klog.SetOutput(writer)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		Controller()
	},
}

func Execute() {
	err := command.Execute()
	util.FailOnError(err)
}
