package transport

import (
	myErrors "consumer/internal/errors"
	my_errors "consumer/internal/errors"
	myLog "consumer/internal/logger"
	"net/http"

	"github.com/valyala/fasthttp"
)

func (hb *HandlersBuilder) GetHtml() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func GetHtml")
		if ctx.IsGet() {
			err1 := hb.templUserGet.Execute(ctx.Response.BodyWriter(), nil)
			if err1 != nil {
				myLog.Log.Errorf("GetHtml error during executing of file get: %v", err1)
				ctx.Response.SetStatusCode(400)
				return
			}
			ctx.Response.Header.Add("content-type", "text/html")
		}
	}, "GetHtml")
}

func (hb *HandlersBuilder) GetHtmlLoginUser() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func GetHtml")
		if ctx.IsGet() {

			err2 := hb.templUserLogin.Execute(ctx.Response.BodyWriter(), nil)

			if err2 != nil {
				myLog.Log.Errorf("GetHtml error during executing of file login: %v", err2)
				ctx.Response.SetStatusCode(400)
				return
			}
			ctx.Response.Header.Add("content-type", "text/html")
		}
	}, "GetHtml")
}

func (hb *HandlersBuilder) GetHtmlRegistr() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func GetHtml")
		if ctx.IsGet() {
			err3 := hb.templUserRegistration.Execute(ctx.Response.BodyWriter(), nil)
			if err3 != nil {
				myLog.Log.Errorf("GetHtml error during executing of file registr: %v", err3)
				ctx.Response.SetStatusCode(400)
				return
			}
			ctx.Response.Header.Add("content-type", "text/html")
		}
	}, "GetHtml")
}

func (hb *HandlersBuilder) Registration() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func Registration: %+v", string(ctx.Request.Body()))
		if ctx.IsPost() {
			user, err := ParseJsonUser(ctx)
			if err != nil {
				myLog.Log.Errorf("err: %v", err.Error())
				WriteJson(ctx, "Error parse")
			} else {
				err := hb.srv.Registration(user)
				if err != nil {
					myLog.Log.Errorf("Registration", err.Error())
					ctx.SetStatusCode(http.StatusBadRequest)
				} else {
					myLog.Log.Debugf("Sucses Registration")
				}
			}
			///сделвать редирект на вход
		} else {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Warnf("message from func Registration %v", myErrors.ErrMethodNotAllowed.Error())
		}
	}, "Registration")
}

func (hb *HandlersBuilder) LoginUser() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func Login: %+v", string(ctx.Request.Body()))
		if ctx.IsPut() {
			phone, err := ParseJsonLogin(ctx)
			if err != nil {
				myLog.Log.Errorf("err: %v", err.Error())
				WriteJson(ctx, "Error parse")
			} else {
				err := hb.srv.CheckRegistration(phone)
				if err != nil {
					if err == myErrors.ErrNotFoundUser {
						myLog.Log.Errorf("Error Login", err.Error())
						ctx.SetStatusCode(fasthttp.StatusNotFound)
						//ctx.Redirect("http://localhost:8080/api/v1/auth/registration", http.StatusSeeOther)
						//ctx.Redirect("http://localhost:8081/api/v1/set", http.StatusSeeOther) // для наглядности
						// сделать редирект на регистрацию
					} else {
						WriteJson(ctx, "Error login")
						ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					}
				} else {
					myLog.Log.Debugf("Sucses Login")
					token, err := hb.srv.GenerateRandomToken(phone)
					if err != nil {
						myLog.Log.Debugf("Error generate token: %+v", err)
					} else {
						myLog.Log.Debugf("Token: %+v", token)
					}
					//WriteJson(ctx, "Sucses Login")
					WriteJsonToken(ctx, token)
					//сделать редирект на домашнюю стрницу
				}
			}
		} else {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Warnf("message from func Login %v", myErrors.ErrMethodNotAllowed.Error())
		}
	}, "Login")
}

func (hb *HandlersBuilder) GetStatus() func(ctx *fasthttp.RequestCtx) {
	myLog.Log.Infof("Start func GetStatus")
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if ctx.IsGet() {

			auth := string(ctx.QueryArgs().Peek("auth"))
			if auth == "" {
				myLog.Log.Debugf("equql reqeust")
			} else {

				phone, err := hb.srv.ValidateToken(auth)
				if err != nil {
					myLog.Log.Errorf("Invalid token: ", err)
				} else {
					err := hb.srv.FindPhoneUser(phone)
					if err != nil {
						ctx.SetStatusCode(fasthttp.StatusSeeOther)
						// на фронте этот стутус код на редирект на регистрацию
					} else {
						orderUUID := string(ctx.QueryArgs().Peek("order_uid"))
						if orderUUID == "" {
							myLog.Log.Debugf("equql reqeust")
						} else {
							myLog.Log.Debugf("func Get with id %+v", orderUUID)
							status, time, err := hb.srv.GetStatusSrv(orderUUID)
							if err != nil {
								if err == myErrors.ErrNotFoundOrder {
									ctx.SetStatusCode(fasthttp.StatusNotFound)
								} else {
									ctx.SetStatusCode(fasthttp.StatusInternalServerError)
								}
							} else {
								myLog.Log.Debugf("sucsess get status order: %+v", orderUUID)

								WriteJsonStatus(ctx, status, time)
							}
						}
					}
				}
			}
		} else {
			WriteJson(ctx, my_errors.ErrMethodNotAllowed.Error())
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Debugf("MethodNotAllowed")
		}
	}, "GetStatus")
}

func (hb *HandlersBuilder) Get() func(ctx *fasthttp.RequestCtx) {
	myLog.Log.Infof("Start func Get")
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if ctx.IsGet() {
			auth := string(ctx.QueryArgs().Peek("auth"))
			if auth == "" {
				myLog.Log.Debugf("equql reqeust")
			} else {

				phone, err := hb.srv.ValidateToken(auth)
				if err != nil {
					myLog.Log.Errorf("Invalid token: ", err)
				} else {
					err := hb.srv.FindPhoneUser(phone)
					if err != nil {
						ctx.SetStatusCode(fasthttp.StatusSeeOther)
						// на фронте этот стутус код на редирект на регистрацию
					} else {
						orderUUID := string(ctx.QueryArgs().Peek("order_uid"))
						if orderUUID == "" {
							myLog.Log.Debugf("equql reqeust")
						} else {
							myLog.Log.Debugf("func Get with id %+v", orderUUID)
							order, err := hb.srv.GetOrderSrv(orderUUID)
							if err != nil {
								if err == myErrors.ErrNotFoundOrder {
									ctx.SetStatusCode(fasthttp.StatusNotFound)
								} else {
									ctx.SetStatusCode(fasthttp.StatusInternalServerError)
								}
							} else {
								myLog.Log.Debugf("sucsess get: %+v", orderUUID)

								WriteJsonOrder(ctx, order)
							}
						}
					}
				}
			}
		} else {
			WriteJson(ctx, my_errors.ErrMethodNotAllowed.Error())
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Debugf("MethodNotAllowed")
		}
	}, "Get")
}
