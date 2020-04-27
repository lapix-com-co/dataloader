package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/lapix-com-co/dataloader/pkg"
)

var provider *string
var output *string
var input = pkg.LoaderInput{}

func main() {
	args := flag.Args()
	if len(args) == 0 {
		// Default: proccess the package in the current directory.
		args = []string{"."}
	}

	output = flag.String("output", "", "output file name; default srcdir/<type>_dataloader.go")
	provider = flag.String("provider", "", "dataloader provider name; default mysql")

	flag.StringVar(&input.Type, "type", "", "type name; must be set")
	flag.StringVar(&input.TableName, "table", "", "database table name; default type snakecased")
	flag.StringVar(&input.OrderKey, "okey", "", "pagination sort key; default id")
	flag.Parse()
	log.SetPrefix("dataloader:")

	input.Provider = pkg.Provider(*provider)
	input.Pattern = args

	content, err := pkg.BuildLoader(input)
	if err != nil {
		log.Fatal(err)
	}

	if *output == "" {
		output = &args[0]
	}

	outputName := filepath.Join(*output, strings.ToLower(input.Type)+"_dataloader.go")
	if err := ioutil.WriteFile(outputName, content, 0644); err != nil {
		log.Fatal(err)
	}
}
