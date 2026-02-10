package handler

import (
	"myblog/database"
	"myblog/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoginResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Uid   int    `json:"uid"`
	Token string `json:"token"`
}

type RegisterResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.PostForm("user")
		pass := ctx.PostForm("pass")
		if len(name) == 0 {
			ctx.JSON(
				http.StatusBadRequest,
				&LoginResponse{
					Code:  1,
					Msg:   "must indicate user name",
					Uid:   0,
					Token: "",
				},
			)
			return
		}
		if len(pass) != 32 {
			ctx.JSON(
				http.StatusBadRequest,
				&LoginResponse{
					Code:  2,
					Msg:   "invalid password",
					Uid:   0,
					Token: "",
				},
			)
			return
		}
		user := database.GetUserByName(name)
		if user == nil {
			ctx.JSON(
				http.StatusForbidden,
				&LoginResponse{
					Code:  3,
					Msg:   "user not exist",
					Uid:   0,
					Token: "",
				},
			)
			return
		}
		if user.PassWd != pass {
			ctx.JSON(
				http.StatusForbidden,
				&LoginResponse{
					Code:  4,
					Msg:   "incorrect password",
					Uid:   0,
					Token: "",
				})
			return
		}
		zap.L().Info("user login success", zap.String("name", name), zap.Int("uid", user.Id))

		header := &util.JwtHeader{}
		payload := &util.JwtPayload{
			Issue:       "blog",
			IssueAt:     time.Now().Unix(),
			Expiration:  time.Now().Add(database.TOKEN_EXPIRE).Add(24 * time.Hour).Unix(),
			UserDefined: map[string]any{"uid": user.Id},
		}
		jwtToken, err := util.GenJwt(header, payload, util.JWT_SECRET)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				&LoginResponse{
					Code:  5,
					Msg:   "generate jwtToken failed",
					Uid:   0,
					Token: "",
				},
			)
			return
		}
		refreshToken := util.RandStringRunes(20)
		database.SetToken(refreshToken, jwtToken)
		ctx.SetCookie(
			"refresh_token",
			refreshToken,
			int(database.TOKEN_EXPIRE.Seconds()),
			"/",
			"",
			false,
			true,
		)
		ctx.JSON(
			http.StatusOK,
			&LoginResponse{
				Code:  0,
				Msg:   "success",
				Uid:   user.Id,
				Token: jwtToken,
			},
		)

	}
}

func GetAuthToken(ctx *gin.Context) {
	refreshToken := ctx.PostForm("refresh_token")
	authToken := database.GetToken(refreshToken)
	ctx.String(http.StatusOK, authToken)
}

func NewRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.PostForm("user")
		pass := ctx.PostForm("pass")

		if len(name) == 0 {
			ctx.JSON(
				http.StatusBadRequest,
				&RegisterResponse{
					Code: 1,
					Msg:  "must indicate user name",
				},
			)
			return
		}

		if len(pass) != 32 {
			ctx.JSON(
				http.StatusBadRequest,
				&RegisterResponse{
					Code: 2,
					Msg:  "invalid password",
				},
			)
			return
		}

		if database.GetUserByName(name) != nil {
			ctx.JSON(
				http.StatusConflict,
				&RegisterResponse{
					Code: 3,
					Msg:  "user already exist",
				},
			)
			return
		}

		err := database.CreateUser(name, pass)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				&RegisterResponse{
					Code: 4,
					Msg:  "create user failed",
				},
			)
			return
		}

		ctx.JSON(
			http.StatusOK,
			&RegisterResponse{
				Code: 0,
				Msg:  "success",
			},
		)
	}
}
