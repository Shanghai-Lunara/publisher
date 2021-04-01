package scheduler

import (
	"fmt"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

type Login struct {
}

func (l *Login) LoginHandler(c *gin.Context) {
	a, err := c.Cookie("test-cookies")
	if err != nil {
		zaplogger.Sugar().Error(err)
	}
	zaplogger.Sugar().Infow("print cookie", "value", a)

	c.SetCookie("test-cookies", fmt.Sprintf("random-cookie-%d", rand.Intn(9999999)), 1000, "/", "", false, false)
	c.JSON(http.StatusOK, "121212121")
}

func (l *Login) LogoutHandler(c *gin.Context) {

}
