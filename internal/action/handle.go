// type FetchState uint8

// const (
// 	SUCCESS FetchState = 1 << iota
// 	TIME_LIMIT
// 	FAIL
// )

// type ChannelStorage struct {
// 	state FetchState
// 	log   string
// }

// func FetchRegistry(source *registry.RegistryDataSource, args []string) int {

// 	keys := make([]string, 0)

// 	var wg sync.WaitGroup

// 	goCount := 5

// 	ch := make(chan ChannelStorage)

// 	if len(args) == 0 {
// 		keys = append(keys, source.Keys...)
// 	} else {
// 		keys = append(keys, args[0])
// 	}
// 	if len(keys) == 1 {
// 		if _, ok := source.Registry[keys[0]]; !ok {
// 			logger.Warn(internal.StringJoin("[Grm]: warning! can't found alias", keys[0], "please check it exist."))
// 			return 1
// 		}
// 	}

// 	for i := 0; i < goCount; i++ {
// 		go printFetchResult(&wg, ch)
// 	}
// 	for i := 0; i < len(keys); i++ {
// 		key := keys[i]
// 		fetchImpl := func() (FetchState, string) {
// 			url := source.Registry[key]
// 			log := internal.StringJoin("[Grm]: fetch", key)
// 			res := internal.Fetch(url)

// 			if res.IsTimeout {
// 				log = internal.StringJoin(log, "state", res.Status)
// 			} else {
// 				log = internal.StringJoin(log, fmt.Sprintf("%.2f%s", res.Time, "s"), "state:", res.Status)
// 			}
// 			log = internal.StringJoin(log)

// 			if res.IsTimeout {
// 				return TIME_LIMIT, log
// 			}

// 			if res.StatusCode != 200 {
// 				return FAIL, log
// 			}
// 			return SUCCESS, log
// 		}
// 		wg.Add(1)
// 		sendFetchResult(fetchImpl, ch)

// 	}
// 	wg.Wait()
// 	return 0
// }

// func printFetchResult(wg *sync.WaitGroup, ch chan ChannelStorage) {
// 	for m := range ch {
// 		switch m.state {
// 		case TIME_LIMIT:
// 			logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
// 				return fmt.Sprintf("%s%s%s", c.Dim, m.log, c.Reset)
// 			})
// 		case SUCCESS:
// 			logger.Success(m.log)
// 		case FAIL:
// 			logger.Error(m.log)
// 		}

// 		wg.Done()
// 	}

// }

// func sendFetchResult(f func() (FetchState, string), ch chan ChannelStorage) {
// 	go func() {
// 		state, log := f()
// 		ch <- ChannelStorage{
// 			state,
// 			log,
// 		}
// 	}()
// }

package action

type RegistryDataSource struct {
	Name     string
	Uri      string
	Internal bool
}
