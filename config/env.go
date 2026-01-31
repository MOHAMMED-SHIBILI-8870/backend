package config

import "os"

func GetEnv(key, def string) string {
	if v:=os.Getenv(key);v != key{
		return  v
	}
	return def
}