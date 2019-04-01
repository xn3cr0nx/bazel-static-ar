package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, dir, output string

var bazelarCmd = &cobra.Command{
	Use:   "bazelar",
	Short: "Tools to archive together static library files produced by bazel",
	Long: `Bazelar
parses bazel-bin folder to extract all .a static library files and archive them up together
	`,
	// PersistentPreRun: func(cmd *cobra.Command, args []string) {	},
	Run: func(cmd *cobra.Command, args []string) {
		dir := viper.GetString("dir")
		log.Printf("Looking for bazel-bin in directory %s", dir)
		content, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatalf(err.Error())
		}
		var bazelBin os.FileInfo
		for _, f := range content {
			if f.Name() == "bazel-bin" {
				bazelBin = f
			}
		}
		if bazelBin == nil {
			log.Fatalf("bazel-bin folder not found in directory %s", dir)
		}

		time.Sleep(1 * time.Second)
		log.Print("Extracting files from bazel-bin...")
		time.Sleep(1 * time.Second)

		staticFiles, err := extractStaticFiles(filepath.Join(dir, bazelBin.Name()))
		if err != nil {
			log.Fatalf(err.Error())
		}
		files := strings.Join(staticFiles, " ")

		log.Print("Merging the archive")
		time.Sleep(1 * time.Second)

		command := strings.Fields(fmt.Sprintf("ar rsv %s %s", viper.GetString("out"), files))
		fmt.Println("what cmd?", command)
		output, err := exec.Command(command[0], command[1:]...).Output()
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("%s", output)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the bazelarCmd.
func Execute() {
	if err := bazelarCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// // Adds subdirectories command
	// bazelarCmd.AddCommand(block.BlockCmd)

	// Adds root flags and persistent flags
	// bazelarCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Sets logging level to Debug")
	bazelarCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	bazelarCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bazelar.yaml)")
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf(fmt.Sprintf("Bazelar %v", err))
	}
	bazelarCmd.PersistentFlags().StringVarP(&dir, "dir", "d", wd, "Sets the path to the directory containing bazel-bin folder")
	viper.SetDefault("blocksDir", wd)
	bazelarCmd.PersistentFlags().StringVarP(&output, "out", "o", fmt.Sprintf("%s/%s", wd, "output.a"), "Sets the path to the output where the .a file will be saved")
	viper.SetDefault("out", fmt.Sprintf("%s/%s", wd, "output.a"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// viper.SetDefault("debug", false)

	viper.BindPFlag("dir", bazelarCmd.PersistentFlags().Lookup("dir"))
	viper.BindPFlag("out", bazelarCmd.PersistentFlags().Lookup("out"))

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
		// Search config in home directory with name ".bazelar" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bazelar")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func extractStaticFiles(dir string) ([]string, error) {
	content, err := ioutil.ReadDir(dir)
	var files []string
	if err != nil {
		return nil, err
	}

	for _, c := range content {
		if c.IsDir() {
			subFiles, err := extractStaticFiles(filepath.Join(dir, c.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
			// } else if strings.Contains(c.Name(), ".a") {
		} else if strings.Contains(c.Name(), ".o") {
			// re := regexp.MustCompile("\\w+.a$")
			re := regexp.MustCompile("\\w+.o$")
			match := re.FindString(c.Name())
			if strings.Compare(match, c.Name()) == 0 {
				fmt.Println(c.Name())
				time.Sleep(100 * time.Millisecond)
				files = append(files, fmt.Sprintf("%s/%s", dir, c.Name()))
			}
		}
	}

	return files, nil
}
