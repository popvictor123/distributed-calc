# Распределённый Калькулятор

Этот проект реализует распределённый калькулятор, который может вычислять арифметические выражения. Он состоит из двух основных компонентов:

*   **Оркестратор:** Управляет распределением вычислительных задач и отслеживает ход выполнения выражений. Он получает выражения, разбивает их на более мелкие задачи и назначает эти задачи агентам.
*   **Агент:** Выполняет фактические вычисления. Агенты запрашивают задачи у оркестратора, выполняют их и возвращают результаты.

## Структура Проекта

Проект организован в следующие каталоги:

*   `cmd/`: Содержит точки входа для оркестратора и агента.
    *   `orchestrator/`:  Основное приложение для оркестратора.
    *   `agent/`: Основное приложение для агента.
*   `internal/`: Содержит основную логику как для оркестратора, так и для агента.
    *   `orchestrator/`:  Содержит код, специфичный для оркестратора.
        *   `api/`:  Определяет обработчики HTTP API для взаимодействия с оркестратором (как публичные, так и внутренние).
        *   `calculator/`:  Реализует логику разбора выражений, генерации AST (абстрактного синтаксического дерева) и создания задач.
        *   `models/`: Определяет структуры данных, используемые оркестратором (например, `Expression`, `Task`).
        *   `repository/`:  Предоставляет хранилище в памяти для выражений и задач.  Обрабатывает сохранение и извлечение.
        *   `service/`:  Реализует бизнес-логику, координируя работу между API, репозиторием и калькулятором.
    *   `agent/`: Содержит код, специфичный для агента.
        *   `agent.go`: Реализует основной цикл агента, получение задач, обработку и отправку результатов.

## Оркестратор

### Функциональность

1.  **Отправка выражений:** Принимает арифметические выражения через POST-запрос к `/api/v1/calculate`.
2.  **Разбор выражений:** Разбирает выражение в абстрактное синтаксическое дерево (AST).
3.  **Генерация задач:** Преобразует AST в набор меньших, независимых задач. Эти задачи представляют собой отдельные арифметические операции (сложение, вычитание, умножение, деление) или получение значений. Задачи создаются с зависимостями, чтобы операции выполнялись в правильном порядке.
4.  **Управление задачами:** Хранит задачи в репозитории в памяти и отслеживает их статус (Pending, Processing, Completed, Error).
5.  **Назначение задач:** Предоставляет задачи агентам через GET-запрос к `/internal/task`. Оркестратор отдает приоритет задачам, зависимости которых были разрешены.
6.  **Агрегация результатов:** Получает результаты задач от агентов через POST-запрос к `/internal/task`. Обновляет статус задач и, когда все задачи для выражения завершены, вычисляет окончательный результат.
7.  **Статус выражения:** Позволяет получать статус выражения и результаты через GET-запросы к `/api/v1/expressions` и `/api/v1/expressions/{id}`.

### Конечные точки API

*   **Публичный API (`/api/v1`)**

    *   `POST /calculate`: Отправляет новое выражение для вычисления.
        *   Тело запроса: `{"expression": "2 + 2 * 2"}`
        *   Ответ: `{"id": "<uuid>"}` (ID выражения)
    *   `GET /expressions`: Получает список всех выражений и их статус.
        *   Ответ: `{"expressions": [{"id": "<uuid>", "status": "COMPLETED", "result": 6}, ...]}`
    *   `GET /expressions/{id}`: Получает подробную информацию о конкретном выражении.
        *   Ответ: `{"expression": {"id": "<uuid>", "status": "COMPLETED", "result": 6}}`
*   **Внутренний API (`/internal`)**

    *   `GET /task`: Получает следующую доступную задачу для агента.
        *   Ответ: `{"task": {"id": "<uuid>", "arg1": 2, "arg2": 2, "operation": "MULTIPLICATION", "operation_time": 2000}}`
    *   `POST /task`: Отправляет результат выполненной задачи.
        *   Тело запроса: `{"id": "<uuid>", "result": 4}`
        *   Ответ: `{"status": "success"}`

### Переменные окружения

*   `TIME_ADDITION_MS` (по умолчанию: 1000): Имитируемое время выполнения операций сложения (в миллисекундах).
*   `TIME_SUBTRACTION_MS` (по умолчанию: 1000): Имитируемое время выполнения операций вычитания (в миллисекундах).
*   `TIME_MULTIPLICATIONS_MS` (по умолчанию: 2000): Имитируемое время выполнения операций умножения (в миллисекундах).
*   `TIME_DIVISIONS_MS` (по умолчанию: 2000): Имитируемое время выполнения операций деления (в миллисекундах).

## Агент

### Функциональность

1.  **Запрос задач:** Периодически запрашивает задачи у оркестратора.
2.  **Выполнение задач:** Выполняет полученную задачу, выполняя указанную арифметическую операцию. Имитирует время обработки на основе поля `operation_time` и настроенных переменных окружения.
3.  **Отправка результатов:** Отправляет результат задачи обратно оркестратору.
4.  **Пул рабочих процессов:** Использует настраиваемое количество рабочих горутин для параллельной обработки задач.

### Переменные окружения

*   `ORCHESTRATOR_URL` (по умолчанию: `http://localhost:8080`): URL-адрес оркестратора.
*   `COMPUTING_POWER` (по умолчанию: 3): Количество рабочих горутин, используемых для обработки задач.

## Запуск проекта

1.  **Установите зависимости:**

    ```bash
    go mod tidy
    ```

2.  **Запустите оркестратор:**

    ```bash
    go run ./cmd/orchestrator/main.go
    ```

3.  **Запустите один или несколько агентов:**

    ```bash
    go run ./cmd/agent/main.go
    ```
    Вы можете запустить несколько агентов, чтобы имитировать распределенную среду. Вы также можете изменить переменную окружения `COMPUTING_POWER`, чтобы изменить количество рабочих процессов на агент. Например:
    ```bash
    COMPUTING_POWER=5 go run ./cmd/agent/main.go
    ```

3.  **Отправьте выражение:**

    *   **Успешный запрос:**
        ```bash
        curl -X POST -H "Content-Type: application/json" -d '{"expression": "2 + 2 * (3 - 1) / 2"}' http://localhost:8080/api/v1/calculate
        ```
         Ожидаемый ответ (HTTP статус 201 Created):
        ```json
        {"id": "<uuid>"}
        ```

    *   **Неверный формат выражения (отсутствует закрывающая скобка):**
        ```bash
        curl -X POST -H "Content-Type: application/json" -d '{"expression": "2 + 2 * (3 - 1 / 2"}' http://localhost:8080/api/v1/calculate
        ```
        Ожидаемый ответ (HTTP статус 500 Internal Server Error):
        ```json
        {"error": "Failed to process expression"}
        ```
         В логах оркестратора будет более детальная ошибка: `missing closing parenthesis`.

    *   **Пустое выражение:**
        ```bash
        curl -X POST -H "Content-Type: application/json" -d '{"expression": ""}' http://localhost:8080/api/v1/calculate
        ```
        Ожидаемый ответ (HTTP статус 422 Unprocessable Entity):
        ```json
        {"error": "Expression is required"}
        ```

    *   **Невалидный JSON:**

    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"expression": "2 + 2" ' http://localhost:8080/api/v1/calculate
    ```

    Ожидаемый ответ (HTTP статус 422 Unprocessable Entity):

    ```json
    {"error": "Invalid request payload"}
    ```

    * **Деление на ноль (ошибка во время выполнения задачи агентом):**
       Хотя деление на ноль обрабатывается в коде `agent.go`, и агент не падает, сама задача переходит в состояние ошибки. Оркестратор *не* помечает выражение как ошибочное автоматически, если только *все* задачи не завершатся с ошибкой, или если ошибка произойдет на этапе *парсинга*.
       Чтобы увидеть ошибку, вам нужно будет:
        1.  Отправить выражение с делением на ноль:
            ```bash
            curl -X POST -H "Content-Type: application/json" -d '{"expression": "2 / 0"}' http://localhost:8080/api/v1/calculate
            ```
        2.  Получить ID выражения из ответа.
        3.  Подождать некоторое время (агент выполняет задачу с задержкой).
        4.  Запросить информацию об этом выражении:
             ```bash
             curl http://localhost:8080/api/v1/expressions/<expression_id>
             ```
             Вы *не* увидите `"status": "ERROR"` на уровне выражения (потому что другие задачи, например, извлечение числа `2`, могли выполниться успешно).  Чтобы увидеть ошибку на уровне *задачи*, вам бы понадобился эндпоинт, который возвращает список *задач*, связанных с выражением (такого эндпоинта в текущем коде нет).  Но в логах *агента* вы увидите сообщение об ошибке: `division by zero`.

4.  **Проверьте статус выражения:**

    ```bash
     curl http://localhost:8080/api/v1/expressions/<expression_id>
    ```
   Замените `<expression_id>` на ID, возвращенный из запроса `calculate`.  Вы также можете вывести список всех выражений с помощью
   ```bash
   curl http://localhost:8080/api/v1/expressions
   ```
*   **Выражение не найдено:**

    ```bash
    curl http://localhost:8080/api/v1/expressions/invalid-uuid
    ```
      Ожидаемый ответ (HTTP статус 422 Unprocessable Entity):
    ```json
    {"error": "Invalid expression ID"}
    ```


    ```bash
    curl http://localhost:8080/api/v1/expressions/00000000-0000-0000-0000-000000000000
    ```
    Ожидаемый ответ (HTTP статус 404 Not Found):
    ```json
    {"error": "Expression not found"}
    ```
