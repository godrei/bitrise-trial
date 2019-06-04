package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	contentsFolder = "contents"
	payloadFolder  = contentsFolder + "/Payload"
	plistFilePath  = "/Info.plist"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("No .ipa file specified. Usage: `go run bitrise-trial path/to/file.ipa`")
	}
	filePath := args[1]

	unzip(filePath, contentsFolder)
	plistPath := plistPath(payloadFolder)

	result := parsePlist(plistPath)
	printResults(result)

	os.RemoveAll(contentsFolder)
}

func printResults(result plistInfo) {
	fmt.Println("Bundle Identifier:" + result.bundleID)
	fmt.Println("Version number:" + result.versionNumber)
	fmt.Println("Build number:" + result.buildNumber)
	if len(result.bundleIcons) > 0 {
		fmt.Println("Icon files:")
		for _, path := range result.bundleIcons {
			fmt.Println("- " + path)
		}
	} else {
		fmt.Println("No icon files specified")
	}
}

func plistPath(bundleFolderPath string) string {
	files, err := ioutil.ReadDir(bundleFolderPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		log.Fatal("No files found in bundle")
	}

	appName := files[0].Name()

	absFilePath, err := filepath.Abs(bundleFolderPath + "/" + appName)
	if err != nil {
		log.Fatal(err)
	}

	absPlistPath := absFilePath + plistFilePath

	return absPlistPath
}

func unzip(src string, dest string) {
	r, err := zip.OpenReader(src)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			log.Fatal(err)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Fatal(err)
		}

		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			log.Fatal(err)
		}
	}
}
