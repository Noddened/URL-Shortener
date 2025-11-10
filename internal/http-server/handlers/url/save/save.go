package save

import (
	resp "github.com/Noddened/URL-Shortener/internal/lib/api/response"
	"github.com/Noddened/URL-Shortener/internal/lib/random"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(URL, alias string) (int64, error)
}

// TODO: перекинуть это в конфиг когда надо будет
const aliasLength = 6

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriterm r * http.Request) {
		const op = "handlers.url.save.New"

		// Для наглядности добавляем к объекту логгера op и request_id
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Создаем объект запроса и парсим запрос
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Если запрос с пустым телом
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))
			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		// Ну и на всякий случай еще логи, лучше больше, чем меньше
		log.Info("request body decoded", slog.Any("req", req))

		// Объект валидатора, в него передаем структуру
		if err := validator.New().Struct(req); err != nil {
			// приводим ошибку к типу ошибки валидации
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.Error(validateErr.Error()))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			// Запись с таким Alias уже есть
			log.Info("url already exist", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exist"))
			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias: alias,
	})
}