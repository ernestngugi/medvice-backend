package services

import (
	"testing"

	"github.com/ernestngugi/medvice-backend/internal/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCacheController(t *testing.T) {

	redisProvider := mocks.NewMockRedisProvider()

	cacheService := NewTestCacheService(redisProvider)

	Convey("TestCacheController", t, func() {

		Convey("can add something to cache", func() {

			err := cacheService.CacheValue("key", "value")
			So(err, ShouldBeNil)
		})

		Convey("can get a cached value", func() {

			err := cacheService.CacheValue("key1", "value1")
			So(err, ShouldBeNil)

			var value string

			err = cacheService.GetCachedValue("key1", &value)
			So(err, ShouldBeNil)

			So(value, ShouldEqual, "value1")
		})

		Convey("can check if a key already exists", func() {

			err := cacheService.CacheValue("key", "value")
			So(err, ShouldBeNil)

			exist, err := cacheService.Exists("key")
			So(err, ShouldBeNil)
			So(exist, ShouldBeTrue)
		})

		Convey("can remove key from cache", func() {

			key := "key"

			err := cacheService.CacheValue(key, "value1")
			So(err, ShouldBeNil)

			err = cacheService.RemoveFromCache(key)
			So(err, ShouldBeNil)

			exists, err := cacheService.Exists(key)
			So(err, ShouldBeNil)
			So(exists, ShouldBeFalse)
		})
	})
}
