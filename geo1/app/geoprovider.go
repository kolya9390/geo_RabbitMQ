package app

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kolya9390/gRPC_GeoProvider/server_rpc/service/dadata"
	"github.com/kolya9390/gRPC_GeoProvider/server_rpc/service/rebbit"
	"github.com/kolya9390/gRPC_GeoProvider/server_rpc/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GeoProviderService struct {
	storege   storage.GeoRepository
	limiter   *RateLimiter
	service_d dadata.DadataService
	service_r rebbit.MessageQueuer
}

type RateLimiter struct {
	requests map[string]chan time.Time
	mu       sync.Mutex
}

func NewRateLimiter(rateLimit int, per time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]chan time.Time),
	}
	go func() {
		ticker := time.NewTicker(per)
		defer ticker.Stop()
		for range ticker.C {
			rl.mu.Lock()
			for user, reqs := range rl.requests {
				for len(reqs) > 0 {
					<-reqs
				}
				delete(rl.requests, user)
			}
			rl.mu.Unlock()
		}
	}()
	return rl
}

func NewGeoProvider(storageDB storage.GeoRepository, service_d dadata.DadataService, service_r rebbit.MessageQueuer) *GeoProviderService {
	return &GeoProviderService{
		storege:   storageDB,
		service_d: service_d,
		service_r: service_r,
		limiter:   NewRateLimiter(5, time.Minute)}
}

func (gp *GeoProviderService) AddressSearch(input string) ([]*dadata.Address, error) {
	var result []*dadata.Address

	// Проверка кэша
	if ok, err := gp.storege.CheckAvailability(input); ok {
		if err != nil {
			return nil, err
		}
		addresses, err := gp.storege.Get(input)
		if err != nil {
			log.Printf("Get error: %s", err)
			return nil, err
		}

		for _, address := range addresses {
			result = append(result, &dadata.Address{
				GeoLat: address.GeoLat,
				GeoLon: address.GeoLon,
				Result: address.Region,
			})
		}
		return result, nil
	}

	requests := gp.limiter.Take("userTest")
select {
case requests <- time.Now():
    // Лимит запросов не превышен, продолжаем выполнение
default:
    // Если метод Take() заблокирован, значит, лимит запросов превышен
    // Отправляем сообщение в RabbitMQ
    go func() {
        err := gp.service_r.Publish("notification", []byte("пользователь превысил лимит TEST!"))
        if err != nil {
            log.Printf("Error publishing message to RabbitMQ: %v", err)
        }
    }()

    // Возвращаем ошибку 429 Too Many Requests
    return nil, status.Errorf(codes.ResourceExhausted, "429 Too Many Requests")
}
	// Если данные в кэше устарели или их нет, обращаемся к сервису Dadata
	respData, err := gp.service_d.AddressSearch(input)
	if err != nil {
		log.Printf("AddressSearch error: %s", err)
		return nil, err
	}

	if len(respData) < 1 {
		return nil, fmt.Errorf("address was not found")
	}

	// Сохранение данных в кэше
	err = gp.storege.Add(input, respData[0].Result, respData[0].GeoLat, respData[0].GeoLon)
	if err != nil {
		log.Printf("Add error: %s", err)
	}

	for _, address := range respData {
		result = append(result, &address)
	}

	return result, nil
}

func (gp *GeoProviderService) GeoCode(lat, lng string) ([]*dadata.Address, error) {

	var result []*dadata.Address

	geocode := fmt.Sprintf("%s %s", lat, lng)

	// Проверка кэша
	if ok, err := gp.storege.CheckAvailability(geocode); ok {
		if err != nil {
			return nil, err
		}
		addresses, err := gp.storege.Get(geocode)
		if err != nil {
			log.Printf("Get error: %s", err)
			return nil, err
		}

		for _, address := range addresses {
			result = append(result, &dadata.Address{
				GeoLat: address.GeoLat,
				GeoLon: address.GeoLon,
				Result: address.Region,
			})
		}
		return result, nil
	}

	requests := gp.limiter.Take("userTest")
	select {
	case requests <- time.Now():
		// Лимит запросов не превышен, продолжаем выполнение
	default:
		// Если метод Take() заблокирован, значит, лимит запросов превышен
		// Отправляем сообщение в RabbitMQ
		go func() {
			err := gp.service_r.Publish("notification", []byte("пользователь превысил лимит TEST!"))
			if err != nil {
				log.Printf("Error publishing message to RabbitMQ: %v", err)
			}
		}()
	
		// Возвращаем ошибку 429 Too Many Requests
		return nil, status.Errorf(codes.ResourceExhausted, "429 Too Many Requests")
	}

	// Если данные в кэше устарели или их нет, обращаемся к сервису Dadata
	respData, err := gp.service_d.GeoCode(lat, lng)
	if err != nil {
		log.Printf("AddressSearch error: %s", err)
		return nil, err
	}

	if len(respData) < 1 {
		return nil, fmt.Errorf("address was not found")
	}

	// Сохранение данных в кэше
	err = gp.storege.Add(geocode, respData[0].Result, respData[0].GeoLat, respData[0].GeoLon)
	if err != nil {
		log.Printf("Add error: %s", err)
	}

	for _, address := range respData {
		result = append(result, &address)
	}

	return result, nil
}

func (rl *RateLimiter) Take(user string) chan time.Time {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	reqs, ok := rl.requests[user]
	if !ok {
		reqs = make(chan time.Time, 5)
		rl.requests[user] = reqs
	}
	return reqs
}
