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
2. 403 не прошел авторизацию

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
2. 403 Unauthorized

## 2. Профиль
### 2.1 Получение информации профиля

Запрос: `/profile` типа `GET`

Ответ:
1. 200 ok
```json
{
    "username": "string",
    "email": "string",
    "firstname": "string",
    "lastname": "string",
    "avatar": "string",
    "password": ""
}
```
| Ключ          | Значение                 |
| ------------- | ------------------------ |
| `username`    | Логин                    |
| `email`       | Адрес электронной почты  |
| `firstname`   | Имя                      |
| `lastname`    | Фамилия                  |
| `avatar`      | Ссылка на аватарку (url) |
2. 401 unauthorized
### 2.2 Загрузить аватарку

Запрос: `/avatar` типа `PUT`

Ответ:
1. 200 ok
```html
<form enctype="multipart/form-data">
    <input name="avatar" type="file" />
</form>
```
2. 401 unauthorized