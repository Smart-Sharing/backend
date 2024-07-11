package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var data struct {
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}

	bodyBin, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read r.Body in Login", slog.String("error", err.Error()))
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500 during parse login data",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	err = json.Unmarshal(bodyBin, &data)
	if err != nil {
		slog.Error("failed to unmarshal login data", slog.String("error", err.Error()))
		if err = utils.RespondWith400(w, "login data is wrong"); err != nil {
			slog.Error("failed to respond with 500 during wrong login data",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)

		}
		return
	}

	user, err := h.service.GetUserByPhoneNumber(data.PhoneNumber)
	if err != nil {
		slog.Error("failed to get user from db by phone number",
			slog.String("phone_number", data.PhoneNumber),
			slog.String("error", err.Error()),
		)

		if err = utils.RespondWith400(w, "user with such phone number doesn't exists"); err != nil {
			msg := fmt.Sprintf("failed to respond with 400 during not found user with phone_number=%s", data.PhoneNumber)
			slog.Error(msg,
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("phone_number", data.PhoneNumber),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	if data.Password != user.Password {
		if err = utils.RespondWith400(w, "user password is not correct"); err != nil {
			slog.Error("failed to respond with 500 during user password is not correct",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("phone_number", data.PhoneNumber),
				slog.String("right_password", user.Password),
				slog.String("user_input_password", data.Password),
				slog.String("error", err.Error()),
			)

		}
		return
	}
	token, err := h.service.GenerateToken(*user, h.cfg.Secret, h.cfg.TokenTTL)
	if err != nil {
		slog.Error("failed to generate JWT token", slog.String("error", err.Error()))
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500 during generate jwt token",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
	}

	response := struct {
		Token string `json:"token"`
	}{Token: token}

	if err := utils.RespondWithJSON(w, 200, response); err != nil {
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with JSON with JWT token",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("token", token),
				slog.String("error", err.Error()),
			)
		}
	}
}
