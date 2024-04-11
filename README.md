# avito_banners
Для запуска контейнера:
docker-compose up -d

Сделать:
- Кэш баннеров
- Историчность баннеров
- Запуск сервиса в контейнере докера
- интеграционный тест
- нагрузочный тест
- Описание readme
- описание проблем и их решений:
  (в частности решение с чекэкзистом и с токенами авторизации)
- Пооставлять комментарии в коде

Вопросы:
- По условию задачи предположил, что токены админа и юзера - константы. 
Хотя логичнее использовать уникальные токены для каждого юзера и админа с набором полей (админ или юзер) и хранить их например в редиске
- Предполагается, что один и тот же тэг не может находиться в списке тэгов двух разных баннеров с одинаковыми фичами
- Предполагается, что при обновлении баннера изменению подлежит только содержимое баннера, а тэг и фича остаются неизменными, 
иначе при изменении тэга или фичи фактически получим новый баннер, т.к. по условию задачи тэги и фича однозначно определяют баннер
- При удалении баннера необходимо ли удалять все его версии? Преположу, что необходимо удалить конкретный баннер по его id, без удаления его версий

Использую golangci-lint linters:
- errcheck: errcheck is a program for checking for unchecked errors in Go code. These unchecked errors can be critical bugs in some cases [fast: false, auto-fix: false]
- gosimple (megacheck): Linter for Go source code that specializes in simplifying code [fast: false, auto-fix: false]
- govet (vet, vetshadow): Vet examines Go source code and reports suspicious constructs. It is roughly the same as 'go vet' and uses its passes. [fast: false, auto-fix: false]
- ineffassign: Detects when assignments to existing variables are not used [fast: true, auto-fix: false]
- staticcheck (megacheck): It's a set of rules from staticcheck. It's not the same thing as the staticcheck binary. The author of staticcheck doesn't support or approve the use of staticcheck as a library inside golangci-lint. [fast:
false, auto-fix: false]
- unused (megacheck): Checks Go code for unused constants, variables, functions and types [fast: false, auto-fix: false]

