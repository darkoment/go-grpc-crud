package main

import (
	"flag"
	"log"
	"net/http"

	pb "github.com/darkoment/go-grpc-crud-api/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

type Book struct {
	BookID  uint32 `json:"book_id"`
	Name    string `json:"name"`
	Year    string `json:"year"`
	Edition string `json:"edition"`
}

type Author struct {
	AuthorID  uint32 `json:"author_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewBookServiceClient(conn)

	r := gin.Default()

	r.GET("/book", func(ctx *gin.Context) {
		res, err := client.GetBooks(ctx, &pb.ReadBooksRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"Book": res.Books,
		})
	})

	r.GET("/book/:book_id", func(ctx *gin.Context) {
		book_id := ctx.Param("book_id")
		res, err := client.GetBook(ctx, &pb.ReadBookRequest{Bookid: book_id})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message ": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"book": res.Book,
		})
	})

	r.POST("/book", func(ctx *gin.Context) {
		var book Book
		err := ctx.ShouldBind(&book)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		data := &pb.Book{
			Bookid:  book.BookID,
			Name:    book.Name,
			Year:    book.Year,
			Edition: book.Edition,
		}
		res, err := client.CreateBook(ctx, &pb.CreateBookRequest{
			Book: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"book": res.Book,
		})
	})

	r.PUT("/book/:book_id", func(ctx *gin.Context) {
		var book Book
		err := ctx.ShouldBind(&book)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.UpdateBook(ctx, &pb.UpdateBookRequest{
			Book: &pb.Book{
				Bookid:  book.BookID,
				Name:    book.Name,
				Year:    book.Year,
				Edition: book.Edition,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"book": res.Book,
		})
		return
	})

	r.DELETE("book/:book_id", func(ctx *gin.Context) {
		book_id := ctx.Param("book_id")
		res, err := client.DeleteBook(ctx, &pb.DeleteBookRequest{Bookid: book_id})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if res.Success == true {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Book deleted successfully",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error deleting book",
			})
			return
		}
	})

	r.GET("/author", func(ctx *gin.Context) {
		res, err := client.GetAuthors(ctx, &pb.ReadAuthorsRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"Authors": res.Authors,
		})
	})

	r.GET("/author/:author_id", func(ctx *gin.Context) {
		author_id := ctx.Param("author_id")
		res, err := client.GetAuthor(ctx, &pb.ReadAuthorRequest{Authorid: author_id})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message ": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"author": res.Author,
		})
	})

	r.POST("/author", func(ctx *gin.Context) {
		var author Author
		err := ctx.ShouldBind(&author)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		data := &pb.Author{
			Authorid:  author.AuthorID,
			FirstName: author.FirstName,
			LastName:  author.LastName,
		}
		res, err := client.CreateAuthor(ctx, &pb.CreateAuthorRequest{
			Author: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"author": res.Author,
		})
	})

	r.PUT("/author/:author_id", func(ctx *gin.Context) {
		var author Author
		err := ctx.ShouldBind(&author)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.UpdateAuthor(ctx, &pb.UpdateAuthorRequest{
			Author: &pb.Author{
				Authorid:  author.AuthorID,
				FirstName: author.FirstName,
				LastName:  author.LastName,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"author": res.Author,
		})
		return
	})

	r.DELETE("author/:author_id", func(ctx *gin.Context) {
		author_id := ctx.Param("author_id")
		res, err := client.DeleteAuthor(ctx, &pb.DeleteAuthorRequest{Authorid: author_id})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if res.Success == true {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Author deleted successfully",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error deleting author",
			})
			return
		}
	})

	r.Run(":5000")

}
