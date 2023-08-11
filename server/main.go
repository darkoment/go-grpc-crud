/*
Сторона сервера:
Этот модуль содержит код для реализации на стороне сервера,
а также подключение к базе данных и операций с ней
*/
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/darkoment/go-grpc-crud-api/proto"
	//"github.com/google/uuid" - убрана за отсутствие необходимости в использовании

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Функция инициализирующая клиент базы данных
func init() {
	DatabaseConnection()
}

var DB *gorm.DB
var err error

// type Test struct {
// 	Name string `gorm:"primarykey"`
// }

// Структура для созадния таблицы Book
// С ее атрибутами
type Book struct {
	BookID  uint32    `gorm:"primarykey;autoIncrement;not null"`
	Name    string    `gorm:"not null"`
	Year    string    `gorm:"not null"`
	Edition string    `gorm:"not null"`
	Authors []*Author `gorm:"many2many:book_author;"` //References:author_id
}

// Структура для создания таблицы author
// С ее атрибутами
type Author struct {
	AuthorID  uint32  `gorm:"primarykey;autoIncrement;not null"`
	FirstName string  `gorm:"not null"`
	LastName  string  `gorm:"not null"`
	Books     []*Book `gorm:"many2many:book_author;"` //References:book_id
}

// Структура для создания смежной таблицы bookauthor
// С ее атрибутами

// Возможно эта структура не нужна

type BookAuthor struct {
	BookID   uint32 `gorm:"primarykey;not null;"` //foreignKey:fk_book_id
	AuthorID uint32 `gorm:"primarykey;not null;"` //foreignKey:fk_author_id
}

// Интерфейс для названия таблиц
// т.к. GORM добавляет окончание s к названиям таблиц
type Tabler interface {
	TableName() string
}

func (Book) TableName() string {
	return "book"
}

func (Author) TableName() string {
	return "author"
}

func (BookAuthor) TableName() string {
	return "book_author"
}

//Функция для установления соединения

func DatabaseConnection() {
	host := "localhost"
	port := "3306"
	dbName := "test"
	dbUser := "TestAdmin"
	password := "TestOnGo"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		password,
		host,
		port,
		dbName,
	)

	// Передаем dsn в gorm
	// подробнее https://gorm.io/ru_RU/docs/connecting_to_the_database.html
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(Book{}, Author{})
	//DB.Migrator().RenameTable("books", "book")
	//DB.Migrator().RenameTable("authors", "author")
	//DB.Migrator().CurrentDatabase()
	if err != nil {
		log.Fatal("Error connection to the database ", err)
	}
	fmt.Println("Database connection successful...")
}

// Создание gRPC сервера
var (
	port = flag.Int("port", 50051, "gRPC server port")
)

type server struct {
	pb.UnimplementedBookServiceServer
}

// Метод отвечающий за создание записи о книге в таблицу Book
func (*server) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	fmt.Println("Create Book")
	book := req.GetBook()
	//book.Bookid = uuid.New().String()

	data := Book{
		BookID:  book.GetBookid(),
		Name:    book.GetName(),
		Year:    book.GetYear(),
		Edition: book.GetEdition(),
	}

	res := DB.Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("Book creation unsuccessful")
	}
	return &pb.CreateBookResponse{
		Book: &pb.Book{
			Bookid:  book.GetBookid(),
			Name:    book.GetName(),
			Year:    book.GetYear(),
			Edition: book.GetEdition(),
		},
	}, nil
}

// Метод отвечающий за получение записи о книге из таблицы Book
func (*server) GetBook(ctx context.Context, req *pb.ReadBookRequest) (*pb.ReadBookResponse, error) {
	fmt.Println("Read Book", req.GetBookid())
	var book Book
	res := DB.Find(&book, "book_id = ?", req.GetBookid())
	if res.RowsAffected == 0 {
		return nil, errors.New("Book not found")
	}
	return &pb.ReadBookResponse{
		Book: &pb.Book{
			Bookid:  book.BookID,
			Name:    book.Name,
			Year:    book.Year,
			Edition: book.Edition,
		},
	}, nil
}

// Метод отвечающий за получения записей книг из таблицы Book
func (*server) GetBooks(ctx context.Context, req *pb.ReadBooksRequest) (*pb.ReadBooksResponse, error) {
	fmt.Println("Read Books")
	var books []Book
	res := DB.Find(&books)
	if res.RowsAffected == 0 {
		return nil, errors.New("Books not found")
	}
	allBooks := make([]*pb.Book, len(books))
	for i := range books {
		allBooks[i] = &pb.Book{
			Bookid:  books[i].BookID,
			Name:    books[i].Name,
			Year:    books[i].Year,
			Edition: books[i].Edition,
		}
	}
	return &pb.ReadBooksResponse{
		Books: allBooks,
	}, nil
}

// Метод отвечающий за обновление записи о книге в таблице Book
func (*server) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	fmt.Println("Update Book")
	var book Book
	reqBook := req.GetBook()

	res := DB.Model(&book).Where("book_id=?", reqBook.Bookid).Updates(
		Book{Name: reqBook.Name, Year: reqBook.Year, Edition: reqBook.Edition})

	if res.RowsAffected == 0 {
		return nil, errors.New("Books not found")
	}

	return &pb.UpdateBookResponse{
		Book: &pb.Book{
			Bookid:  book.BookID,
			Name:    book.Name,
			Year:    book.Year,
			Edition: book.Edition,
		},
	}, nil
}

// Метод отвечающий за удаление записи о книге из таблицы Book (удаление происходит по индексу)
func (*server) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	fmt.Println("Delete Book")
	var book Book
	res := DB.Where("book_id=?", req.GetBookid()).Delete(&book)
	if res.RowsAffected == 0 {
		return nil, errors.New("Book not found")
	}

	return &pb.DeleteBookResponse{
		Success: true,
	}, nil
}

// Метод отвечающий за создание записи о книге в таблицу Author
func (*server) CreateAuthor(ctx context.Context, req *pb.CreateAuthorRequest) (*pb.CreateAuthorResponse, error) {
	fmt.Println("Create Author")
	author := req.GetAuthor()
	//book.Authorid = uuid.New().String()

	data := Author{
		AuthorID:  author.GetAuthorid(),
		FirstName: author.GetFirstName(),
		LastName:  author.GetLastName(),
	}

	res := DB.Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("Author creation unsuccessful")
	}
	return &pb.CreateAuthorResponse{
		Author: &pb.Author{
			Authorid:  author.GetAuthorid(),
			FirstName: author.GetFirstName(),
			LastName:  author.GetLastName(),
		},
	}, nil
}

// Метод отвечающий за получение записи о книге из таблицы Author
func (*server) GetAuthor(ctx context.Context, req *pb.ReadAuthorRequest) (*pb.ReadAuthorResponse, error) {
	fmt.Println("Read Author", req.GetAuthorid())
	var author Author
	res := DB.Find(&author, "author_id = ?", req.GetAuthorid())
	if res.RowsAffected == 0 {
		return nil, errors.New("Author not found")
	}
	return &pb.ReadAuthorResponse{
		Author: &pb.Author{
			Authorid:  author.AuthorID,
			FirstName: author.FirstName,
			LastName:  author.LastName,
		},
	}, nil
}

// Метод отвечающий за получения записей книг из таблицы Author
func (*server) GetAuthors(ctx context.Context, req *pb.ReadAuthorsRequest) (*pb.ReadAuthorsResponse, error) {
	fmt.Println("Read Authors")
	var authors []Author
	res := DB.Find(&authors)
	if res.RowsAffected == 0 {
		return nil, errors.New("Authors not found")
	}
	allAuthors := make([]*pb.Author, len(authors))
	for i := range authors {
		allAuthors[i] = &pb.Author{
			Authorid:  authors[i].AuthorID,
			FirstName: authors[i].FirstName,
			LastName:  authors[i].LastName,
		}
	}
	return &pb.ReadAuthorsResponse{
		Authors: allAuthors,
	}, nil
}

// Метод отвечающий за обновление записи о книге в таблице Author
func (*server) UpdateAuthor(ctx context.Context, req *pb.UpdateAuthorRequest) (*pb.UpdateAuthorResponse, error) {
	fmt.Println("Update Author")
	var author Author
	reqAuthor := req.GetAuthor()

	res := DB.Model(&author).Where("author_id=?", reqAuthor.Authorid).Updates(
		Author{FirstName: reqAuthor.FirstName, LastName: reqAuthor.LastName})

	if res.RowsAffected == 0 {
		return nil, errors.New("Authors not found")
	}

	return &pb.UpdateAuthorResponse{
		Author: &pb.Author{
			Authorid:  author.AuthorID,
			FirstName: author.FirstName,
			LastName:  author.LastName,
		},
	}, nil
}

// Метод отвечающий за удаление записи о книге из таблицы Author (удаление происходит по индексу)
func (*server) DeleteAuthor(ctx context.Context, req *pb.DeleteAuthorRequest) (*pb.DeleteAuthorResponse, error) {
	fmt.Println("Delete Author")
	var author Author
	res := DB.Where("author_id=?", req.GetAuthorid()).Delete(&author)
	if res.RowsAffected == 0 {
		return nil, errors.New("Author not found")
	}

	return &pb.DeleteAuthorResponse{
		Success: true,
	}, nil
}

func (*server) CreateBookAuthor(ctx context.Context, req *pb.CreateBookAuthorRequest) (*pb.CreateBookAuthorResponse, error) {
	fmt.Println("Create related Book and Author")
	book_author := req.GetBookauthor()

	data := BookAuthor{
		BookID:   book_author.GetBookid(),
		AuthorID: book_author.GetAuthorid(),
	}

	res := DB.Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("relation creation unsuccessful")
	}
	return &pb.CreateBookAuthorResponse{
		Bookauthor: &pb.BookAuthor{
			Bookid:   book_author.GetBookid(),
			Authorid: book_author.GetAuthorid(),
		},
	}, nil
}

// Метод отвечающий за получение записи о книге из таблицы BookAuthor
// func (*server) GetBookAuthor(ctx context.Context, req *pb.ReadBookAuthorRequest) (*pb.ReadBookAuthorResponse, error) {
// 	fmt.Println("Read relates Book Author", req.GetAuthorid())

// }

// Метод отвечающий за получения записей книг из таблицы BookAuthor
func (*server) GetBooksAuthors(ctx context.Context, req *pb.ReadBooksAuthorsRequest) (*pb.ReadBooksAuthorsResponse, error) {
	fmt.Println("Read related Books Authors")
	var book_author []BookAuthor
	res := DB.Find(&book_author)
	if res.RowsAffected == 0 {
		return nil, errors.New("relates not found")
	}
	allRelated := make([]*pb.BookAuthor, len(book_author))
	for i := range book_author {
		allRelated[i] = &pb.BookAuthor{
			Authorid: book_author[i].AuthorID,
			Bookid:   book_author[i].BookID,
		}
	}
	return &pb.ReadBooksAuthorsResponse{
		Booksauthors: allRelated,
	}, nil
}

// Метод отвечающий за обновление записи о книге в таблице BookAuthor
//func (*server) UpdateBookAuthor(ctx context.Context, req *pb.UpdateBookAuthorRequest) (*pb.UpdateBookAuthorResponse, error) {
// 	fmt.Println("Update BookAuthor")

// }

// // Метод отвечающий за удаление записи о книге из таблицы BookAuthor (удаление происходит по индексу)
// func (*server) DeleteBookAuthor(ctx context.Context, req *pb.DeleteBookAuthorRequest) (*pb.DeleteBookAuthorResponse, error) {
// 	fmt.Println("Delete BookAuthor")

// }

func main() {
	fmt.Println("gRPC server running ...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterBookServiceServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
