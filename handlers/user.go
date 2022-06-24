package handlers

import (
	"agchavez/go/rest-ws/models"
	"agchavez/go/rest-ws/repository"
	"agchavez/go/rest-ws/server"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type SingUpRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SingUpResponse struct {
	Menssage string `json:"message"`
	Status   bool   `json:"status"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type UserResponse struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

type GetUsersResponse struct {
	Users []UserResponse `json:"users"`
	Count int            `json:"count"`
}

type ListError struct {
	Errors []string `json:"errors"`
}

func SingUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SingUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var user = models.User{
			LastName:  strings.ToUpper(request.LastName),
			Email:     request.Email,
			Password:  string(hashedPassword),
			FirstName: strings.ToUpper(request.FirstName),
		}

		err = repository.InsertUser(r.Context(), &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(
			SingUpResponse{
				Menssage: "User created successfully",
				Status:   true,
				Email:    user.Email,
			})
	}
}

func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = LoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//Validations data
		var listErrors ListError

		if request.Email == "" {
			listErrors.Errors = append(listErrors.Errors, "El email es obligatorio")
		} else {
			re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
			if re.FindAllStringSubmatch(request.Email, -1) == nil {
				listErrors.Errors = append(listErrors.Errors, "El email no es valido")
			}
		}
		if len(request.Password) < 5 {
			listErrors.Errors = append(listErrors.Errors, "El password debe tener al menos 5 caracteres")
		}

		if len(listErrors.Errors) > 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(listErrors)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
		if err != nil {
			http.Error(w, "Invalid data", http.StatusUnauthorized)
			return
		}

		claims := models.AppClaims{
			UserID: user.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(
			LoginResponse{
				Email: user.Email,
				Name:  user.FirstName + " " + user.LastName,
				Token: tokenString,
			})

	}
}

// List all users
func ListUsersHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params = models.ParamsQuery{}
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")
		if limit == "" {
			limit = "15"
		}
		params.Limit, _ = strconv.Atoi(limit)
		if offset == "" {
			offset = "0"
		}
		params.Offset, _ = strconv.Atoi(offset)
		users, err := repository.GetUsers(r.Context(), params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		usersResponse := make([]UserResponse, len(users))
		for i, user := range users {
			usersResponse[i] = UserResponse{
				ID:        uint64(user.ID),
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetUsersResponse{
			Users: usersResponse,
			Count: len(users),
		})
	}
}

func GetUserHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := repository.GetUserByID(r.Context(), int(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		userResponse := UserResponse{
			ID:        uint64(user.ID),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userResponse)
	}
}
