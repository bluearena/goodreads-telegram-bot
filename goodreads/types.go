package goodreads

import "encoding/xml"

type UserShelves struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
}

type GoodreadsUserShelves struct {
	GoodreadsResponse xml.Name      `xml:"GoodreadsResponse"`
	Shelves           []UserShelves `xml:"shelves"`
}
