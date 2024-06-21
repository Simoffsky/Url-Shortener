package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"url-shorter/internal/jwt"
	"url-shorter/internal/models"
)

func (s *LinkServer) handleLink(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = removeTrailingSlash(r.PathValue("short"))
	switch r.Method {
	case http.MethodGet:
		s.handleRedirect(w, r)
	case http.MethodDelete:
		s.handleRemoveLink(w, r)
	case http.MethodPut:
		s.handleEditLink(w, r)
	default:
		s.writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}
func (s *LinkServer) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	var link models.Link

	token := r.Header.Get("Authorization")
	if token != "" {
		login, err := jwt.ParseJWT(token, s.config.JwtSecret)
		if err != nil {
			s.handleError(w, err)
			return
		}
		link.CreatorLogin = login
	}

	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.linkService.CreateLink(link); err != nil {
		s.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *LinkServer) handleRemoveLink(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Path

	creatorLogin, err := getUserLogin(r.Header.Get("Authorization"), s.config.JwtSecret)
	if err != nil {
		s.handleError(w, err)
		return
	}

	s.logger.Debug("Removing link: " + short)

	if err := s.linkService.RemoveLink(creatorLogin, short); err != nil {
		s.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *LinkServer) handleEditLink(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Path

	var link models.Link
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	creatorLogin, err := getUserLogin(r.Header.Get("Authorization"), s.config.JwtSecret)
	if err != nil {
		s.handleError(w, err)
		return
	}
	link.CreatorLogin = creatorLogin
	s.logger.Debug("Editing link: " + short)

	if err := s.linkService.EditLink(creatorLogin, short, link); err != nil {
		s.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *LinkServer) handleRedirect(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Path
	s.logger.Debug("Redirecting to: " + short)

	url, err := s.linkService.GetLink(short)
	if err != nil {
		s.handleError(w, err)
		return
	}

	statsVisitor := &models.LinkStatVisitor{
		LinkShort: short,
		Visitor: models.LinkVisitor{
			VisitorIP:    r.RemoteAddr,
			VisitorAgent: r.UserAgent(),
		},
		TimeAt: time.Now().Unix(),
	}

	err = s.statService.SendStat(statsVisitor)
	if err != nil {
		s.handleError(w, err)
		return
	}

	http.Redirect(w, r, url.Url, http.StatusMovedPermanently)
}

func (s *LinkServer) handleQRCode(w http.ResponseWriter, r *http.Request) {
	var qrSize int

	short := getLink(r.URL.Path)

	_, err := s.linkService.GetLink(short)
	if err != nil {
		s.handleError(w, err)
		return
	}

	link := getFullUrl(r)
	sizeQuery := r.URL.Query().Get("size")

	if sizeQuery != "" {
		qrSize, err = strconv.Atoi(sizeQuery)
		if err != nil {
			s.logger.Error("Invalid size query parameter: " + sizeQuery)
			qrSize = 0
		}
	}

	s.logger.Debug("Getting QR code for: " + link)

	qr, err := s.linkService.GetQRCode(link, qrSize)
	if err != nil {
		s.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	if _, err := w.Write(qr); err != nil {
		s.handleError(w, err)
	}
}

func (s *LinkServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	s.logger.Debug("senging gRPC request with login: " + user.Login)
	if err := s.authService.Register(user); err != nil {
		s.handleError(w, err)
		return
	}
	fmt.Println("Registered successfully")
	w.WriteHeader(http.StatusCreated)
}

func (s *LinkServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	token, err := s.authService.Login(user)
	if err != nil {
		s.handleError(w, err)
		return
	}

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)

}

// post with
func (s *LinkServer) handleStats(w http.ResponseWriter, r *http.Request) {
	short := getLink(r.URL.Path)
	stats, err := s.statService.GetStatForLink(short)
	if err != nil {
		s.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		s.handleError(w, err)
	}
}

func (s *LinkServer) handleError(w http.ResponseWriter, err error) {
	var modelErr models.Error
	if !errors.As(err, &modelErr) {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}
	s.writeError(w, modelErr.StatusCode, err)
}

func (s *LinkServer) writeError(w http.ResponseWriter, statusCode int, err error) {
	s.logger.Error(fmt.Sprintf("HTTP error(%d): %s", statusCode, err.Error()))
	http.Error(w, err.Error(), statusCode)
}

func getFullUrl(r *http.Request) string {
	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}
	return scheme + "://" + r.Host + r.RequestURI
}

func getLink(path string) string {
	trimmedPath := strings.TrimPrefix(path, "/")
	parts := strings.Split(trimmedPath, "/")
	lastPart := parts[len(parts)-1]
	cleanLastPart := strings.Split(lastPart, "?")[0]
	return cleanLastPart
}
func removeTrailingSlash(short string) string {
	if short == "" {
		return short
	}
	if short[len(short)-1] == '/' {
		return short[:len(short)-1]
	}
	return short
}

func getUserLogin(token, jwtSecret string) (string, error) {
	if token == "" {
		return "", nil
	}
	login, err := jwt.ParseJWT(token, jwtSecret)
	if err != nil {
		return "", err
	}
	return login, nil
}
