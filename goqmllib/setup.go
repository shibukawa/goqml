package goqmllib

import (
	"bytes"
	"github.com/fatih/color"
	"github.com/shibukawa/configdir"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"fmt"
)

func Setup() {
	qtdir := FindQt()

	if qtdir == "" {
		color.Red("\nCan't find Qt. if you don't use standard install location, set QTDIR environment variable.\n")
		os.Exit(1)
	}

	var pkgconfigPath string
	useExtraEnv := false
	if runtime.GOOS == "darwin" && strings.HasSuffix(qtdir, "clang_64") {
		cache := configdir.New("", "goqml").QueryCacheFolder()
		pkgconfigPath = filepath.Join(cache.Path, "pkgconfig")
		os.MkdirAll(pkgconfigPath, 0755)
		files, err := AssetDir("packageconfig")
		if err != nil {
			panic(err)
		}
		variable := map[string]string{
			"Prefix": qtdir,
		}
		for _, file := range files {
			templateSrc := MustAsset("packageconfig/" + file)
			var buffer bytes.Buffer
			tmp := template.Must(template.New(file).Delims("[[", "]]").Parse(string(templateSrc)))
			err = tmp.Execute(&buffer, variable)
			if err != nil {
				panic(err)
			}
			ioutil.WriteFile(filepath.Join(pkgconfigPath, file), buffer.Bytes(), 0644)
		}
		useExtraEnv = true
	}
	get := Command("go", ".", "get", "gopkg.in/qml.v1")
	genqrc := Command("go", ".", "get", "-ldflags", fmt.Sprintf(`-s -w -r %s`, filepath.Join(qtdir, "lib")), "gopkg.in/qml.v1/cmd/genqrc")
	if useExtraEnv {
		get.AddEnv("PKG_CONFIG_PATH=" + pkgconfigPath)
		genqrc.AddEnv("PKG_CONFIG_PATH=" + pkgconfigPath)
	}
	get.Run()
	genqrc.Run()
}

func FindQt() string {
	// from environment variable
	env := os.Getenv("QTDIR")
	if env != "" {
		return env
	}

	// from qmake
	cmd := exec.Command("qmake", "-query", "QT_INSTALL_PREFIX")
	output, err := cmd.CombinedOutput()
	if err == nil {
		return strings.TrimSpace(string(output))
	}

	// Search default install folder
	var paths []string
	if runtime.GOOS == "windows" {
		paths = []string{
			os.Getenv("USERPROFILE"),
			"C:\\", "D:\\",
			os.Getenv("ProgramFiles"),
			os.Getenv("ProgramFiles(x86)"),
		}
	} else {
		paths = []string{
			os.Getenv("HOME"),
			"/",
		}
	}
	for _, path := range paths {
		versions, err := ioutil.ReadDir(filepath.Join(path, "Qt"))
		if err != nil {
			continue
		}
		var biggestDir string
		for _, version := range versions {
			if strings.HasPrefix(version.Name(), "5.") {
				stat, _ := os.Stat(filepath.Join(path, "Qt", version.Name()))
				if !stat.IsDir() {
					continue
				}
				if version.Name() > biggestDir {
					biggestDir = version.Name()
				}
			}
		}
		if biggestDir == "" {
			continue
		}
		targets, err := ioutil.ReadDir(filepath.Join(path, "Qt", biggestDir))
		var candidate string
		for _, target := range targets {
			name := target.Name()
			if hasPrefixes(name, "ios", "android", "winphone", "winrt") {
				continue
			}
			if strings.HasPrefix(name, "mingw") {
				// mingw has higher priority than MSVC by default
				// because Qt bundles mingw. It is good for default behavior to make it easy
				candidate = name
				break
			} else if name > candidate {
				// Higher version of MSVC has higher priority
				candidate = name
			}
		}
		if candidate != "" {
			return filepath.Join(path, "Qt", biggestDir, candidate)
		}
	}
	return ""
}

func hasPrefixes(str string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}
