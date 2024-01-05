package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"AppSistemas/models"

)

// RegistrarEquipoHandler maneja la solicitud de registro de equipos
func RegistrarEquipoHandler(w http.ResponseWriter, r *http.Request, modelo *models.Model) {
	if r.Method == http.MethodPost {
		equipo := models.Equipo{
			CodigoSistema: r.FormValue("CodigoSistema"),
			Marca:         r.FormValue("Marca"),
			REF:           r.FormValue("REF"),
			TipoEquipo:    r.FormValue("TipoEquipo"),
			Modelo:        r.FormValue("Modelo"),
			Serial:        r.FormValue("Serial"),
			Estado:        r.FormValue("Estado"),
		}

		err := models.NewModel().InsertarEquipo(&equipo)
		if err != nil {
			http.Error(w, "Error al registrar el equipo", http.StatusInternalServerError)
			return
		}

		// Redirigir a una página de éxito o hacer lo que necesites después del registro
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Si no es una solicitud POST, renderizar el formulario de registro
	tmpl, err := template.ParseFiles("templates/registerdevice.html")
	if err != nil {
		http.Error(w, "Error al cargar el formulario de registro", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// RegistrarUsuarioHandler maneja la solicitud de registro de usuarios
func RegistrarUsuarioHandler(w http.ResponseWriter, r *http.Request, modelo *models.Model) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		err := modelo.CreateUser(username, password)
		if err != nil {
			http.Error(w, "Error al registrar el usuario", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Renderizar el formulario de registro
	// (puedes utilizar un paquete de plantillas como html/template aquí)
	fmt.Fprint(w, "Formulario de registro")
}

// IniciarSesionHandler maneja la solicitud de inicio de sesión
func IniciarSesionHandler(w http.ResponseWriter, r *http.Request, modelo *models.Model) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := modelo.GetUserByUsername(username)
		if err != nil {
			http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Error(w, "Contraseña incorrecta", http.StatusUnauthorized)
			return
		}

		// Iniciar sesión y redirigir a la página principal
		modelo.SetSession(w, username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Renderizar el formulario de inicio de sesión
	// (puedes utilizar un paquete de plantillas como html/template aquí)
	fmt.Fprint(w, "Formulario de inicio de sesión")
}

// MostrarEquiposHandler maneja la solicitud para mostrar equipos con paginación y un número mínimo de registros por página
func MostrarEquiposHandler(w http.ResponseWriter, r *http.Request, modelo *models.Model) {
	// Obtener el número de página desde los parámetros de la URL
	pagina := obtenerNumeroPagina(r)

	// Obtener los equipos para la página actual
	equipos, err := modelo.ObtenerEquiposPaginados(pagina)
	if err != nil {
		http.Error(w, "Error al obtener equipos", http.StatusInternalServerError)
		return
	}

	// Realizar la conversión de la página a un número entero con un valor predeterminado de 1
	numPagina, err := strconv.Atoi(fmt.Sprintf("%v", pagina))
	if err != nil || numPagina <= 0 {
		numPagina = 1
	}

	// Calcular la página anterior y pasarla a la plantilla
	paginaAnterior := numPagina - 1

	// Crear el mapa de funciones
	funcMap := template.FuncMap{
		"sumar": func(a, b int) int {
			return a + b
		},
		"restar": func(a, b int) int {
			return a - b
		},
	}

	// Crear una nueva instancia de template con el mapa de funciones
	tmpl := template.New("").Funcs(funcMap)

	// Renderizar la vista con los equipos utilizando tu función renderTemplate y el nuevo mapa de funciones
	renderTemplate(w, tmpl, "equipos.html", map[string]interface{}{
		"Equipos":        equipos,
		"Pagina":         numPagina,
		"PaginaAnterior": paginaAnterior,
	})
}

// obtenerNumeroPagina obtiene el número de página de los parámetros de la URL
func obtenerNumeroPagina(r *http.Request) int {
	// Implementa la lógica para obtener el número de página de los parámetros de la URL
	// Puedes utilizar mux.Vars(r) si estás utilizando gorilla/mux.
	// En este ejemplo, simplemente devuelve 1 si no se proporciona un número de página válido.
	// Puedes extender esta lógica según tus necesidades.
	// Se recomienda validar y manejar adecuadamente los errores en un entorno de producción.
	pagina := 1
	if p, ok := r.URL.Query()["pagina"]; ok {
		// Convierte el valor del parámetro "pagina" a un número entero
		// y maneja posibles errores, por ejemplo, si no es un número válido.
		// En este ejemplo, simplemente se establece a 1 si no es válido.
		numPagina, err := strconv.Atoi(p[0])
		if err == nil && numPagina > 0 {
			pagina = numPagina
		}
	}
	return pagina
}

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

var funcMap = template.FuncMap{
	"sumar": func(a, b int) int {
		return a + b
	},
	"restar": func(a, b int) int {
		return a - b
	},
}
