package controllers

import (
	"book/database"
	"book/graph/model"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateBookListing(bookInfo model.NewBookListing) *model.BookListing {
	booksCollection := database.GetDB().Collection("books")
	currentTime := time.Now()
	currentTimeFloat := float64(currentTime.Unix())

	result, err := booksCollection.InsertOne(context.Background(), bson.M{"title": bookInfo.Title, "bookname": bookInfo.Bookname, "description": bookInfo.Description, "author": bookInfo.Author, "addedOn": currentTimeFloat})
	if err != nil {
		log.Fatalf("Failed to insert book: %v", err)
	}

	// Create a BookListing with the inserted data
	insertedID := result.InsertedID.(primitive.ObjectID).Hex()
	fmt.Println(insertedID)
	return &model.BookListing{
		ID:          insertedID,
		Title:       bookInfo.Title,
		Author:      bookInfo.Author,
		Bookname:    bookInfo.Bookname,
		Description: bookInfo.Description,
		AddedOn:     currentTimeFloat,
		// Add other fields accordingly
	}
}

func GetBookByID(id string) (*model.BookListing, error) {
	booksCollection := database.GetDB().Collection("books")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var book database.Book
	err = booksCollection.FindOne(context.Background(), filter).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return nil, fmt.Errorf("book not found")
		}

		log.Printf("Error fetching book by ID: %v", err)
		return nil, err
	}

	bookListing := &model.BookListing{
		ID:          objectID.Hex(),
		Title:       book.Title,
		Author:      book.Author,
		Description: book.Description,
		AddedOn:     book.AddedOn,
		Bookname:    book.Bookname,
	}

	return bookListing, nil
}
func GetBooks() ([]*model.BookListing, error) {
	booksCollection := database.GetDB().Collection("books")

	// Empty filter to fetch all documents
	filter := bson.M{}

	// Execute the query to get all books
	cursor, err := booksCollection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Error fetching books: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var books []*model.BookListing
	for cursor.Next(context.Background()) {
		var book model.BookListingDB
		err := cursor.Decode(&book)
		if err != nil {
			log.Printf("Error decoding book: %v", err)
			continue
		}

		// Convert the book to a BookListing and append to the list
		bookListing := &model.BookListing{
			ID:          book.ID.Hex(),
			Title:       book.Title,
			Author:      book.Author,
			Description: book.Description,
			AddedOn:     book.AddedOn,
			// Add other fields accordingly
		}
		books = append(books, bookListing)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Error processing books cursor: %v", err)
		return nil, err
	}

	return books, nil
}

func UpdateBooks(bookID string, bookInfo model.UpdateBookListing) (*model.BookListing, error) {
	booksCollection := database.GetDB().Collection("books")

	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		log.Printf("Invalid bookID format: %v", err)
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	update := bson.M{}
	if bookInfo.Title != nil {
		update["title"] = *bookInfo.Title
	}
	if bookInfo.Bookname != nil {
		update["bookname"] = *bookInfo.Bookname
	}
	if bookInfo.Author != nil {
		update["author"] = *bookInfo.Author
	}
	if bookInfo.Description != nil {
		update["description"] = *bookInfo.Description
	}

	result, err := booksCollection.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		log.Printf("Error updating book: %v", err)
		return nil, err
	}

	fmt.Println("ModifiedCount:", result.ModifiedCount)

	return GetBookByID(bookID)
}

func FindBooksByAuthor(author string) ([]*model.BookListing, error) {
	booksCollection := database.GetDB().Collection("books")

	filter := bson.M{"author": author}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := booksCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []*model.BookListing
	for cursor.Next(context.Background()) {
		var book model.BookListingDB
		err := cursor.Decode(&book)
		if err != nil {
			log.Printf("Error decoding book: %v", err)
			continue
		}

		// Convert the book to a BookListing and append to the list
		bookListing := &model.BookListing{
			ID:          book.ID.Hex(),
			Title:       book.Title,
			Author:      book.Author,
			Description: book.Description,
			AddedOn:     book.AddedOn,
			// Add other fields accordingly
		}
		books = append(books, bookListing)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Error processing books cursor: %v", err)
		return nil, err
	}

	return books, nil
}

func DeleteBookById(BookId string) *model.DeleteBook {
	booksCollection := database.GetDB().Collection("books")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(BookId)
	filter := bson.M{"_id": _id}
	_, err := booksCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Println(err)
	}
	return &model.DeleteBook{DeletedBookID: BookId}
}
