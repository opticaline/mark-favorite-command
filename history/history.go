package history

import (
	"os"
	"strings"
)

type History interface {
	GetHistory() []string
}

func OpenHistory() (History, error) {
	var history History
	if home, ok := os.LookupEnv("HOME"); ok {
		if shell, ok := os.LookupEnv("SHELL"); ok {
			switch {
			case strings.Contains(shell, "zsh"):
				history = ZshHistory(home + "/.zsh_history")
			case strings.Contains(shell, "bash"):
				history = BashHistory(home + "/.bash_history")
			}
		}
	}
	return history, nil
}
