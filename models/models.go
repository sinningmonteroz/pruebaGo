package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// Usuario representa la estructura de la tabla de usuarios en la base de datos
type Usuario struct {
	ID       int
	Username string
	Password string
}

// Equipo representa la estructura de datos del equipo
type Equipo struct {
	ID            int
	CodigoSistema string
	Marca         string
	REF           string
	TipoEquipo    string
	Modelo        string
	Serial        string
	Estado        string
}

const (
	DBUsername = "root"
	DBPassword = ""
	DBHost     = "localhost"
	DBPort     = "3306"
	DBName     = "wurthco_crm2016"
)

// Modelo representa la estructura del modelo en la arquitectura MVC
type Model struct {
	db *sql.DB
}

// NewModelo crea una nueva instancia del modelo
func NewModel() *Model {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUsername, DBPassword, DBHost, DBPort, DBName)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatal(err)
	}
	return &Model{db: db}
}

// GetUserByUsername obtiene un usuario por nombre de usuario
func (m *Model) GetUserByUsername(username string) (Usuario, error) {
	var user Usuario
	query := "SELECT rowid, login, pass FROM llx_user WHERE login = ?"
	row := m.db.QueryRow(query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	return user, err
}

// CreateUser crea un nuevo usuario
func (m *Model) CreateUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = m.db.Exec("INSERT INTO llx_user (login, pass) VALUES (?, ?)", username, hashedPassword)
	return err
}

// GetUserFromSession obtiene un usuario de la sesión
func (m *Model) GetUserFromSession(r *http.Request) Usuario {
	username, _ := m.GetSession(r)
	user, _ := m.GetUserByUsername(username)
	return user
}

// SetSession establece la sesión del usuario
func (m *Model) SetSession(w http.ResponseWriter, username string) {
	expiration := time.Now().Add(30 * time.Minute)
	cookie := http.Cookie{Name: "session", Value: username, Expires: expiration, Path: "/"}
	http.SetCookie(w, &cookie)
}

// ClearSession limpia la sesión del usuario
func (m *Model) ClearSession(w http.ResponseWriter) {
	cookie := http.Cookie{Name: "session", Value: "", Expires: time.Now(), Path: "/"}
	http.SetCookie(w, &cookie)
}

// GetSession obtiene la sesión del usuario
func (m *Model) GetSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// InsertarEquipo inserta un nuevo equipo en la base de datos
func (m *Model) InsertarEquipo(equipo *Equipo) error {
	query := "INSERT INTO equipos (CodigoSistema, Marca, REF, TipoEquipo, Modelo, Serial, Estado) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := m.db.Exec(query, equipo.CodigoSistema, equipo.Marca, equipo.REF, equipo.TipoEquipo, equipo.Modelo, equipo.Serial, equipo.Estado)
	if err != nil {
		return fmt.Errorf("Error al insertar equipo en la base de datos: %v", err)
	}
	return nil
}

// ObtenerEquiposPaginados obtiene una página específica de equipos con paginación y 25 registros por página
func (m *Model) ObtenerEquiposPaginados(pagina int) ([]Equipo, error) {
	elementosPorPagina := 25
	offset := (pagina - 1) * elementosPorPagina
	query := "SELECT * FROM equipos LIMIT ? OFFSET ?"
	rows, err := m.db.Query(query, elementosPorPagina, offset)
	if err != nil {
		return nil, fmt.Errorf("error al obtener equipos desde la base de datos: %v", err)
	}
	defer rows.Close()

	var equipos []Equipo
	for rows.Next() {
		var equipo Equipo
		err := rows.Scan(&equipo.ID, &equipo.CodigoSistema, &equipo.Marca, &equipo.REF, &equipo.TipoEquipo, &equipo.Modelo, &equipo.Serial, &equipo.Estado)
		if err != nil {
			return nil, fmt.Errorf("error al escanear equipo desde la base de datos: %v", err)
		}
		equipos = append(equipos, equipo)
	}

	return equipos, nil
}
