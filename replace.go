package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func replace(dir string, from *regexp.Regexp, to string) error {
	children, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, child := range children {
		if strings.HasPrefix(child.Name(), ".") {
			continue
		}
		p := filepath.Join(dir, child.Name())
		if child.IsDir() {
			if err := replace(p, from, to); err != nil {
				return err
			}
		} else {
			b, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(p, from.ReplaceAll(b, []byte(to)), child.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	from := flag.String("from", "", "What to replace.")
	to := flag.String("to", "", "What to replace it with.")

	flag.Parse()

	if *from == "" {
		flag.Usage()
		os.Exit(1)
	}

	fromReg, err := regexp.Compile(*from)
	if err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := replace(cwd, fromReg, *to); err != nil {
		panic(err)
	}
}
