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
package cmd

import (
	"fmt"

	"github.com/dkyaorui/gpt-cli/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "your openai token",
	Long: `Your openai token.
	It will be used to call the OpenAI API.
	
	gpt-cli config token
	gpt-cli config token your_token.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			token := viper.GetString(config.TokenConfigKey)
			if token == "" {
				fmt.Println("token not set")
				return
			}
			fmt.Printf("token=%s\n", viper.GetString(config.TokenConfigKey))
			return
		}
		token := args[0]
		viper.Set(config.TokenConfigKey, token)
		if err := viper.WriteConfig(); err != nil {
			fmt.Println("set token fail")
			fmt.Println(err)
			return
		}
		fmt.Printf("set token success, token=%s\n", token)
	},
}

func init() {
	configCmd.AddCommand(tokenCmd)
}
