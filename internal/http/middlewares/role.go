package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

// role - the minimal role level which will have access to resource
func RoleBasedAccess(secret string, requiredJob entities.UserJob, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header["Authorization"]

		if len(token) == 0 || token[0] == "" {
			if err := utils.RespondWith400(w, "missing authorization token"); err != nil {
				if err := utils.RespondWith500(w); err != nil {
					slog.Error("failed to respond with 500 on missing auth token",
						slog.String("path", r.URL.Path),
						slog.String("method", r.Method),
						slog.String("error", err.Error()),
					)
				}
			}
			return
		}

		tokenData := strings.Split(token[0], " ")
		if len(tokenData) != 2 {
			if err := utils.RespondWith400(w, "wrong auth token format"); err != nil {
				if err := utils.RespondWith500(w); err != nil {
					slog.Error("failed to respond with 500 on wrong auth token format",
						slog.String("path", r.URL.Path),
						slog.String("method", r.Method),
						slog.Any("token", token),
						slog.String("error", err.Error()),
					)
				}
			}
		}

		claims, ok := extractClaims(tokenData[1], secret)
		if !ok {
			if err := utils.RespondWith500(w); err != nil {
				slog.Error("failed to respond with 500 on error with parse token claims",
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.Any("token", token),
					slog.String("error", err.Error()),
				)
			}
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			slog.Error("failed to extract	`exp` from token claims", slog.Any("exp", exp))
			if err := utils.RespondWith500(w); err != nil {
				slog.Error("failed to respond with 500 on error with extract `exp` from token claims",
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.Any("token", token),
					slog.String("error", err.Error()),
				)
			}
			return
		}

		expTime := int64(exp)
		// Function validating token
		if tokenExpire(expTime) {
			if err := utils.RespondWith401(w, "token is expired"); err != nil {
				if err = utils.RespondWith500(w); err != nil {
					slog.Error("failed to respond with 500 on error with token expired",
						slog.String("path", r.URL.Path),
						slog.String("method", r.Method),
						slog.Any("token", token),
						slog.Any("exp", exp),
						slog.String("error", err.Error()),
					)

				}
			}
			return
		}

		userJob, ok := claims["job_position"]
		if !ok {
			if err := utils.RespondWith500(w); err != nil {
				slog.Error("failed to respond with 500 on failed get job_position from claims",
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.Any("token", token),
					slog.String("error", err.Error()),
				)
			}
			return
		}

		job, ok := userJob.(entities.UserJob)
		if !ok || !hasUserPermission(job, requiredJob) {
			if err := utils.RespondWith400(w, "user have no access to this resource"); err != nil {
				if err := utils.RespondWith500(w); err != nil {
					slog.Error("failed to respond with 500 on user have no permissions to resourse",
						slog.String("path", r.URL.Path),
						slog.String("method", r.Method),
						slog.Any("token", token),
						slog.String("error", err.Error()),
					)
				}
			}
			return
		}

		// User is admin, the access to resources is allowed
		slog.Info("handle request with auth",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.Any("token", token),
			slog.Any("token_claims", claims),
		)
		userId, ok := claims["id"]
		if !ok {
			slog.Error("failed to get 'id' from claims", slog.Any("claims", claims))
			next.ServeHTTP(w, r)
			return
		}
		id, ok := userId.(int)
		if !ok {
			slog.Error("failed to cast useId to int", slog.Any("userId", userId))
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(context.Background(), "user_id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
