# Tender Managemant API

Этот репозиторий содержит опиcание написанного API в рамках тестового задания на стажировку в компанию "Авито" на позицию backend разработчика.

*Небольшое примечание: так как я не получил ответы на свои вопросы и не мог задать остальные, в спорных по моему мнению местах я действовал исходя из своих мыслей и логики. **Задание не является завершенным, так как ~5 дней оказалось недостаточно для такого объема работы.*** 

# Table Of Contents

- [Tender Managemant API](#tender-managemant-api)
- [Table Of Contents](#table-of-contents)
  - [Описание задания](#описание-задания)
  - [Проектирование базы данных](#проектирование-базы-данных)
  - [Кодогенерация из openapi](#кодогенерация-из-openapi)
  - [Проектирование архитектуры сервера](#проектирование-архитектуры-сервера)
  - [Конфигурация сервера и создание миграций](#конфигурация-сервера-и-создание-миграций)
  - [База данных](#база-данных)
  - [Валидации](#валидации)
  - [Логгирование](#логгирование)
  - [Слой API](#слой-api)
  - [Запуск](#запуск)
  - [Тестирование](#тестирование)
    - [GET api/ping](#get-apiping)
    - [POST api/tenders/new](#post-apitendersnew)
    - [GET api/tenders/my](#get-apitendersmy)
    - [GET api/tenders/:tenderId/status](#get-apitenderstenderidstatus)
    - [PUT api/tenders/:tenderId/status](#put-apitenderstenderidstatus)
    - [PATCH api/tenders/:tenderId/edit](#patch-apitenderstenderidedit)
    - [PUT api/tenders/:tenderId/rollback/:version](#put-apitenderstenderidrollbackversion)
    - [POST api/bids/new](#post-apibidsnew)
    - [GET api/bids/my](#get-apibidsmy)
    - [GET api/bids/:tenderId/list](#get-apibidstenderidlist)
    - [GET api/bids/:bidId/status](#get-apibidsbididstatus)
    - [PUT api/bids/:bidId/status](#put-apibidsbididstatus)
    - [PATCH api/bids/:bidId/edit](#patch-apibidsbididedit)
    - [PUT api/bids/:bidId/sumbit\_decision](#put-apibidsbididsumbit_decision)
    - [PUT api/bids/:bidId/feedback](#put-apibidsbididfeedback)
    - [PUT api/bids/:bidId/rollback/:version](#put-apibidsbididrollbackversion)
    - [GET api/bids/:tenderId/reviews](#get-apibidstenderidreviews)
  - [Использованные технологии:](#использованные-технологии)
  - [Вывод](#вывод)


## Описание задания

В задании требовалось разработать API для менеджмента тендеров и предложений с возможностями редактирования, изменения статуса, роллбэков к старым версиям и тд. **Все возможности будут описаны далее.**

## Проектирование базы данных

Было решено начать выполнение с проектирования базы данных. Прямо на этот моменте я впал в ступор, так как небольшая часть базы данных была приведена в `задание/README.md`, где, почему-то, все `id` были в форматах `int`, когда в `задание/openapi.yml`, например, `organizationId` был в формате `UUID`.

Подход был следующий: мы уже имеем таблицы `employee`, `organization`, `organization_responsible`. Нам необходимо хранить **предложения** (`bids`) и **тендеры** (`tenders`). Поэтому создаем таблицы для предложений и тендеров (`bids`, `tenders`). Внимательно просмотрев спецификацию `задание/openapi.yml` стало понятно, что необходимо так же хранить информацию о прошлых версиях предложений и тендеров, чтобы можно было делать откат к предыдущим версиям. Так появились таблицы `tender_versions` и `bid_versions`.

Далее стало понятно, что надо хранить решения о предложениях и фидбэки. Появились таблицы `bid_decisions` и `bid_feedbacks`. 

Некоторые атрибуты таблиц имеют кастомные типы данных, а именно список ENUM. Всего их 6:

```sql
CREATE TYPE organization_type AS ENUM ('IE', 'LLC', 'JSC');
CREATE TYPE tender_service_type AS ENUM ('Construction', 'Delivery', 'Manufacture');
CREATE TYPE tender_status AS ENUM ('Created', 'Published', 'Closed');
CREATE TYPE bid_status AS ENUM ('Created', 'Published', 'Canceled', 'Approved', 'Rejected');
CREATE TYPE bid_author_type AS ENUM ('Organization', 'User');
CREATE TYPE bid_decision AS ENUM ('Approved', 'Rejected');
```

Все эти типы данных были указаны в спецификации.

Дополнительно было создано 2 индекса для оптимизации поиска текущей версии:

```sql
CREATE INDEX idx_tender_versions_current ON tender_versions(tender_id, is_current);
CREATE INDEX idx_bid_versions_current ON bid_versions(bid_id, is_current);
```

Итоговый вариант полученной базы данных в виде диаграммы приведен ниже:

![Database Diagram](resources/db-diagram.png)

## Кодогенерация из openapi

Когда диаграмма базы данных была получена, было принято решение о написании самого сервера. Конечно же, если имееется спецификация в формате openapi, мы можем сгенерировать из нее код.

Для такой задачи был выбран `oapi-codegen`, которые поддерживает выбор фреймворка для написания http-сервера и который поддерживает `v3 openapi`.

## Проектирование архитектуры сервера

Когда код был сгенрирован, он весь содержался в одном файле, что, конечно же, не является хорошей практикой. 

Изначально планировалась Чистая архитектура по всем канонам дядюшки Боба, однако я все еще помню споры в разных чатах на эту тему, 

Чтобы не допустить крупной ошибки, была выбрана более простая архитектура с разделением на 3 слоя: 

1. **API**. Здесь находятся все http обработчики и по сути сам сервер. Когда приходит на какую-либо ручку запрос, из него извлекаются все необходимые параметры и формируется структура параметров, которая передается на следующий слой.
2. **Service**. Этот слой отвечает за обработку переданных от API параметров. Сначала идет валидация параметров, чтобы предотварить некоторые ошибки, не обращаясь к следующему слою, после валидации идет вызов следующего слоя:
3. **Database**. По сути, это центр нашего сервера. Здесь лежит бизнес логика, модели и работа с базой данных. 

Таким образом, получаем следующую модель работы сервера:

![Server Diagram](resources/server-diagram.png)

## Конфигурация сервера и создание миграций

В задании требовалось, чтобы для конфигурации сервера использовались переменные окружения. Удобным и простым вариантом является использование `.env` файла. 

Так как нельзя хранить секреты в репозиториях, был создан `.env.example`, где описана структура итогового `.env` файла. ***(не запушен в ветку)***

Перейдем к миграциям. Всего есть 2 версии миграций. Первая создает основные таблицы и типы данных + индексы. Вторая версия выполняет запросы для заполнения данными некоторых таблиц: `employee`, `organization`, `organization_responsible`. Это сделано, потому что в спецификации никак не описан процесс создания пользователей, организаций и ответственных за организации. К тому же, этот сервис не предназначен для такого функционала, потому что он управляет тендерами и предложениями.

## База данных

В задании требовалось использовать `Postgres`, что и было сделано, однако остается место для маневра, так как используется интерфейс `Database`, который содержит все методы для работы с БД.

Следование принципам **SOLID** помогло спроектировать базу данных так, что нет зависимости от конкретной реализации.

Сейчас имеется реализация только Postgres, однако можно добавить реализацию, например, `MySQL` и добавить переменную окрудения `DB_DRIVER`, где указывать, какой драйвер использовать.

Это поможет легко сменить базу данных при необходимости.

## Валидации

Чтобы слои были отделены друг от друга, было принято решение добавить между `API` и `Database` слой сервисов, которые получают параметры от `API`, валидирует их и выполняют запросы к БД.

Добавление слоя сервисов помогает отделить проверку получаемых параметров от базы данных и от слоя `API`.

## Логгирование

По схожему принципу с пунктом "База данных" для логгера был написан интерфейс и реализован логгер через `zap`. Есть возможность сменить логгер посредсвом имплементации интерфейса `Logger` 

## Слой API

Это слой, куда приходят все запросы извне. Именно здесь формируются параметры для запросов к БД, которые далее передаются в сервисный слой, а затем в слой БД.

Все поддерживаемые `routes`:

```go
func (s *server) RegisterHandlers() {
	s.r.GET("/api/bids/my", s.GetUserBids)
	s.r.POST("/api/bids/new", s.CreateBid)
	s.r.PATCH("/api/bids/:bidId/edit", s.EditBid)
	s.r.PUT("/api/bids/:bidId/feedback", s.SubmitBidFeedback)
	s.r.PUT("/api/bids/:bidId/rollback/:version", s.RollbackBid)
	s.r.GET("/api/bids/:bidId/status", s.GetBidStatus)
	s.r.PUT("/api/bids/:bidId/status", s.UpdateBidStatus)
	s.r.PUT("/api/bids/:bidId/submit_decision", s.SubmitBidDecision)
	s.r.GET("/api/bids/:tenderId/list", s.GetBidsForTender)
	s.r.GET("/api/bids/:tenderId/reviews", s.GetBidReviews)
	s.r.GET("/api/ping", s.CheckServer)
	s.r.GET("/api/tenders", s.GetTenders)
	s.r.GET("/api/tenders/my", s.GetUserTenders)
	s.r.POST("/api/tenders/new", s.CreateTender)
	s.r.PATCH("/api/tenders/:tenderId/edit", s.EditTender)
	s.r.PUT("/api/tenders/:tenderId/rollback/:version", s.RollbackTender)
	s.r.GET("/api/tenders/:tenderId/status", s.GetTenderStatus)
	s.r.PUT("/api/tenders/:tenderId/status", s.UpdateTenderStatus)
}
```

Все эти пути были описани в спецификации. Реализованы обработчики были в сторогом следовании спецификации.

## Запуск

Для запуска сервера необходимо наличие:
1. Docker
2. docker-compose
3. Golang 1.23+ версии

Шаги следующие:

1. Клонируем репозиторий:
```sh
git clone https://github.com/0x0FACED/tender-service.git
```
2. Делаем `cd` в корневую дирекорию проекта:
```sh
cd .../tender-service
```
1. Создаем `.env` файл по спецификации:
```sh
SERVER_ADDRESS=localhost:8080

POSTGRES_CONN=postgres://yourusername:yourpassword@localhost:5432/yourdatabase?sslmode=disable
POSTGRES_USERNAME=yourusername
POSTGRES_PASSWORD=yourpassword
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DATABASE=yourdatabase
```
4. Создаем базу данных, если ее нет и меняем данные в docker-compose файле на свои
5. Запускаем через `docker-compose`:
```sh
sudo docker-compose up --build
```
ИЛИ
```bash
sudo docker-compose up -d
```
Для запуска в фоновом режиме

*Далее будут описаны все обработчики и приведены скриншоты тестирования при помощи `Postman`.*

## Тестирование

### GET api/ping

Я не дождался ответа по почте, где я интересовался, а стоит ли пинговать дополнительно базу данных. Решил, что было бы неплоохо проверять, работает ли база данных, принимает ли она запросы, поэтому такая проверка была добавлена.

![api/ping](resources/ping.png)

Если сервер готов принимать запросы, в том числе установлено соединение с базой данных, тогда сервер вернет `"OK"` и `StatusCode = 200`

Если же, например, подключение было разорвано, тогда сервер вернет ошибку с кодом 500 и сообщением, что база данных не отвечает:

![api/ping](resources/ping_err.png)

*Все дальнейшие запросы будут выполняться к серверу, в базе данных которого нет данных, кроме тех, которые применяются во второй миграции.*

### POST api/tenders/new

Начнем с создания тендера.

На скриншоте ниже показан пример успешного создания тендера, если все данные правильно переданы.

![api/tenders/new](resources/tender_new.png)

На скриншоте ниже показан пример с возвратом 403 от сервера, так как переданный organizationId не совпал с ответственным за организацию в БД.

![api/tenders/new_fail](resources/tender_new_err.png)

### GET api/tenders/my

Создадим еще несколько тендеров и попытаемся их получить в ответе:

![api/tenders/my](resources/tenders_my.png)

Все тендеры были созданы от лица `john_doe`, теперь если мы попробуем получить тендеры от лица другого `employee`, то ничего не получим:

![api/tenders/my_null](resources/tenders_my_null.png)

### GET api/tenders/:tenderId/status

Просто получаем статус тендера. Ничего больше.

![api/tenders/:tenderId/status](resources/tenders_get_status.png)

Если мы изменим `username`, то будет ответ, что такой тендер не найден, потому что проверяется причастность к созданию тендера по username, а точнее по ответственному за организацию.

### PUT api/tenders/:tenderId/status

Изменяем статус и возращаем обновленный тендер.

![api/tenders/:tenderId/status](resources/tenders_put_status.png)

### PATCH api/tenders/:tenderId/edit

Передаем в теле запроса новые имя, описание и `serviceType`.

![api/tenders/:tenderId/edit](resources/tenders_edit.png)

### PUT api/tenders/:tenderId/rollback/:version

Перед тем как сделать rollback, сделаем еще несколько изменений тендера, который был приведен в предыдущем пункте.

В результате в таблице версий тендеров будут такие записи:

![tender_versions](resources/tenders_versions.png)

Теперь сделаем откат к, например, 3-й версии:

![api/tenders/:tenderId/rollback/:version](resources/tenders_rollback_v3.png)

В спецификации указано, что такой откат - новая версия, следовательно теперь последняя версия под номером 6.

### POST api/bids/new

Создавать предложения можно только к существующим тендерам.

![api/bids/new](resources/bids_new.png)

### GET api/bids/my

![api/bids/my](resources/bids_my.png)

### GET api/bids/:tenderId/list

![api/bids/:tenderId/list](resources/bids_list.png)

### GET api/bids/:bidId/status

Просто получаем статус бида :D

![api/bids/:bidId/status](resources/bids_get_status.png)

*Важно: этот запрос был выполнен после апдейта статуса.*

### PUT api/bids/:bidId/status

Меняем статус бида.

![api/bids/:bidId/status](resources/bids_put_status.png)

### PATCH api/bids/:bidId/edit

Редактировать бид может только его создатель.

![api/bids/:bidId/edit](resources/bids_edit.png)

### PUT api/bids/:bidId/sumbit_decision

Принятие решение - инкремент версии. 

![api/bids/:bidId/submit_decision](resources/bids_submit_decision.png)

### PUT api/bids/:bidId/feedback

![api/bids/:bidId/feedback](resources/bids_feedback.png)

![bid_feedbacks](resources/bid_feedbacks.png)

### PUT api/bids/:bidId/rollback/:version

По аналогии с роллбэком у тендера работает.

### GET api/bids/:tenderId/reviews

![api/bids/:tenderId/feedback](resources/bids_reviews.png)

Не работает, так как я допустил ошибку при проектировании базы данных. Надо было добавить поле `bid_id` в `bid_feedbacks`, чтобы можно было связать отзывы и предложения, чего я не сделал, увы, а сил и времени уже нет.

## Использованные технологии:

1. [golang-migrate](https://github.com/golang-migrate/migrate)
2. [godotenv](https://github.com/joho/godotenv)
3. [echo](https://github.com/labstack/echo)
4. [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
5. [zap](https://github.com/uber-go/zap)
6. [Postman](https://www.postman.com/)
7. [Docker](https://www.docker.com/)
8. [docker-compose](https://docs.docker.com/compose/)

## Вывод

В целом, задание было частично выполнено. Многие моменты в коде было не доделаны, например, так и остался `context.TODO()` при вызове методов сервисов, а можно сделать, например, хотя бы контекст с таймаутом на 10сек. Это условно. 

Многие комментарии остались, потому что они была расставлены `oapi-codegen`, я решил не чистить их. Свои комментарии почти все почистил, но некоторые целенаправленно оставил.

Также были допущены ошибки при проектировании базы данных. Правятся они созданием 3-й версии миграции.

