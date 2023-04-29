api:

1) POST /api/register - регистрация
2) POST /api/login - логин
3) GET /api/posts/ - список всех постов
4) POST /api/posts/ - добавление поста - обратите внимание - есть с урлом, а есть с текстом
5) GET /api/posts/{CATEGORY_NAME} - список постов конкретной категории
6) GET /api/post/{POST_ID} - детали поста с комментами
7) POST /api/post/{POST_ID} - добавление коммента
8) DELETE /api/post/{POST_ID}/{COMMENT_ID} - удаление коммента
9) GET /api/post/{POST_ID}/upvote - рейтинг постп вверх
10) GET /api/post/{POST_ID}/downvote - рейтинг поста вниз
11) GET /api/post/{POST_ID}/unvote - отмена голоса 
12) DELETE /api/post/{POST_ID} - удаление поста
13) GET /api/user/{USER_LOGIN} - получение всех постов конкртеного пользователя

запуск на localhost:8091
```
go run redditclone/cmd/redditclone/main.go 
```

