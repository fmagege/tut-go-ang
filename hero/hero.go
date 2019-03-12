package hero

import (
	"errors"
	"github.com/rs/xid"
	"sync"
)

var (
	list []Hero
	mtx sync.RWMutex
	once sync.Once
)

type Hero struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

func init() {
	once.Do(initializeList)
}

func initializeList() {
	list = []Hero{}

	//	// add mock data
	Add("Jason Bourne")
	Add("Clark Kent")
}

func Get() []Hero {
	return list
}

func Add(name string) string {
	h := newHero(name)
	mtx.Lock()
	list = append(list, h)
	mtx.Unlock()
	return h.ID
}

func Delete(id string) error {
	location, err := findHeroLocation(id)
	if err != nil {
		return err
	}
	removeElementByLocation(location)
	return nil
}

func findHeroLocation(id string) (int, error) {
	mtx.RLock()
	defer mtx.RUnlock()
	for i, h := range list {
		if isMatchingID(h.ID, id) {
			return i, nil
		}
	}
	return 0, errors.New("could not find Hero based on id")
}

func newHero(name string) Hero {
	return Hero{
		ID: xid.New().String(),
		Name: name,
	}
}

func removeElementByLocation(i int) {
	mtx.Lock()
	list = append(list[:i], list[i+1:]...)
	mtx.Unlock()
}

func isMatchingID(a string, b string) bool {
	return a == b
}
