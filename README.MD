# GraphQL Posts API

## Описание
Система для создания и чтения постов и комментариев с использованием GraphQL.
Поддерживаются:
- Создание/получение пользователей, постов и комментариев.
- Возможность отключения комментариев автором поста.
- GraphQL Subscriptions для получения новых комментариев в реальном времени.
- Пагинация комментариев (в дальнейшем).

## Технические требования
- Go, Gin, GORM, gqlgen.
- База данных: PostgreSQL (с возможностью использования in-memory хранилища).
- Docker, docker-compose.

## Запуск проекта
### Локально
1. Установите зависимости: `go mod download`
2. Создайте файл `.env` с необходимыми переменными.
```
STORAGE_TYPE= postgres или memory // по умолчанию postgres
```
* Если выбрано postgres, то необходимо указать следующие переменные:
```
DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME
PGADMIN_DEFAULT_EMAIL
PGADMIN_DEFAULT_PASSWORD
```
3. Далее необходимо создать базу данных с указанными переменными в файле `.env`.
4. Запустите сервер: `go run cmd/main.go`

### Через Docker
1. Соберите и запустите контейнеры: `docker-compose up --build -d`
2. API будет доступен на `http://localhost:8080/query`
3. Так же можно подключиться к бд через `http://localhost:5050`

## Тестирование
Запуск тестов: `go test ./graph -v`

## Примеры запросов
Примеры GraphQL запросов, мутаций и подписок приведены ниже.
Для проверки запросов использовалась программа `Insomnia`
### Запросы

##### Созданеи пользователя
```graphql
mutation {
  createUser(username: "test_user") {
    id
    username
    createdAt
  }
}
```

##### Создание поста
```graphql
mutation {
  createPost(title: "Мой первый пост", content: "Это тестовый пост", authorID: 1) {
    id
    title
    content
    createdAt
  }
}
```

##### Получение поста
```graphql
query {
  getPost (id: "1") {
    title
    content
    createdAt
  }
```

##### Получение поста и его комментариев
```graphql
query {
  getPost (id: "1") {
    title
    content
    createdAt
  }
  
  getComments(postID: "1", limit: 5, offset: 0) {
    id
    content
    createdAt
    children {
      id
      content
      children {
        id
        content
      }
    }
  }
}
```

##### Получение постов
```graphql
query {
  getPosts {
    id
    title
    content
    createdAt
  }
}
```

##### Создание комментария (обычные и вложенные)
```graphql
mutation {
  createComment(
    postID: "1"
    authorID: "1"
    content: "Parent comment"
  ) {
    id
    content
  }
}
----------------------------
mutation {
  createComment(
    postID: "1"
    parentID: "1"
    authorID: "1"
    content: "Reply to comment"
  ) {
    id
    content
    parent {
      id
      content
    }
  }
}
```

##### Получение комментариев с пагинацией
```graphql
query {
  getComments(
    postID: "1"
    limit: 5
    offset: 0
  ) {
    id
    content
    createdAt
    children {
      id
      content
    }
  }
}
```

##### Отключение комментария
```graphql
mutation {
  toggleComments(postID: 3, disable: true) {
    id
    disableComments
  }
}
```

##### Подписка на новые комментарии
```graphql
subscription {
  onNewComment(postID: "2") {
    id
    content
    createdAt
  }
}
```
