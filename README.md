# Task Management API

REST API для управления задачами, построенный на Go с использованием Gin, pgx и PostgreSQL (Supabase).

## Особенности

- ✅ CRUD операции для задач и подзадач
- ✅ Управление пользователями и их назначение на задачи
- ✅ Сортировка задач по дате выполнения и статусу
- ✅ Фильтрация по статусам: `all`, `not-started`, `completed`, `in-progress`
- ✅ Интеграция с Supabase PostgreSQL
- ✅ Современный pgx/v5 драйвер для высокой производительности
- ✅ CORS поддержка для фронтенда
- ✅ Валидация данных

## Установка и настройка

### Требования

- Go 1.24+
- PostgreSQL (Supabase)

### Шаги установки

1. **Клонируйте репозиторий**

```bash
git clone <your-repo>
cd event_app
```

2. **Установите зависимости**

```bash
go mod download
```

3. **Настройте переменные окружения**

Создайте файл `.env` в корне проекта:

**Вариант 1: Используя DATABASE_URL (рекомендуется для Supabase)**

```env
# Supabase Database URL (найдите в Project Settings > Database > Connection string > URI)
DATABASE_URL=postgresql://postgres:5I9QJ8ruCgqKiPvW@db.your-project-ref.supabase.co:5432/postgres

# Server Configuration
PORT=8080

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:3000
```

**Вариант 2: Используя отдельные параметры**

```env
# Supabase Database Configuration
DB_HOST=db.your-project-ref.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=5I9QJ8ruCgqKiPvW
DB_NAME=postgres
DB_SSLMODE=require

# Server Configuration
PORT=8080

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:3000
```

**Важно:**

- Замените `your-project-ref` на реальный reference вашего Supabase проекта
- Используйте ваш реальный пароль `5I9QJ8ruCgqKiPvW` или создайте новый в Supabase

### Где найти DATABASE_URL в Supabase:

1. Откройте ваш проект в Supabase Dashboard
2. Перейдите в **Project Settings** → **Database**
3. В разделе **Connection string** выберите **URI**
4. Скопируйте строку подключения

5. **Запустите приложение**

```bash
go run cmd/app/main.go
```

Сервер запустится на `http://localhost:8080`

## Технологический стек

- **Go 1.24+** - Язык программирования
- **Gin** - HTTP веб-фреймворк
- **pgx/v5** - Современный PostgreSQL драйвер с высокой производительностью
- **Supabase** - PostgreSQL как сервис
- **Clean Architecture** - Архитектурный паттерн

### Преимущества pgx/v5:

- 🚀 **Высокая производительность** - быстрее lib/pq
- 🎯 **Нативная поддержка PostgreSQL** - все специфичные типы
- 🔧 **Контекстная поддержка** - лучшая отмена операций
- 📦 **Современное API** - более удобное использование

## API Endpoints

### Health Check

- `GET /health` - Проверка состояния сервера

### Задачи (Tasks)

- `POST /api/v1/tasks` - Создать задачу
- `GET /api/v1/tasks` - Получить список задач (с фильтрацией и сортировкой)
- `GET /api/v1/tasks/:id` - Получить задачу по ID
- `PUT /api/v1/tasks/:id` - Обновить задачу
- `DELETE /api/v1/tasks/:id` - Удалить задачу

### Назначение пользователей

- `POST /api/v1/tasks/:id/users/:user_id` - Назначить пользователя на задачу
- `DELETE /api/v1/tasks/:id/users/:user_id` - Убрать пользователя с задачи

### Подзадачи (SubTasks)

- `POST /api/v1/tasks/:id/subtasks` - Создать подзадачу
- `GET /api/v1/tasks/:id/subtasks` - Получить подзадачи задачи
- `PUT /api/v1/subtasks/:id` - Обновить подзадачу
- `DELETE /api/v1/subtasks/:id` - Удалить подзадачу

### Пользователи (Users)

- `POST /api/v1/users` - Создать пользователя
- `GET /api/v1/users` - Получить всех пользователей
- `GET /api/v1/users/:id` - Получить пользователя по ID
- `PUT /api/v1/users/:id` - Обновить пользователя
- `DELETE /api/v1/users/:id` - Удалить пользователя

## Примеры запросов

### Создание задачи

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Выполнить проект",
    "icon_name": "code",
    "due_date": "2024-12-31T23:59:59Z",
    "status": "not-started",
    "user_ids": ["user-1", "user-2"]
  }'
```

### Получение задач с фильтрацией

```bash
# Все задачи
curl "http://localhost:8080/api/v1/tasks"

# Только завершенные задачи, отсортированные по дате выполнения
curl "http://localhost:8080/api/v1/tasks?status=completed&sort_by=due_date&sort_type=asc"

# Задачи в работе с пагинацией
curl "http://localhost:8080/api/v1/tasks?status=in-progress&limit=10&offset=0"
```

### Создание пользователя

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван Иванов",
    "src": "https://example.com/avatar.jpg"
  }'
```

### Создание подзадачи

```bash
curl -X POST http://localhost:8080/api/v1/tasks/task-id/subtasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Подзадача 1",
    "description": "Описание подзадачи",
    "status": "not-started"
  }'
```

## Структура данных

### Task

```json
{
  "id": "uuid",
  "title": "string",
  "icon_name": "string",
  "start_time": "string (optional)",
  "end_time": "string (optional)",
  "due_date": "datetime",
  "progress": "number",
  "status": "not-started|completed|in-progress",
  "comments": "number",
  "attachments": "number",
  "links": "number",
  "users": [{"id": "string", "name": "string", "src": "string"}],
  "sub_tasks": [SubTask],
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### SubTask

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "title": "string",
  "description": "string (optional)",
  "status": "not-started|completed|in-progress",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### User

```json
{
  "id": "uuid",
  "name": "string",
  "src": "string"
}
```

## Параметры фильтрации и сортировки

### Параметры запроса для GET /api/v1/tasks:

- `status` - Фильтр по статусу: `all`, `not-started`, `completed`, `in-progress`
- `sort_by` - Поле для сортировки: `due_date`, `status`, `created_at` (default)
- `sort_type` - Тип сортировки: `asc`, `desc` (default)
- `limit` - Лимит записей (default: 50)
- `offset` - Смещение для пагинации (default: 0)

## База данных

Приложение автоматически создает необходимые таблицы при запуске:

- `tasks` - Основные задачи
- `sub_tasks` - Подзадачи
- `users` - Пользователи
- `task_user_assignments` - Связь задач и пользователей

## Разработка

### Структура проекта

```
cmd/app/           # Точка входа приложения
internal/
├── config/        # Конфигурация БД
├── env/           # Утилиты для env переменных
├── handlers/      # HTTP обработчики
├── models/        # Модели данных
├── repository/    # Слой работы с БД
└── service/       # Бизнес логика
```

### Архитектура

Проект использует Clean Architecture:

1. **Handlers** - HTTP слой (контроллеры)
2. **Services** - Бизнес логика
3. **Repository** - Слой данных
4. **Models** - Структуры данных

## Лицензия

MIT
