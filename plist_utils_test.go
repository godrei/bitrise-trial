package main

import "testing"

func TestBundleIDParsing(t *testing.T) {
	result := parsePlist("./sample/Info.plist").bundleID
	expected := "hu.jozsefvesza.example"
	if result != expected {
		t.Errorf("Build number is %s, wanted %s", result, expected)
	}
}

func TestVersionNumberParsing(t *testing.T) {
	result := parsePlist("./sample/Info.plist").versionNumber
	expected := "1.0"
	if result != expected {
		t.Errorf("Build number is %s, wanted %s", result, expected)
	}
}

func TestBuildNumberParsing(t *testing.T) {
	result := parsePlist("./sample/Info.plist").buildNumber
	expected := "1.0"
	if result != expected {
		t.Errorf("Build number is %s, wanted %s", result, expected)
	}
}

func TestBundleIconsParsing(t *testing.T) {
	result := parsePlist("./sample/Info.plist").bundleIcons
	expected := []string{"Default.png", "Default@2x.png"}
	for i := 0; i < 2; i++ {
		if result[i] != expected[i] {
			t.Errorf("Bundle icon path is %s, wanted %s", result[i], expected[i])
		}
	}
}
