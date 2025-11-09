package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)
		log.Info("logger middleware enabled")

		// Обработчик:
		fn := func(w http.ResponseWriter, r *http.Request) {
			// собираем исходную информацию о запросе
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			// Создаем обертку http.ResponseWriter
			// для получния сведений об ответе
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// момент получения запроса
			t1 := time.Now()

			// Запись идет в лог в defer
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			// Передаем управление следующему обработчику

			next.ServeHTTP(ww, r)
		}

		// Возвращаем созданный выше обработчик
		return http.HandlerFunc(fn)
	}
}
