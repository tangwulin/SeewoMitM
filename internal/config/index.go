package config

var ResourceServerAddr string

func GetResourceServerAddr() string {
	return ResourceServerAddr
}

func SetResourceServerAddr(addr string) {
	ResourceServerAddr = addr
}
