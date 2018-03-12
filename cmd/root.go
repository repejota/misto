// Copyright 2018 The misto Authors. All rights reserved.
//
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with this
// work for additional information regarding copyright ownership.  The ASF
// licenses this file to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
// License for the specific language governing permissions and limitations
// under the License.

package cmd

import (
	"fmt"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/repejota/misto"
)

var (
	inputFlag   string
	outputFlag  string
	verboseFlag bool
	versionFlag bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "misto",
	Short: "Tail aggregated logs",
	Long:  `Misto tails and aggregates logs`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.FatalLevel)

		formatter := &log.TextFormatter{
			FullTimestamp: true,
		}
		log.SetFormatter(formatter)

		// --verbose
		if verboseFlag {
			log.SetLevel(log.DebugLevel)
		}

		m := misto.NewMisto()

		// --version
		if versionFlag {
			versionInformation := m.ShowVersion()
			fmt.Println(versionInformation)
			os.Exit(2)
		}

		// graceful shutdown signal
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt)

		m.Start()

		<-shutdown
		fmt.Println()
		m.Stop()

	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.SetUsageFunc(UsageFunc)
	RootCmd.Flags().StringVarP(&inputFlag, "input", "i", "", "set an input")
	RootCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "set an output")
	RootCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "enable verbose mode")
	RootCmd.Flags().BoolVarP(&versionFlag, "version", "V", false, "show version number")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Unimplemented
}

// UsageFunc prints command usage help message
func UsageFunc(cmd *cobra.Command) error {
	fmt.Println("Usage:")
	fmt.Println("  misto [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -i, --input=INPUT	set an input")
	fmt.Println("  -o, --output=OUTPUT 	set an output")
	fmt.Println("  -h, --help		help for misto")
	fmt.Println("  -v, --verbose   	enable verbose mode")
	fmt.Println("  -V, --version   	show version number")
	fmt.Println()
	fmt.Println("Inputs:")
	fmt.Println("  dummy    		Generate dummy data")
	fmt.Println()
	fmt.Println("Outputs:")
	fmt.Println("  stdout    		Prints log events to STDOUT")
	return nil
}
