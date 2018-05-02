package main

import "encoding/xml"

type UserShelf struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
}

type GoodreadsUserShelves struct {
	GoodreadsResponse xml.Name `xml:"GoodreadsResponse"`
	Shelves           struct {
		UserShelves []UserShelf `xml:"user_shelf"`
	} `xml:"shelves"`
}

type Book struct {
	ID                 int     `xml:"best_book>id"`
	BooksCount         int     `xml:"books_count"`
	RatingsCount       int     `xml:"ratings_count"`
	TextReviewsCount   int     `xml:"text_review_counts"`
	OriginalPublicYear int     `xml:"original_publication_year"`
	AverageRating      float64 `xml:"average_rating"`
	Title              string  `xml:"best_book>title"`
	Author             string  `xml:"best_book>author>name"`
}

type BookSearchResponse struct {
	GoodreadsResponse xml.Name `xml:"GoodreadsResponse"`
	Books             []Book   `xml:"search>results>work"`
}
