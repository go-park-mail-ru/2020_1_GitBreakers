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
```json
{
    "message": "sometext"
}
```

### 1.2 Регистрация

Запрос: `/signup` типа `POST`

Тело запроса:

```json
{
    "login": "string",
    "password": "string",
    "email": "string"
}
```
Ответ:  
1. 200 OK+поставит куки session_id
2. 409 Conflict(уже есть такой юзер)
```json
{
    "message": "already exsist user with"
}
```

### 1.3 Логаут

Запрос: `/logout` типа `GET`

Ответ:  
1. 200 OK+уберет куки session_id
2. 403 Unauthorized
```json
{
    "message": "you must to be logged in"
}
```

## 2. Профиль
### 2.1 Получение информации профиля

Запрос: `/profile` типа `GET`

Ответ:
```json
{
    "username": "string",
    "email": "string",
    "firstname": "string",
    "lastname": "string",
    "avatar": "string"
}
```
| Ключ          | Значение                 |
| ------------- | ------------------------ |
| `username`    | Логин                    |
| `email`       | Адрес электронной почты  |
| `firstname`   | Имя                      |
| `lastname`    | Фамилия                  |
| `avatar`      | Ссылка на аватарку (url) |

### 2.2 Изменение информации профиля

Запрос: `/profile` типа `PUT`

Тело запроса:
```json
{
    "username": "string",
    "email": "string",
    "firstname": "string",
    "lastname": "string"
}
```
| Ключ          | Значение                 |
| ------------- | ------------------------ |
| `username`    | Логин                    |
| `email`       | Адрес электронной почты  |
| `firstname`   | Имя                      |
| `lastname`    | Фамилия                  |

### 2.3 Удаление профиля

Запрос: `/profile` типа `DELETE`

### 2.4 Изменение пароля

Запрос: `/password` типа `PUT`

Тело запроса:
```json
{
    "old_password": "string",
    "new_password": "string"
}
```
| Ключ           | Значение      |
| -------------- | ------------- |
| `old_password` | Старый пароль |
| `new_password` | Новый пароль  |

### 2.5 Изменение аватарки

Запрос: `/avatar` типа `PUT`
         |
