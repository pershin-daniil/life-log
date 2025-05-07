# life-log
Life Log is an app for tracking and recording activities, habits, and tasks. It helps capture meaningful life moments, organize daily routines, and build productive habits.

## Конфигурация

Приложение использует два механизма конфигурации:

1. **YAML-конфигурация (config.yaml)** - для настроек приложения и параметров, специфичных для окружения (порты, таймауты и т.д.)
2. **Переменные окружения (ENV)** - для секретов и чувствительной информации (пароли, токены, ключи)

### YAML-конфигурация (настройки приложения)

Создайте файл `config.yaml` в корне проекта или используйте переменную окружения `CONFIG_PATH` для указания пути к файлу.

Пример содержимого `config.yaml`:
```yaml
# Конфигурация сервера
server:
  # HTTP порт
  port: 8080
  
  # Таймаут в секундах для graceful shutdown
  shutdown_timeout_seconds: 3
  
  # Таймаут в секундах для чтения заголовков запроса
  read_header_timeout_seconds: 1
```

### Переменные окружения (секреты)

| Переменная      | Описание                               | Пример значения              |
|-----------------|----------------------------------------|------------------------------|
| VERSION_TAG     | Версия приложения                      | "v1.0.0"                     |
| DATABASE_URL    | URL подключения к базе данных          | "postgres://user:pass@db:5432/life_log" |
| POSTGRES_USER   | Имя пользователя PostgreSQL            | "postgres"                   |
| POSTGRES_PASSWORD | Пароль пользователя PostgreSQL       | "secret"                     |
| POSTGRES_DB     | Имя базы данных PostgreSQL             | "life_log"                   |
| CONFIG_PATH     | Путь к файлу конфигурации (опционально)| "configs/dev.yaml"           |
| GHCR_USERNAME   | Имя пользователя GitHub Container Registry | "username"               |

## Запуск

1. Клонируйте репозиторий
2. Установите зависимости: `go mod download`
3. Создайте файл `config.yaml` на основе `config.yaml.example`
4. Экспортируйте необходимые переменные окружения или создайте файл `.env`
5. Запустите приложение: `go run cmd/life-log/main.go`
