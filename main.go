package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/article.html"))

var blog Blog

func main() {
	loadedBlog, err := readJSON("articles.json")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier JSON:", err)
		return
	}

	blog = loadedBlog

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/article/", articleHandler)
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/admin/add/", addArticleHandler)
	http.HandleFunc("/admin/delete/", deleteArticleHandler)

	// Démarrez le serveur web
	http.ListenAndServe(":8080", nil)
}

type Article struct {
	ID        int    `json:"id"`
	Categorie string `json:"categorie"`
	Titre     string `json:"titre"`
	Contenu   string `json:"contenu"`
	Images    Image  `json:"images"`
}

type Image struct {
	URL string `json:"url"`
}

// Blog représente une collection d'articles
type Blog struct {
	Articles []Article `json:"articles"`
}

func readJSON(filePath string) (Blog, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Blog{}, fmt.Errorf("erreur lors de l'ouverture du fichier JSON: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var blog Blog
	err = decoder.Decode(&blog)
	if err != nil {
		return Blog{}, fmt.Errorf("erreur lors du décodage du fichier JSON: %v", err)
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
func indexHandler(w http.ResponseWriter, r *http.Request) {
	categorie := r.URL.Query().Get("categorie")

	if categorie == "" {
		templates.ExecuteTemplate(w, "index.html", blog.Articles)
	} else {
		var filteredArticles []Article
		for _, article := range blog.Articles {
			if strings.ToLower(article.Categorie) == strings.ToLower(categorie) {
				filteredArticles = append(filteredArticles, article)
			}
		}
		templates.ExecuteTemplate(w, "index.html", filteredArticles)
	}
}

// Gestionnaire pour la page des catégories
func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "categories.html", blog.Articles)
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de l'article de l'URL
	id := strings.TrimPrefix(r.URL.Path, "/article/")
	articleID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// Rechercher l'article par ID
	var foundArticle *Article
	for _, article := range blog.Articles {
		if article.ID == articleID {
			foundArticle = &article
			break
		}
	}

	// Vérifier si l'article a été trouvé
	if foundArticle == nil {
		http.NotFound(w, r)
		return
	}

	// Exécuter le template spécifique pour l'article
	templates.ExecuteTemplate(w, "article.html", foundArticle)
}

// Gestionnaire pour la recherche
func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire le terme de recherche de l'URL
	term := strings.TrimPrefix(r.URL.Path, "/search/")
	var results []Article
	for _, article := range blog.Articles {
		if strings.Contains(strings.ToLower(article.Titre), strings.ToLower(term)) {
			results = append(results, article)
		}
	}
	templates.ExecuteTemplate(w, "index.html", results)
}

// Gestionnaire pour la partie administration
func adminHandler(w http.ResponseWriter, r *http.Request) {
	// Implémenter les fonctionnalités d'administration
	// ...
}

// Gestionnaire pour l'ajout d'article
func addArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Implémenter l'ajout d'article
	// ...
}

// Gestionnaire pour la suppression d'article
func deleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Implémenter la suppression d'article
	// ...
}
