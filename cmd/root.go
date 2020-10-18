package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	bucket  string
	prefix  string
	region  string
	dryrun  bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "s3-restore",
	Short: "Restores deleted objects of an S3 version enabled bucket.",
	Long:  `Restores deleted objects of an S3 version enabled bucket.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.s3-restore.yaml)")
	rootCmd.PersistentFlags().StringVar(&bucket, "bucket", "", "S3 bucket name to restore the objects from")
	rootCmd.PersistentFlags().StringVar(&prefix, "prefix", "", "S3 prefix to look for delete objects")
	rootCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "AWS region where S3 bucket is located")
	rootCmd.PersistentFlags().BoolVar(&dryrun, "dryrun", false, "Whether to just to print the actions being performed or to execute them")

	rootCmd.MarkPersistentFlagRequired("bucket")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".s3-restore" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".s3-restore")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
