package main

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"AppSistemas/controllers"
	"AppSistemas/models"

)

// IMPORTAR LOS MODELOS

var modelo *models.Model

var templates *template.Template

func init() {
	// Cargar los templates al inicio de la aplicación
	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error al cargar las plantillas:", err)
	}
}

func main() {
	modelo = models.NewModel()

	r := mux.NewRouter()

	// Ruta para los archivos estáticos (CSS, imágenes, etc.)
	r.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))

	r.HandleFunc("/", checkSession(homeHandler)).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("GET")
	r.HandleFunc("/login", authHandler).Methods("POST")
	r.HandleFunc("/create_user", createUserHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")
	r.HandleFunc("/auth", authHandler).Methods("POST")
	// Rutas de la aplicación
	r.HandleFunc("/registrar", registerDeviceHandler).Methods("GET")
	r.HandleFunc("/registerdevice", func(w http.ResponseWriter, r *http.Request) {
		controllers.RegistrarEquipoHandler(w, r, modelo)
	})

	// Resto de tus rutas...

	r.HandleFunc("/", homeHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

//FUNCIONES

func homeHandler(w http.ResponseWriter, r *http.Request) {
	user := modelo.GetUserFromSession(r)
	renderTemplate(w, "index", map[string]interface{}{"Username": user.Username})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := modelo.GetUserByUsername(username)
	if err != nil {
		log.Println(err)
		renderTemplate(w, "login", map[string]interface{}{"Error": "Error de autenticación"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		renderTemplate(w, "login", map[string]interface{}{"Error": "Credenciales incorrectas"})
		return
	}

	modelo.SetSession(w, user.Username)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		err := modelo.CreateUser(username, password)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
			return
		}

		modelo.SetSession(w, username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	renderTemplate(w, "create_user", nil)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	modelo.ClearSession(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func checkSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar si hay una sesión activa
		if !isSessionActive(r) {
			// Si no hay una sesión activa, redirigir a la página de inicio de sesión
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Si hay una sesión activa, pasar al siguiente manejador
		next(w, r)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	// Eliminar la cookie de sesión (o realiza cualquier otra acción necesaria para cerrar sesión)
	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})

	// Redirigir a la página de inicio de sesión
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func isSessionActive(r *http.Request) bool {
	// Implementa la lógica para verificar si hay una sesión activa
	// Puedes usar cookies, tokens JWT, o algún otro mecanismo.
	// En este ejemplo, simplemente devuelve true si la cookie "session" existe.
	_, err := r.Cookie("session")
	return err == nil
}

func registerDeviceHandler(w http.ResponseWriter, r *http.Request) {
	// Manejar la ruta "/register"
	err := templates.ExecuteTemplate(w, "registerdevice.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CODIGO PARA CORRER TEMPLANTES

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, templateName string, data map[string]interface{}) {
	if tmpl == nil {
		tmpl = template.New("").Funcs(template.FuncMap{
			"sumar": func(a, b int) int {
				return a + b
			},
			"restar": func(a, b int) int {
				return a - b
			},
		})
	}

	err := tmpl.ExecuteTemplate(w, templateName+".html", data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
