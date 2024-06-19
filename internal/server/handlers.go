package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"url-shorter/internal/models"
)

func (s *LinkServer) handler(w http.ResponseWriter, r *http.Request) {
	var link models.Link

	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.linkService.CreateLink(link); err != nil {
		if errors.Is(err, models.ErrLinkAlreadyExists) {
			s.writeError(w, http.StatusBadRequest, err)
		} else {
			s.writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *LinkServer) handleRedirect(w http.ResponseWriter, r *http.Request) {
	short := removeTrailingSlash(r.URL.Path[len("/"):])

	s.logger.Debug("Redirecting to: " + short)

	url, err := s.linkService.GetLink(short)
	if err != nil {
		if errors.Is(err, models.ErrLinkNotFound) {
			s.writeError(w, http.StatusNotFound, err)
		} else {
			s.writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	http.Redirect(w, r, url.Url, http.StatusMovedPermanently)
}

func (s *LinkServer) handleQRCode(w http.ResponseWriter, r *http.Request) {

	link := getFullUrl(r)

	sizeQuery := r.URL.Query().Get("size")
	size, err := strconv.Atoi(sizeQuery)

	if err != nil {
		s.logger.Error("Invalid size query parameter: " + sizeQuery)
		size = 0 // 0 means default size
	}

	s.logger.Debug("Getting QR code for: " + link)

	qr, err := s.linkService.GetQRCode(link, size)
	if err != nil {
		if errors.Is(err, models.ErrLinkNotFound) {
			s.writeError(w, http.StatusNotFound, err)
		} else {
			s.writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Set("Content-Type", "image/png")
	if _, err := w.Write(qr); err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
	}
}

func (s *LinkServer) writeError(w http.ResponseWriter, errCode int, err error) {
	if err == nil {
		err = errors.New("(WARNING)!: writeError() called with nil error")
	}
	s.logger.Error(fmt.Sprintf("HTTP error(%d): %s", errCode, err.Error()))
	http.Error(w, err.Error(), errCode)
}

type Request struct {
	Url      string `json:"url"`
	ShortUrl string `json:"short_url"`
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
func removeTrailingSlash(short string) string {
	if short == "" {
		return short
	}
	if short[len(short)-1] == '/' {
		return short[:len(short)-1]
	}
	return short
}
