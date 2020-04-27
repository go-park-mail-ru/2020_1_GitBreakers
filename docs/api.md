**Во всех запросах кроме, `"GET", "HEAD", "OPTIONS", "TRACE"`,
 необходимо ставить заголовок `X-Csrf-Token` с CSRF токеном**

## Выдача CSRF токена
Запрос: `/api/v1/csrftoken` типа `GET`

Ответ:
1. 200 - ok
    - Токен приходит в **заголовке** `X-Csrf-Token`
2. 401 - unauthorized
    - токен может получить только авторизованный пользователь

## 1. Авторизация
### 1.1 Логин

Запрос: `/session` типа `POST`
required: login: string(>3 симв), password:string(>5 симв)  

Тело запроса:
```json
{
    "login": "string",
    "password": "string"
}
```
Ответ:  
1. 200 OK+поставит куки session_id
2. 400 невалдиный json или невалидные поля
3. 401 неверный пароль
4. 404 нет такого юзера
5. 406 уже авторизован

### 1.2 Регистрация

Запрос: `/user/profile` типа `POST`

Тело запроса:  
required email, login(>3 симв), password(>5 симв)  
```json
{
  "login": "string",
  "email": "string",
  "name": "string",
  "password": "string"
}
```
Ответ:  
1. 201 Created+поставит куки session_id  
2. 400 Невалидныые данные(плохой json или поля не прошли валидацию)  
3. 409 Conflict(уже есть такой юзер)  

### 1.3 Логаут

Запрос: `/session` типа `DELETE`

Ответ:  
1. 200 OK+уберет куки session_id
2. 401 Unauthorized

## 2. Профиль
### 2.1 Получение информации профиля

Запрос: `/user/profile` типа `GET`

Ответ:
1. 200 ok
```json
{
    "id": "int",
    "login": "string",
    "email": "string",
    "name": "string",
    "image": "string"
}
```
| Ключ          | Значение                 |
| ------------- | ------------------------ |
| `login`       | Логин                    |
| `email`       | Адрес электронной почты  |
| `name`        | Имя+Фамилия              |
| `avatar`      | Ссылка на аватарку (url) |
2. 401 unauthorized  
3. 404 юзера не существует  
### 2.2 Обновить данные юзера
Запрос: `/user/profile` типа `PUT`  
Тело запроса:  
Все поля опциональны  
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
4. 409 уже есть юзер с такими данными  
### 2.3 Получить инфу по логину
Запрос: `/user/profile/{login}` типа `GET`  

Ответ:  
1. 200 ok
```json
{
    "id": "int",
    "name": "string",
    "login": "string",
    "image": "http://localhost:8080/static/image/avatar/default.jpg",
    "email": "string"
}
```
2. 404 не найден такой юзер  
### 2.4 Загрузить аватарку

Запрос: `/user/avatar` типа `PUT`  
Запрос: Картинка(6MB max size), имя поля **avatar**  
```html
<form enctype="multipart/form-data">
    <input name="avatar" type="file" />
</form>
```
Ответ:
1. 200 ok
2. 400 плохой файл(недопустимый формат или большой размер)  
2. 401 не авторизован  
## 3. Репозиторий
### 3.1 Создать новый репозиторий
Запрос: `/user/repo` типа `POST`  
Required: name(alphanumeric), 
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
2. 400 невалидный json или сами данные
3. 401 unauthorized
4. 409 есть репак с таким названием  
### 3.2 Получить репозиторий по имени юзера и названию
Запрос: `/repo/{username}/{reponame}` типа `GET`  
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
2. 403 нет прав на просмотр(приватный)  
3. 404 не найден репозиторий с таким username+reponame  
### 3.3 Получить список своих репозиториев
Запрос: `/user/repo` типа `GET`  
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
Запрос: `/user/repo/{username}` типа `GET`  
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
2. 404 не найден такой юзер  
## 4. Ветки и коммиты  
### 4.1 Получить список веток по логину и названию репозитория  
Запрос: `/repo/{username}/{reponame}/branches` типа `GET`  
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
2. 403 (нет прав на просмотр)  
3. 404 (нет такого юзера или репозитория)  
### 4.2 Получить список коммитов 
Запрос: `/repo/{username}/{reponame}/commits/hash/{hash}` типа `GET`  
{commithash} - **хеш коммита** или последний коммит ветки(передается при получении списка веток)
Образец:  
`89.208.198.186:8080/repo/logggers/hefherser/commits/hash/23c70a09237681d7a0d908220a1a1af44ee74229?offset=2&limit=5`
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
2. 403 (нет прав на просмотр)  
3. 404 (нет такого юзера или репозитория или коммита)  
### 4.3 Получить список коммитов ветки (аналог 4.2)
Запрос: `/repo/{username}/{reponame}/commits/branch/{branchname}` типа `GET`  
{branchname} - обычное название ветки (master, dev,prod ...)  
Образец:  
`89.208.198.186:8080/localhost:8080/repo/lox5000/testname/commits/branch/master`  
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
2. 403 (нет прав на просмотр)  
3. 404 (нет такого юзера или репозитория или ветки)  
### 4.4 Получить список файлов по коммиту 
Запрос: `/repo/{username}/{reponame}/files/{hashcommits}` типа `GET`
с параметрами:
- `path` - путь до папки с файлами, например `./`, или `files/`.
Если параметр пустой, то считается, что  `path=./` 
   
Образец:  
`89.208.198.186:8080/repo/logggers/hefherser/files/07818d5fa5aef7dd7dca1d556f59c7a146a9b00c?path=docker/s6/crond`
Ответ:  
1. 200 ok
```json
[
    {
        "name": "LICENSE",
        "file_type": "blob",
        "file_mode": "regular",
        "file_size": 1054,
        "is_binary": false,
        "content_type": "text/plain; charset=utf-8",
        "entry_hash": "0640c41d4b3b4633682d839f980bcc33fca6e970"
    },
    {
        "name": "conf",
        "file_type": "tree",
        "file_mode": "dir",
        "file_size": -1,
        "is_binary": true,
        "content_type": "",
        "entry_hash": "9161b4edeb4e928405650145c6ca85f2131a6cff"
    }
]
```
2. 403 (нет прав на просмотр)  
3. 404 (нет такого юзера, репозитория, файла или коммита)   
### 4.5 Просмотр одного файла  
Запрос: `/repo/{username}/{reponame}/files/{hashcommits}` типа `GET`  
Образец:  
`89.208.198.186:8080/repo/logggers/hefherser/files/07818d5fa5aef7dd7dca1d556f59c7a146a9b00c?path=docker/main.go`
Ответ:  
1. 200 ok  
```json
{
    "file_info": {
        "name": "Dockerfile.rpihub",
        "file_type": "blob",
        "file_mode": "regular",
        "file_size": 1413,
        "is_binary": false,
        "content_type": "text/plain; charset=utf-8",
        "entry_hash": "d83a9f5ab53796e6c71491b466840d36184b0376"
    },
    "content": "text from files create table if not exists users\n(\n    nickname text, "
} 
```
- is_binary всегда false, иначе ошибка (нельзя 
    посмотреть бинарный файл)
2. 403 (нет прав на просмотр)  
3. 404 (нет такого юзера, репозитория, файла или коммита)   
## 5. Issues  
### 5.1 Создать issues  
Запрос: `/func/repo/{repoID}/issues` типа `POST`  
Required: author_id, repo_id,title(>0 symbol),message(>0symbol)    
Тело: 
```json
{
  "author_id": 54,
  "repo_id": 24,
  "title": "sfsfdsfsd",
  "message": "kekek",
  "label": "resolved",
  "is_closed": true
}
``` 
Ответ:  
1. 201 created  
2. 400 невалидный json
3. 401 (не авторизован)  
4. 403 (нет прав на создание в этом репозитории)  
5. 404 (нет такого юзера или репозитория)  
### 5.2 Получить список issues  
Запрос: `/func/repo/{repoID}/issues` типа `GET` 
Ответ:  
1.200 ok
```json
[{
  "id": 43,
  "author_id": 534,
  "repo_id": 24,
  "title": "sfsfdsfsd",
  "message": "kekekfafafasfasfasfasfasfas",
  "label": "resolved",
  "is_closed": true,
  "created_at": "2020-04-22T17:34:07.529Z"
},
{
  "id": 45,
  "author_id": 5,
  "repo_id": 24,
  "title": "sfsfdsfsd",
  "message": "kekek",
  "label": "resolved",
  "is_closed": true,
  "created_at": "2020-04-22T17:34:07.529Z"
}]
``` 
 2. 400 указана строка на месте repo_id  
 3. 403 (нет прав на просмтор в этом репозитории)  
 4. 404 (нет такого юзера или репозитория) 
### 5.3 Обновить issues  
Запрос: `/func/repo/{repoID}/issues` типа `POST`  
Можно обновить: message,title,label  
Required: id
Тело: 
```json
{
  "id": 242,
  "title": "sfsfdsfsd",
  "message": "kekek",
  "label": "resolved"
}
``` 
Ответ:  
1. 200 ок   
2. 400 невалидный json
3. 401 (не авторизован)  
4. 403 (нет прав на апдейт)  
5. 404 (нет такого вопроса или репозитория)  
### 5.4 Закрыть issues  
Запрос: `/func/repo/{repoID}/issues` типа `DELETE`   
Тело: 
```json
{
  "id": 242
}
``` 
Ответ:  
1. 200 ок  
2. 400 невалидный json
3. 401 (не авторизован)  
4. 403 (нет прав на апдейт)  
5. 404 (нет такого вопроса или репозитория)  
## 6. Stars 
### 6.1 Добавить/удалить звезду  
Запрос: `/func/repo/{repoID}/stars` типа `PUT`  
Описание: true добавит звезду/false уберет  
Тело: 
```json
{
  "vote": true
}
``` 
Ответ:  
1. 200 ок
2. 400 невалидный json
3. 401 (не авторизован)  
### 6.2 Список избранных репозиториев  
Запрос: `/func/repo/{login}/stars?limit=100&offset=2` типа `GET`  
Тело: 
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
Ответ:  
1. 200 ок  
2. 404 нет юзера с таким логином  
### 6.3 Список юзеров которые лайкнули репозиторий  
Запрос: `/func/repo/{repoID}/stars/users?limit=100&offset=2` типа `GET`  
Тело: 
```json
[{
    "id": "int",
    "login": "string",
    "email": "string",
    "name": "string",
    "image": "string"
},
{
    "id": "int",
    "login": "string",
    "email": "string",
    "name": "string",
    "image": "string"
}]
``` 
Ответ:  
1. 200 ок  
2. 400 вместо repoID  скидывают строку  
3. 403 нет прав на просмотр    
4. 404 нет такого repoID  
## 7. News  
### 7.1 Список новостей   
Запрос: `/func/repo/{repoID}/news?limit=100&offset=2` типа `GET`  
limit и offset опциональные параметры  
Тело: 
```json
[{
    "id": 5,
    "author_id": 5,
    "repo_id": 654433,
    "message": "dwrwrwsfdfrfe",
    "date": "2020-04-26T19:02:10.073Z"
},
{
    "id": 5,
    "author_id": 5,
    "repo_id": 654433,
    "message": "dwrwrwrfsffe",
    "date": "2020-04-26T19:02:10.073Z"
}]
``` 
Ответ:  
1. 200 ок   
2. 400 невалдиный repoid(строка)  
3. 401 неавторизован  
4. 403 нет доступа на просмотр новостей в данной репке  
5. 404 нет репки с таким repoid    
