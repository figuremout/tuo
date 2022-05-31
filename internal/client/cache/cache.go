package cache

import (
	"errors"
	"time"

	"github.com/allegro/bigcache/v3"
	jsonUtil "github.com/githubzjm/tuo/internal/pkg/utils/json"
)

var Cache *bigcache.BigCache

const (
	KeyUser  = "User"  // the login user
	KeyUsers = "Users" // all related users
	KeyToken = "Token"
)

var KeyDescription = map[string]string{
	KeyUser:  "Current user",
	KeyUsers: "All related Users",
	KeyToken: "Current user's token",
}

type User struct {
	UserID   uint   `json:"userid"`
	Username string `json:"username"`
}

// Should not rely on the cache's expire and clean mechanism, the kv will be removed due to many reasons,
// like expiration time or no space left for the new entry, or because delete was called.
// if kv exist, use it, if get new value or find the value invalid then reset it
// if kv not exist, try to get new value

// to make cached token survice longer than token valid duration can make full use of token, reduce the re-login times
func InitCache(lifeWindow, cleanWindow time.Duration) error {
	var err error
	config := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// lifeWindow is the expire interval, cleanWindow is the clean interval,
		// the kv can still be got after expire but before clean,
		// the shorter cleanWindow than lifeWindow, the sooner expired kv cleaned.
		// the longest kv survive time is lifeWindow + cleanWindow

		// time after which entry can be evicted
		LifeWindow: lifeWindow,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: cleanWindow, // the value smaller, the exipred entry will be deleted sooner

		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,

		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 500,

		// prints information about additional memory allocation
		Verbose: true,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 8192,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		OnRemoveWithReason: nil,
	}

	Cache, err = bigcache.NewBigCache(config)
	if err != nil {
		return err
	}
	return nil
}

// will be called when CleanWindow time is reached
// func onExpire(key string, extry []byte, reason bigcache.RemoveReason) {
// 	if reason == bigcache.Expired {
// 		// entry removed because of exipiration
// 		if key == KeyToken {
// 			// auto login to refresh token
// 		}
// 		if key == "k" {
// 			fmt.Print("expire")
// 		}
// 	}
// }

func Set(key string, value interface{}) error {
	return Cache.Set(key, jsonUtil.Dumps(value))
}

// if exist, append;if not exist, create
func Append(key string, value interface{}) error {
	return Cache.Append(key, jsonUtil.Dumps(value))
}

// may get expired entry
// return []byte{} if key not exist or other error occurs
func Get(key string) ([]byte, error) {
	bytes, err := Cache.Get(key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return []byte{}, nil
		}
		return []byte{}, err
	}
	return bytes, nil
}

func Del(key string) error {
	return Cache.Delete(key)
}
func Reset() error {
	return Cache.Reset()
}
func Close() error {
	return Cache.Close()
}

func Capacity() int {
	return Cache.Capacity()
}
func Len() int {
	return Cache.Len()
}

func Range(f func(entry *bigcache.EntryInfo)) {
	iterator := Cache.Iterator()
	for iterator.SetNext() {
		entry, err := iterator.Value()
		if err != nil {
			return
		}
		f(&entry)
	}
}
