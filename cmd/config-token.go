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
	"github.com/fzdwx/infinite"
	"github.com/fzdwx/infinite/components/input/text"
	"github.com/fzdwx/infinite/theme"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var showTokenFlag bool

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Config openai token.",
	Long: `Config openai token.
	It will be used to call the OpenAI API.
	
	Demo:

	gpt-cli config token
	gpt-cli config token [-s, --show].`,
	Run: func(cmd *cobra.Command, args []string) {
		if showTokenFlag {
			token := viper.GetString(config.TokenConfigKey)
			if token == "" {
				fmt.Println("token not set")
				return
			}
			fmt.Printf("token=%s\n", viper.GetString(config.TokenConfigKey))
			return
		}

		inf := infinite.NewText(
			text.WithPrompt("Input your token: "),
			text.WithPromptStyle(theme.DefaultTheme.PromptStyle),
		)
		input, err := inf.Display()
		if err != nil {
			return
		}
		viper.Set(config.TokenConfigKey, input)
		if err := viper.WriteConfig(); err != nil {
			fmt.Println("set token fail")
			fmt.Println(err)
			return
		}
		fmt.Printf("set token success, token=%s\n", input)
	},
}

func init() {
	configCmd.AddCommand(tokenCmd)
	tokenCmd.PersistentFlags().BoolVarP(&showTokenFlag, "show", "s", false, "show current token")
}
