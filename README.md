# Инструкция по запуску
***

## Начало работы
1. Перейдите в папку, в которой вы хотите работать с репозиторием и откройте ее в терминале.
2. Клонируйте репозиторий на свой компьютер, используя команду 
`git clone git@github.com:IlyaChgn/merch_shop.git`
***

## Развертывание приложения
1. Создайте в корне проекта файл .env, в котором будут содержаться переменные окружения.
2. Добавьте в .env следующие переменные:
    - *CONFIG_PATH*=./internal/pkg/config/config.yaml
    - *SECRET_KEY*=api-secret-key
    - *POSTGRES_HOST*=postgres
    - *DATABASE_NAME*=shop
    - *POSTGRES_USERNAME*=postgres
    - *POSTGRES_PASSWORD*=postgres
    - *POSTGRES_PORT*=5432
   
   **Примечание**: все значения переменных за исключением *CONFIG_PATH*, *POSTGRES_HOST* и *DATABASE_NAME* могут быть 
изменены.
3. Перейдите в терминал и запустите проект командой `docker compose up -d`.
***

## Запуск тестов
1. Для корректной работы тестов добавьте в .env следующие переменные:
    - *TEST_SECRET_KEY*=test-secret-key
    - *TEST_DB_HOST*=localhost
    - *TEST_DB_NAME*=shop
    - *TEST_POSTGRES_USER*=postgres
    - *TEST_POSTGRES_PASSWORD*=postgres
    - *TEST_DB_PORT*=5432

    **Примечание**: все значения переменных за исключением *TEST_DB_HOST* и *TEST_DB_NAME* могут быть
      изменены.
2. Для запуска тестов с выводом результатов используйте команду `make test`.
3. Для запуска тестов с отображением покрытия по каждому модулю и всему проекту используйте команду `make cover`.
