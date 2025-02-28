package ostools

import "os"

func NewFile(dirName string) error {
	return os.Mkdir(dirName, 0755)
}
