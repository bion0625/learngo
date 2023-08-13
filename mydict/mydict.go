package mydict

import "errors"

// Dictionary type export
type Dictionary map[string]string

var (
	errNotFount = errors.New("Not Fount")
	errCantUpdate = errors.New("Cant Update non-existing word")
	errCantDelete = errors.New("Cant Delete non-existing word")
	errWordExists = errors.New("That word is already exists")
)

// Search for word
func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word]
	if exists {
		return value, nil
	}else{
		return "", errNotFount
	}
}

// Add a word to dictionary
func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)
	if err == errNotFount {
		d[word] = def
	}else if err == nil {
		return errWordExists
	}
	return nil
}

// Update a word
func (d Dictionary) Update(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = def
	case errNotFount:
		return errCantUpdate
	}
	return nil
}

// Delete a word
func (d Dictionary) Delete(word string) error {
	_, err := d.Search(word)
	if err == errNotFount {
		return errCantDelete
	}
	delete(d, word)
	return nil
}