# TZ Parser

Этот проект представляет собой веб-скрапинг парсер для сбора данных о товарах с сайта kuper.ru. Данные извлекаются с помощью headless-браузера Chrome и сохраняются в формате CSV.

## Требования
- Go 1.20+
- Google Chrome
- chromedp (Go-библиотека для управления Chrome)

## Установка

1. Клонируйте репозиторий:
   ```sh
   git clone https://github.com/EmelinDanila/product-parser.git
   cd product-parser
   ```

2. Установите зависимости:
   ```sh
   go mod tidy
   ```

## Использование

Запустите парсер с помощью команды:
```sh
   go run main.go
```

После выполнения скрипта данные о товарах будут сохранены в файле `products.csv`.

## Конфигурация
- По умолчанию браузер работает в headless-режиме (не отображается на экране).
- Используется прокси-сервер для обхода возможных блокировок.
- Данные сохраняются в CSV в удобочитаемом формате.

## Формат CSV
Файл `products.csv` содержит три столбца:
- `Название` — наименование товара
- `Цена` — цена со скидкой
- `Ссылка` — URL-адрес страницы товара