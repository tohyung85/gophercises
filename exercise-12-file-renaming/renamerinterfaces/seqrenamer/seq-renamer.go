package seqrenamer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type SequenceRenamer struct {
	renameMap map[string]int
}

func (b *SequenceRenamer) InitRenamer(path string) error {
	renameMap, err := b.createRenameMap(path)
	if err != nil {
		return err
	}
	b.renameMap = renameMap

	return nil
}

func (b *SequenceRenamer) IsRenamable(path string) (bool, error) {
	re, err := regexp.Compile(`_[0-9]+.`)
	if err != nil {
		return false, err
	}
	matched := re.Match([]byte(path))
	if !matched {
		return false, nil
	}
	return true, nil
}

func (b *SequenceRenamer) RenamePath(path string) (string, error) {
	basePath := getBasePath(path)
	seq := getSeqNumber(path)
	maxSeq, inMap := b.renameMap[basePath]
	if !inMap {
		return path, fmt.Errorf("Error: BasePath not in rename map! Something wrong with initialization")
	}
	newPath := fmt.Sprintf("%s (%d of %d)%s\n", basePath, seq, maxSeq, filepath.Ext(path))
	return newPath, nil
}

func (b *SequenceRenamer) createRenameMap(path string) (map[string]int, error) {
	renameMap := make(map[string]int)

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Unable to access folder/file: %s\n", path)
			return nil
		}
		renamable, err := b.IsRenamable(path)
		if renamable {
			baseName := getBasePath(path)
			seqNumber := getSeqNumber(path)
			maxSeq, inMap := renameMap[baseName]
			if !inMap || maxSeq < seqNumber {
				renameMap[baseName] = seqNumber
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return renameMap, nil
}

func getBasePath(fullPath string) string {
	lastUnderScoreIdx := strings.LastIndex(fullPath, "_")
	return fullPath[:lastUnderScoreIdx]
}

func getSeqNumber(fullPath string) int {
	re := regexp.MustCompile(`_[0-9]+.`)
	reDigitsOnly := regexp.MustCompile(`[0-9]+`)
	numString := reDigitsOnly.Find(re.Find([]byte(fullPath)))
	num, err := strconv.Atoi(string(numString))
	if err != nil {
		log.Fatal("something wrong with string to int conversion")
	}
	return num
}
