# Auth service golang
## Тестовое задание от компании Medods

Ссылка на работающее приложение: https://auth-service-golang.herokuapp.com/

Эндпоинты реализованного приложения представлены ниже:

1) Метод: GET, Путь: /auth/user/{id} - Получение пары access/refresh токенов для пользователя. Id пользователя указывается в ссылке запроса.

Пример запроса:
```
curl -X GET https://auth-service-golang.herokuapp.com/auth/user/...
```
**Где ... - id пользователя.**

2) Метод: POST, Путь: /auth/tokens/refresh - Осуществление refresh операции по паре access/refresh токенов.

Пример запроса:
```
curl -X POST -d '{"access_token":"...","refresh_token":"..."}' https://auth-service-golang.herokuapp.com/auth/tokens/refresh
```

**Где ... - access и refresh токены.**

3) Метод: DELETE, Путь: /auth/refresh - Удаление конкретного refresh токена из базы данных. Refresh токен передается в теле запроса.

Пример запроса:
```
curl -X DELETE -d '{"refresh_token":"..."}' https://auth-service-golang.herokuapp.com/auth/refresh
```
**Где ... - полученный ранее refresh токен.**

4) Метод: DELETE, Путь: /auth/user/refresh - Удаление всех refresh токенов из базы данных для конкретного пользователя. Id пользователя передается в теле запроса.
Пример запроса:
```
curl -X DELETE -d '{"user_id":"..."}' https://auth-service-golang.herokuapp.com/auth/user/refresh
```
**Где ... - id пользователя.**