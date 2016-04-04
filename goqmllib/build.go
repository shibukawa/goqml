package goqmllib

import (
	"bytes"
	"fmt"
	"github.com/Kodeworks/golang-image-ico"
	"github.com/fatih/color"
	"github.com/shibukawa/configdir"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"github.com/k0kubun/pp"
)

func Build() {
	os.MkdirAll("Resources", 0755)
	os.MkdirAll("_build", 0755)

	buildApplication()
}

func createIcon(outputDir string) {
	target := Target()

	var iconImage image.Image
	resourceDir := "Resources"
	os.MkdirAll(resourceDir, 0755)

	workDir := "_build"

	resultPath1 := filepath.Join(outputDir, "WindowsAppIcon.ico")
	resultPath2 := filepath.Join(outputDir, "MacOSXAppIcon.icns")
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
	println("@3")

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
		err = SequentialRun(".").
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

func buildApplication() {
	qtdir := FindQt()

	if qtdir == "" {
		color.Red("\nCan't find Qt. if you don't use standard install location, set QTDIR environment variable.\n")
		os.Exit(1)
	}

	if runtime.GOOS == "darwin" && strings.HasSuffix(qtdir, "clang_64") {
		config, _ := LoadConfig()

		variable := map[string]string{
			"Target":       config.Name,
			"TargetSmall":  strings.ToLower(config.Name),
			"Version":      fmt.Sprintf("%d.%d.%d", config.Version[0], config.Version[1], config.Version[2]),
			"ShortVersion": fmt.Sprintf("%d.%d", config.Version[0], config.Version[1]),
			"Copyright":    config.Copyright(),
		}
		var buffer bytes.Buffer
		src := MustAsset("templates/Info.plist")
		tmp := template.Must(template.New("Info.plist").Delims("[[", "]]").Parse(string(src)))
		err := tmp.Execute(&buffer, variable)
		if err != nil {
			panic(err)
		}

		binDir := fmt.Sprintf("%s.app/Contents/MacOS", config.Name)
		os.MkdirAll(binDir, 0755)

		// Write Info.plist
		ioutil.WriteFile(fmt.Sprintf("%s.app/Contents/Info.plist", config.Name), buffer.Bytes(), 0644)

		// Write Icon.
		rscDir := fmt.Sprintf("%s.app/Contents/Resources", config.Name)
		os.MkdirAll(rscDir, 0755)
		createIcon(rscDir)

		// Write executable file
		cmd := Command("go", ".", "generate")
		err = cmd.Run()
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}

		binPath := filepath.Join(binDir, strings.ToLower(config.Name))
		buffer.Reset()
		for _, framework := range config.AdditionalFrameworks {
			fmt.Fprintf(&buffer, " -framework %s", framework)
		}
		cmd = Command("go", ".", "build", "-ldflags", fmt.Sprintf(`-s -w -r %s%s`, filepath.Join(qtdir, "lib"), buffer.String()), "-o", binPath)
		cache := configdir.New("", "goqml").QueryCacheFolder()
		pp.Println(cmd.command.Args)
		pkgconfigPath := filepath.Join(cache.Path, "pkgconfig")
		cmd.AddEnv("PKG_CONFIG_PATH=" + pkgconfigPath)
		err = cmd.Run()
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
	}
}
