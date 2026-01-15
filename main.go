package main

import (
	"flag"
	"fmt"
	"go-zipfile/zipfile"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	useDeflate := flag.Bool("d", false, "compress archive with deflate algorithm")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "tool needs at least two arguments")
		flag.Usage()
		os.Exit(2)
	}
	out := args[0]

	zip := zipfile.NewZip()
	if *useDeflate {
		zip.SetCompressionMethod(zipfile.CompressionMethodDeflated)
	}

	for _, arg := range args[1:] {
		if err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
			if path == "." {
				return nil
			}
			if info.IsDir() && !strings.HasSuffix(path, "\\") {
				path = path + "\\"
			}
			fmt.Println(path)
			return zip.Add(path)
		}); err != nil {
			panic(err)
		}
	}

	f, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	if err = zip.Marshal(f); err != nil {
		panic(err)
	}
}
