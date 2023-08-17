package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userRespone struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.Users) userRespone {
	return userRespone{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPassword, err := util.HashPassWord(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		fmt.Println("sssss")
		errCode := db.ErrorCode(err)

		if errCode == db.Uniqueviolation {
			ctx.JSON(http.StatusForbidden, errResponse(err))
		}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, user)
	}
	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type logInUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type logInUserResponse struct {
	SessionID             uuid.UUID   `json:"session_id"`
	AccessToken           string      `json:"access_token"`
	AccessTokenExpiredAt  time.Time   `json:"access_token_expired_at"`
	RefreshToken          string      `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time   `json:"refresh_token_expired_at"`
	User                  userRespone `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req logInUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = util.CheckPassWord(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	refreshtoken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	sessions, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     refreshPayload.Username,
		RefreshToken: refreshtoken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	rsp := logInUserResponse{
		SessionID:             sessions.ID,
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshtoken,
		RefreshTokenExpiredAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
