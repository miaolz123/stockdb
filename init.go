package main

const (
	version = "0.0.1"
)

func init() {
	log(logInfo, "stockdb version "+version)
	loadConfig("")
}
