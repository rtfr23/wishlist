
# Wishlist API

## Что делает сервис

REST API для управления вишлистами и подарками. 
Позволяет:
-   Регистрироваться и авторизоваться
-   Создавать, читать, обновлять и удалять вишлисты
-   Управлять подарками в вишлистах (добавлять, обновлять, удалять)
-   Делиться вишлистами по публичной ссылке 
-   Бронировать подарки без авторизации по публичной ссылке
## Как запустить

1.  Создать файл `.env` в корне проекта
2.  Запустить сервис и базу данных:
    docker compose up --build -d
3.  Сервис будет доступен по адресу: `http://localhost:8080`
## Примеры запросов

### 1. Авторизация

**Регистрация**
`curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"user@mail.com","password":"123123"}'`
  **Аутентификация**
`curl -X POST http://localhost:8080/auth/signin \
  -H "Content-Type: application/json" \
  -d '{"email":"user@mail.com","password":"123123"}'`
  Вернется токен, необходимо сохранить его:
  `export TOKEN="..."`
 
  ### 2. Вишлисты

**Создать**
`curl -X POST http://localhost:8080/wishlists \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"event":"День рождения","date":"2026-12-31T23:59:59Z"}'`
  **Получить вишлист**
  `curl -X GET http://localhost:8080/wishlists \
  -H "Authorization: Bearer $TOKEN"`
  **Обновить**
  `curl -X PATCH http://localhost:8080/wishlists/$WISHLIST_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"event":"Новое название"}'`
  **Удалить**
  `curl -X DELETE http://localhost:8080/wishlists/$WISHLIST_ID \
  -H "Authorization: Bearer $TOKEN"`
  ### 3. Подарки
  **Добавить**
  `curl -X POST http://localhost:8080/wishlists/$WISHLIST_ID/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Наушники","priority":5,"url":"https://URL.com"}'`
  **Обновить**
  `curl -X PATCH http://localhost:8080/wishlists/$WISHLIST_ID/items/$ITEM_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"priority":3}'`
  **Удалить**
  `curl -X DELETE http://localhost:8080/wishlists/$WISHLIST_ID/items/$ITEM_ID \
  -H "Authorization: Bearer $TOKEN"`
  **Получить все подарки вишлиста**
  `curl -X GET http://localhost:8080/wishlists/$WISHLIST_ID/items \
  -H "Authorization: Bearer $TOKEN"`
  **Получить конкретный подарок**
  `curl -X GET http://localhost:8080/wishlists/$WISHLIST_ID/items/$ITEM_ID \
  -H "Authorization: Bearer $TOKEN"`
  ### 4. Публичный доступ
  **Просмотр вишлиста**
  `curl -X GET http://localhost:8080/public/wishlists/$WISHLIST_TOKEN`
  **Забронировать подарок**
  `curl -X POST http://localhost:8080/public/wishlists/$WISHLIST_TOKEN/items/$ITEM_ID/reserve \
  -H "Content-Type: application/json" \
  -d '{}'`
  