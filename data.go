package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
)

// Article représente une entité d'article
type Article struct {
	ID        int    `json:"id"`
	Categorie string `json:"categorie"`
	Titre     string `json:"titre"`
	Contenu   string `json:"contenu"`
	Images    string `json:"images"`
}

// Blog représente une collection d'articles
type Blog struct {
	Articles []Article `json:"articles"`
}

// Fonction pour lire les articles depuis un fichier CSV
func readCSV(filePath string) (Blog, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Blog{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return Blog{}, err
	}

	var blog Blog
	for _, record := range records {
		article := Article{
			ID:       len(blog.Articles) + 1,
			Category: record[0],
			Title:    record[1],
			Content:  record[2],
		}
		blog.Articles = append(blog.Articles, article)
	}

	return blog, nil
}

// Fonction pour écrire les articles dans un fichier JSON
func writeJSON(filePath string, blog Blog) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(blog)
}
