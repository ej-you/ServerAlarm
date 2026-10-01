# Server Alarm

Система уведомления пользователей о превышении порога температуры в серверной.
Основана на [ntfy](https://ntfy.sh/) сервисе.

## Описание

Приложение состоит из двух сервисов:

- nfty
- api

`ntfy` используется для отправки уведомлений на устройства.

`api` опрашивает базу данных, в которую пишет датчик, и посылает запрос `ntfy` на отправку сообщения.

## Конфигурация

### Файл [./api/config.yml](./api/config.yml)

Все параметры подробно описаны в самом файле.

### Файл [./ntfy/.env](./ntfy/.env)

```dotenv
NTFY_BASE_URL='https://external.com'

NTFY_CACHE_FILE='/var/lib/ntfy/cache.db'
NTFY_AUTH_FILE='/var/lib/ntfy/auth.db'
NTFY_ATTACHMENT_CACHE_DIR='/var/lib/ntfy/attachments'

# is reverse-proxy in use
NTFY_BEHIND_PROXY=true
NTFY_AUTH_DEFAULT_ACCESS='deny-all'

# admin password
NTFY_PASSWORD='your-password'
```

### Файл [./api/.env](./api/.env)

```dotenv
DB_PASSWORD="test_password"

# ntfy admin token
NTFY_TOKEN="tk_pmxwhj94glo0kiwutzwa0zp8lb056"
```

## Запуск и настройка

### Запуск `nfty` сервиса

> ! _Вначале запустим только `ntfy`_
> _После его настройки запустим `api`_

```shell
docker compose up -d ntfy
```

### Добавление админа

> ! _Далее для админа используется имя `root`._
> _Если хотите изменить имя админа, то используйте другое имя везде, где указан `root`._
>
> ! _Пароль админа можно поменять в [.env](./ntfy/.env) файле для `ntfy`_
> _Для этого нужно изменить значение ключа `NTFY_PASSWORD`_

```shell
docker compose exec -it ntfy ntfy user add --role=admin root
```

### Создание токена для админа

```shell
docker compose exec -it ntfy ntfy token add root
```

### Прописываем токен в [.env](./api/.env) для `api`

```shell
nano ./api/.env
```

Ключу `NTFY_TOKEN` зададим значение в виде полученного токена

```dotenv
NTFY_TOKEN='tk_pmxwhj94glo0kiwutzwa0zp8lb056'
```

### Запуск `api`

```shell
docker compose up --build -d api
```

## Администрирование

### Создание пользователя

```shell
docker compose exec -it ntfy ntfy user add your-username
```

> ! _Также для работы с пользователями есть другие функции._
>
> ! _Подсказку по ним можно получить командой_
> \- `docker compose exec -it ntfy ntfy user --help`

### Изменение конфига [./api/config.yml](./api/config.yml)

Для применения изменений в этом файле при уже запущенном сервисе `api` его необходимо перезапустить.
Для этого используйте команду:

```shell
docker compose restart api
```
