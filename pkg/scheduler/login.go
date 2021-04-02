package scheduler

import (
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Login struct {
}

func NewLogin(rbacPath string) *Login {
	l := &Login{}
	go func() {
		//a, err := gormadapter.NewAdapter("mysql", dao.Get().Mysql.MasterDsn())
		//if err != nil {
		//	zaplogger.Sugar().Fatal(err)
		//}
		//e, err := casbin.NewEnforcer(rbacPath, a)
		//if err != nil {
		//	zaplogger.Sugar().Fatal(err)
		//}
		//
		//// Load the policy from DB.
		//e.LoadPolicy()
		//
		//// Check the permission.
		//e.Enforce("alice", "data1", "read")
		//
		//// Modify the policy.
		//// e.AddPolicy(...)
		//// e.RemovePolicy(...)
		//
		//// Save the policy back to DB.
		//e.SavePolicy()
	}()

	return l
}

func (l *Login) LoginHandler(c *gin.Context) {
	zaplogger.Sugar().Infow("LoginHandler print token", "value", c.Request.Header.Get("Token"))
	c.JSON(http.StatusOK, "121212121")
}

func (l *Login) LogoutHandler(c *gin.Context) {

}
