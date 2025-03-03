package ostools

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func NewFile(filePath, projectName string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	out := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		textLine := scanner.Text()
		if strings.Contains(scanner.Text(), "{PROJECTNAME}") {
			textLine = strings.ReplaceAll(scanner.Text(), "{PROJECTNAME}", projectName)
		}
		out = out + textLine + "\n"
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return out, nil
}
