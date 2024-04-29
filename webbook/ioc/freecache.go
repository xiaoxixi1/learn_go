package ioc

import (
	"github.com/coocood/freecache"
)

func InitFreeCache() *freecache.Cache {
	cacheSize := 100 * 1024 * 1024
	return freecache.NewCache(cacheSize)

}
