package main

import (
	"io/ioutil"
	"log"
	"reflect"

	"howett.net/plist"
)

const (
	bundleIDKey        = "CFBundleIdentifier"
	versionNumberKey   = "CFBundleShortVersionString"
	buildNumberKey     = "CFBundleVersion"
	bundleIconFilesKey = "CFBundleIconFiles"
)

type plistInfo struct {
	bundleID      string
	versionNumber string
	buildNumber   string
	bundleIcons   []string
}

func (e *plistInfo) UnmarshalPlist(unmarshal func(interface{}) error) error {

	var plistDict map[string]interface{}

	if err := unmarshal(&plistDict); err != nil {
		return err
	}

	bundleID := reflect.ValueOf(plistDict[bundleIDKey]).String()
	versionNumber := reflect.ValueOf(plistDict[versionNumberKey]).String()
	buildNumber := reflect.ValueOf(plistDict[buildNumberKey]).String()

	iconPathsValue := reflect.ValueOf(plistDict[bundleIconFilesKey])
	var iconPaths []string

	if iconPathsValue.Kind() == reflect.Slice {
		for i := 0; i < iconPathsValue.Len(); i++ {
			value := iconPathsValue.Index(i).Interface().(string)
			iconPaths = append(iconPaths, value)
		}
	}

	*e = plistInfo{
		bundleID:      bundleID,
		versionNumber: versionNumber,
		buildNumber:   buildNumber,
		bundleIcons:   iconPaths}

	return nil
}

func parsePlist(plistPath string) plistInfo {
	bytes, _ := ioutil.ReadFile(plistPath)

	var decoded plistInfo
	_, err := plist.Unmarshal(bytes, &decoded)
	if err != nil {
		log.Fatal(err)
	}

	return decoded
}
