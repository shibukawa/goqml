package goqmllib

import (
	"bytes"
	"github.com/Kodeworks/golang-image-ico"
	"github.com/fatih/color"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Build() {
	os.MkdirAll("Resources", 0755)
	os.MkdirAll("_build", 0755)

	CreateIcon()
}

func CreateIcon() {
	target := Target()

	var iconImage image.Image
	resourceDir := "Resources"
	os.MkdirAll(resourceDir, 0755)

	workDir := "_build"

	resultPath1 := filepath.Join(resourceDir, "WindowsAppIcon.ico")
	resultPath2 := filepath.Join(resourceDir, "MacOSXAppIcon.icns")
	pngPath := filepath.Join("Resources", "icon.png")

	file, err := os.Open(pngPath)
	var oldestUnixTime int64

	if target == "windows" {
		oldestUnixTime = unixTime(resultPath1)
	} else if target == "darwin" {
		oldestUnixTime = unixTime(resultPath2)
	} else {
		return
	}

	if err == nil {
		inputStat, err := file.Stat()
		if err == nil && oldestUnixTime > inputStat.ModTime().Unix() {
			return
		}
		iconImage, _, err = image.Decode(file)
		if err != nil {
			iconImage = nil
		}
	} else if oldestUnixTime > 0 {
		return

	}
	sourceFile := filepath.Join("Resources", "icon.png")
	if iconImage == nil {
		sourceFile = "default image"
		pngPath = filepath.Join(resourceDir, "icon.png")
		ioutil.WriteFile(pngPath, MustAsset("resources/qt-logo.png"), 0644)
		iconImage, _, err = image.Decode(bytes.NewReader(MustAsset("resources/qt-logo.png")))
		if err != nil {
			panic(err) // qtpm should be able to read fallback icon anytime
		}
	}

	if target == "windows" {
		icon, err := os.Create(resultPath1)
		defer icon.Close()
		if err != nil {
			panic(err)
		}
		ico.Encode(icon, iconImage)
		color.Magenta("Wrote: %s from %s\n", filepath.Join("Resources", "WindowsAppIcon.ico"), sourceFile)
	} else {
		os.MkdirAll(filepath.Join(workDir, "MacOSXAppIcon.iconset"), 0755)
		err = SequentialRun(resourceDir).
			Run("sips", "-z", "16", "16", pngPath, "--out", filepath.Join(workDir, "MacOSXAppIcon.iconset", "icon_16x16.png")).
			Run("sips", "-z", "32", "32", pngPath, "--out", filepath.Join(workDir, "MacOSXAppIcon.iconset", "icon_16x16@2x.png")).
			Run("sips", "-z", "32", "32", pngPath, "--out", filepath.Join(workDir, "MacOSXAppIcon.iconset", "icon_32x32.png")).
			Run("sips", "-z", "64", "64", pngPath, "--out", filepath.Join(workDir, "MacOSXAppIcon.iconset", "icon_32x32@2x.png")).
			Run("sips", "-z", "128", "128", pngPath, "--out", filepath.Join(workDir, "MacOSXAppIcon.iconset", "icon_128x128.png")).
			Run("sips", "-z", "256", "256", pngPath, "--out", filepath.Join(workDir, "MacOSXAppIcon.iconset", "icon_128x128@2x.png")).
			Run("sips", "-z", "256", "256", pngPath, "--out", filepath.Join(workDir, "MacOSXAppIcon.iconset", "icon_256x256.png")).
			Run("sips", "-z", "512", "512", pngPath, "--out", filepath.Join(workDir, "MacOSXAppIcon.iconset", "icon_256x256@2x.png")).
			Run("sips", "-z", "512", "512", pngPath, "--out", filepath.Join(workDir, "MacOSXAppIcon.iconset", "icon_512x512.png")).
			Run("sips", "-z", "1024", "1024", pngPath, "--out", filepath.Join(workDir, "MacOSXAppIcon.iconset", "icon_512x512@2x.png")).
			Run("iconutil", "-c", "icns", "--output", resultPath2, filepath.Join(workDir, "MacOSXAppIcon.iconset")).Finish()
		if err != nil {
			panic(err)
		}
		color.Magenta("Wrote: %s from %s\n", filepath.Join("Resources", "MacOSXAppIcon.icns"), sourceFile)
	}
}
