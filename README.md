# filmoteka

Приложение запускается командой:
```
docker-compose up
```
### Авторизация.
#### POST /signin
Результатом успешной авторизации является отдача cookie. Пример запроса:
```
{
    "login":"andrey",
    "passsword": "andrey"
}
```
### Регистрация
#### POST /signup
Результатом успешной регистрация является создание нового бользователя в БД. Пример запроса:
```
{
    "login":"andrey",
    "passsword": "andrey"
}
```

### Выход
#### DELETE /logout
Для выхода из аккаунта необходима кука session_id, которая была получена при авторизации.

### Проверка авторизации
#### GET /authcheck
Аутентификация пользователя. Проверка просходи по куке session_id.

### Получение списка актёров
#### GET /api/v1/actors

### Добавление нового актёра
#### POST /api/v1/actors/add

### Редактирование информации об актёре
#### PATCH /api/v1/actors/update

### Удаление актёра
#### DELETE /api/v1/actors/delete

### Получение списка фильмов
#### GET /api/v1/films

### Поиск фильмов
#### GET /api/v1/films/search

### Добавление фильма
#### POST /api/v1/films/add

### Редактирование информации о фильме
#### PATCH /api/v1/films/update

### Удаление фильма
#### DELETE /api/v1/films/delete
