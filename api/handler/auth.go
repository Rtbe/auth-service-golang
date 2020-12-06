package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/auth-service-go/api/model"
	"example.com/auth-service-go/internal/entity"
	"example.com/auth-service-go/internal/repository"
	"github.com/go-chi/chi"
)

//InitAuthRoutes initializes /auth subrouter
func (h *Handler) InitAuthRoutes(repo repository.Token) {
	h.Router.Route("/auth", func(r chi.Router) {
		r.Get("/user/{userID}", get(h.Context, repo))
		r.Post("/tokens/refresh", refreshTokens(h.Context, repo))
		r.Delete("/refresh", deleteRefreshToken(h.Context, repo))
		r.Delete("/user/refresh", deleteUserRefreshTokens(h.Context, repo))
	})
}

func get(ctx context.Context, repo repository.Token) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "userID")
		if id == "" {
			respondWithError("Error parsing user id from request url", http.StatusBadRequest, w)
			return
		}

		tokenPair, err := entity.CreateTokenPair(id)
		if err != nil {
			respondWithError("Error creating tokens", http.StatusInternalServerError, w)
			return
		}

		respTokens := model.TokenPair{
			AccessToken: tokenPair.AccessToken.Token,
			//Convert refresh token into base64 string before sending it to the user.
			RefreshToken: entity.EncodeToken64(tokenPair.RefreshToken.Token),
		}

		err = repo.Insert(ctx, tokenPair)
		if err != nil {
			respondWithError("Error inserting tokens into DB", http.StatusInternalServerError, w)
			return
		}

		respondWithJSON("data", respTokens, http.StatusOK, w)
	}
}

func refreshTokens(ctx context.Context, repo repository.Token) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokens := &model.TokenPair{}
		err := json.NewDecoder(r.Body).Decode(tokens)
		if err != nil {
			respondWithError("Error parsing pair of tokens", http.StatusBadRequest, w)
			return
		}
		if tokens.AccessToken == "" {
			respondWithError("Access token is empty", http.StatusBadRequest, w)
			return
		}
		if tokens.RefreshToken == "" {
			respondWithError("Refresh token is empty", http.StatusBadRequest, w)
			return
		}

		refreshToken, err := entity.DecodeToken64(tokens.RefreshToken)
		if err != nil {
			respondWithError(fmt.Sprintf("Error decoding refresh token %s", tokens.RefreshToken), http.StatusInternalServerError, w)
			return
		}

		claimsRefreshToken, err := entity.ParseRefreshToken(refreshToken)
		if err != nil {
			respondWithError(fmt.Sprintf("Error parsing refresh token %s. %s", refreshToken, err.Error()), http.StatusInternalServerError, w)
			return
		}
		claimsAccessToken, err := entity.ParseAccessToken(tokens.AccessToken)
		if err != nil {
			respondWithError(fmt.Sprintf("Error parsing refresh token %s. %s", refreshToken, err.Error()), http.StatusInternalServerError, w)
			return
		}
		//Check bind between access and refresh token.
		if claimsAccessToken.Refresh_uuid != claimsRefreshToken.UUID {
			respondWithError("Access token does not belong to refresh token", http.StatusInternalServerError, w)
			return
		}

		isRefreshTokenInDB := repo.IsRefreshTokenInDB(ctx, claimsRefreshToken.UUID)
		if !isRefreshTokenInDB {
			respondWithError("There is no such refresh token", http.StatusNotFound, w)
			return
		}

		err = repo.RefreshTokenSetIsUsed(ctx, claimsAccessToken.Refresh_uuid)
		if err != nil {
			respondWithError(err.Error(), http.StatusInternalServerError, w)
			return
		}

		tokenPair, err := entity.CreateTokenPair(claimsRefreshToken.User_id)
		if err != nil {
			respondWithError("Error creating tokens", http.StatusInternalServerError, w)
			return
		}
		respTokens := model.TokenPair{
			AccessToken: tokenPair.AccessToken.Token,
			//Convert refresh token into base64 string before sending it to the user.
			RefreshToken: entity.EncodeToken64(tokenPair.RefreshToken.Token),
		}

		err = repo.Insert(ctx, tokenPair)
		if err != nil {
			respondWithError("Error inserting tokens into DB", http.StatusInternalServerError, w)
			return
		}

		respondWithJSON("data", respTokens, http.StatusOK, w)
	}
}

func deleteRefreshToken(ctx context.Context, repo repository.Token) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestRefreshToken := model.RefreshToken{}

		err := json.NewDecoder(r.Body).Decode(&requestRefreshToken)
		if err != nil {
			respondWithError("Error parsing refresh token", http.StatusBadRequest, w)
			return
		}
		if requestRefreshToken.Token == "" {
			respondWithError("Refresh token is empty", http.StatusBadRequest, w)
			return
		}

		refreshToken, err := entity.DecodeToken64(requestRefreshToken.Token)
		if err != nil {
			respondWithError(fmt.Sprintf("Error decoding refresh token %s", requestRefreshToken.Token), http.StatusInternalServerError, w)
			return
		}
		claimsRefreshToken, err := entity.ParseRefreshToken(refreshToken)
		if err != nil {
			respondWithError(fmt.Sprintf("Error parsing refresh token: %s. %s", refreshToken, err.Error()), http.StatusInternalServerError, w)
			return
		}

		userID := claimsRefreshToken.User_id
		isUserInDB := repo.IsUserInDB(ctx, userID)
		if !isUserInDB {
			respondWithError("There is no such user", http.StatusNotFound, w)
			return
		}

		refreshTokenUUID := claimsRefreshToken.UUID
		isRefreshTokenInDB := repo.IsRefreshTokenInDB(ctx, refreshTokenUUID)
		if !isRefreshTokenInDB {
			respondWithError("There is no such refresh token", http.StatusNotFound, w)
			return
		}

		err = repo.DeleteRefreshToken(ctx, userID, refreshTokenUUID)
		if err != nil {
			respondWithError("Error deleting refresh token", http.StatusInternalServerError, w)
			return
		}

		respondWithJSON("message", "Refresh token was successfully deleted", http.StatusOK, w)
	}
}

func deleteUserRefreshTokens(ctx context.Context, repo repository.Token) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		u := &model.User{}
		err := json.NewDecoder(r.Body).Decode(u)
		if err != nil {
			respondWithError("Error parsing user id", http.StatusBadRequest, w)
			return
		}
		if u.UserID == "" {
			respondWithError("User id is empty", http.StatusBadRequest, w)
			return
		}

		isUserInDB := repo.IsUserInDB(ctx, u.UserID)
		if !isUserInDB {
			respondWithError("There is no such user", http.StatusNotFound, w)
			return
		}

		err = repo.DeleteUserRefreshTokens(ctx, u.UserID)
		if err != nil {
			respondWithError("Error deleting refresh tokens", http.StatusInternalServerError, w)
			return
		}

		respondWithJSON("message", "User refresh tokens was successfully deleted", http.StatusOK, w)
	}
}
