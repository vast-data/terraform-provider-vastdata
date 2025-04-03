package metadata

var clusterConfig = map[string]string{}

func SetClusterConfig(key, value string) {
	clusterConfig[key] = value
}

func GetClusterConfig(key string) (string, bool) {
	value, exists := clusterConfig[key]
	return value, exists
}
