## 1. Авторизация
### 1.1 Логин

Запрос: `/login` типа `POST`
login: string, password:string

Тело запроса:
```json
{
    "login": "string",
    "password": "string"
}
```
Ответ:  
1. 200 OK+поставит куки session_id
2. 401 неверный пароль или логин
3. 400 уже авторизован или невалидный json

### 1.2 Регистрация

Запрос: `/signup` типа `POST`

Тело запроса:  
required email, login, password
```json
{
  "login": "string",
  "email": "string",
  "name": "string",
  "password": "string"
}
```
Ответ:  
1. 200 OK+поставит куки session_id
2. 409 Conflict(уже есть такой юзер)

### 1.3 Логаут

Запрос: `/logout` типа `GET`

Ответ:  
1. 200 OK+уберет куки session_id
2. 401 Unauthorized

## 2. Профиль
### 2.1 Получение информации профиля

Запрос: `/whoami` типа `GET`

Ответ:
1. 200 ok
```json
{
    "id": "int",
    "login": "string",
    "email": "string",
    "name": "string",
    "avatar": "string"
}
```
| Ключ          | Значение                 |
| ------------- | ------------------------ |
| `login`       | Логин                    |
| `email`       | Адрес электронной почты  |
| `name`        | Имя+Фамилия              |
| `avatar`      | Ссылка на аватарку (url) |
2. 401 unauthorized  
### 2.2 Обновить данные юзера
Запрос: `/profile` типа `PUT`
```json
{
    "email": "string",
    "name": "string",
    "password": "string"
}
```
Ответ:
1. 200 ok
2. 401 unauthorized
3. 400 json невалидный
### 2.3 Получить инфу по логину
Запрос: `/profile/{login}` типа `GET`

Ответ:
1. 200 ok
```json
{
    "id": "int",
    "name": "string",
    "login": "string",
    "image": "./static/image/avatar/default.jpg",
    "email": "string"
}
```
2. 404 не найден такой юзер
### 2.4 Загрузить аватарку

Запрос: `/avatar` типа `PUT`

Ответ:
1. 200 ok
```html
<form enctype="multipart/form-data">
    <input name="avatar" type="file" />
</form>
```
2. 401 unauthorized