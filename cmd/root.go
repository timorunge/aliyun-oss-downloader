// Copyright © 2019 Timo Runge <me@timorunge.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package cmd

import (
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/timorunge/aliyun-oss-downloader/app"
)

const (
	defaultCreateDestinationDir bool   = false
	defaultRegion               string = "eu-central-1"
	defaultThreads              int    = 5
	defaultMaxKeys              int    = 250
)

var (
	accessKeyID          string
	accessKeySecret      string
	bucket               string
	createDestinationDir bool = defaultCreateDestinationDir
	cfgFile              string
	destinationDir       string
	marker               string
	maxKeys              int    = defaultMaxKeys
	region               string = defaultRegion
	threads              int    = defaultThreads
)

var rootCmd = &cobra.Command{
	Use:   "aliyun-oss-downloader",
	Short: "Downlaods all objects from an Aliyun OSS bucket.",
	Long: `aliyun-oss-downloader downloads all objects from an Aliyun OSS bucket to a defined local destination.
It copies new and modified files.`,
	Run: func(cmd *cobra.Command, args []string) { app.Download() },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&accessKeyID, "accessKeyID", accessKeyID, "Your access key ID")
	viper.BindPFlag("accessKeyID", rootCmd.PersistentFlags().Lookup("accessKeyID"))
	rootCmd.MarkFlagRequired("accessKeyID")

	rootCmd.PersistentFlags().StringVar(&accessKeySecret, "accessKeySecret", accessKeySecret, "Your access key secret")
	viper.BindPFlag("accessKeySecret", rootCmd.PersistentFlags().Lookup("accessKeySecret"))
	rootCmd.MarkFlagRequired("accessKeySecret")

	rootCmd.PersistentFlags().StringVarP(&bucket, "bucket", "b", bucket, "The name of the OSS bucket which should be downloaded")
	viper.BindPFlag("bucket", rootCmd.PersistentFlags().Lookup("bucket"))
	rootCmd.MarkFlagRequired("bucket")

	rootCmd.PersistentFlags().BoolVarP(&createDestinationDir, "createDestinationDir", "", createDestinationDir, "Create the (local) destination directory if not existing")
	viper.BindPFlag("createDestinationDir", rootCmd.PersistentFlags().Lookup("createDestinationDir"))

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "Config file (default \"$HOME/.aliyun-oss-downloader.yaml\")")

	rootCmd.PersistentFlags().StringVarP(&destinationDir, "destinationDir", "", destinationDir, "The (local) destination directory")
	viper.BindPFlag("destinationDir", rootCmd.PersistentFlags().Lookup("destinationDir"))
	rootCmd.MarkFlagRequired("destinationDir")

	rootCmd.PersistentFlags().StringVarP(&marker, "marker", "", marker, "The marker to start the download")
	viper.BindPFlag("marker", rootCmd.PersistentFlags().Lookup("marker"))

	rootCmd.PersistentFlags().IntVarP(&maxKeys, "maxKeys", "", maxKeys, "The amount of objects which are fetched in a single request")
	viper.BindPFlag("maxKeys", rootCmd.PersistentFlags().Lookup("maxKeys"))

	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", region, "The name of the OSS region in which you have stored your bucket")
	viper.BindPFlag("region", rootCmd.PersistentFlags().Lookup("region"))
	rootCmd.MarkFlagRequired("region")

	rootCmd.PersistentFlags().IntVarP(&threads, "threads", "", threads, "The amount of threads to use")
	viper.BindPFlag("threads", rootCmd.PersistentFlags().Lookup("threads"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".aliyun-oss-downloader")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error while loading config file: %s", viper.ConfigFileUsed())
	}
}
