# coinService

Это сервис для управления монетами и мерчем в Avito.

## Запуск с использованием Docker

1. Должны быть установлены [Docker](https://docs.docker.com/get-docker/) и [Docker Compose](https://docs.docker.com/compose/install/).

2. Клонируем репозиторий:
   ```bash
   git clone https://github.com/your-username/coinService.git
   cd coinService
3. Запускаем проект:
   ```bash
   docker-compose up -d
4. Сервис будет доступен по адресу:

    API: http://localhost:8080
    База данных PostgreSQL: localhost:5432
5. Остановить:
   ```bash
   docker-compose down
