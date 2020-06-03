package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"regexp"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

var sessionKey = RandStringBytesMask(12)

var store = sessions.NewCookieStore([]byte(sessionKey))

type Users struct {
	ID       string
	Login    string
	Password string
	Phone    string
	Email    string
}

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

func userSave(w http.ResponseWriter, r *http.Request) (string, error) {
	u := *&Users{}
	var newPass string
	var dat map[string]string
	var err error
	var settingErr string = ""
	r.ParseMultipartForm(1024)
	if !open {
		return "db must be opened before saving!", err
	}

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("user"))

		b.ForEach(func(k, v []byte) error {
			if err := json.Unmarshal(v, &dat); err != nil {
				panic(err)
			}

			if string(dat["Email"]) == string(r.FormValue("email")) {
				settingErr = "Email is in data base"
				return err
			}
			return nil
		})
		return nil
	})
	if settingErr != "" {
		fmt.Fprintf(w, settingErr)
		return settingErr, err
	}

	if !validateEmail(r.FormValue("email")) {
		settingErr = "Email address is invalid"
		fmt.Fprintf(w, settingErr)
		return settingErr, err
	}

	err = db.Batch(func(tx *bolt.Tx) error {
		users, err := tx.CreateBucketIfNotExists([]byte("user"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		id, err := users.NextSequence()
		if err != nil {
			log.Fatal(err)
		}
		u.ID = strconv.FormatUint(id, 10)
		u.Email = r.FormValue("email")

		newPass = stringWithCharset(13, charset)

		password := []byte(newPass)
		// Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		u.Password = string(hashedPassword)

		enc, err := json.Marshal(u)
		if err != nil {
			return err
		}
		err = users.Put([]byte(u.ID), enc)
		return err
	})

	sendEmail(string(r.FormValue("email")), newPass)

	return "1", err
}

func sendEmail(email, code string) {
	auth := smtp.PlainAuth("", "iisusnawin@yandex.ru", "testAuthJob02", "smtp.yandex.ru")
	body := "Your Login: " + email + " \nYour password: " + code
	header := make(map[string]string)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
	err := smtp.SendMail("smtp.yandex.ru:25", auth, "iisusnawin@yandex.ru", []string{email}, []byte(message))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Email Sent!")
}

func GetUser(login string) (*Users, error) {
	if !open {
		return nil, fmt.Errorf("db must be opened before saving!")
	}
	var u *Users
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("user"))
		k := []byte(login)
		u, err = decode(b.Get(k))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Could not get Person ID %s", login)
		return nil, err
	}
	return u, nil
}

func decode(data []byte) (*Users, error) {
	var u *Users
	err := json.Unmarshal(data, &u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func checkAuth(w http.ResponseWriter, r *http.Request) {
	var u *Users
	var err error
	var dat map[string]string
	var settingErr string = ""
	r.ParseMultipartForm(1024)

	if err != nil {
		log.Fatal(err)
	}

	if !validateEmail(r.FormValue("email")) {
		settingErr = "Email address is invalid"
		fmt.Fprintf(w, settingErr)
		return
	}

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("user"))

		b.ForEach(func(k, v []byte) error {
			if err := json.Unmarshal(v, &dat); err != nil {
				panic(err)
			}
			if dat["Email"] == r.FormValue("email") {
				u, _ = GetUser(string(dat["ID"]))
				return err
			}
			return nil
		})
		return nil
	})

	if len(r.FormValue("password")) != 13 {
		settingErr = "Change login or password"
		fmt.Fprintf(w, settingErr)
		return
	}
	if u.Email != string(r.FormValue("email")) {
		settingErr = "Change login or password"
		fmt.Fprintf(w, settingErr)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(r.FormValue("password")))
	if err != nil {
		fmt.Fprintf(w, "Change login or password")
		return
	}

	session, err := store.Get(r, "user_id")
	if err != nil {
		fmt.Println(err)
	}
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 4,
		HttpOnly: true,
	}
	session.Values["user_id"] = string(u.ID)

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "1")
	return

}

func logut(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_id")
	if err != nil {
		fmt.Println(err)
	}
	session.Values["user_id"] = nil
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(session.Values["user_id"])
	http.Redirect(w, r, "http://localhost:4888", 301)
	return
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func registration(w http.ResponseWriter, r *http.Request) {
	userSave(w, r)
}

func goForget(w http.ResponseWriter, r *http.Request) {
	var uId string
	var err error
	var dat map[string]string
	var settingErr string = ""
	r.ParseMultipartForm(1024)

	if !validateEmail(r.FormValue("email")) {
		settingErr = "Email address is invalid"
		fmt.Fprintf(w, settingErr)
		return
	}

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("user"))

		b.ForEach(func(k, v []byte) error {
			if err := json.Unmarshal(v, &dat); err != nil {
				panic(err)
			}
			if dat["Email"] == r.FormValue("email") {
				uId, _ = dat["ID"]
				return err
			}
			return nil
		})
		return nil
	})
	delete("user", uId)
	resId, _ := userSave(w, r)
	fmt.Fprintf(w, resId)
}
