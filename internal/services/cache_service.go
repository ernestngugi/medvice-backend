package services

import (
	"encoding/json"

	"github.com/ernestngugi/medvice-backend/internal/providers"
)

type (
	CacheService interface {
		CacheValue(key string, value any) error
		Exists(key string) (bool, error)
		GetCachedValue(key string, result any) error
		RemoveFromCache(key string) error
	}

	cacheService struct {
		redisProvider providers.Redis
	}
)

func NewCacheService(
	redisProvider providers.Redis,
) CacheService {
	return &cacheService{
		redisProvider: redisProvider,
	}
}

func NewTestCacheService(
	redisProvider providers.Redis,
) *cacheService {
	return &cacheService{
		redisProvider: redisProvider,
	}
}

func (s *cacheService) CacheValue(
	key string,
	value any,
) error {

	cacheData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = s.redisProvider.Set(key, cacheData)
	if err != nil {
		return err
	}

	return nil
}

func (s *cacheService) Exists(
	key string,
) (bool, error) {

	exists, err := s.redisProvider.Exists(key)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *cacheService) GetCachedValue(
	key string,
	result any,
) error {

	payload, err := s.redisProvider.Get(key)
	if err != nil {
		return err
	}

	data, ok := payload.([]byte)
	if !ok {
		return err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	return nil
}

func (s *cacheService) RemoveFromCache(
	key string,
) error {

	err := s.redisProvider.Del(key)
	if err != nil {
		return err
	}

	return nil
}
