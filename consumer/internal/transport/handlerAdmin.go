package transport

import (
	"fmt"

	myErrors "consumer/internal/errors"
	my_errors "consumer/internal/errors"
	myLog "consumer/internal/logger"

	"net/http"
	"strconv"

	"github.com/valyala/fasthttp"
)

func (hb *HandlersBuilder) GetHtmlLoginAdmin() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func GetHtmlLoginAdmin")
		if ctx.IsGet() {

			err2 := hb.templAdminLogin.Execute(ctx.Response.BodyWriter(), nil)

			if err2 != nil {
				myLog.Log.Errorf("GetHtmlLoginAdmin error during executing of file login: %v", err2)
				ctx.Response.SetStatusCode(400)
				return
			}
			ctx.Response.Header.Add("content-type", "text/html")
		}
	}, "GetHtmlLoginAdmin")
}

func (hb *HandlersBuilder) GetHtmlWorkAdmin() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func GetHtmlWorkAdmin")
		if ctx.IsGet() {

			err2 := hb.templAdminAllMethods.Execute(ctx.Response.BodyWriter(), nil)

			if err2 != nil {
				myLog.Log.Errorf("GetHtmlWorkAdmin error during executing of file login: %v", err2)
				ctx.Response.SetStatusCode(400)
				return
			}
			ctx.Response.Header.Add("content-type", "text/html")
		}
	}, "GetHtmlWorkAdmin")
}

func (hb *HandlersBuilder) UpdateStatus() func(ctx *fasthttp.RequestCtx) {
	myLog.Log.Infof("Start func UpdateStatus")
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if ctx.IsPut() {
			auth := string(ctx.QueryArgs().Peek("auth"))
			if auth == "" {
				myLog.Log.Debugf("equql reqeust")
			} else {

				id, err := hb.srv.ValidateTokenAdmin(auth)
				if err != nil {
					myLog.Log.Errorf("Invalid token: ", err)
				} else {
					err := hb.srv.CheckAdmin(id)
					if err != nil {
						ctx.SetStatusCode(fasthttp.StatusSeeOther)
						// на фронте этот стутус код на редирект на регистрацию
					} else {
						orderID := string(ctx.QueryArgs().Peek("order_id"))
						if orderID == "" {
							myLog.Log.Debugf("equql reqeust: order id")
						} else {
							order_id, err := strconv.Atoi(orderID)
							if err != nil {
								myLog.Log.Errorf("Invalid order id: %+v", err.Error())
							} else {
								status := string(ctx.QueryArgs().Peek("status"))
								if status == "" {
									myLog.Log.Debugf("equql reqeust: new status")
								} else {
									if (status != "create") && (status != "assembly") && (status != "delivery") {
										myLog.Log.Errorf("invalid status")
									} else {
										myLog.Log.Debugf("func Get with id %+v", orderID)
										err = hb.srv.UpdateStatusSrv(order_id, status)
										if err != nil {
											if err == myErrors.ErrNotFoundOrder {
												ctx.SetStatusCode(fasthttp.StatusNotFound)
											} else {
												ctx.SetStatusCode(fasthttp.StatusInternalServerError)
											}
										} else {
											myLog.Log.Debugf("sucsess update status: %+v", order_id)

										}
									}
								}

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
	}, "UpdateStatus")
}

func (hb *HandlersBuilder) GiveOrderToDeliveryMan() func(ctx *fasthttp.RequestCtx) {
	myLog.Log.Infof("Start func GiveOrderToDeliveryMan")
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if ctx.IsPut() {
			auth := string(ctx.QueryArgs().Peek("auth"))
			if auth == "" {
				myLog.Log.Debugf("equql reqeust")
			} else {

				id, err := hb.srv.ValidateTokenAdmin(auth)
				if err != nil {
					myLog.Log.Errorf("Invalid token: ", err)
				} else {
					err := hb.srv.CheckAdmin(id)
					if err != nil {
						ctx.SetStatusCode(fasthttp.StatusSeeOther)
						// на фронте этот стутус код на редирект на регистрацию
					} else {
						orderID := string(ctx.QueryArgs().Peek("order_id"))
						if orderID == "" {
							myLog.Log.Debugf("equql reqeust: order id")
						} else {
							order_id, err := strconv.Atoi(orderID)
							if err != nil {
								myLog.Log.Errorf("Invalid order id: %+v", err.Error())
							} else {
								delivery_man_id_ := string(ctx.QueryArgs().Peek("delivery_man_id"))
								if delivery_man_id_ == "" {
									myLog.Log.Debugf("equql reqeust: new status")
								} else {
									delivery_man_id, err := strconv.Atoi(delivery_man_id_)
									if err != nil {
										myLog.Log.Errorf("Invalid order id: %+v", err.Error())
									} else {
										myLog.Log.Debugf("func GiveOrderToDeliveryMan with id %+v", orderID)

										err = hb.srv.GiveOrderDelivery(order_id, delivery_man_id) //
										if err != nil {
											if err == myErrors.ErrNotFoundOrder {
												ctx.SetStatusCode(fasthttp.StatusNotFound)
											} else {
												ctx.SetStatusCode(fasthttp.StatusInternalServerError)
											}
										} else {
											myLog.Log.Debugf("sucsess update status: %+v", order_id)

										}
									}

								}
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
	}, "GiveOrderToDeliveryMan")
}

func (hb *HandlersBuilder) CreateDeliveryMan() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func CreateDeliveryMan: %+v", string(ctx.Request.Body()))
		if ctx.IsPost() {
			fmt.Println(string(ctx.Request.Header.Method()))
			delivery_man, err := ParseJsonDEliveryMan(ctx)
			if err != nil {
				myLog.Log.Errorf("err: %v", err.Error())
				WriteJson(ctx, "Error parse")
			} else {
				id, err := hb.srv.CreateDeliveryMan(delivery_man)
				if err != nil {
					myLog.Log.Errorf("AddDeliveryMan", err.Error())
				} else {
					myLog.Log.Debugf("Sucses AddDeliveryMan with id: %+v", id)
				}
			}

		} else {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Warnf("message from func CreateDeliveryMan %v", myErrors.ErrMethodNotAllowed.Error())
		}
	}, "CreateDeliveryMan")
}

// func (hb *HandlersBuilder) CheckDeliveryStart() func(ctx *fasthttp.RequestCtx) {
// 	return metrics(func(ctx *fasthttp.RequestCtx) {
// 		myLog.Log.Debugf("Start func CheckDeliveryStart: %+v", string(ctx.Request.Body()))
// 		if ctx.IsPut() {
// 			delivery_man_id_ := string(ctx.QueryArgs().Peek("delivery_man_id"))
// 			if delivery_man_id_ == "" {
// 				myLog.Log.Debugf("equql reqeust: new status")
// 			} else {
// 				delivery_man_id, err := strconv.Atoi(delivery_man_id_)
// 				if err != nil {
// 					myLog.Log.Errorf("Invalid order id: %+v", err.Error())
// 				} else {
// 					result, err := hb.srv.CheckDeliveryStart(delivery_man_id)
// 					if err != nil {
// 						if result {
// 							myLog.Log.Debugf("Start delivery")
// 						} else {
// 							myLog.Log.Debugf("Not ready for delivery")
// 						}
// 					} else {
// 						myLog.Log.Errorf("Error CheckDeliveryStart: %+v", err.Error())
// 					}
// 				}
// 			}
// 		} else {
// 			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
// 			myLog.Log.Warnf("message from func CheckDeliveryStart %v", myErrors.ErrMethodNotAllowed.Error())
// 		}
// 	}, "CheckDeliveryStart")
// }

func (hb *HandlersBuilder) LoginAdmin() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func Login: %+v", string(ctx.Request.Body()))
		if ctx.IsPut() {
			id_ := string(ctx.QueryArgs().Peek("id"))
			if id_ == "" {
				myLog.Log.Debugf("equql reqeust")
				WriteJson(ctx, "equql reqeust")
				ctx.SetStatusCode(http.StatusBadRequest)
			} else {
				id, err := strconv.Atoi(id_)
				if err != nil {
					myLog.Log.Debugf("bad reqeust")
					WriteJson(ctx, "bad reqeust")
					ctx.SetStatusCode(http.StatusBadRequest)
				} else {
					err := hb.srv.CheckAdmin(id)
					if err != nil {
						if err == myErrors.ErrNotFoundUser {
							myLog.Log.Errorf("Error Login", err.Error())
							ctx.SetStatusCode(fasthttp.StatusNotFound)
						} else {
							WriteJson(ctx, "Error login")
							ctx.SetStatusCode(fasthttp.StatusInternalServerError)
						}
					} else {
						myLog.Log.Debugf("Sucses Login")
						token, err := hb.srv.GenerateAdminToken(id)
						if err != nil {
							myLog.Log.Debugf("Error generate token: %+v", err)
						} else {
							myLog.Log.Debugf("Token: %+v", token)
						}
						//WriteJson(ctx, "Sucses Login")
						WriteJsonToken(ctx, token)
					} //сделать редирект на домашнюю стрницу
				}
			}

		} else {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Warnf("message from func Login %v", myErrors.ErrMethodNotAllowed.Error())
		}
	}, "Login")
}
