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

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dkyaorui/gpt-cli/config"
	"github.com/fzdwx/infinite"
	"github.com/fzdwx/infinite/color"
	"github.com/fzdwx/infinite/components"
	"github.com/fzdwx/infinite/components/selection/singleselect"
	"github.com/fzdwx/infinite/style"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var showModelFlag bool

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Config openai's model.",
	Long: `Config openai's model.
	
	Demo:

	gpt-cli config model
	gpt-cli config model [-s, --show].`,
	Run: func(cmd *cobra.Command, args []string) {
		if showModelFlag {
			currentModel := viper.GetString(config.ModelConfigKey)
			fmt.Printf("model=%s\n", currentModel)
			return
		}

		input := components.NewInput()
		input.Prompt = "Filtering: "
		input.PromptStyle = style.New().Bold().Italic().Fg(color.LightBlue)

		choiceIndex, err := infinite.NewSingleSelect(
			config.ModelList,
			singleselect.WithKeyBinding(
				singleselect.KeyMap{
					Up: key.NewBinding(
						key.WithKeys("up"),
						key.WithHelp("↑", "move up"),
					),
					Down: key.NewBinding(
						key.WithKeys("down"),
						key.WithHelp("↓", "move down"),
					),
					Choice: key.NewBinding(
						key.WithKeys("enter"),
						key.WithHelp("enter", "choice and finish selection"),
					),
					Confirm: key.NewBinding(
						key.WithKeys("enter"),
						key.WithHelp("enter", "choice and finish selection"),
					),
					NextPage: key.NewBinding(
						key.WithKeys(tea.KeyPgDown.String()),
						key.WithHelp("pageup", "next page"),
					),
					PrevPage: key.NewBinding(
						key.WithKeys(tea.KeyPgUp.String()),
						key.WithHelp("pagedown", "prev page"),
					),
					Quit: key.NewBinding(
						key.WithKeys("ctrl+c"),
						key.WithHelp("^C", "kill program"),
					),
				},
			),
			singleselect.WithFilterInput(input),
		).Display("select model: ")
		if err != nil {
			fmt.Println(err)
			return
		}

		model := config.ModelList[choiceIndex]
		viper.Set(config.ModelConfigKey, model)
		if err := viper.WriteConfig(); err != nil {
			fmt.Println("set model fail")
			fmt.Println(err)
			return
		}
		fmt.Printf("set model success, model=%s\n", model)
	},
}

func init() {
	configCmd.AddCommand(modelCmd)
	modelCmd.PersistentFlags().BoolVarP(&showModelFlag, "show", "s", false, "show current model")
}
