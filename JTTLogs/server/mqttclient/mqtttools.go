package mqttclient

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var mqttCache *cache.Cache

func init()  {
	mqttCache = cache.New(45*time.Second, 10*time.Second)
}

