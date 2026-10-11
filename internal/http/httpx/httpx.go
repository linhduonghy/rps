package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	apperror "rps/internal/errors"
	"rps/internal/http/response"
)

func WriteJSON[T any](
	writer http.ResponseWriter,
	resp response.Response[T],
) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(resp.Code)
	return json.NewEncoder(writer).Encode(resp)
}

func WriteSuccess[T any](
	writer http.ResponseWriter,
	data T,
) error {
	resp := response.Success(data)
	return WriteJSON(writer, resp)
}

func WriteSuccessDetail[T any](
	writer http.ResponseWriter,
	status int,
	msg string,
	data T,
) error {
	resp := response.SuccessWithCodeAndMsg(status, msg, data)
	return WriteJSON(writer, resp)
}

func WriteErrorDetail(
	w http.ResponseWriter,
	status int,
	msg string,
) {
	_ = WriteJSON(
		w,
		response.Error(status, msg),
	)
}

func WriteError(
	w http.ResponseWriter,
	err error,
) {
	switch {
	case errors.Is(err, apperror.NotFoundErr):
		WriteErrorDetail(w, http.StatusNotFound, err.Error())

	case errors.Is(err, apperror.AlreadyExistsErr):
		WriteErrorDetail(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, apperror.InvalidInputErr):
		WriteErrorDetail(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, apperror.UnauthorizedErr):
		WriteErrorDetail(w, http.StatusUnauthorized, err.Error())

	case errors.Is(err, apperror.ForbiddenErr):
		WriteErrorDetail(w, http.StatusForbidden, err.Error())

	case errors.Is(err, apperror.InternalServerErr):
		WriteErrorDetail(w, http.StatusInternalServerError, err.Error())

	default:
		WriteErrorDetail(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
	}
}
