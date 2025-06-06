# Personal-data-API
Тестовое задание Junior Golang Developer
Effective Mobile
Реализовать сервис, который будет получать по API ФИО, из открытых API обогащать
ответ наиболее вероятными возрастом, полом и национальностью и сохранять данные в
БД. По запросу выдавать инфу о найденных людях.
Необходимо реализовать следующее:

Выставить rest методы
- [x] Для получения данных с различными фильтрами и пагинацией
  Параметры запроса:
    - Name (имя человека для поиска)
    - Surname (фамилия человека для поиска)
    - Patronymic (отчество человека для поиска)
    - Limit (максимальное количество выдаваемых записей)
    - Offset (сдвиг по записям)
    - sortBy (сортировка по параметру. Возможные поля для сортировки: id, name, surname, patronymic, age, gender, nationality)
    - sortOrder (сортировка ASC и DESC)
      Ответ сервиса:
    ```json
  {
    "data": [
      {
        "id": "int",
        "name": "string",
        "surname": "string",
        "patronymic": "string",
        "age": "int",
        "gender": "string",
        "nationality": "string"
      } 
    ],
    "total": "int"
  }
  ```
- [x] Для удаления по идентификатору
  Параметры запроса:
    - ID (индентификатор записи)
- [x] Для изменения сущности
  Параметры запроса:
    - ID (индентификатор записи)
    ```json
     {
      "age": 0,
      "gender": "string",
      "id": 0,
      "name": "string",
      "nationality": "string",
      "patronymic": "string",
      "surname": "string"
    }
  ```
- [x] Для добавления новых людей в формате
  Параметры запроса:
    ```json
    {
      "name": "string",
      "patronymic": "string",
      "surname": "string"
    }
  ```
Корректное сообщение обогатить
- [x] Возрастом - https://api.agify.io/?name=Dmitriy
- [x] Полом - https://api.genderize.io/?name=Dmitriy
- [x] Национальностью - https://api.nationalize.io/?name=Dmitriy
- [x] Обогащенное сообщение положить в БД postgres (структура БД должна быть создана
  путем миграций)
- [x] Покрыть код debug- и info-логами
- [x] Вынести конфигурационные данные в .env
  Пример .env:
```
Postgres_User=postgres
Postgres_Password=7321968
Postgres_DbName=test
Postgres_Host=localhost
Postgres_Port=5432
Server_Port=8080
DB_Host=db
Migration_Flag=yes
Migration_Path=file://internal/database/migrations
Agify_URL=https://api.agify.io/?name=
Genderize_URL=https://api.genderize.io/?name=
Nationalize_URL=https://api.nationalize.io/?name=
```
- [x] Сгенерировать сваггер на реализованное API
  Доступен по хосту http://localhost:8080/swagger/index.html

## Запуск сервиса
Заполнить .env файл переменными
```
Postgres_User=testUser
Postgres_Password=testPassword
Postgres_DbName=testDBName
Postgres_Host=localhost
Postgres_Port=5432
Server_Port=8080
DB_Host=db
Migration_Flag=yes
Migration_Path=file://internal/database/migrations
```
``` 
docker-compose up --build
 ```