package main

import (
	"flag"
	_ "fmt"
	"github.com/tohyung85/gophercises/exercise-12-file-renaming/filerenamer"
	"github.com/tohyung85/gophercises/exercise-12-file-renaming/renamerinterfaces/seqrenamer"
)

func main() {
	filePathPtr := flag.String("f", "./sample1", "Filepath to start renaming")
	flag.Parse()
	seqRenamer := &seqrenamer.SequenceRenamer{}
	filerenamer.RenameFiles(*filePathPtr, seqRenamer)
}
