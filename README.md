# avito_banners
Для сборки и запуска сервиса выполнить команду:
`docker-compose up --build` и дождаться в логах сервиса строчки: `HTTP server started on addr 0.0.0.0:8000`

Для выключения:
`docker-compose down`

# Проблемы, с которыми столкнулся и их решения
- По условию задачи предположил, что токены админа и юзера - константы, следственно положил их в конфиги.
Хотя логичнее использовать уникальные токены для каждого юзера и админа с набором полей (админ или юзер) и хранить их например в редиске.
- Предполагаю, что один и тот же тэг не может находиться в списке тэгов двух разных баннеров с одинаковыми фичами, т.к. по API в запросе `/get-banner`
приходит один `tag_id` и `feature_id`, по которым необходимо выдать ОДИН баннер, иначе их могло бы быть несколько.
- Предполагаю, что при обновлении баннера изменению подлежит только содержимое баннера, а тэг и фича остаются неизменными,
иначе при изменении тэга или фичи фактически получим новый баннер, т.к. по условию задачи тэги и фича однозначно определяют баннер.
- При удалении баннера необходимо ли удалять все его версии? Предположу, что необходимо удалить конкретный баннер по его id, без удаления его версий,
иначе мы можем потерять содержимое баннеров других версий.
- Также написал небольшой скрипт для генерации баннеров в базе [filldata_test.go](internal%2Ftests%2Ffilldata_test.go).

Примеры запросов из Postman [requestsexample.md](requestsexample.md)

# По доп. заданиям

1) Сервис адаптирован для увеличения количества фичей и тэгов, добавлял >10000, прожует и 100000 и больше, но и время ответа увеличится при запросах в базу,
хотя при работе с кэшем время измениться не должно.
2) Результаты нагрузочных тестов см. в [loadtests.md](loadtests.md)
3) Реализовал версионность баннеров - при обновлении баннера создаем новую версию с обновленным `updated_at`,
актуальной версией считается та, у которой последний `updated_at`. Всего хранится 3 версии. При выборе более старой версии - обновляем `updated_at` у целевой версии.
4) Эндпоинт реализован, однако не стал реализовывать механизм отложенного действия - можно запустить удаление в горутине, но тогда в респонс ничего не отдам
и будет неизвестен результат удаления, и время выполнения запроса < 50 мс без нагрузки.
5) Написаны интеграционные тесты дял каждого эндпоинта в [banner_test.go](internal%2Ftests%2Fbanner_test.go). Необходимо запускать весь пакет тестов,
т.к. во время выполнения тестов выполняется цикл создания, обновления, получения, удаления и т.п. Перед запуском тестов таблицу в базе пересоздать.
6) Использую golangci-lint linters:
- errcheck: errcheck is a program for checking for unchecked errors in Go code. These unchecked errors can be critical bugs in some cases [fast: false, auto-fix: false]
- gosimple (megacheck): Linter for Go source code that specializes in simplifying code [fast: false, auto-fix: false]
- govet (vet, vetshadow): Vet examines Go source code and reports suspicious constructs. It is roughly the same as 'go vet' and uses its passes. [fast: false, auto-fix: false]
- ineffassign: Detects when assignments to existing variables are not used [fast: true, auto-fix: false]
- staticcheck (megacheck): It's a set of rules from staticcheck. It's not the same thing as the staticcheck binary. The author of staticcheck doesn't support or approve the use of staticcheck as a library inside golangci-lint. [fast:
false, auto-fix: false]
- unused (megacheck): Checks Go code for unused constants, variables, functions and types [fast: false, auto-fix: false]