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
	"context"
	"fmt"
	"time"

	"github.com/dkyaorui/gpt-cli/common"
	"github.com/dkyaorui/gpt-cli/config"

	"github.com/fzdwx/infinite"
	"github.com/fzdwx/infinite/components"
	"github.com/fzdwx/infinite/components/input/text"
	"github.com/fzdwx/infinite/components/spinner"
	"github.com/fzdwx/infinite/theme"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Container struct {
	client   *openai.Client
	messages []openai.ChatCompletionMessage
}

const (
	interruptedCmd = 0
	finishCmd      = 1
)

func (c *Container) InitContainer() {
	token := viper.GetString(config.TokenConfigKey)
	c.client = openai.NewClient(token)
	c.resetMessages()
}

func (c *Container) resetMessages() {
	c.messages = []openai.ChatCompletionMessage{
		{
			Role:    common.ChatRoleSystem,
			Content: "Show your response in Markdown format.",
		},
		{
			Role:    common.ChatRoleAssistant,
			Content: "Okay.",
		},
	}
}

func (c *Container) Run() {

	common.PrintToMarkdown(common.CliHelpInfo)

	for {
		inf := infinite.NewText(
			text.WithPrompt("[In]:"),
			text.WithPromptStyle(theme.DefaultTheme.PromptStyle),
			text.WithKeyMap(components.InputDefaultKeyMap()),
		)
		input, err := inf.Display()
		if err != nil {
			return
		}
		switch input {
		case "":
			continue
		case "-h":
			common.PrintToMarkdown(common.CliHelpInfo)
			continue
		case "-exit":
			fmt.Printf("\nBey.\n")
			return
		case "-reset":
			c.resetMessages()
			fmt.Printf("\nChat context cleared.\n\n")
		default:
			if err := c.request(input); err != nil {
				panic(err)
			}
		}
	}
}

func (c *Container) request(in string) error {
	return infinite.NewSpinner(
		spinner.WithShape(components.Dot),
	).Display(func(spinner *spinner.Spinner) {

		var total = int64(6000)
		var response openai.ChatCompletionResponse
		var err error
		var cmdCh = make(chan int)
		go func() {
			response, err = c.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo0301,
				Messages: c.generateMessage(common.ChatRoleUser, in),
			})
			cmdCh <- finishCmd
		}()

		for i := int64(0); i < total-1; i++ {
			select {
			case cmd := <-cmdCh:
				switch cmd {
				case finishCmd:

					spinner.Finish("Done")

					close(cmdCh)

					if err != nil {
						fmt.Println(err)
						return
					}

					message := response.Choices[0].Message
					result := message.Content

					c.addMessage(common.ChatRoleUser, in)
					c.addMessage(message.Role, result)

					common.PrintToMarkdown(result)
					return
				}
			default:
				spinner.Refresh("Wait...")
				time.Sleep(time.Millisecond * 50)
			}
		}

		spinner.Failed("time out, please retry")
	})

}

func (c *Container) addMessage(role, content string) {
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	})
}

func (c *Container) generateMessage(role, content string) []openai.ChatCompletionMessage {
	requestMessages := c.messages
	requestMessages = append(requestMessages, openai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	})
	return requestMessages
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run GPT CLI program",
	Long:  `Run GPT CLI program, start chat with gpt. Please config your key first`,
	Run: func(cmd *cobra.Command, args []string) {
		var container = Container{}
		container.InitContainer()
		container.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
