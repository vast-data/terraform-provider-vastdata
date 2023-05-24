package metadata

var cluster_config map[string]string = map[string]string{}

func SetClusterConfig(key, value string) {
	cluster_config[key] = value
}

func GetClusterConfig(key string) (string, bool) {
	value, exists := cluster_config[key]
	return value, exists
}
