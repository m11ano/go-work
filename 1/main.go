package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {

	return dirTreeEngine(out, path, printFiles, nil)
}

func dirTreeEngine(out io.Writer, path string, printFiles bool, stackParents []bool) error {

	dirItems, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	lastItemIndex := len(dirItems) - 1
	if len(dirItems) > 0 && !printFiles {
		for i := len(dirItems) - 1; i >= 0; i-- {
			if dirItems[i].IsDir() {
				lastItemIndex = i
				break
			}
		}
	}

	for i, item := range dirItems {
		isLastItem := i == lastItemIndex
		itemName := item.Name()
		line := ""
		for _, p := range stackParents {
			if p {
				line += "\t"
			} else {
				line += "│\t"
			}
		}

		if isLastItem {
			line += "└───"
		} else {
			line += "├───"
		}

		if item.IsDir() {
			line += itemName
			_, _ = fmt.Fprintln(out, line)
			childStack := slices.Concat(stackParents, []bool{isLastItem})
			_ = dirTreeEngine(out, filepath.Join(path, itemName), printFiles, childStack)
		} else if printFiles {
			info, err := item.Info()
			if err != nil {
				continue
			}
			line += itemName
			sizeTextValue := "empty"
			if info.Size() > 0 {
				sizeTextValue = strconv.FormatInt(info.Size(), 10) + "b"
			}
			line += fmt.Sprintf(" (%s)", sizeTextValue)
			_, _ = fmt.Fprintln(out, line)
		}
	}

	return nil
}
