package tests

// Функция для генерации любого количества баннеров по эндпоинту create-banner
// Создание выполняется параллельно, но с задержкой в 20 миллисекунд для избежания одинаковых created_at
// Пробовал генерировать больше 10000 баннеров и по ним тестировал все эндпоинты
// Баннеры с неповторяющимися feature_ids
/*func TestCreateBannersManyFeatureId(t *testing.T) {

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

	wg.Wait()
}*/

// Генерация баннеров с одним feature_id, но с разными tag_ids
/*func TestCreateBannersOneFeatureId(t *testing.T) {

	uri := "/create-banner"
	headers := map[string]string{
		"Authorization": "tokenforadmin",
	}

	var wg sync.WaitGroup

	for i := 99999; i <= 100100; i++ {
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
	}

	wg.Wait()
}*/
