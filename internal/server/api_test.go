// Пакет для работы с сервером и обработчиками API.

package server

import (
	"GoExamCensor/internal/config"
	"GoExamCensor/internal/logger"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var offensive = []string{"qwerty", "йцукен", "zxvbnm"}

func Test_isOffensive(t *testing.T) {

	tests := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "Not_offensive",
			text: "Пакет для работы с сервером и обработчиками API",
			want: false,
		},
		{
			name: "Offensive",
			text: "Пакет для работы qwerty и обработчиками API",
			want: true,
		},
		{
			name: "Offensive_without_spaces",
			text: "Пакетдляработыqwertyи обработчиками API",
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := isOffensive(tt.text, offensive); got != tt.want {
				t.Errorf("isOffensive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCensor(t *testing.T) {
	logger.Discard()
	cfg := &config.Config{
		CensorList: offensive,
	}

	tests := []struct {
		name     string
		text     string
		wantCode int
	}{
		{
			name:     "Not_offensive",
			text:     `{"content":"Пакет для работы с сервером и обработчиками API"}`,
			wantCode: 200,
		},
		{
			name:     "Offensive",
			text:     `{"content":"Пакет для работы qwerty и обработчиками API"}`,
			wantCode: 400,
		},
		{
			name:     "Offensive_without_spaces",
			text:     `{"content":"Пакетдляработыqwertyи обработчиками API"}`,
			wantCode: 400,
		},
		{
			name:     "Offensive_with_fields",
			text:     `{"postId":1234,"content":"Пакетдляработыqwertyи обработчиками API"}`,
			wantCode: 400,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			mux.HandleFunc("POST /", Censor(cfg))
			srv := httptest.NewServer(mux)
			defer srv.Close()

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.text))
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("Censor() code = %v, want %v", rr.Code, tt.wantCode)
			}
		})
	}
}
