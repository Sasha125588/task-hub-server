# 🚀 Быстрый старт

## 1. Настройте Supabase

1. Откройте ваш проект в [Supabase Dashboard](https://supabase.com/dashboard)
2. Перейдите в **Project Settings** → **Database**
3. В разделе **Connection string** выберите **URI**
4. Скопируйте строку подключения (будет выглядеть примерно так):
   ```
   postgresql://postgres.your-project-ref:your-password@aws-0-region.pooler.supabase.com:6543/postgres
   ```

## 2. Создайте .env файл

Создайте файл `.env` в корне проекта:

```env
DATABASE_URL=postgresql://postgres:5I9QJ8ruCgqKiPvW@db.kpVIQlol7giFTJkn.supabase.co:5432/postgres
PORT=8080
ALLOWED_ORIGINS=http://localhost:3000
```

**⚠️ Важно:** Замените строку подключения на вашу реальную из Supabase!

## 3. Протестируйте подключение

```bash
go run test_connection.go
```

Вы должны увидеть:

```
✅ Successfully connected to database!
PostgreSQL version: PostgreSQL 15.1 on x86_64-pc-linux-gnu...
✅ Inserted test record with ID: 1
🎉 Database connection test completed successfully!
```

## 4. Запустите API сервер

```bash
go run cmd/app/main.go
```

Вы должны увидеть:

```
Successfully connected to database with pgx driver
Successfully created tables
Server starting on port 8080
```

## 5. Протестируйте API

### Health check:

```bash
curl http://localhost:8080/health
```

### Создайте пользователя:

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Тестовый пользователь", "src": "https://example.com/avatar.jpg"}'
```

### Создайте задачу:

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Моя первая задача",
    "icon_name": "star",
    "due_date": "2024-12-31T23:59:59Z",
    "status": "not-started"
  }'
```

### Получите все задачи:

```bash
curl http://localhost:8080/api/v1/tasks
```

## 🎉 Готово!

Ваш Task Management API работает и готов к подключению фронтенда на `http://localhost:8080`

## Проблемы?

1. **Ошибка подключения к БД**: Проверьте DATABASE_URL в .env файле
2. **CORS ошибки**: Убедитесь что ALLOWED_ORIGINS включает адрес вашего фронтенда
3. **Порт занят**: Измените PORT в .env файле

Полная документация в [README.md](README.md)
