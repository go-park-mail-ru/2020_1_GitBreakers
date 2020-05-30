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

Запрос: `/api/v1/session` типа `POST`
required: login: string(>3 симв), password:string(>5 симв)

Тело запроса:
```json
{
    "login": "string",
    "password": "string"
}
```
Примечания:
- полем login может быть как логин пользователя, так и его email

Ответ:
1. 200 OK+поставит куки session_id
2. 400 невалдиный json или невалидные поля
3. 401 неверный пароль
4. 404 нет такого юзера
5. 406 уже авторизован

### 1.2 Регистрация

Запрос: `/api/v1/user/profile` типа `POST`

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

Запрос: `/api/v1/session` типа `DELETE`

Ответ:
1. 200 OK+уберет куки session_id
2. 401 Unauthorized

## 2. Профиль
### 2.1 Получение информации профиля

Запрос: `/api/v1/user/profile` типа `GET`

Ответ:
1. 200 ok
```json
{
    "id": 75,
    "name": "kek",
    "login": "kek2101",
    "image": "/static/image/avatar/default.png",
    "email": "alexloh500@mail.ru",
    "created_at": "2020-01-01T00:00:00Z"
}
```
| Ключ          | Значение                 |
| ------------- | ------------------------ |
| `login`       | Логин                    |
| `email`       | Адрес электронной почты  |
| `name`        | Имя+Фамилия              |
| `image`       | Ссылка на аватарку (url) |
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
### 2.3 Получить информацию о пользователе по логину
Запрос: `/api/v1/user/profile/{login or email}` типа `GET`
можно почту скинуть или логин

Ответ:
1. 200 ok
```json
{
    "id": "int",
    "name": "string",
    "login": "string",
    "image": "/static/image/avatar/default.jpg",
    "email": "string"
}
```
2. 404 не найден такой юзер

### 2.4 Получить информацию о пользователе по ID
Запрос: `/api/v1/user/id/{id}/profile` типа `GET`

Ответ:
1. 200 ok
```json
{
    "id": "int",
    "name": "string",
    "login": "string",
    "image": "/static/image/avatar/default.jpg",
    "email": "string"
}
```
2. 400 - некорректный id (не удалось преобразовать в число)
3. 404 не найден такой юзер

### 2.5 Загрузить аватарку

Запрос: `/api/v1/user/avatar` типа `PUT`
Запрос: Картинка(6MB max size), имя поля **avatar**
```html
<form enctype="multipart/form-data">
    <input name="avatar" type="file" />
</form>
```
Ответ:
1. 200 ok
2. 400 плохой файл(недопустимый формат)
2. 401 не авторизован
3. 413 слишком большой файл для загрузки
## 3. Репозиторий
### 3.1 Создать новый репозиторий
Запрос: `/api/v1/user/repo` типа `POST`
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
Запрос: `/api/v1/repo/{username}/{reponame}` типа `GET`
Ответ:
1. 200 ok  (поле parent_repository_info может быть пустым, если is_fork = false)
```json
{
    "id": "int",
    "owner_id": "int",
    "name": "string",
    "description": "string",
    "is_public": "bool",
    "is_fork": "bool",
    "created_at": "date",
    "parent_repository_info": {
        "id": "int",
        "owner_id": "int",
        "name": "string",
        "author_login": "string"
    }
}
```
2. 403 нет прав на просмотр(приватный)
3. 404 не найден репозиторий с таким username+reponame
### 3.3 Получить список своих репозиториев
Запрос: `/api/v1/user/repo` типа `GET`
Ответ:
1. 200 ok
```json
[
    {
        "id": 67,
        "owner_id": 75,
        "name": "pes",
        "description": "",
        "is_fork": true,
        "created_at": "2020-05-17T20:08:50.029927Z",
        "is_public": true,
        "stars": 0,
        "forks": 0,
        "merge_requests_open": 0,
        "author_login": "mudila2101",
        "parent_repository_info": {
            "id": 37,
            "owner_id": 51,
            "name": "horse",
            "author_login": "cheburek111"
        }
    },
    {
        "id": 68,
        "owner_id": 75,
        "name": "pess",
        "description": "",
        "is_fork": false,
        "created_at": "2020-05-17T20:08:50.029927Z",
        "is_public": true,
        "stars": 0,
        "forks": 0,
        "merge_requests_open": 0,
        "author_login": "mudila2101",
        "parent_repository_info": {}
    }
]
```
### 3.4 Получить список репозиториев юзера его логину
Запрос: `/api/v1/user/repo/{username}` типа `GET`
Ответ:
1. 200 ok
```json
[
    {
        "id": 67,
        "owner_id": 75,
        "name": "pes",
        "description": "",
        "is_fork": true,
        "created_at": "2020-05-17T20:08:50.029927Z",
        "is_public": true,
        "stars": 0,
        "forks": 0,
        "merge_requests_open": 0,
        "author_login": "mudila2101",
        "parent_repository_info": {
            "id": 37,
            "owner_id": 51,
            "name": "horse",
            "author_login": "cheburek111"
        }
    }
]
```
2. 404 не найден такой юзер

### 3.5 Получить список репозиториев юзера его ID
Запрос: `/api/v1/user/id/{id}/repo` типа `GET`
Ответ:
1. 200 ok
```json
[
    {
        "id": 67,
        "owner_id": 75,
        "name": "pes",
        "description": "",
        "is_fork": true,
        "created_at": "2020-05-17T20:08:50.029927Z",
        "is_public": true,
        "stars": 0,
        "forks": 0,
        "merge_requests_open": 0,
        "author_login": "mudila2101",
        "parent_repository_info": {
            "id": 37,
            "owner_id": 51,
            "name": "horse",
            "author_login": "cheburek111"
        }
    }
]
```
2. 400 - некорректный id (не удалось преобразовать в число)
3. 404 не найден такой юзер

### 3.6 Удалить репозиторий
Запрос: `/api/v1/user/repo` типа `DELETE`
```json
{
    "name": "string"
}
```
Ответ:
1. 200 - репозиторий успешно удалён
2. 400 - некорректные данные
3. 401 - пользователь не авторизован
4. 404 - не найден целевой репозиторий

## 4. Ветки и коммиты
### 4.1 Получить список веток по логину и названию репозитория
Запрос: `/api/v1/repo/{username}/{reponame}/branches` типа `GET`
Ответ:
1. 200 ok
```json
[
    {
        "name": "master",
        "commit": {
            "commit_hash": "89944a4f685579117fb8a36649cdd4b99b3d56e6",
            "commit_author_name": "Deiklov",
            "commit_author_email": "romanov408g@mail.ru",
            "commit_author_when": "2020-04-27T17:58:11+03:00",
            "committer_name": "Deiklov",
            "committer_email": "romanov408g@mail.ru",
            "committer_when": "2020-04-27T17:58:11+03:00",
            "tree_hash": "9e60190702cadececd04cd8faf82aa8659e57ada",
            "commit_parents": [
                "1bc2f16b52e5cbbb3b64e1f050fc25009e9a4404"
            ]
        }
    }
]
```
2. 403 (нет прав на просмотр)
3. 404 (нет такого юзера или репозитория)
### 4.2 Получить список коммитов
Запрос: `/api/v1/repo/{username}/{reponame}/commits/hash/{hash}` типа `GET`
{commithash} - **хеш коммита** или последний коммит ветки(передается при получении списка веток)
Образец:
`89.208.198.186:8080/api/v1/repo/logggers/hefherser/commits/hash/23c70a09237681d7a0d908220a1a1af44ee74229?offset=2&limit=5`
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
Запрос: `/api/v1/repo/{username}/{reponame}/commits/branch/{branchname}` типа `GET`
{branchname} - обычное название ветки (master, dev,prod ...)
Образец:
`89.208.198.186:8080/api/v1/localhost:8080/repo/lox5000/testname/commits/branch/master`
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
Запрос: `/api/v1/repo/{username}/{reponame}/files/{hashcommits}` типа `GET`
с параметрами:
- `path` - путь до папки с файлами, например `./`, или `files/`.
Если параметр пустой, то считается, что  `path=./`

Образец:
`89.208.198.186:8080/api/v1/repo/logggers/hefherser/files/07818d5fa5aef7dd7dca1d556f59c7a146a9b00c?path=docker/s6/crond`
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
4. 406 (неприемлимо) - файл либо бинарный, либо мы не можем отправить его в формате json.
для получения бинарного контента см методы ниже
### 4.5 Просмотр одного файла
Запрос: `/api/v1/repo/{username}/{reponame}/files/{hashcommits}` типа `GET`
Образец:
`89.208.198.186:8080/api/v1/repo/logggers/hefherser/files/07818d5fa5aef7dd7dca1d556f59c7a146a9b00c?path=docker/main.go`
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

### 4.6 Получение HEAD репозитория (его состояния по умолчанию)
Запрос: `/api/v1/repo/{username}/{reponame}/head` типа `GET`
Образец: `89.208.198.186:8080/api/v1/repo/logggers/hefherser/head`
Ответ:
```json
{
    "name": "master",
    "commit": {
        "commit_hash": "1e16727dfaf12c018cb4c069762d9ab9d62c0814",
        "commit_author_name": "UlianaBespalova",
        "commit_author_email": "43138516+UlianaBespalova@users.noreply.github.com",
        "commit_author_when": "2020-04-27T20:06:36+03:00",
        "committer_name": "GitHub",
        "committer_email": "noreply@github.com",
        "committer_when": "2020-04-27T20:06:36+03:00",
        "tree_hash": "88b5bc91f0d49751be59e28ee5e5e4204ab12733",
        "commit_parents": [
            "4a7bc0859ace4e7c479d3739756d81f0c9cb7bc8"
        ]
    }
}
```
1. 200 ok
2. 204 no content - репозиторий существует,
    но не имеет состояния по умолчанию (ветки master)
3. 403 - нет прав доступа на просмотр
4. 404 - запрашиваемый репозиторий не существует

### 4.7 Получение информации о ветке по имени
Запрос: `/repo/{username}/{reponame}/branch/{branchname}` типа `GET`
Образец: `89.208.198.186:8080/api/v1/repo/nickeskov/codehub/branch/master`
Ответ:
```json
{
    "name": "master",
    "commit": {
        "commit_hash": "1e16727dfaf12c018cb4c069762d9ab9d62c0814",
        "commit_author_name": "UlianaBespalova",
        "commit_author_email": "43138516+UlianaBespalova@users.noreply.github.com",
        "commit_author_when": "2020-04-27T20:06:36+03:00",
        "committer_name": "GitHub",
        "committer_email": "noreply@github.com",
        "committer_when": "2020-04-27T20:06:36+03:00",
        "tree_hash": "88b5bc91f0d49751be59e28ee5e5e4204ab12733",
        "commit_parents": [
            "4a7bc0859ace4e7c479d3739756d81f0c9cb7bc8"
        ]
    }
}
```
1. 200 ok
3. 403 - нет прав доступа на просмотр
4. 404 - запрашиваемый репозиторий или ветка не существует

### 4.9 Получение контента файла в бинарном виде по ветке и пути
Запрос: `/repo/{username}/{reponame}/branch/{branchname}/tree/{path}` типа `GET`
Образец: `89.208.198.186:8080/api/v1/repo/nickeskov/codehub/branch/master/tree/some-file`
Ответ:
```text
Some binary data
```
1. 200 ok
2. 400 - некорректный запрос, т.е. запросили не файл, а директорию или что-то другое
3. 403 - нет прав доступа на просмотр
4. 404 - запрашиваемый пользователь/репозиторий/коммит/файл не существует
5. 413 - файл слишком большой и не может быть отправлен

### 4.10 Получение контента файла в бинарном виде по хешу коммита и пути
Запрос: `/repo/{username}/{reponame}/commit/{hash}/tree/{path}` типа `GET`
Образец: `89.208.198.186:8080/api/v1/repo/nickeskov/codehub/commit/4a7bc0859ace4e7c479d3739756d81f0c9cb7bc8/some-file`
Ответ:
```text
Some binary data
```
1. 200 ok
2. 400 - некорректный запрос, т.е. запросили не файл, а директорию или что-то другое
3. 403 - нет прав доступа на просмотр
4. 404 - запрашиваемый пользователь/репозиторий/коммит/файл не существует
5. 413 - файл слишком большой и не может быть отправлен

## 5. Issues
### 5.1 Создать issues
Запрос: `/api/v1/func/repo/{repoID}/issues` типа `POST`
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
Запрос: `/api/v1/func/repo/{repoID}/issues` типа `GET`
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
  "created_at": "2020-04-22T17:34:07.529Z",
  "author_login": "nickeskov",
  "author_image": "http://localhost:8080/static/image/avatar/default.jpg"
},
{
  "id": 45,
  "author_id": 5,
  "repo_id": 24,
  "title": "sfsfdsfsd",
  "message": "kekek",
  "label": "resolved",
  "is_closed": true,
  "created_at": "2020-04-22T17:34:07.529Z",
  "author_login": "nickeskov",
  "author_image": "http://localhost:8080/static/image/avatar/default.jpg"
}]
```
 2. 400 указана строка на месте repo_id
 3. 403 (нет прав на просмтор в этом репозитории)
 4. 404 (нет такого юзера или репозитория)
### 5.3 Обновить issues
Запрос: `/api/v1/func/repo/{repoID}/issues` типа `POST`
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
Запрос: `/api/v1/func/repo/{repoID}/issues` типа `DELETE`
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
Запрос: `/api/v1/func/repo/{repoID}/stars` типа `PUT`
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
4. 409 double vote
### 6.2 Список избранных репозиториев
Запрос: `/api/v1/func/repo/{login}/stars?limit=100&offset=2` типа `GET`
Тело:
```json
[{
    "id": "int",
    "owner_id": "int",
    "name": "string",
    "description": "string",
    "is_public": "bool",
    "is_fork": "bool",
    "created_at": "date",
    "author_login": "string"
},
{
    "id": "int",
    "owner_id": "int",
    "name": "string",
    "description": "string",
    "is_public": "bool",
    "is_fork": "bool",
    "created_at": "date",
    "author_login": "string"
}]
```
Ответ:
1. 200 ок
2. 404 нет юзера с таким логином
### 6.3 Список юзеров которые лайкнули репозиторий
Запрос: `/api/v1/func/repo/{repoID}/stars/users?limit=100&offset=2` типа `GET`
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
Запрос: `/api/v1/func/repo/{repoID}/news?limit=100&offset=2` типа `GET`
limit и offset опциональные параметры
Тело:
```json
[{
    "id": 5,
    "author_id": 5,
    "repo_id": 654433,
    "message": "dwrwrwsfdfrfe",
    "label": "",
    "date": "2020-04-26T19:02:10.073Z",
    "author_login": "nickeskov",
    "author_image": "http://localhost:8080/static/image/avatar/default.jpg"
},
{
    "id": 5,
    "author_id": 5,
    "repo_id": 654433,
    "message": "dwrwrwrfsffe",
    "label": "",
    "date": "2020-04-26T19:02:10.073Z",
    "author_login": "nickeskov",
    "author_image": "http://localhost:8080/static/image/avatar/default.jpg"
}]
```
Ответ:
1. 200 ок
2. 400 невалдиный repoid(строка)
3. 401 неавторизован
4. 403 нет доступа на просмотр новостей в данной репке
5. 404 нет репки с таким repoid
## 8. Search
### 8.1 Поиск юзеров среди всей базы
Запрос: `/api/v1/func/search/{params}?query=keksik&limit=100&offset=2` типа `GET`
limit и offset опциональные параметры
query это левая часть ника, по которой идет поиск
params=**allusers**
Тело:
Ответ:
```json
[
    {
        "id": 4,
        "owner_id": 1,
        "name": "loshok",
        "description": "",
        "is_fork": false,
        "created_at": "2020-05-17T22:23:25.185682+03:00",
        "is_public": true,
        "stars": 0,
        "forks": 0,
        "author_login": "keksik500",
        "parent_repository_info": {}
    }
]
```
```json
[
    {
        "id": 1,
        "name": "sddd",
        "login": "keksik500",
        "image": "/static/image/avatar/default.jpg",
        "email": "keksik@mail.ru",
        "created_at": "2020-05-17T22:18:18.711999+03:00"
    }
]
```
1. 200 ок
2. 400 невалдиные get параметры
### 8.2 Поиск репозиториев среди всей базы
Аналогично 8.1
params=**allrepo**
### 8.3 Поиск среди своих репозиториев
Аналогично 8.1
params=**myrepo**
### 8.4 Поиск среди избранных репозиториев(на которые юзер повесил звезду)
Аналогично 8.1
params=**starredrepo**
## 9. Fork
### 9.1 Ответвление репозитория к себе
Запрос: `/api/v1/func/repo/fork` типа `POST`
Указать либо id от чего форкаемся(from_repo_id), либо связку логин+название репы
new_name новое имя которое будет отображатся
при указании и id и логин+название репы приоритет будет у id
```json
{
  "from_repo_id": 45,
  "from_author_name": "pesntv",
  "from_repo_name": "muhtar",
  "new_name": "muhtar2"
}
```
Ответ:
1. 201 created
2. 400 невалидный json или сами данные
3. 401 unauthorized
4. 403 форкаем приватный репак
5. 409 уже форкали или форкаем свой же репак
## 10. PullRequest
### 10.1 Создание PullRequest
Запрос: `/api/v1/func/repo/pullrequests` типа `POST`
branch указываем название веток текстовое
```json
{
  "title": "kekemdaa",
  "author_id": 20,
  "from_repo_id": 440,
  "to_repo_id": 550,
  "branch_from": "dev",
  "branch_to": "master",
  "message": "kekeka"
}
```
Ответ:
```json
{
    "author_id": 40,
    "branch_from": "dev",
    "branch_to": "master",
    "closer_user_id": null,
    "created_at": "2020-05-26T20:17:46.776782Z",
    "from_author_login": "nickeskov",
    "from_repo_id": 61,
    "from_repo_name": "codehub",
    "id": 69,
    "is_accepted": false,
    "is_closed": false,
    "label": "",
    "message": "",
    "status": "ok",
    "title": "In Master",
    "to_author_login": "nickeskov",
    "to_repo_id": 61,
    "to_repo_name": "codehub"
}
```
Возможные значения поля status в ответе (по этим полям не стоит определять закрыт pull request или нет):
* "" - Это некорректное значение, если такое значение пришло, что-то явно не так
* "ok" - Сервер сам **может** слить ветки (pull request открыт)
* "error" - Произошла ошибка (pull request закрыт)
* "merged" - Слит (pull request закрыт)
* "rejected" - Отклонён (pull request закрыт)
* "conflict" - Сервер **не может** слить ветки из-за конфликта (pull request открыт)
* "up_to_date" - Сервер **не может** слить ветки, так как эти изменения уже есть в целевой ветке (pull request открыт)
* "no_changes" - Сервер сам **может** слить ветки, но не обнаружено изменений (pull request открыт)
* "bad_to_branch" - Target ветки не существует (pull request закрыт)
* "bad_from_branch" - Source ветки не существует (pull request закрыт)

1. 201 created
2. 400 некорректный json или сами данные
3. 401 unauthorized
4. 403 приватный репозиторий или пользователь не имеет права на это действие
5. 409 возник конфликт при создании, т.е. пытаемся слить ветки, которые
не имеют общего предка или общий предок это ветка, которую мы пытаемся слить в целевую
(например мы master сливаем в dev, хотя в dev уже есть все изменения master)

### 10.2 Получение списка всех PullRequest в наш репо или из нашего репока
Запрос: `/api/v1/func/repo/{repoID}/pullrequests/in?limit=2&offset=0` типа `GET`
Запрос: `/api/v1/func/repo/{repoID}/pullrequests/out?limit=2&offset=0` типа `GET`
limit offset - лимит и смещение
title from_repo to_repo branch_to branch_from обязательны
Ответ:
```json
[
    {
        "author_id": 40,
        "branch_from": "dev",
        "branch_to": "master",
        "closer_user_id": null,
        "created_at": "2020-05-26T20:17:46.776782Z",
        "from_author_login": "nickeskov",
        "from_repo_id": 61,
        "from_repo_name": "codehub",
        "id": 69,
        "is_accepted": false,
        "is_closed": false,
        "label": "",
        "message": "",
        "status": "ok",
        "title": "In Master",
        "to_author_login": "nickeskov",
        "to_repo_id": 61,
        "to_repo_name": "codehub"
    }
]
```
Возможные значения поля status в ответе (по этим полям не стоит определять закрыт pull request или нет):
* "" - Это некорректное значение, если такое значение пришло, что-то явно не так
* "ok" - Сервер сам **может** слить ветки (pull request открыт)
* "error" - Произошла ошибка (pull request закрыт)
* "merged" - Слит (pull request закрыт)
* "rejected" - Отклонён (pull request закрыт)
* "conflict" - Сервер **не может** слить ветки из-за конфликта (pull request открыт)
* "up_to_date" - Сервер **не может** слить ветки, так как эти изменения уже есть в целевой ветке (pull request открыт)
* "no_changes" - Сервер сам **может** слить ветки, но не обнаружено изменений (pull request открыт)
* "bad_to_branch" - Target ветки не существует (pull request закрыт)
* "bad_from_branch" - Source ветки не существует (pull request закрыт)

Nullable поля:
* author_id
* closer_user_id
* from_repo_id
* from_author_login
* from_repo_name
* to_repo_id
* to_author_login
* to_repo_name

1. 200 ok
2. 400 невалидный json или сами данные
3. 401 unauthorized
4. 403 лезем в приватный репак
### 10.3 Принять PullRequest
Запрос: `/api/v1/func/repo/pullrequests` типа `PUT`

Ответ:
```json
{
  "id": 20,
  "to_repo_id": 30
}
```
1. 200 ok
2. 400 невалидный json или сами данные(id не существуют)
3. 401 unauthorized
4. 403 лезем в чужой реквест
5. 409 имеется конфликт и сервер не может слить изменения самостоятельно
### 10.4 Закрыть PullRequest
Запрос: `/api/v1/func/repo/pullrequests` типа `DELETE`
Ответ:
```json
{
  "id":20,
  "to_repo_id": 30
}
```
1. 200 ok
2. 400 невалидный json или сами данные(id не существуют)
3. 401 unauthorized
4. 403 лезем в какой-то чужой реквест
### 10.5 Все PullRequest юзера
Запрос: `/api/v1/user/pullrequests?limit=1&offset=0` типа `GET`
Ответ:
```json
[
    {
        "author_id": 40,
        "branch_from": "dev",
        "branch_to": "master",
        "closer_user_id": null,
        "created_at": "2020-05-26T20:17:46.776782Z",
        "from_author_login": "nickeskov",
        "from_repo_id": 61,
        "from_repo_name": "codehub",
        "id": 69,
        "is_accepted": false,
        "is_closed": false,
        "label": "",
        "message": "",
        "status": "ok",
        "title": "In Master",
        "to_author_login": "nickeskov",
        "to_repo_id": 61,
        "to_repo_name": "codehub"
    }
]
```
Возможные значения поля status в ответе (по этим полям не стоит определять закрыт pull request или нет):
* "" - Это некорректное значение, если такое значение пришло, что-то явно не так
* "ok" - Сервер сам **может** слить ветки (pull request открыт)
* "error" - Произошла ошибка (pull request закрыт)
* "merged" - Слит (pull request закрыт)
* "rejected" - Отклонён (pull request закрыт)
* "conflict" - Сервер **не может** слить ветки из-за конфликта (pull request открыт)
* "up_to_date" - Сервер **не может** слить ветки, так как эти изменения уже есть в целевой ветке (pull request открыт)
* "no_changes" - Сервер сам **может** слить ветки, но не обнаружено изменений (pull request открыт)
* "bad_to_branch" - Target ветки не существует (pull request закрыт)
* "bad_from_branch" - Source ветки не существует (pull request закрыт)

Nullable поля:
* author_id
* closer_user_id
* from_repo_id
* from_author_login
* from_repo_name
* to_repo_id
* to_author_login
* to_repo_name

1. 200 ok
2. 400 невалидный json или сами данные(id не существуют)
3. 401 unauthorized

### 10.6 Получение PullRequest по id
Запрос: `/api/v1/func/repo/pullrequest/{id}/diff` типа `GET`

Параметры:
- **id**  - id pullrequest

Ответ:
```json
[
    {
        "author_id": 40,
        "branch_from": "dev",
        "branch_to": "master",
        "closer_user_id": null,
        "created_at": "2020-05-26T20:17:46.776782Z",
        "from_author_login": "nickeskov",
        "from_repo_id": 61,
        "from_repo_name": "codehub",
        "id": 69,
        "is_accepted": false,
        "is_closed": false,
        "label": "",
        "message": "",
        "status": "ok",
        "title": "In Master",
        "to_author_login": "nickeskov",
        "to_repo_id": 61,
        "to_repo_name": "codehub"
    }
]
```
Возможные значения поля status в ответе (по этим полям не стоит определять закрыт pull request или нет):
* "" - Это некорректное значение, если такое значение пришло, что-то явно не так
* "ok" - Сервер сам **может** слить ветки (pull request открыт)
* "error" - Произошла ошибка (pull request закрыт)
* "merged" - Слит (pull request закрыт)
* "rejected" - Отклонён (pull request закрыт)
* "conflict" - Сервер **не может** слить ветки из-за конфликта (pull request открыт)
* "up_to_date" - Сервер **не может** слить ветки, так как эти изменения уже есть в целевой ветке (pull request открыт)
* "no_changes" - Сервер сам **может** слить ветки, но не обнаружено изменений (pull request открыт)
* "bad_to_branch" - Target ветки не существует (pull request закрыт)
* "bad_from_branch" - Source ветки не существует (pull request закрыт)

Nullable поля:
* author_id
* closer_user_id
* from_repo_id
* from_author_login
* from_repo_name
* to_repo_id
* to_author_login
* to_repo_name

1. 200 ok
2. 400 невалидный json или сами данные(id не существуют)
3. 401 unauthorized

### 10.7 PullRequest diff
Запрос: `/api/v1/func/repo/pullrequest/{id}/diff` типа `GET`

Параметры:
- **id**  - id pull request

Ответ:
```json
{
  "diff": "string",
  "status":"ok"
}
```
Примечания:
- Поле diff может быть пустым

Возможные значения поля status в ответе (по этим полям не стоит определять закрыт pull request или нет):
* "" - Это некорректное значение, если такое значение пришло, что-то явно не так
* "ok" - Сервер сам **может** слить ветки (pull request открыт)
* "error" - Произошла ошибка (pull request закрыт)
* "merged" - Слит (pull request закрыт)
* "rejected" - Отклонён (pull request закрыт)
* "conflict" - Сервер **не может** слить ветки из-за конфликта (pull request открыт)
* "up_to_date" - Сервер **не может** слить ветки, так как эти изменения уже есть в целевой ветке (pull request открыт)
* "no_changes" - Сервер сам **может** слить ветки, но не обнаружено изменений (pull request открыт)
* "bad_to_branch" - Target ветки не существует (pull request закрыт)
* "bad_from_branch" - Source ветки не существует (pull request закрыт)

1. 200 ok
2. 400 - некорректный id (не удалось преобразовать в число)
3. 404 не найден такой pull request
