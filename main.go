package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/article.html", "templates/category1.html", "templates/category2.html", "templates/category3.html"))

var blog Blog

func (a *Article) URL() template.URL {
	return template.URL(fmt.Sprintf("/%s/article/%d", strings.ToLower(a.Categorie), a.ID))
}

func main() {
	loadedBlog, err := readJSON("articles.json")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier JSON:", err)
		return
	}

	blog = loadedBlog

	// Gère la route "/static/" pour servir des fichiers statiques depuis le dossier "static"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/category1/", category1Handler)
	http.HandleFunc("/category2/", category2Handler)
	http.HandleFunc("/category3/", category3Handler) // Ajout du gestionnaire de catégorie
	http.HandleFunc("/article/", articleHandler)
	http.HandleFunc("/search/", searchHandler)

	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/admin/add/", addArticleHandler)
	http.HandleFunc("/admin/delete/", deleteArticleHandler)

	// Démarrez le serveur web
	http.ListenAndServe(":8080", nil)
}

type Article struct {
	ID           int    `json:"id"`
	Categorie    string `json:"categorie"`
	Titre        string `json:"titre"`
	Auteur       string `json:"auteur"`
	Contenu      string `json:"contenu"`
	Images       Image  `json:"images"`
	ContenuCourt string
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

	// Ajouter le début du contenu
	for i := range blog.Articles {
		const maxContentLength = 300
		if len(blog.Articles[i].Contenu) > maxContentLength {
			blog.Articles[i].ContenuCourt = blog.Articles[i].Contenu[:maxContentLength] + "..."
		} else {
			blog.Articles[i].ContenuCourt = blog.Articles[i].Contenu
		}
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
	categorie := r.URL.Query().Get("category")

	var articlesToDisplay []Article

	if categorie == "" {
		// If no category specified, select 10 random articles from the blog
		articlesToDisplay = getRandomArticles(blog.Articles, 10)
	} else {
		// If category specified, filter articles by category
		var filteredArticles []Article
		for _, article := range blog.Articles {
			if strings.ToLower(article.Categorie) == strings.ToLower(categorie) {
				filteredArticles = append(filteredArticles, article)
			}
		}
		// Select 10 random articles from the filtered list
		articlesToDisplay = getRandomArticles(filteredArticles, 10)
	}

	templates.ExecuteTemplate(w, "index.html", articlesToDisplay)
}

// Function to get n random articles from a given list
func getRandomArticles(articles []Article, n int) []Article {
	if n >= len(articles) {
		return articles
	}

	// Shuffle the articles randomly
	shuffledArticles := make([]Article, len(articles))
	perm := rand.Perm(len(articles))
	for i, randIndex := range perm {
		shuffledArticles[i] = articles[randIndex]
	}

	// Select the first n articles from the shuffled list
	return shuffledArticles[:n]
}

func category1Handler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les articles de la catégorie 1
	categoryArticles := getArticlesByCategory("TOPS 10")
	templates.ExecuteTemplate(w, "category1.html", categoryArticles)
}

func category2Handler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les articles de la catégorie 2
	categoryArticles := getArticlesByCategory("Tutoriels")
	templates.ExecuteTemplate(w, "category2.html", categoryArticles)
}

func category3Handler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les articles de la catégorie 3
	categoryArticles := getArticlesByCategory("Nouveautes")
	templates.ExecuteTemplate(w, "category3.html", categoryArticles)
}

// Fonction utilitaire pour récupérer les articles d'une catégorie spécifique
func getArticlesByCategory(category string) []Article {
	var categoryArticles []Article
	for _, article := range blog.Articles {
		if strings.ToLower(article.Categorie) == strings.ToLower(category) {
			categoryArticles = append(categoryArticles, article)
		}
	}
	return categoryArticles
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
	term := r.URL.Query().Get("term")
	var results []Article
	for _, article := range blog.Articles {
		if strings.Contains(strings.ToLower(article.Titre), strings.ToLower(term)) {
			results = append(results, article)
		}
	}
	templates.ExecuteTemplate(w, "search.html", results)
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
