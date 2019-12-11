Написал redditclone

Структура проекта
.
├── bin
│   └── redditclone
├── cmd
│   └── redditclone
│       └── main.go
├── pkg
│   ├── handlers
│   │   ├── posts.go
│   │   └── user.go
│   ├── items
│   │   ├── comment.go
│   │   ├── posts.go
│   │   └── repo.go
│   ├── middleware
│   │   ├── accesslog.go
│   │   ├── auth.go
│   │   └── panic.go
│   ├── session
│   │   ├── manager.go
│   │   └── session.go
│   └── user
│       └── user.go
├── readme.md
└── template
    ├── index.html
    └── static

ID поста и пользователя для простоты сделал автоинкремент (позже доработаю при подключении бд)

