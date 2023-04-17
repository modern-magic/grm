package shell

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

func MakeConfirm(msg string) bool {
	prompt := promptui.Prompt{
		Label:       fmt.Sprintf("%s [y/n]", msg),
		IsConfirm:   true,
		HideEntered: true,
	}
	result, err := prompt.Run()
	return err == nil && strings.ToLower(result) == "y"
}

func MakePrompt(msg string, validate func(input string) error) (string, error) {
	prompt := promptui.Prompt{
		Label:       msg,
		Validate:    validate,
		HideEntered: true,
	}
	result, err := prompt.Run()
	return result, err
}
