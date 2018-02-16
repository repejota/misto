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

	"github.com/repejota/misto"
	"github.com/spf13/cobra"
)

var (
	verboseFlag bool
	versionFlag bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "misto",
	Short: "Tail docker logs",
	Long:  `Misto tails logs from a docker daemon`,
	Run: func(cmd *cobra.Command, args []string) {

		if versionFlag {
			showVersion()
		}

		hub, err := misto.NewHub()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		hub.ListenAndServe()
	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "enable verbose mode")
	RootCmd.Flags().BoolVarP(&versionFlag, "version", "V", false, "show version number")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}

// snowVersion shows the program build and version information.
func showVersion() {
	// TODO:
	// Show the real version information
	fmt.Println("misto v.0.0.0-0498275295")
	os.Exit(2)
}
