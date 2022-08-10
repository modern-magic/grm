package action

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/modern-magic/grm/internal"
	"github.com/modern-magic/grm/internal/logger"
	"github.com/modern-magic/grm/internal/registry"
)

func getCurrent() string {
	return registry.ReadNpm()
}

func ShowSources(source *registry.RegistryDataSource) {

	outLen := len(source.Keys) + 3

	cur := getCurrent()

	for _, key := range source.Keys {
		prefix := ""
		uri := source.Registry[key]
		if cur == uri {
			prefix = "* "
		}

		log := internal.StringJoin(prefix, key, getDashLine(key, outLen), uri, registry.Eol())

		if prefix == "" {
			fmt.Printf("%s", log)
		} else {
			logger.Success(log)

		}
	}

}

// show current registry uri and alias

func ShowCurrent() {
	cur := getCurrent()
	logger.Info(internal.StringJoin("[Grm]: you are using", cur))
}

func SetCurrent(source *registry.RegistryDataSource, args []string) int {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn("[Grm]: Plese pass an alias.")
			return
		}
	}()
	name := internal.PickArgs(args, 0)
	uri, err := getRegistryMeta(name, source.Registry, func(n string) (string, error) {
		return "", errors.New(internal.StringJoin("[Grm]: Can't found alias", name, "in your .nrmrc file. Please check it exist."))
	})
	if err != nil {
		logger.Error(internal.StringJoin(err.Error(), registry.Eol()))
		return 1
	}
	err = registry.WriteNrm("", uri, "")
	if err != nil {
		logger.Error(internal.StringJoin("[Grm]: error with", err.Error(), registry.Eol()))
		return 1
	}
	logger.Success(internal.StringJoin("[Grm]: use", name, "success~", registry.Eol()))
	return 0

}

// del .nrm file registry alias

func DelRegistry(source *registry.RegistryDataSource, args []string) int {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn("[Grm]: Plese pass an alias.")
			return
		}
	}()
	name := internal.PickArgs(args, 0)
	_, err := getRegistryMeta(name, source.UserRegistry, func(n string) (string, error) {
		return "", errors.New(internal.StringJoin("[Grm]: Can't found alias", name, "in your .nrmrc file. Please check it exist."))
	})
	if err != nil {
		logger.Error(internal.StringJoin(err.Error(), registry.Eol()))
		return 1
	}
	err = registry.DelNrm(name)
	if err != nil {
		logger.Error(internal.StringJoin("[Grm]: del registry fail", err.Error(), registry.Eol()))
		return 1
	}
	logger.Success(internal.StringJoin("[Grm]: del registry", name, "success!", registry.Eol()))
	return 0
}

func getRegistryMeta(name string, source map[string]string, callback func(name string) (string, error)) (string, error) {
	meta, ok := source[name]
	if !ok {
		return callback(name)
	}
	return meta, nil
}

func AddRegistry(source *registry.RegistryDataSource, args []string) int {

	name := ""
	home := ""
	uri := ""

	if len(args) <= 1 {
		logger.Error(internal.StringJoin("[Grm]: name and registry url is must be entry", registry.Eol()))
		return 1
	}
	name = args[0]

	if _, ok := source.UserRegistry[name]; ok {
		logger.Error("[Grm]: can't be the same as the default source name!")
		return 1
	}

	uri = args[1]
	if len(args) == 2 {
		home = uri
	}
	if len(args) >= 3 {
		home = args[2]
	}

	if !internal.IsUri(uri) && !internal.IsUri(home) {
		logger.Error("[Grm]: please verify the uri address you entered.")
		return 1
	}

	err := addRegistryImpl(name, uri, home)

	if err != nil {
		logger.Error(internal.StringJoin("[Grm]: add registry fail", err.Error(), registry.Eol()))
		return 1
	}
	logger.Success(internal.StringJoin("[Grm]: add registry success!", registry.Eol()))
	return 0
}

type FetchState uint8

const (
	SUCCESS FetchState = 1 << iota
	TIME_LIMIT
	FAIL
)

type ChannelStorage struct {
	state FetchState
	log   string
}

func FetchRegistry(source *registry.RegistryDataSource, args []string) int {

	keys := make([]string, 0)

	var wg sync.WaitGroup

	goCount := 5

	ch := make(chan ChannelStorage)

	if len(args) == 0 {
		keys = append(keys, source.Keys...)
	} else {
		keys = append(keys, args[0])
	}
	if len(keys) == 1 {
		if _, ok := source.Registry[keys[0]]; !ok {
			logger.Warn(internal.StringJoin("[Grm]: warning! can't found alias", keys[0], "please check it exist.", registry.Eol()))
			return 1
		}
	}

	for i := 0; i < goCount; i++ {
		go printFetchResult(&wg, ch)
	}
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		fetchImpl := func() (FetchState, string) {
			url := source.Registry[key]
			log := internal.StringJoin("[Grm]: fetch", key)
			res := internal.Fetch(url)

			if res.IsTimeout {
				log = internal.StringJoin(log, "state", res.Status)
			} else {
				log = internal.StringJoin(log, fmt.Sprintf("%.2f%s", res.Time, "s"), "state:", res.Status)
			}
			log = internal.StringJoin(log, registry.Eol())

			if res.IsTimeout {
				return TIME_LIMIT, log
			}

			if res.StatusCode != 200 {
				return FAIL, log
			}
			return SUCCESS, log
		}
		wg.Add(1)
		sendFetchResult(fetchImpl, ch)

	}
	wg.Wait()
	return 0
}

func printFetchResult(wg *sync.WaitGroup, ch chan ChannelStorage) {
	for m := range ch {
		switch m.state {
		case TIME_LIMIT:
			logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
				return fmt.Sprintf("%s%s%s", c.Dim, m.log, c.Reset)
			})
		case SUCCESS:
			logger.Success(m.log)
		case FAIL:
			logger.Error(m.log)
		}

		wg.Done()
	}

}

func sendFetchResult(f func() (FetchState, string), ch chan ChannelStorage) {
	go func() {
		state, log := f()
		ch <- ChannelStorage{
			state,
			log,
		}
	}()
}

func addRegistryImpl(name, uri, home string) error {
	return registry.WriteNrm(name, uri, home)

}

func getDashLine(key string, length int) string {
	final := math.Max(2, (float64(length) - float64(len(key))))
	bar := make([]string, int(final))
	for i := range bar {
		bar[i] = "-"
	}
	return strings.Join(bar[:], "-")
}
