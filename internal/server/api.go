// Пакет для работы с сервером и обработчиками API.
package server

import (
	"GoExamCensor/internal/config"
	"GoExamCensor/internal/logger"
	"GoExamCensor/internal/middleware"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

// Request - структура запроса на цензурирование комментария.
type Request struct {
	Content string `json:"content"`
}

// Censor возвращает на запрос код 200, если комментарий успешно прошел
// проверку цензора, либо код 400, если не прошел проверку и содержит
// недопустимые выражения.
func Censor(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.Censor"

		log := slog.Default().With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("request to censor comment")

		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("cannot decode request", logger.Err(err))
			http.Error(w, "cannot decode request", http.StatusBadRequest)
			return
		}
		log.Debug("request body decoded")

		if isOffensive(req.Content, cfg.CensorList) {
			log.Info("comment contains offensive words")
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		log.Info("comment is allowed")

		w.WriteHeader(http.StatusOK)
		log.Info("request served successfuly")
	}
}

// isOffensive проверяет контент комментария на содержание недопустимых
// выражений.
func isOffensive(text string, words []string) bool {
	for _, word := range words {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}
