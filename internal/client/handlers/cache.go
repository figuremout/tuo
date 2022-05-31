package handlers

import (
	"fmt"

	"github.com/allegro/bigcache/v3"
	"github.com/githubzjm/tuo/internal/client/cache"
	"github.com/githubzjm/tuo/internal/client/common"
)

func CacheInfo() {
	fmt.Printf("%-8s %-8s\n", "Bytes", "Entries")
	fmt.Printf("%-8d %-8d\n", cache.Capacity(), cache.Len())
}

func CacheReset() {
	if err := cache.Reset(); err != nil {
		cacheError(err)
	}
}

func CacheList() {
	if cache.Len() > 0 {
		fmt.Printf("%-8s %s\n", "Key", "Value")
	}
	cache.Range(func(entry *bigcache.EntryInfo) {
		fmt.Printf("%-8s %s\n", entry.Key(), string(entry.Value()))
		// fmt.Printf(common.NormalColor("KEY: ")+"%s"+common.NormalColor("\nVALUE: ")+"%s\n", entry.Key(), string(entry.Value()))
	})
}

func cacheError(err error) {
	fmt.Printf(common.ErrorColor("Cache Error: %s\n"), err.Error())
}
