# Auth service golang
## Тестовое задание от компании Medods

Ссылка на работающее приложение: 

Эндпоинты реализованного приложения представлены ниже:

1) Метод: GET, Путь: /auth/{id} - Получение пары access/refresh токенов для пользователя. Id пользователя указывается в ссылке запроса.

Пример запроса:
```
curl -X GET http://localhost:8080/auth/user/...
```
**Где ... - id пользователя.**

2) Метод: GET, Путь: /auth/{id} - Осуществление refresh операции по паре access/refresh токенов.

Пример запроса:
```
curl -X POST -d '{"access_token":"...","refresh_token":"..."}' http://localhost:8080/auth/tokens/refresh
```

**Где ... - access и refresh токены.**

3) Метод: DELETE, Путь: /user/refresh - Удаление конкретного refresh токена из базы данных. Refresh токен передается в теле запроса.

Пример запроса:
```
curl -X DELETE -d '{"refresh_token":"..."}' http://localhost:8080/auth/refresh
```
**Где ... - полученный ранее refresh токен.**

4) Метод: GET, Путь: /user/refresh- Удаление всех refresh токенов из базы данных для конкретного пользователя. Id пользователя передается в теле запроса.
Пример запроса:
```
curl -X DELETE -d '{"user_id":"..."}' http://localhost:8080/auth/user/refresh
```
**Где ... - id пользователя.**