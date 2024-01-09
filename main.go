package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/article.html"))

var blog Blog

func main() {
	// Chargez les données depuis un fichier CSV
	blog, _ = readCSV("articles.csv")

	// Configuration des gestionnaires de routage
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/article/", articleHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/admin/add/", addArticleHandler)
	http.HandleFunc("/admin/delete/", deleteArticleHandler)

	// Démarrez le serveur web
	http.ListenAndServe(":8080", nil)
}

// Gestionnaire pour la page d'accueil
func indexHandler(w http.ResponseWriter, r *http.Request) {
	categorie := r.URL.Query().Get("categorie")

	var filteredArticles []Article
	if categorie != "" {
		for _, article := range blog.Articles {
			if strings.ToLower(article.Categorie) == strings.ToLower(categorie) {
				filteredArticles = append(filteredArticles, article)
			}
		}
	} else {
		filteredArticles = blog.Articles
	}

	templates.ExecuteTemplate(w, "index.html", filteredArticles)
}

// Gestionnaire pour la page d'article
func articleHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de l'article de l'URL
	id := strings.TrimPrefix(r.URL.Path, "/article/")
	for _, article := range blog.Articles {
		if id == strconv.Itoa(article.ID) {
			templates.ExecuteTemplate(w, "article.html", article)
			return
		}
	}
	http.NotFound(w, r)
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
