package utils

import (
	"io/ioutil"
)

func ReadFile(path string) string {

	ctx, _ := ioutil.ReadFile(path)

	return string(ctx)
}

func WriteFile(content, path string) error {

	//err := ioutil.WriteFile(path, []byte(content), os.ModePerm)
	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
