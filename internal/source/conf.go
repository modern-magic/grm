package source

import (
	"bufio"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/modern-magic/grm/internal/fs"
)

type S uint8

const (
	Npm S = iota
	Yarn
	Tencent
	NpmMirror
	System
)

var DefaultSource = map[string]S{
	"https://registry.npmjs.org/":            Npm,
	"https://registry.yarnpkg.com/":          Yarn,
	"https://mirrors.cloud.tencent.com/npm/": Tencent,
	"https://registry.npmmirror.com/":        NpmMirror,
}

var DefaultKey = map[S]string{
	Npm:       "https://registry.npmjs.org/",
	Yarn:      "https://registry.yarnpkg.com/",
	Tencent:   "https://mirrors.cloud.tencent.com/npm/",
	NpmMirror: "https://registry.npmmirror.com/",
}

func EnsureDefaultKey(input string) S {
	var s S
	switch input {
	case Npm.String():
		s = Npm
	case Yarn.String():
		s = Yarn
	case Tencent.String():
		s = Tencent
	case NpmMirror.String():
		s = NpmMirror
	default:
		s = System
	}
	return s
}

var SourceToString = []string{
	"npm",
	"yarn",
	"tencet",
	"npmMirror",
	"system",
}

func (s S) String() string {
	return SourceToString[s]
}

func readConf(path, alias string, c chan string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	sb := strings.Builder{}
	sb.WriteString(alias)
	sb.WriteString("->")
	if scanner.Scan() {
		sb.WriteString(scanner.Text())
		c <- sb.String()
	}
}

func ReadConf(path string) (url string, err error) {
	ch := make(chan string)
	go readConf(path, "*", ch)
	line := <-ch
	expr := strings.Split(line, "->")
	if len(expr) == 0 {
		return url, errors.New("Internal Error")
	}
	return expr[1], nil
}

type GrmConfig struct {
	BaseDir  string
	ConfPath string
	Paths    []string
	files    map[string]string // user conf
	parse    *GrmIni
}

func NewGrmConf() *GrmConfig {

	conf := &GrmConfig{
		BaseDir:  path.Join(fs.SystemPreffix, "grm"),
		ConfPath: path.Join(fs.SystemPreffix, ".npmrc"),
	}
	conf.parse = NewGrmIniParse(conf)
	return conf
}

func (g *GrmConfig) ListAllPath() {
	files := g.scanner()
	list := make([]string, 0, len(g.files)+len(SourceToString))
	for _, s := range SourceToString {
		if s != System.String() {
			list = append(list, s)
		}
	}
	for k := range files {
		list = append(list, k)
	}
	g.files = files
	g.Paths = list
}

func (g *GrmConfig) GetCurrentPath() string {
	g.parse.Get("registry")
	return g.parse.Path
}

func (g *GrmConfig) SetCurrentPath(target string) bool {
	return g.parse.Set("registry", target)
}

func (g *GrmConfig) GetCurrentConf() string {
	return g.parse.ToString()
}

func (g *GrmConfig) scanner() (files map[string]string) {
	if _, err := os.Stat(g.BaseDir); os.IsNotExist(err) {
		return files
	}
	fd, err := os.ReadDir(g.BaseDir)
	if err != nil {
		return files
	}
	files = make(map[string]string, len(fd))
	for _, file := range fd {
		if !file.IsDir() {
			n := file.Name()
			p := filepath.Join(g.BaseDir, n)
			files[n] = p
		}
	}
	return files
}

func (g *GrmConfig) ScannerUserConf() (source, key map[string]string) {
	if len(g.files) == 0 {
		return nil, nil
	}
	source = make(map[string]string)
	key = make(map[string]string)
	var wg sync.WaitGroup
	c := make(chan string)
	for name, file := range g.files {
		wg.Add(1)
		go func(path string, name string) {
			readConf(path, name, c)
			wg.Done()
		}(file, name)
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	for line := range c {
		if line == "" {
			continue
		}
		expr := strings.Split(line, "->")
		k := expr[0]
		v := expr[1]
		source[v] = k
		key[k] = v
	}

	return source, key
}

func (g *GrmConfig) MergePaths(userConf map[string]string) map[string]string {
	merged := make(map[string]string, len(userConf)+len(DefaultKey))
	for k, v := range userConf {
		merged[k] = v
	}
	for k, v := range DefaultKey {
		merged[k.String()] = v
	}
	return merged
}

func (g *GrmConfig) Files() map[string]string {
	return g.files
}
