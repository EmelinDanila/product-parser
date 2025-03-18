package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/chromedp/chromedp"
)

// Структура для хранения данных о продукте
type Product struct {
	Name  string // Название товара
	Price string // Цена товара (со скидкой)
	URL   string // Ссылка на товар
}

func main() {
	// Настройка Chrome в headless-режиме (без GUI)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer("http://185.10.129.14:3128"), // Используем прокси
		chromedp.Flag("headless", true),                   // Запуск в фоне
		chromedp.Flag("ignore-certificate-errors", true),  // Отключение проверки SSL
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
		// Открываем страницу
		chromedp.Navigate("https://kuper.ru/metro/c/myaso-ptitsa-40ebce3/konservi-iz-myasa?ads_identity.ads_promo_identity.placement_uid=cns3h249u8b43pmgrklg&ads_identity.ads_promo_identity.site_uid=c9qep2jupf8ugo3scn10&sid=1"),

		// Ждём, пока загрузится блок с товарами
		chromedp.WaitVisible("#__next > div.body > div.Category_root__rJZG0 > main > section > div > div > div > div > div > main", chromedp.ByQuery),

		// Извлекаем данные о товарах
		chromedp.Evaluate(`
			(() => {
				const products = [];
				const productBlocks = document.querySelectorAll('#__next > div.body > div.Category_root__rJZG0 > main > section > div > div > div > div > div > main > div > div.ProductsGrid_grid__VTD1X > div');

				productBlocks.forEach(item => {
					const nameElement = item.querySelector('.ProductCard_titleContainer__Kh_kg h3');
					const discountPriceElement = item.querySelector('.ProductCardPriceWrapper_root__0aldJ .ProductCard_price__LnWjd > div:first-child'); // Цена со скидкой
					const urlElement = item.querySelector('a');

					const name = nameElement ? nameElement.innerText.trim() : '';
					const price = discountPriceElement ? discountPriceElement.innerText.trim() : '';
					const url = urlElement ? urlElement.href.trim() : '';

					if (name && price && url) {
						products.push({ name, price, url });
					}
				});
				return products;
			})()`, &products,
		),
	)

	if err != nil {
		log.Fatalf("Ошибка при выполнении chromedp: %v", err)
	}

	// Сохраняем данные в CSV
	saveToCSV(products)
}

// Функция сохранения данных в CSV-файл
func saveToCSV(products []Product) {
	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatalf("Ошибка при создании CSV файла: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Записываем заголовки
	writer.Write([]string{"Название товара", "Цена (₽)", "Ссылка"})

	// Записываем данные о товарах
	for _, product := range products {
		writer.Write([]string{product.Name, product.Price, product.URL})
	}

	fmt.Println("Данные сохранены в products.csv")
}
