package source

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/edsrzf/mmap-go"
	"github.com/modern-magic/grm/internal/fs"
)

type S uint8

const (
	Npm S = iota
	Yarn
	HuaWei
	Tencent
	NpmMirror
	System
)

var DefaultSource = map[string]S{
	"https://registry.npmjs.org/":                  Npm,
	"https://registry.yarnpkg.com/":                Yarn,
	"https://repo.huaweicloud.com/repository/npm/": HuaWei,
	"https://mirrors.cloud.tencent.com/npm/":       Tencent,
	"https://registry.npmmirror.com/":              NpmMirror,
}

var DefaultKey = map[S]string{
	Npm:       "https://registry.npmjs.org/",
	Yarn:      "https://registry.yarnpkg.com/",
	HuaWei:    "https://repo.huaweicloud.com/repository/npm/",
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
	case HuaWei.String():
		s = HuaWei
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
	"huawei",
	"tencet",
	"npmMirror",
	"system",
}

func (s S) String() string {
	return SourceToString[s]
}

func readConf(path, alias string, c chan string) error {
	f, err := os.Open(path)
	if err != nil {
		c <- ""
		return err
	}
	defer f.Close()

	data, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		c <- ""
		return err
	}
	defer data.Unmap()

	scanner := bufio.NewScanner(os.NewFile(uintptr(f.Fd()), ""))
	if scanner.Scan() {
		c <- fmt.Sprintf("%s->%s", alias, scanner.Text())
	} else {
		c <- ""
	}
	return err
}

type GrmConfig struct {
	BaseDir  string
	ConfPath string
	Paths    []string
	files    []string // user conf
	aliases  []string // user alias
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
	aliases, files := g.scanner()
	list := make([]string, 0, len(aliases)+len(SourceToString))
	list = append(list, SourceToString...)
	list = append(list, aliases...)
	g.files = files
	g.aliases = aliases
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

func (g *GrmConfig) scanner() (aliases []string, files []string) {
	if _, err := os.Stat(g.BaseDir); os.IsNotExist(err) {
		return nil, nil
	}
	fd, err := os.ReadDir(g.BaseDir)
	if err != nil {
		return nil, nil
	}
	var b strings.Builder
	aliases = make([]string, 0, len(fd))
	files = make([]string, 0, len(fd))
	for _, file := range fd {
		if !file.IsDir() {
			aliases = append(aliases, file.Name())
			b.Reset()
			b.WriteString(g.BaseDir)
			b.WriteByte('/')
			b.WriteString(file.Name())
			files = append(files, b.String())
		}
	}
	return aliases, files
}

func (g *GrmConfig) ScannerUserConf() (source, key map[string]string) {
	if len(g.files) == 0 {
		return nil, nil
	}

	source = make(map[string]string)
	key = make(map[string]string)
	var wg sync.WaitGroup
	c := make(chan string)
	for pos, file := range g.files {
		wg.Add(1)
		go func(path string, pos int) {
			err := readConf(path, g.aliases[pos], c)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read config file: %v\n", err)
			}
			wg.Done()
		}(file, pos)
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
