package shell

import (
	"errors"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

func MakeConfirm(msg string) bool {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("%s [y/n]", msg),
		IsConfirm: true,
	}
	result, err := prompt.Run()
	return err == nil && strings.ToLower(result) == "y"
}

func MakePrompt(msg string) string {
	prompt := promptui.Prompt{
		Label: msg,
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("Can't be empty")
			}
			return nil
		},
	}
	result, _ := prompt.Run()
	return result
}
