package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"

	"github.com/boltdb/bolt"
)

var db *bolt.DB
var open bool
var pathSite string

func Open() error {
	var err error
	dbfile := path.Join(pathSite + "/db/bolt.db")
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(dbfile, 0600, config)
	if err != nil {
		log.Fatal(err)
	}
	open = true
	return nil
}

func Close() {
	open = false
	db.Close()
}

type Person struct {
	ID   string
	Name string
	Age  string
	Job  string
}

func (p *Person) save() error {
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}
	err := db.Batch(func(tx *bolt.Tx) error {
		people, err := tx.CreateBucketIfNotExists([]byte("people"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		enc, err := p.encode()
		if err != nil {
			return fmt.Errorf("could not encode Person %s: %s", p.ID, err)
		}
		err = people.Put([]byte(p.ID), enc)
		return err
	})
	return err
}

func delete(buket, id string) error {
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}
	err := db.Batch(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(buket)).Delete([]byte(id))
		return err
	})
	return err
}

func (p *Person) gobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func gobDecode(data []byte) (*Person, error) {
	var p *Person
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Person) encode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func GetPerson(id string) (*Users, error) {
	if !open {
		return nil, fmt.Errorf("db must be opened before saving!")
	}
	var u *Users
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("user"))
		k := []byte(id)
		u, err = decode(b.Get(k))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Could not get Person ID %s", id)
		return nil, err
	}
	return u, nil
}

func List(bucket string) {
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucket)).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})
}
func RandStringBytesMask(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}

func ListPrefix(bucket, prefix string) {
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucket)).Cursor()
		p := []byte(prefix)
		for k, v := c.Seek(p); bytes.HasPrefix(k, p); k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})
}

func ListRange(bucket, start, stop string) {
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucket)).Cursor()
		min := []byte(start)
		max := []byte(stop)
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			fmt.Printf("%s: %s\n", k, v)
		}
		return nil
	})
}

func cut(text string, limit int) string {
	runes := []rune(text)
	countStr := len(runes) - limit
	if len(runes) >= countStr {
		return string(runes[:countStr])
	}
	return text
}

func main() {
	pathSiteGo, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	//testAuthJob02
	pathSite = cut(pathSiteGo, 3)

	Open()
	defer Close()

	r := mux.NewRouter()

	r.HandleFunc("/registration", registration)
	r.HandleFunc("/login", checkAuth)
	r.HandleFunc("/forget", goForget)
	r.HandleFunc("/logout", logut)
	r.HandleFunc("/adm/allPhones", allPhones)
	r.HandleFunc("/adm/myPhones", myPhone)
	r.HandleFunc("/adm/favorite", favorite)
	r.HandleFunc("/adm/addPhone", addPhone)
	r.HandleFunc("/adm/editPhone/{id}", editPhone)
	r.HandleFunc("/go/addPhone", goAddPhone)
	r.HandleFunc("/go/editPhone", goEditPhone)
	r.HandleFunc("/go/f/{phoneId}/{res}/{name}/{phone}", goAddFavorite)
	r.HandleFunc("/go/delete/{id}", goDelete)
	r.HandleFunc("/", getIndex)

	fs := http.FileServer(http.Dir(pathSite + "/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":4888", nil))

}
