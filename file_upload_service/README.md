# ЛОКАЛЬНЫЙ ЗАПУСКЕ(БЕЗ ДОКЕРА) File_Uploader Service

## 📌 Описание

Этот сервис отвечает за загрузку файла в объектное хранилище MinIO и предоставление url файла для просмотра. Сервис использует JWT для авторизации

---
## Запуск проекта

### 1️⃣ Настройка переменных окружения
Создайте файл `.env` в корневой директории и добавьте следующие параметры, пример:

```ini
HTTP_SERVER_ADDRESS=localhost:8081
JWT_PUBLIC_KEY_PATH=public_key.pem
MINIO_PORT=localhost:9000
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
MINIO_USE_SSL=false
MINIO_URL_LIFETIME=48h
REDIS_URL_LIFETIME=1h
REDIS_PORT=localhost:6379
```

---

### 2️⃣ Создайте публичный ключ, для валидации JWT токена, если лень, можно просто скопировать ключ, который лежит в auth-service

---

### 3️⃣ Скачайте MinIO из официального сайта

### 4️⃣ Скачайте Redis, для винды можно скачать отсюда [RedisForWindows](https://github.com/tporadowski/redis/releases), далее просто запустить файл redis-server.exe

---
> **Важно:** Сначало нужно запустить MinIO потом только сервер и не забыть запустить redis

### 5️⃣ Запустите 2 команды 

``` make minio-run ```
``` make server-run ```

---

> После запуска сервер будет доступен по адресу `http://localhost:8081`

---

## 🛠 Технологии
- **Go** – версия 1.23.1
- **MinIO** – объектное хранилище
- **JWT** – аутентификация
- **Makefile** – автоматизация задач
- **OpenSSL** – генерация ключей
- **Redis** – Для временного хранения url

---
