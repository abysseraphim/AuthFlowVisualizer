package input

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func LocalHandler(path string) ([]string, error) {
	input, err := os.Stat(path)
	jsFiles := []string{}

	if os.IsNotExist(err) {
		fmt.Println("[e]Path Does Not Seem To Exist.")
		return nil, err

	}
	if err != nil {
		fmt.Println("[e]Error Occured:", err)
		return nil, err

	}

	if !input.IsDir() {
		fmt.Println("[e]You Have to Enter a Directory Path.")
		return nil, errors.New("[e]path is not a directory")
	}

	// walkDir takes two inputs: starting path, a callback function. go sais i walk the directory and whenever reaching a file or directory, ill call this function.
	// the callback function takes 3 parameters: file's full path, an object that gives file information (inteface), an error in case anything enexpected happens.
	// but why should callback return an error? consider it a termination signal. if you return nil, operation continures and if you return an error, func will be terminated.
	// and finally why do we assign all of this WalkDir's output to an error? because path might not exist at all, there is a permission error, callback returns an error.
	walkErr := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".js" {
			jsFiles = append(jsFiles, path)
		}
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}

	return jsFiles, nil
}
