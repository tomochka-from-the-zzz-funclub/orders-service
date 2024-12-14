package transport

import (
	"consumer/internal/config"
	myLog "consumer/internal/logger"
	"consumer/internal/service"
	"fmt"
	"html/template"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
)

type HandlersBuilder struct {
	srv                   service.InterfaceService
	rout                  *router.Router
	templUserGet          *template.Template
	templUserLogin        *template.Template
	templUserRegistration *template.Template
	templAdminLogin       *template.Template
	templAdminAllMethods  *template.Template
}

func HandleCreate(cfg config.Config, s service.InterfaceService) {
	tmplg, err1 := template.ParseFiles("../app/getUser.html")
	tmpll, err2 := template.ParseFiles("../app/loginUser.html")
	tmplr, err3 := template.ParseFiles("../app/registrationUser.html")
	templAl, err4 := template.ParseFiles("../app/loginAdmin.html")
	templAm, err5 := template.ParseFiles("../app/workAdmin.html")

	if err1 != nil {
		myLog.Log.Errorf("GetHtml error during parsing of file getUser: %v", err1)
		return
	}
	if err2 != nil {
		myLog.Log.Errorf("GetHtml error during parsing of file loginUser: %v", err2)
		return
	}
	if err3 != nil {
		myLog.Log.Errorf("GetHtml error during parsing of file registrUser: %v", err3)
		return
	}
	if err4 != nil {
		myLog.Log.Errorf("GetHtml error during parsing of file loginAdmin: %v", err4)
		return
	}
	if err5 != nil {
		myLog.Log.Errorf("GetHtml error during parsing of file workAdmin: %v", err5)
		return
	}
	hb := HandlersBuilder{
		srv:                   s,
		rout:                  router.New(),
		templUserGet:          tmplg,
		templUserLogin:        tmpll,
		templUserRegistration: tmplr,
		templAdminLogin:       templAl,
		templAdminAllMethods:  templAm,
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8090", nil)
	}()
	//user html
	hb.rout.GET("/api/v1/user", hb.GetHtml())
	hb.rout.GET("/api/v1/user/auth/login", hb.GetHtmlLoginUser())
	hb.rout.GET("/api/v1/user/auth/registration", hb.GetHtmlRegistr())

	//user work methods
	hb.rout.GET("/api/v1/user/get/order", hb.Get())
	hb.rout.GET("/api/v1/user/get/status", hb.GetStatus())

	//user auth
	hb.rout.POST("/api/v1/user/registration", hb.Registration())
	hb.rout.PUT("/api/v1/user/login", hb.LoginUser())

	//admin html
	hb.rout.GET("/api/v1/admin", hb.GetHtmlWorkAdmin())
	hb.rout.GET("/api/v1/admin/auth/login", hb.GetHtmlLoginAdmin())

	//admin work methods
	hb.rout.PUT("/api/v1/admin/update/status", hb.UpdateStatus())
	hb.rout.POST("/api/v1/admin/adition/delivery_man", hb.CreateDeliveryMan())
	hb.rout.PUT("/api/v1/admin/create/delivery", hb.GiveOrderToDeliveryMan())

	//admin auth
	hb.rout.PUT("/api/v1/admin/login", hb.LoginAdmin())

	//hb.rout.PUT("/api/v1/check_delivery", hb.CheckDeliveryStart())

	fmt.Println(fasthttp.ListenAndServe(":8080", hb.rout.Handler))
}
