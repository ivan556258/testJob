package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

type PhoneBook struct {
	ID       string
	Name     string
	Phone    string
	Owner    string
	Favorite string
}

var IsLetter = regexp.MustCompile(`[a-zA-Z]+$`).MatchString

func allPhones(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "user_id")
	if err != nil {
		fmt.Println(err)
	}
	if session.Values["user_id"] == nil {
		http.Redirect(w, r, "http://localhost:4888", 301)
		return
	}
	var parsedData []PhoneBook
	var dat map[string]string
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("phoneBook")).Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {

			if err := json.Unmarshal(v, &dat); err != nil {
				panic(err)
			}
			if dat["Name"] != "" || dat["Phone"] != "" {

				id, _ := json.Marshal(dat["ID"])
				name, _ := json.Marshal(dat["Name"])
				phone, _ := json.Marshal(dat["Phone"])

				parsedData = append(parsedData, PhoneBook{
					ID:    strings.Trim(string(id), "\""),
					Name:  strings.Trim(string(name), "\""),
					Phone: strings.Trim(string(phone), "\""),
				})

			}

		}

		t, _ := template.ParseFiles(pathSite + "/html/lk/members.html")

		type n map[string]interface{}

		t.Execute(w, n{
			"data": parsedData,
		})

		return nil
	})

}

func getIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(pathSite + "/html/i.html")
	t.Execute(w, nil)
}

func myPhone(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_id")
	if err != nil {
		fmt.Println(err)
	}

	if session.Values["user_id"] == nil {
		http.Redirect(w, r, "http://localhost:4888", 301)
		return
	}
	var parsedData []PhoneBook
	var dat map[string]string
	idOwner := session.Values["user_id"].(string)
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("phoneBook")).Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {

			if err := json.Unmarshal(v, &dat); err != nil {
				panic(err)
			}
			if dat["Owner"] == idOwner {

				id, _ := json.Marshal(dat["ID"])
				name, _ := json.Marshal(dat["Name"])
				phone, _ := json.Marshal(dat["Phone"])
				favorite, _ := json.Marshal(dat["Favorite"])

				parsedData = append(parsedData, PhoneBook{
					ID:       strings.Trim(string(id), "\""),
					Name:     strings.Trim(string(name), "\""),
					Phone:    strings.Trim(string(phone), "\""),
					Favorite: strings.Trim(string(favorite), "\""),
				})

			}
		}

		t, _ := template.ParseFiles(pathSite + "/html/lk/myPhone.html")

		type n map[string]interface{}

		t.Execute(w, n{
			"data": parsedData,
		})

		return nil
	})
}

func favorite(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_id")
	if err != nil {
		fmt.Println(err)
	}
	if session.Values["user_id"] == nil {
		http.Redirect(w, r, "http://localhost:4888", 301)
		return
	}
	var parsedData []PhoneBook
	var dat map[string]string
	idOwner := session.Values["user_id"].(string)
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("phoneBook")).Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {

			if err := json.Unmarshal(v, &dat); err != nil {
				panic(err)
			}
			if dat["Owner"] == idOwner && dat["Favorite"] == "1" {

				id, _ := json.Marshal(dat["ID"])
				name, _ := json.Marshal(dat["Name"])
				phone, _ := json.Marshal(dat["Phone"])
				favorite, _ := json.Marshal(dat["Favorite"])

				parsedData = append(parsedData, PhoneBook{
					ID:       strings.Trim(string(id), "\""),
					Name:     strings.Trim(string(name), "\""),
					Phone:    strings.Trim(string(phone), "\""),
					Favorite: strings.Trim(string(favorite), "\""),
				})

			}
		}

		t, _ := template.ParseFiles(pathSite + "/html/lk/favorite.html")

		type n map[string]interface{}

		t.Execute(w, n{
			"data": parsedData,
		})

		return nil
	})
}

func goEditPhone(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_id")
	if err != nil {
		fmt.Println(err)
	}
	if session.Values["user_id"] == nil {
		http.Redirect(w, r, "http://localhost:4888", 301)
		return
	}
	changePhone(w, r)
}

func addPhone(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_id")
	if err != nil {
		fmt.Println(err)
	}
	if session.Values["user_id"] == nil {
		http.Redirect(w, r, "http://localhost:4888", 301)
		return
	}

	t, _ := template.ParseFiles(pathSite + "/html/lk/addPhone.html")

	type n map[string]interface{}

	t.Execute(w, nil)
}

func editPhone(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_id")
	if err != nil {
		fmt.Println(err)
	}
	if session.Values["user_id"] == nil {
		http.Redirect(w, r, "http://localhost:4888", 301)
		return
	}
	var parsedData PhoneBook
	var dat map[string]string
	vars := mux.Vars(r)
	idOwner := session.Values["user_id"].(string)
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("phoneBook")).Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {

			if err := json.Unmarshal(v, &dat); err != nil {
				panic(err)
			}
			if dat["Owner"] == idOwner && dat["ID"] == vars["id"] {

				id, _ := json.Marshal(dat["ID"])
				name, _ := json.Marshal(dat["Name"])
				phone, _ := json.Marshal(dat["Phone"])
				favorite, _ := json.Marshal(dat["Favorite"])
				parsedData = PhoneBook{
					ID:       strings.Trim(string(id), "\""),
					Name:     strings.Trim(string(name), "\""),
					Phone:    strings.Trim(string(phone), "\""),
					Favorite: strings.Trim(string(favorite), "\""),
				}

			}
		}

		t, _ := template.ParseFiles(pathSite + "/html/lk/editPhone.html")

		type n map[string]interface{}

		t.Execute(w, n{
			"data": parsedData,
		})

		return nil
	})
}

func goAddPhone(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_id")
	if err != nil {
		fmt.Println(err)
	}
	if session.Values["user_id"] == nil {
		http.Redirect(w, r, "http://localhost:4888", 301)
		return
	}
	newPhone(w, r)
}

func goAddFavorite(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_id")
	if err != nil {
		fmt.Println(err)
	}
	if session.Values["user_id"] == nil {
		http.Redirect(w, r, "http://localhost:4888", 301)
		return
	}
	changeFavorite(w, r)
	http.Redirect(w, r, "/adm/myPhones", http.StatusFound)

}

func newPhone(w http.ResponseWriter, r *http.Request) error {
	a := *&PhoneBook{}
	var err error
	session, _ := store.Get(r, "user_id")
	r.ParseMultipartForm(1024)
	if !open {
		fmt.Errorf("db must be opened before saving!")
	}
	err = db.Batch(func(tx *bolt.Tx) error {
		articale, err := tx.CreateBucketIfNotExists([]byte("phoneBook"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		if r.FormValue("id") != "" {
			a.ID = string(r.FormValue("id"))
		} else {
			id, err := articale.NextSequence()
			if err != nil {
				log.Fatal(id)
			}
			a.ID = strconv.FormatUint(id, 10)
		}
		if len(r.FormValue("name")) > 100 || len(r.FormValue("phone")) > 25 {
			fmt.Fprintf(w, "Very long text. Need little text.")
			return nil
		}
		isPhone := IsLetter(r.FormValue("phone"))
		if isPhone == true {
			fmt.Fprintf(w, "Word don't use in number phone")
			return nil
		}
		a.Name = r.FormValue("name")
		a.Phone = r.FormValue("phone")
		a.Owner = session.Values["user_id"].(string)
		if r.FormValue("favorite") == "on" {
			a.Favorite = "1"
		} else {
			a.Favorite = "0"
		}

		enc, err := json.Marshal(a)
		if err != nil {
			return err
		}
		err = articale.Put([]byte(a.ID), enc)
		return err
	})
	fmt.Fprintf(w, "1")
	return err
}

func changeFavorite(w http.ResponseWriter, r *http.Request) error {
	a := *&PhoneBook{}
	var err error
	vars := mux.Vars(r)
	session, _ := store.Get(r, "user_id")
	if !open {
		fmt.Errorf("db must be opened before saving!")
	}
	err = db.Batch(func(tx *bolt.Tx) error {
		articale, err := tx.CreateBucketIfNotExists([]byte("phoneBook"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		a.ID = vars["phoneId"]
		if vars["res"] == "0" {
			a.Favorite = "1"
		} else {
			a.Favorite = "0"
		}

		a.Name = vars["name"]
		a.Phone = vars["phone"]
		a.Owner = session.Values["user_id"].(string)

		enc, err := json.Marshal(a)
		if err != nil {
			return err
		}
		err = articale.Put([]byte(a.ID), enc)
		return err
	})
	return err
}

func changePhone(w http.ResponseWriter, r *http.Request) error {
	a := *&PhoneBook{}
	var err error
	session, _ := store.Get(r, "user_id")
	r.ParseMultipartForm(1024)
	if !open {
		fmt.Errorf("db must be opened before saving!")
	}
	err = db.Batch(func(tx *bolt.Tx) error {
		articale, err := tx.CreateBucketIfNotExists([]byte("phoneBook"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		a.ID = r.FormValue("idPhone")
		if r.FormValue("favorite") == "on" {
			a.Favorite = "1"
		} else {
			a.Favorite = "0"
		}
		if len(r.FormValue("name")) > 100 || len(r.FormValue("phone")) > 25 {
			fmt.Fprintf(w, "Very long text. Need little text.")
			return nil
		}
		isPhone := IsLetter(r.FormValue("phone"))
		if isPhone == true {
			fmt.Fprintf(w, "Word dom't use in number phone")
			return nil
		}
		a.Name = r.FormValue("name")
		a.Phone = r.FormValue("phone")
		a.Owner = session.Values["user_id"].(string)

		enc, err := json.Marshal(a)
		if err != nil {
			return err
		}
		err = articale.Put([]byte(a.ID), enc)
		return err
	})
	fmt.Fprintf(w, "1")
	return err
}

func goDelete(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_id")
	if err != nil {
		fmt.Println(err)
	}
	vars := mux.Vars(r)
	if session.Values["user_id"] == nil {
		http.Redirect(w, r, "http://localhost:4888", 301)
		return
	}
	delete("phoneBook", vars["id"])
	http.Redirect(w, r, "/adm/myPhones", http.StatusFound)
}
