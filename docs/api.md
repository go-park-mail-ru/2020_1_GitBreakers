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
        "name": "xp/git-lfs",
        "commit": {
            "commit_hash": "23c70a09237681d7a0d908220a1a1af44ee74229",
            "commit_author_name": "ᴜɴᴋɴᴡᴏɴ",
            "commit_author_email": "u@gogs.io",
            "commit_author_when": "2020-03-30T00:09:37+08:00",
            "committer_name": "ᴜɴᴋɴᴡᴏɴ",
            "committer_email": "u@gogs.io",
            "committer_when": "2020-03-30T00:09:37+08:00",
            "tree_hash": "c79a6098241e27d82de8f3a916dfa3d6ce0d9b7d",
            "commit_parents": [
                "5164d782afd860a5642c9bf71fea5f1723151ea6"
            ]
        }
    },
    {
        "name": "master",
        "commit": {
            "commit_hash": "07818d5fa5aef7dd7dca1d556f59c7a146a9b00c",
            "commit_author_name": "ᴜɴᴋɴᴡᴏɴ",
            "commit_author_email": "u@gogs.io",
            "commit_author_when": "2020-04-05T06:36:08+08:00",
            "committer_name": "GitHub",
            "committer_email": "noreply@github.com",
            "committer_when": "2020-04-05T06:36:08+08:00",
            "tree_hash": "14c89609a04f269123413f676a8cbe68c197de07",
            "commit_parents": [
                "bae1d6ccd81cd427382a2456e7c3646bdac9cf46"
            ]
        }
    }
]
```
2. 404 (нет такого юзера или репозитория)  
### 4.2 Получить список коммитов 
Запрос: `/{username}/{reponame}/commits/{branchname}` типа `GET`  
{branchname} - **хеш коммита** ветки(передается при получении списка веток)
Образец:  
`89.208.198.186:8080/logggers/hefherser/commits/23c70a09237681d7a0d908220a1a1af44ee74229?offset=2&limit=5`
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
### 4.3 Получить список коммитов ветки (аналог 4.2)
Запрос: `/{username}/{reponame}/{branchname}/commits` типа `GET`  
{branchname} - обычное название ветки (master, dev,prod ...)  
Образец:  
`89.208.198.186:8080/localhost:8080/lox5000/testname/bmstu/commits`  
Ответ:  
1. 200 ok
```json
[
    {
        "commit_hash": "2ef55ce2af5701880f2d165e6dbac49ca60d7e3f",
        "commit_author_name": "Deiklov",
        "commit_author_email": "romanov408g@mail.ru",
        "commit_author_when": "2020-04-05T21:23:55+03:00",
        "committer_name": "Deiklov",
        "committer_email": "romanov408g@mail.ru",
        "committer_when": "2020-04-05T21:23:55+03:00",
        "tree_hash": "ab52364eca9a07eaa7c70458a91edff748c243b7",
        "commit_parents": [
            "610307817a81acad346201f97548e72d6b061607"
        ]
    },
    {
        "commit_hash": "610307817a81acad346201f97548e72d6b061607",
        "commit_author_name": "Deiklov",
        "commit_author_email": "romanov408g@mail.ru",
        "commit_author_when": "2020-04-05T21:17:53+03:00",
        "committer_name": "Deiklov",
        "committer_email": "romanov408g@mail.ru",
        "committer_when": "2020-04-05T21:17:53+03:00",
        "tree_hash": "d1eb67922e1eb2d21191100a669ab6336a0681ff",
        "commit_parents": [
            "47695708f45d379f4608db11cc2b4b26c8c517b2"
        ]
    }
]
```
2. 404 (нет такого юзера или репозитория или коммита)  
### 4.4 Получить список файлов по коммиту 
Запрос: `/{username}/{reponame}/files/{commithash}` типа `GET`  
Образец:  
`89.208.198.186:8080/logggers/hefherser/files/07818d5fa5aef7dd7dca1d556f59c7a146a9b00c?path=docker/s6/crond`
Ответ:  
1. 200 ok
```json
[
    {
        "Name": ".s6-svscan",
        "FileType": "tree",
        "FileMode": "dir",
        "ContentType": "",
        "EntryHash": "720d59ef11a18b3c177e4e85854542b57060f232"
    },
    {
        "Name": "crond",
        "FileType": "tree",
        "FileMode": "regular",
        "ContentType": "",
        "EntryHash": "609789152b89b2d1d6a3fc892bb98115db1b8234"
    }
]
```
2. 404 (нет такого юзера или репозитория или коммита)