package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
)

func replace(dir string, from *regexp.Regexp, to string, dryrun bool) error {
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
			if err := replace(p, from, to, dryrun); err != nil {
				return err
			}
		} else {
			b, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}
			if dryrun {
				diff := difflib.UnifiedDiff{
					A:        difflib.SplitLines(string(b)),
					B:        difflib.SplitLines(from.ReplaceAllString(string(b), to)),
					FromFile: "Original",
					ToFile:   "Replacement",
					Context:  3,
				}
				text, _ := difflib.GetUnifiedDiffString(diff)
				fmt.Printf("%v\t%v\n", p, text)
			} else {
				if err := ioutil.WriteFile(p, from.ReplaceAll(b, []byte(to)), child.Mode()); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func main() {
	from := flag.String("from", "", "What to replace.")
	to := flag.String("to", "", "What to replace it with.")
	dryrun := flag.Bool("dryrun", false, "Don't perform any replacement, just output the ones that would have been made.")

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
	if err := replace(cwd, fromReg, *to, *dryrun); err != nil {
		panic(err)
	}
}
