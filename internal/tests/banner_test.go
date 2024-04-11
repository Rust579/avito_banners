package tests

import (
	"avito_banners/internal/model"
	"fmt"
	"github.com/valyala/fasthttp"
	"sync"
	"testing"
	"time"
)

// Тест эндпоинта get-banner
func TestGetBanner(t *testing.T) {

	var (
		tagId           = "7"
		featureId       = "3"
		useLastRevision = "false"
	)

	uri := "/user_banner?tag_id=" + tagId + "&feature_id=" + featureId + "&use_last_revision=" + useLastRevision
	headers := map[string]string{
		"Authorization": "tokenforadmin",
	}

	resp, err := SendHTTPRequest(uri, fasthttp.MethodGet, headers, 200, nil)
	if err != nil {
		t.Errorf("send request error: %v", err)
	}

	fmt.Println(resp)
}

// Функция для генерации любого количества баннеров через эндпоинт create-banner
// Создание выполняется параллельно, но с задержкой в 20 миллисекунд для избежания одинаковых created_at
// Пробовал генерировать больше 10000 баннеров и по ним тестировал все эндпоинты
func TestCreateBanners(t *testing.T) {

	uri := "/create-banner"
	headers := map[string]string{
		"Authorization": "tokenforadmin",
	}

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		time.Sleep(20 * time.Millisecond)
		go func(i int) {
			defer wg.Done()
			body := model.BannerAddRequest{
				TagIds:    []int{i*2 - 1, i * 2},
				FeatureId: i,
				IsActive:  true,
				BannerItem: map[string]interface{}{
					"text":  i,
					"title": i,
					"url":   i,
				},
			}
			_, err := SendHTTPRequest(uri, fasthttp.MethodGet, headers, 201, body)
			if err != nil {
				t.Errorf("send request error: %v", err)
			}
		}(i)
	}

	/*for i := 99999; i <= 100100; i++ {
		wg.Add(1)
		time.Sleep(20 * time.Millisecond)
		go func(i int) {
			defer wg.Done()
			body := model.BannerAddRequest{
				TagIds:    []int{i, i + 1},
				FeatureId: 1,
				IsActive:  true,
				BannerItem: map[string]interface{}{
					"text":  99999,
					"title": 99999,
					"url":   99999,
				},
			}
			_, err := SendHTTPRequest(uri, fasthttp.MethodGet, headers, 201, body)
			if err != nil {
				t.Errorf("send request error: %v", err)
			}
		}(i)
	}*/

	wg.Wait()
}
