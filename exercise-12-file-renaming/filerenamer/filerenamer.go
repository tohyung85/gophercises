package filerenamer

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileRenamer interface {
	RenamePath(path string) (string, error)
	IsRenamable(path string) (bool, error)
}

type WithInit interface { //Optional Interface
	InitRenamer(path string) error
}

func RenameFiles(rootPath string, renamer FileRenamer) error {
	if renamerWithInit, ok := renamer.(WithInit); ok { //Initialize the renamer if required
		renamerWithInit.InitRenamer(rootPath)
	}
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path: %s\n", path)
			return nil
		}
		renamable, err := renamer.IsRenamable(path)
		if err != nil {
			return err
		}
		if !renamable {
			fmt.Printf("File/folder %s is not renamable.\n", path)
			return nil
		}
		newPath, err := renamer.RenamePath(path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = changeFileName(path, newPath)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func changeFileName(old string, new string) error {
	fmt.Printf("Renaming file from %s to %s", old, new)
	err := os.Rename(old, new)
	return err
}
