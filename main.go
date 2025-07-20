package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func connectDb() *sql.DB {
	file, _ := os.Executable()
	db, err := sql.Open("sqlite3", filepath.Dir(file)+"/db/Dictionaries.db")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return db
}

func getFromDb(word string, languageFrom string, languageTo string) []string {
	db := connectDb()
	defer db.Close()

	var fromId string
	err := db.QueryRow("SELECT id FROM "+languageFrom+" WHERE word = ?", word).Scan(&fromId)
	if err != nil {
		return nil
	}

	rows, err := db.Query("SELECT translation_id FROM "+languageFrom+"_"+languageTo+" WHERE word_id = ?", fromId)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	var toId []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			fmt.Fprintln(os.Stdout, err)
		}
		toId = append(toId, id)
	}

	var translations []string
	for _, to := range toId {
		var word string
		err = db.QueryRow("SELECT word FROM "+languageTo+" WHERE id = ?", to).Scan(&word)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil
		}
		translations = append(translations, word)
	}
	return translations
}

func main() {
	word := os.Args[1]
	transaltions := getFromDb(word, "en", "ru")

	if len(transaltions) > 0 {
		fmt.Println(word)
		fmt.Println("-------------------")
		for _, translation := range transaltions {
			fmt.Fprintln(os.Stdout, translation)
		}
	}
}
