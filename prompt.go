package main

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

func (c *client) gpt(ctx context.Context, prompt string) (string, error) {
	res, err := c.oc.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("create chat completion: %w", err)
	}

	return res.Choices[0].Message.Content, nil
}
