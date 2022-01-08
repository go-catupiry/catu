package helpers

import (
	"bytes"
	"io/ioutil"
	"os"
)

func CleanCSVFile(filePath string) error {
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	outFile := filePath + "_out"
	output := bytes.Replace(input, []byte(`"`), []byte(""), -1)

	if err = ioutil.WriteFile(outFile, output, 0666); err != nil {
		return err
	}

	err = os.Rename(outFile, filePath)
	if err != nil {
		return err
	}

	return nil
}
