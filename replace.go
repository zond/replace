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

func replacePath(path string, from *regexp.Regexp, to string, dryrun bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	replacement := from.ReplaceAll(b, []byte(to))
	if dryrun {
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(string(b)),
			B:        difflib.SplitLines(string(replacement)),
			FromFile: "Original",
			ToFile:   "Replacement",
			Context:  3,
		}
		text, err := difflib.GetUnifiedDiffString(diff)
		if err != nil {
			return err
		}
		if strings.TrimSpace(text) != "" {
			fmt.Printf("*** %v ***\n%v\n", path, text)
		}
	} else {
		if err := ioutil.WriteFile(path, replacement, info.Mode()); err != nil {
			return err
		}
	}
	return nil
}

func replaceDir(dir string, from *regexp.Regexp, to string, dryrun bool) error {
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
			if err := replaceDir(p, from, to, dryrun); err != nil {
				return err
			}
		} else {
			replacePath(p, from, to, dryrun)
		}
	}
	return nil
}

func replaceGlob(glob string, from *regexp.Regexp, to string, dryrun bool) error {
	matches, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	for _, match := range matches {
		if err := replacePath(match, from, to, dryrun); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	from := flag.String("from", "", "What to replace.")
	to := flag.String("to", "", "What to replace it with.")
	dryrun := flag.Bool("dryrun", false, "Don't perform any replacement, just output the ones that would have been made.")
	glob := flag.String("glob", "", "Which files to replace. If empty, all files are recursively selected.")

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
	if *glob == "" {
		if err := replaceDir(cwd, fromReg, *to, *dryrun); err != nil {
			panic(err)
		}
	} else {
		if err := replaceGlob(*glob, fromReg, *to, *dryrun); err != nil {
			panic(err)
		}
	}
}
