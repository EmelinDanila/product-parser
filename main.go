package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

// Структура для хранения данных о продукте
type Product struct {
	Name  string // Название товара
	Price string // Цена товара (со скидкой)
	URL   string // Ссылка на товар
}

// Проверяет, работает ли прокси
func isProxyAlive(proxyURL string) bool {
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return false
	}

	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
		Timeout:   5 * time.Second, // Таймаут 5 секунд
	}

	resp, err := client.Get("https://www.google.com") // Пробный запрос
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func main() {
	// Список прокси
	proxies := []string{
		"http://185.10.129.14:3128",
		"http://116.203.203.208:3128",
		"http://45.140.143.77:18080",
	}

	// Выбираем рабочий прокси
	var selectedProxy string
	for i := 0; i < len(proxies); i++ {
		randomProxy := proxies[rand.Intn(len(proxies))]
		if isProxyAlive(randomProxy) {
			selectedProxy = randomProxy
			fmt.Println("Выбран рабочий прокси:", selectedProxy)
			break
		}
	}

	if selectedProxy == "" {
		log.Fatal("Нет доступных прокси!")
	}

	// Настройка Chrome в headless-режиме (без GUI)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer(selectedProxy), // Используем рабочий прокси
		chromedp.Flag("headless", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	// Создаём контекст Chrome с указанными настройками
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Список для хранения результатов парсинга
	var products []Product

	// Выполняем парсинг
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://kuper.ru/metro/c/myaso-ptitsa-40ebce3/konservi-iz-myasa"),
		chromedp.WaitVisible("#__next > div.body > div.Category_root__rJZG0 > main > section > div > div > div > div > div > main", chromedp.ByQuery),
		chromedp.Evaluate(`
			(function() {
				const products = [];
				const productBlocks = document.querySelectorAll('#__next > div.body > div.Category_root__rJZG0 > main > section > div > div > div > div > div > main > div > div.ProductsGrid_grid__VTD1X > div');

				productBlocks.forEach(item => {
					const nameElement = item.querySelector('.ProductCard_titleContainer__Kh_kg h3');
					const discountPriceElement = item.querySelector('.ProductCardPriceWrapper_root__0aldJ .ProductCard_price__LnWjd > div:first-child'); 
					const urlElement = item.querySelector('a');

					const name = nameElement ? nameElement.innerText.trim() : '';
					const price = discountPriceElement ? discountPriceElement.innerText.trim().replace('₽', '').trim() : '';
					const url = urlElement ? urlElement.href.trim() : '';

					if (name && price && url) {
						products.push({ name, price, url });
					}
				});
				return products;
			})();
		`, &products),
	)

	if err != nil {
		log.Println("Ошибка при выполнении chromedp:", err)
		return
	}

	if len(products) == 0 {
		log.Println("Данные не найдены. Возможно, разметка сайта изменилась или доступ ограничен.")
		return
	}

	// Сохраняем данные в CSV
	saveToCSV(products)
}

// Функция сохранения данных в CSV-файл
func saveToCSV(products []Product) {
	file, err := os.Create("products.csv")
	if err != nil {
		log.Println("Ошибка при создании CSV файла:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	writer.Write([]string{"Название товара", "Цена", "Ссылка"})
	for _, product := range products {
		writer.Write([]string{product.Name, product.Price, product.URL})
	}

	fmt.Println("Данные сохранены в products.csv")
}
