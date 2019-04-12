package history

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type ZshHistory string

func (h ZshHistory) GetHistory() []string {
	f, err := os.OpenFile(string(h), os.O_RDONLY, 0666)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	result := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tmp := strings.Split(scanner.Text(), ";")
		result = append(result, tmp[len(tmp)-1])
	}
	return result
}
