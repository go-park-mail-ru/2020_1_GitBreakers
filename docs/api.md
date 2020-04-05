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
## 3. Репозиторий
### 3.1 Создать новый репозиторий
Запрос: `/repo` типа `POST`
```json
{
    "name": "string",
    "description": "string",
    "is_public": "bool",
    "is_fork": "bool"
}
```
Ответ:
1. 201 created  
2. 401 unauthorized
3. 400 невалидный json 
### 3.2 Получить репозиторий по имени юзера и названию
Запрос: `/{username}/{reponame}` типа `GET`  
Ответ:
1. 200 ok  
```json
{
    "id": "int",
    "owner_id": "int",
    "name": "string",
    "description": "string",
    "is_public": "bool",
    "is_fork": "bool",
    "created_at": "date"
}
```
2. 401 unauthorized
### 3.3 Получить список своих репозиториев
Запрос: `/repolist` типа `GET`  
Ответ:
1. 200 ok  
```json
   [{
    "id": "int",
    "owner_id": "int",
    "name": "string",
    "description": "string",
    "is_public": "bool",
    "is_fork": "bool",
    "created_at": "date"
},
    {
    "id": "int",
    "owner_id": "int",
    "name": "string",
    "description": "string",
    "is_public": "bool",
    "is_fork": "bool",
    "created_at": "date"
}]
```
### 3.4 Получить список репозиториев юзера его логину
Запрос: `/{username}` типа `GET`  
Ответ:
1. 200 ok  
```json
   [{
    "id": "int",
    "owner_id": "int",
    "name": "string",
    "description": "string",
    "is_public": "bool",
    "is_fork": "bool",
    "created_at": "date"
},
{
    "id": "int",
    "owner_id": "int",
    "name": "string",
    "description": "string",
    "is_public": "bool",
    "is_fork": "bool",
    "created_at": "date"
}]
```
## 4. Ветки и коммиты  
### 4.1 Получить список веток по логину и названию репозитория  
Запрос: `/{username}/{reponame}/branches` типа `GET`  
Ответ:  
1. 200 ok  
```json
[
    {
        "name": "refs/heads/bmstu",
        "commit": {
            "commit_hash": "47695708f45d379f4608db11cc2b4b26c8c517b2",
            "commit_author_name": "Deiklov",
            "commit_author_email": "romanov408g@mail.ru",
            "commit_author_when": "2020-04-05T16:30:33+03:00",
            "committer_name": "Deiklov",
            "committer_email": "romanov408g@mail.ru",
            "committer_when": "2020-04-05T16:30:33+03:00",
            "tree_hash": "6d557828f228e61a303344841b199b4454b87c5e",
            "commit_parents": [
                "6ba6381462dfc29f0e6fcd1be049ddefe0ecda33"
            ]
        }
    },
    {
        "name": "refs/heads/hehe",
        "commit": {
            "commit_hash": "6ba6381462dfc29f0e6fcd1be049ddefe0ecda33",
            "commit_author_name": "Deiklov",
            "commit_author_email": "romanov408g@mail.ru",
            "commit_author_when": "2020-04-04T01:59:58+03:00",
            "committer_name": "Deiklov",
            "committer_email": "romanov408g@mail.ru",
            "committer_when": "2020-04-04T01:59:58+03:00",
            "tree_hash": "cc1e204bebe4eacf6964ec3e2d44ce03322e8644",
            "commit_parents": []
        }
    }
]
```
2. 404 (нет такого юзера или репозитория)  
### 4.2 Получить список коммитов 
Запрос: `/{username}/{reponame}/commits/{branchname}` типа `GET`  
{branchname} - хеш коммита ветки(передается при получении списка веток)
Образец:  
`89.208.198.186:8080/lox5000/ahahmfda/commits/47695708f45d379f4608db11cc2b4b26c8c517b2?offset=10&limit=10`
Ответ:  
1. 200 ok
```json  
[
    {
        "commit_hash": "47695708f45d379f4608db11cc2b4b26c8c517b2",
        "commit_author_name": "Deiklov",
        "commit_author_email": "romanov408g@mail.ru",
        "commit_author_when": "2020-04-05T16:30:33+03:00",
        "committer_name": "Deiklov",
        "committer_email": "romanov408g@mail.ru",
        "committer_when": "2020-04-05T16:30:33+03:00",
        "tree_hash": "6d557828f228e61a303344841b199b4454b87c5e",
        "commit_parents": [
            "6ba6381462dfc29f0e6fcd1be049ddefe0ecda33"
        ]
    },
    {
        "commit_hash": "6ba6381462dfc29f0e6fcd1be049ddefe0ecda33",
        "commit_author_name": "Deiklov",
        "commit_author_email": "romanov408g@mail.ru",
        "commit_author_when": "2020-04-04T01:59:58+03:00",
        "committer_name": "Deiklov",
        "committer_email": "romanov408g@mail.ru",
        "committer_when": "2020-04-04T01:59:58+03:00",
        "tree_hash": "cc1e204bebe4eacf6964ec3e2d44ce03322e8644",
        "commit_parents": []
    }
]
```
2. 404 (нет такого юзера или репозитория или коммита)  
### 4.3 Получить список файлов по коммиту 
Запрос: `/{username}/{reponame}/files/{commithash}` типа `GET`  
Образец:  
`localhost:8080/lox5000/testname/files/2ef55ce2af5701880f2d165e6dbac49ca60d7e3f`
Ответ:  
1. 200 ok
```json  

```
2. 404 (нет такого юзера или репозитория или коммита)