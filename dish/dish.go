package dish

import (
	"errors"
	"github.com/rs/xid"
	"sync"
)

type Comment struct {
	Rating  uint8  `json:"rating"`
	Comment string `json:"comment"`
	Author  string `json:"author"`
	Date    string `json:"date"`
}

type Dish struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Image       string    `json"image"`
	Category    string    `json:"category"`
	Featured    string    `json:"featured"`
	Label       string    `json:"label"`
	Price       string    `json:"price"`
	Description string    `json:"description"`
	Comments    []Comment `json:"comments"`
}

var (
	list []Dish
	mtx  sync.RWMutex
	once sync.Once
)

func init() {
	once.Do(initializeList)
}

func initializeList() {
	list = []Dish{}

	//	// add mock data
		Add("Dish Number One")
		Add("Dish Number Two")
}

// Get retrieves all elements from the Dish list
func Get() []Dish {
	return list
}

// Add will add a new Dish with default values and the passed in Dish name
func Add(name string) string {
	d := newDish(name)
	mtx.Lock()
	list = append(list, d)
	mtx.Unlock()
	return d.ID
}

// Delete will remove a Dish from the list
func Delete(id string) error {
	location, err := findDishLocation(id)
	if err != nil {
		return err
	}
	removeElementByLocation(location)
	return nil
}

func removeElementByLocation(i int) {
	mtx.Lock()
	list = append(list[:i], list[i+1:]...)
	mtx.Unlock()
}

func findDishLocation(id string) (int, error) {
	mtx.RLock()
	defer mtx.RUnlock()
	for i, d := range list {
		if isMatchingID(d.ID, id) {
			return i, nil
		}
	}
	return 0, errors.New("could not find Dish based on id")
}

func isMatchingID(a string, b string) bool {
	return a == b
}

func newDish(name string) Dish {
	comment1 := Comment{Rating: 5, Comment: "I love it! Delicious...", Author: "Chubby Boy", Date: "04/13/2018"}
	comment2 := Comment{Rating: 2, Comment: "Very spicy!", Author: "Chauncey Billups", Date: "06/20/2018"}

	return Dish{
		xid.New().String(),
		name,
		"none",
		"Thai",
		"yes",
		"Panang",
		"9.50",
		"Beef Panang",
		[]Comment{
			comment1,
			comment2,
		},
	}
}
