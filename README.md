# Personal-data-API
Тестовое задание Junior Golang Developer
Effective Mobile
Реализовать сервис, который будет получать по API ФИО, из открытых API обогащать
ответ наиболее вероятными возрастом, полом и национальностью и сохранять данные в
БД. По запросу выдавать инфу о найденных людях.
Необходимо реализовать следующее:
1 Выставить rest методы - [x]
- Для получения данных с различными фильтрами и пагинацией
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
- Для удаления по идентификатору
    Параметры запроса:
    - ID (индентификатор записи)
- Для изменения сущности
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
- Для добавления новых людей в формате
    Параметры запроса:
    ```json
    {
      "name": "string",
      "patronymic": "string",
      "surname": "string"
    }
  ```
2 Корректное сообщение обогатить - [x]
- Возрастом - https://api.agify.io/?name=Dmitriy
- Полом - https://api.genderize.io/?name=Dmitriy
- Национальностью - https://api.nationalize.io/?name=Dmitriy
3 Обогащенное сообщение положить в БД postgres (структура БД должна быть создана
путем миграций)
4 Покрыть код debug- и info-логами
5 Вынести конфигурационные данные в .env
6 Сгенерировать сваггер на реализованное API
Доступен по хосту http://localhost:8080
