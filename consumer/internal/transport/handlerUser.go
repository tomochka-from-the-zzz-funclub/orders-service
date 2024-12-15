package transport

import (
	myErrors "consumer/internal/errors"
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
				WriteJsonErr(ctx, ErrorResponse{Message: "Error parse info user", Code: http.StatusBadRequest})
			} else {
				err := hb.srv.Registration(user)
				if err != nil {
					myLog.Log.Errorf("Registration", err.Error())
					WriteJsonErr(ctx, ErrorResponse{Message: "Error registation user", Code: http.StatusInternalServerError})
					//ctx.SetStatusCode(http.StatusBadRequest)
				} else {
					myLog.Log.Debugf("Sucses Registration")
				}
			}
			///сделвать редирект на вход
		} else {
			WriteJsonErr(ctx, ErrorResponse{Message: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
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
				// WriteJson(ctx, "Error parse")
				WriteJsonErr(ctx, ErrorResponse{Message: "Error parse info user", Code: http.StatusBadRequest})
			} else {
				err := hb.srv.CheckRegistration(phone)
				if err != nil {
					if err == myErrors.ErrNotFoundUser {
						myLog.Log.Errorf("Error Login", err.Error())
						WriteJsonErr(ctx, ErrorResponse{Message: "Error Login: not found user", Code: http.StatusNotFound})
						//ctx.SetStatusCode(fasthttp.StatusNotFound)
					} else {
						//WriteJson(ctx, "Error login")
						//ctx.SetStatusCode(fasthttp.StatusInternalServerError)
						WriteJsonErr(ctx, ErrorResponse{Message: "Error Server: login", Code: http.StatusInternalServerError})
					}
				} else {
					myLog.Log.Debugf("Sucses Login")
					token, err := hb.srv.GenerateRandomToken(phone)
					if err != nil {
						myLog.Log.Debugf("Error generate token: %+v", err)
						WriteJsonErr(ctx, ErrorResponse{Message: "Error Server: generate token", Code: http.StatusInternalServerError})
					} else {
						myLog.Log.Debugf("Token: %+v", token)
					}
					//WriteJson(ctx, "Sucses Login")
					WriteJsonToken(ctx, token)
					//сделать редирект на домашнюю стрницу
				}
			}
		} else {
			//ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Warnf("message from func Login %v", myErrors.ErrMethodNotAllowed.Error())
			WriteJsonErr(ctx, ErrorResponse{Message: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
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
				WriteJsonErr(ctx, ErrorResponse{Message: "Error parse info user", Code: http.StatusSeeOther})
			} else {

				phone, err := hb.srv.ValidateToken(auth)
				if err != nil {
					myLog.Log.Errorf("Invalid token: ", err)
					WriteJsonErr(ctx, ErrorResponse{Message: "Invalid token", Code: http.StatusSeeOther})
				} else {
					err := hb.srv.FindPhoneUser(phone)
					if err != nil {
						WriteJsonErr(ctx, ErrorResponse{Message: "Not found user", Code: http.StatusSeeOther})
						//ctx.SetStatusCode(fasthttp.StatusSeeOther)
						// на фронте этот стутус код на редирект на регистрацию
					} else {
						orderUUID := string(ctx.QueryArgs().Peek("order_uid"))
						if orderUUID == "" {
							myLog.Log.Debugf("equql reqeust")
							WriteJsonErr(ctx, ErrorResponse{Message: "Equal request", Code: http.StatusBadRequest})
						} else {
							myLog.Log.Debugf("func Get with id %+v", orderUUID)
							status, time, err := hb.srv.GetStatusSrv(orderUUID)
							if err != nil {
								if err == myErrors.ErrNotFoundOrder {
									WriteJsonErr(ctx, ErrorResponse{Message: "Not Found order", Code: http.StatusNotFound})
									//ctx.SetStatusCode(fasthttp.StatusNotFound)
								} else {
									WriteJsonErr(ctx, ErrorResponse{Message: "Internal server error", Code: http.StatusInternalServerError})
									//	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
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
			//WriteJson(ctx, my_errors.ErrMethodNotAllowed.Error())
			//ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			WriteJsonErr(ctx, ErrorResponse{Message: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
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
				WriteJsonErr(ctx, ErrorResponse{Message: "Equal request", Code: http.StatusSeeOther})
			} else {

				phone, err := hb.srv.ValidateToken(auth)
				if err != nil {
					myLog.Log.Errorf("Invalid token: ", err)
					WriteJsonErr(ctx, ErrorResponse{Message: "Invalid token", Code: http.StatusSeeOther})
				} else {
					err := hb.srv.FindPhoneUser(phone)
					if err != nil {
						//ctx.SetStatusCode(fasthttp.StatusSeeOther)
						WriteJsonErr(ctx, ErrorResponse{Message: "Not found user", Code: http.StatusSeeOther})
						// на фронте этот стутус код на редирект на регистрацию
					} else {
						orderUUID := string(ctx.QueryArgs().Peek("order_uid"))
						if orderUUID == "" {
							myLog.Log.Debugf("equql reqeust")
							WriteJsonErr(ctx, ErrorResponse{Message: "Equql reqeust", Code: http.StatusBadRequest})
						} else {
							myLog.Log.Debugf("func Get with id %+v", orderUUID)
							order, err := hb.srv.GetOrderSrv(orderUUID)
							if err != nil {
								if err == myErrors.ErrNotFoundOrder {
									WriteJsonErr(ctx, ErrorResponse{Message: "Not found order", Code: http.StatusNotFound})
									//ctx.SetStatusCode(fasthttp.StatusNotFound)
								} else {
									WriteJsonErr(ctx, ErrorResponse{Message: "Internal server error", Code: http.StatusInternalServerError})
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
			WriteJsonErr(ctx, ErrorResponse{Message: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Debugf("MethodNotAllowed")
		}
	}, "Get")
}
