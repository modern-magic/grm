package registry

var UserRegistry = make(map[string]RegsitryInfo, 0)

func GetUserRegistryInfo() (map[string]RegsitryInfo, []string) {
	return ReadNrm()
}
