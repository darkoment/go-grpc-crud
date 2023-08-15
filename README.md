# **Описание проекта**
Это тестовое задание для демонстрации работы gRPC на языке GO
Сервис для поиска авторов по названию книге и книг по автору

# **Подготовка к запуску**
## **Установка нужных пакетов**
Перед запуском проекта требуется подключить соответствующие пакеты:

Для получения gRPC
`go get -u google.golang.org/grpc`

Для получения gorm
`go get -u gorm.io/gorm `

Для получения драйвера (mysql) для gorm
`go get -u gorm.io/driver/mysql`

Для получения gin
`go get -u github.com/gin-gonic/gin`

### **Установка protoc**
Если по какой-то причине Compiler protocol buffer не установлен:
`PB_REL="https://github.com/protocolbuffers/protobuf/releases"`

`curl -LO $PB_REL/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip`

разархивировать:
`unzip protoc-3.15.8-linux-x86_64.zip -d $HOME/.local`

обновить окружение:
`export PATH="$PATH:$HOME/.local/bin"`

## **Генерация из file.proto **
Находясь в папке проекта использовать команду:
`protoc --go_out=. --go-grpc_out=. proto/BooksDB.proto`

## **Загрузка тестовых данных в БД из Docker**
Находясь в папке проекта выполнить:
`docker build -t test_mysql .`
затем:
`docker run --detach --name=test_mysql --publish 3303:3306 test_mysql:latest`

## Запуск
Если все пункты выше выполнены, использовать:
`go run server/main.go`

# Тестирование
Сервер работает по localhost:50051

Для отправки запроса по gRPC использовать *insomnia* или *postman*

для книги (таблица book) существуют поля:
`bookid` принимает uint32,
`name` принимает string
`year` принимает string
`edition` принимает string

для Автора (таблицы author) существуют поля:
`authorid` принимает uint32
`first_name` принимает string
`last_name` принимает string

Также есть поля для результирующей (book_author):
`bookid` принимает uint32,
`authorid` принимает uint32

## для взаимодействия с таблицей по отдельности
Для каждой таблицы существуют запросы на Создание, Удаление, Обновление, Получение конкретной позиции, Получение всех позиций. На примере указаны запросы для таблицы book:

`CreateBook` - создать книгу, принимает
```
{
    "book": {
        "name": "Test",
        "edition": "gold_test_edition",
        "year": "2023"
    }
}
```

`UpdateBook` - обновляет книгу, принимает:
```
{
    "book": {
        "bookid": 10,
        "name": "fff",
        "year": "2017",
        "edition": "1"
    }
}
```
`GetBooks` - получает все книги (запрос пустой)
```
{

}
```

`GetBook` - получает книгу по bookid:
```
{
    "bookid": "1"
}
```

`DeleteBook` - удаляет книгу по bookid (Если не существует связи с автором):
```
{
    "bookid": "1"
}
```
## Для получения либо книг по авторову, либо авторов по книге
`GetAuthorBooks` - для поиска всех книг по автору. Необходимо задать его first_name и last_name:
```
{
    "author": {
        "first_name": "Анатолий",
        "last_name": "Громов"
    }
}
```

`GetAuthorsBook` - для поиска всех авторов по книге. Необходимо указать название книги
```
{
    "book": {
        "name": "Сказки громовых"
    }
}
```

