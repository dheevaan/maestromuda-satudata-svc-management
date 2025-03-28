package middleware

import (
	"data-management/src/config"
	"data-management/src/model"
	"data-management/src/service"
	"data-management/src/util/encryption/uaes"
	"log"
	"os"
	"time"

	secret "data-management/src/util/db"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var identityKey = "id"

func InitJwt() (authMiddleware *jwt.GinJWTMiddleware, err error) {
	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "Kuningan",
		Key:         []byte(os.Getenv(config.ENV_KEY_JWT_SECRECT)),
		Timeout:     24 * time.Minute,
		MaxRefresh:  730 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.User{
				Username: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (res interface{}, err error) {
			var param model.User_Login
			if err := c.ShouldBind(&param); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			userRes := model.User{}
			if _, ok := service.NewUserService().GetUserWithValidatePassword("username", param.Username, param.Password, &userRes); ok {
				c.Set("user", userRes)
				res = &userRes
				return
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true

			//? This will limit only superadmin can access the API
			// if v, ok := data.(*model.User); ok && v.Username == "superadmin" {
			// 	return true
			// }
			// return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			if userRaw, found := c.Get("user"); found {
				if user, ok := userRaw.(model.User); ok {
					data := model.Resp_User_Login{
						User: &user,
						Auth: model.Resp_JwtToken{
							Token:  message,
							Expire: time,
						},
					}
					if os.Getenv("PROD_MODE") == "true" {
						SECRET := secret.GenerateRandomString(7)
						log.Println("SECRET", SECRET)
						var Uaes = uaes.NewAES(SECRET)
						resEn, _ := Uaes.Encrypt_Any(data)
						resPon := SECRET + resEn
						c.JSON(code, model.Response{
							Data: &resPon,
						})
					} else {
						c.JSON(code, model.Resp_User_Login{
							User: &user,
							Auth: model.Resp_JwtToken{
								Token:  message,
								Expire: time,
							},
						})
					}
				}
			}
		},

		RefreshResponse: func(c *gin.Context, code int, message string, time time.Time) {
			if userRaw, found := c.Get("user"); found {
				user, _ := userRaw.(model.User)
				c.JSON(code, model.Resp_User_Login{
					User: &user,
					Auth: model.Resp_JwtToken{
						Token:  "message",
						Expire: time,
					},
				})
			}
		},
	})

	if err != nil {
		log.Println("JWT Error:" + err.Error())
	}

	return
}
