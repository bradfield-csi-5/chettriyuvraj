package main

type CacheConf struct {
	ProxyCachePath string
	Server         []Location
}

type Location struct {
	Path string
	// ProxyPass string
	ProxyPath  [4]byte // for simplicity
	ServerPort int
}
