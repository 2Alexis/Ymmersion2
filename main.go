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

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/article.html", "templates/category1.html", "templates/category2.html", "templates/category3.html", "templates/ajout.html", "templates/login.html", "templates/admin.html", "templates/ymmersion2.html", "templates/mentionslegales.html", "templates/contact.html"))

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
	http.HandleFunc("/ymmersion2/", ymmersion2Handler)
	http.HandleFunc("/mentionslegales/", mentionsHandler)
	http.HandleFunc("/contact/", contactHandler)
	http.HandleFunc("/category1/", category1Handler)
	http.HandleFunc("/category2/", category2Handler)
	http.HandleFunc("/category3/", category3Handler) // Ajout du gestionnaire de catégorie
	http.HandleFunc("/article/", articleHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/login/", loginHandler)
	http.Handle("/admin/", adminMiddleware(http.HandlerFunc(adminHandler)))

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

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
	// other fields
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

func authenticateUser(username, password string) (User, error) {
	// Load users from a JSON file (replace "users.json" with your actual file path)
	users, err := loadUsersFromJSON("users.json")
	if err != nil {
		return User{}, fmt.Errorf("error loading users: %v", err)
	}

	// Find the user with the provided username
	var authenticatedUser User
	for _, user := range users {
		if user.Username == username {
			// Check if the password matches
			if strings.EqualFold(user.Password, password) {
				authenticatedUser = user
				break
			} else {
				return User{}, fmt.Errorf("incorrect password for user: %s", username)
			}
		}
	}

	// Check if the user was found
	if authenticatedUser.Username == "" {
		return User{}, fmt.Errorf("user not found: %s", username)
	}

	return authenticatedUser, nil
}
func adminHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authenticated as an admin
	// Check if the user is authenticated as an admin
	user, err := authenticateUser(r.FormValue("username"), r.FormValue("password"))

	// Debug print
	fmt.Printf("Admin Authentication Result - User: %+v, Error: %v\n", user, err)

	if err != nil || !user.Admin {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}
	// Check the request method
	switch r.Method {
	case http.MethodGet:
		// Display the admin page
		templates.ExecuteTemplate(w, "admin.html", nil)
	case http.MethodPost:
		// Handle form submissions (e.g., add or delete articles)
		action := r.FormValue("action")
		if action == "add" {
			// Handle adding article logic here
		} else if action == "delete" {
			// Handle deleting article logic here
		} else {
			http.Error(w, "Invalid action", http.StatusBadRequest)
			return
		}

		// Redirect or show a success message
		http.Redirect(w, r, "/admin/", http.StatusSeeOther)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET
	if r.Method == http.MethodGet {
		// Display the login form
		templates.ExecuteTemplate(w, "login.html", nil)
		return
	}

	// Check if the method is POST (process the login form)
	if r.Method == http.MethodPost {
		// Retrieve form data
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Authenticate the user
		user, err := authenticateUser(username, password)
		if err != nil {
			// Authentication failed, redirect to the login page
			http.Redirect(w, r, "/admin/", http.StatusSeeOther)
			return
		}

		// Check if the user is an admin
		if !user.Admin {
			// User is not an admin, display an error message
			templates.ExecuteTemplate(w, "login.html", "Access denied")
			return
		}

		// Authentication successful, redirect to the admin page
		http.Redirect(w, r, "/admin/", http.StatusSeeOther)
		return
	}
}
func adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract and validate user credentials (for example, from cookies or headers)
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Authenticate the user
		user, err := authenticateUser(username, password)
		if err != nil {
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}

		// Check if the user is an admin
		if !user.Admin {
			http.Error(w, "Not authorized", http.StatusForbidden)
			return
		}

		// Call the next handler if authenticated and authorized
		next.ServeHTTP(w, r)
	})
}

// Function to load users from a JSON file
func loadUsersFromJSON(filePath string) ([]User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var users []User
	err = decoder.Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func ymmersion2Handler(w http.ResponseWriter, r *http.Request) {

	templates.ExecuteTemplate(w, "ymmersion2.html", r)
}

func mentionsHandler(w http.ResponseWriter, r *http.Request) {

	templates.ExecuteTemplate(w, "mentionslegales.html", r)
}
func contactHandler(w http.ResponseWriter, r *http.Request) {

	templates.ExecuteTemplate(w, "contact.html", r)
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
