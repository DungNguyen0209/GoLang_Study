package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RenewAccessTokenUserRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,alphanum"`
}

type RenewAccessTokenUserResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiredAt time.Time `json:"access_token_expired_at"`
}

func (server *Server) RenewAccessTokenUser(ctx *gin.Context) {
	var req RenewAccessTokenUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if session.IsBlocked {
		err = fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err = fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err = fmt.Errorf("missmatch token")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	if time.Now().After(session.ExpiredAt) {
		err = fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	rsp := RenewAccessTokenUserResponse{
		AccessToken:          accessToken,
		AccessTokenExpiredAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
