/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Logrus hooks

// Hook for erroring and exit out on warning
type errorOnWarningHook struct{}

func (errorOnWarningHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.WarnLevel}
}

func (errorOnWarningHook) Fire(entry *logrus.Entry) error {
	logrus.Fatalln(entry.Message)
	return nil
}

var (
	GlobalBundle, GlobalFile, GlobalProvider                    string
	GlobalVerbose, GlobalSuppressWarnings, GlobalErrorOnWarning bool
)

var RootCmd = &cobra.Command{
	Use:   "kompose",
	Short: "A tool helping Docker Compose users move to Kubernetes",
	Long:  `Kompose is a tool to help users who are familiar with docker-compose move to Kubernetes.`,
	// PersitentPreRun will be "inherited" by all children and ran before *every* command unless
	// the child has overridden the functionality. This functionality was implemented to check / modify
	// all global flag calls regardless of app call.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		// Add extra logging when verbosity is passed
		if GlobalVerbose {
			logrus.SetLevel(logrus.DebugLevel)
		}

		// Set the appropriate suppress warnings and error on warning flags
		if GlobalSuppressWarnings {
			logrus.SetLevel(logrus.ErrorLevel)
		} else if GlobalErrorOnWarning {
			hook := errorOnWarningHook{}
			logrus.AddHook(hook)
		}

		// Error out of the user has not chosen Kubernetes or OpenShift
		provider := strings.ToLower(GlobalProvider)
		if provider != "kubernetes" && provider != "openshift" {
			logrus.Fatalf("%s is an unsupported provider. Supported providers are: 'kubernetes', 'openshift'.", GlobalProvider)
		}

	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&GlobalVerbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().BoolVar(&GlobalSuppressWarnings, "suppress-warnings", false, "Suppress all warnings")
	RootCmd.PersistentFlags().BoolVar(&GlobalErrorOnWarning, "error-on-warning", false, "Treat any warning as an error")
	RootCmd.PersistentFlags().StringVarP(&GlobalFile, "file", "f", "docker-compose.yml", "Specify an alternative compose file")
	RootCmd.PersistentFlags().StringVarP(&GlobalBundle, "bundle", "b", "", "Specify a Distributed Application GlobalBundle (DAB) file")
	RootCmd.PersistentFlags().StringVar(&GlobalProvider, "provider", "kubernetes", "Specify a provider. Kubernetes or OpenShift.")
}