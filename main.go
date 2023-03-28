/*
Copyright © 2023 Moye <dkyaorui@163.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dkyaorui/gpt-cli/cmd"
	"github.com/dkyaorui/gpt-cli/config"

	"github.com/spf13/viper"
)

func main() {

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// config file's path: $HOME/.gpt-cli/config.yml
	folderPath := filepath.Join(home, config.ConfigFileFolderName)
	viper.AddConfigPath(folderPath)
	viper.SetConfigName(config.ConfigFileFileName)
	viper.SetConfigType(config.ConfigFileType)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Using config file %s is not exist.", viper.ConfigFileUsed())
		fmt.Println("We are creating by default now.")

		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			if err := os.Mkdir(folderPath, 0755); err != nil {
				fmt.Println("Create config file folder fail.")
				fmt.Println(err)
				os.Exit(1)
			}
		}
		// TODO Set default config
		if err := viper.SafeWriteConfig(); err != nil {
			fmt.Println("Create defualt config file fail.")
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd.Execute()
}
