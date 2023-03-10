package server

import(
	"net/http"

	"rgb/internal/store"

	"github.com/gin-gonic/gin"
)

func signUp(ctx *gin.Context) {
	user := new(store.User)

	if err := ctx.Bind(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := store.AddUser(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return

	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "signed in successfully",
		"jwt": generateJWT(user),
	})
}

func signIn(ctx *gin.Context) {
	user := new(store.User)

	if err := ctx.Bind(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := store.Authenticate(user.Username, user.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"err": "Failed to sign in.",
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "signed in successfully",
		"jwt": generateJWT(user),
	})
	
}