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

	"github.com/c-bata/go-prompt"
	"github.com/charmbracelet/glamour"
	"github.com/dkyaorui/gpt-cli/common"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

type Container struct {
	client *openai.Client

	messages []openai.ChatCompletionMessage
}

const token = ""

func (c *Container) InitContainer() {
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
	fmt.Println("Press 'Ctrl+D' to exit")

	p := prompt.New(
		c.executor,
		func(in prompt.Document) []prompt.Suggest {
			s := []prompt.Suggest{
				{Text: "reset", Description: "Reset the session"},
			}
			return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
		},
		prompt.OptionTitle("Sample Prompt"),
		prompt.OptionPrefix("Input: "),
		prompt.OptionInputTextColor(prompt.Blue),
		prompt.OptionPrefixTextColor(prompt.Green),
		prompt.OptionPreviewSuggestionTextColor(prompt.Green),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSuggestionBGColor(prompt.LightGray),
	)

	p.Run()
}

func (c *Container) addMessage(role, content string) {
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	})
}

func (c *Container) generateMessage(role, content string) []openai.ChatCompletionMessage {
	var requestMessages = c.messages
	requestMessages = append(requestMessages, openai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	})
	return requestMessages
}

func (c *Container) executor(in string) {
	var result string

	switch in {
	case "":
		return
	case "reset":
		c.resetMessages()
		result = "chat session cleared."
	default:
		var response, createErr = c.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo0301,
			Messages: c.generateMessage(common.ChatRoleUser, in),
		})

		if createErr != nil {
			fmt.Println(createErr)
			return
		}

		var message = response.Choices[0].Message

		result = message.Content

		c.addMessage(common.ChatRoleUser, in)
		c.addMessage(message.Role, result)

	}

	var r, newTermRendererErr = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if newTermRendererErr != nil {
		fmt.Println(newTermRendererErr)
		return
	}
	var renderResult, renderErr = r.Render(result)
	if renderErr != nil {
		fmt.Println(newTermRendererErr)
		return
	}
	fmt.Printf("%s\n", string(renderResult))

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
