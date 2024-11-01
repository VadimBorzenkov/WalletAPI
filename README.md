# Тестовое задание 
## Стек технологий
- Go
- PostgreSQL
- Docker
- Docker Compose

## Установка и запуск
1. Клонируйте репозиторий:
   git clone "https://github.com/VadimBorzenkov/WalletAPI"

2. Скопируйте конфигурации:
    cp .env.example .env

3. Замените необходимые переменные окружения в .env-файле

### Запуск через Docker
1. Запустите Docker Compose:
   docker-compose up --build

2. Проверьте, что контейнеры запущены:
   Убедитесь, что контейнеры app и db запущены и работают корректно.
