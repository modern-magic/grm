package shell

import (
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
