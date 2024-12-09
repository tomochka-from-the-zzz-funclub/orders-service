package transport

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"consumer/internal/config"
	myErrors "consumer/internal/errors"
	my_errors "consumer/internal/errors"
	myLog "consumer/internal/logger"
	"consumer/internal/service"

	"net/http"
	"strconv"

	"github.com/fasthttp/router"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
)

type HandlersBuilder struct {
	srv   service.InterfaceService
	rout  *router.Router
	templ *template.Template
}

func HandleCreate(cfg config.Config, s service.InterfaceService) {
	fmt.Println(os.Getwd())
	fmt.Println(os.ReadDir("./src"))
	t, err := os.Getwd()
	absolutePath, err := filepath.Abs(t)
	if err != nil {
		fmt.Println("Ошибка при получении абсолютного пути:", err)
		return
	}

	fmt.Println("Абсолютный путь к директории:", absolutePath)
	tmpl, err := template.ParseFiles("../app/fc.html")
	if err != nil {
		myLog.Log.Fatalf("GetHtml error during parsing of file: %v", err)
		return
	}

	hb := HandlersBuilder{
		srv:   s,
		rout:  router.New(),
		templ: tmpl,
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8090", nil)
	}()

	hb.rout.GET("/api/v1/get", hb.Get())
	hb.rout.PUT("/api/v1/update/status", hb.UpdateStatus())
	hb.rout.POST("/api/v1/adition/delivery_man", hb.CreateDeliveryMan())
	hb.rout.PUT("/api/v1/create/delivery", hb.GiveOrderToDeliveryMan())
	hb.rout.GET("/api/v1", hb.GetHtml())
	fmt.Println(fasthttp.ListenAndServe(":8080", hb.rout.Handler))
}

func (hb *HandlersBuilder) GetHtml() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func GetHtml")
		if ctx.IsGet() {
			err := hb.templ.Execute(ctx.Response.BodyWriter(), nil)
			if err != nil {
				myLog.Log.Errorf("GetHtml error during executing of file: %v", err)
				ctx.Response.SetStatusCode(400)
				return
			}

			ctx.Response.Header.Add("content-type", "text/html")
		}
	}, "GetHtml")
}

func (hb *HandlersBuilder) Get() func(ctx *fasthttp.RequestCtx) {
	myLog.Log.Infof("Start func Get")
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if ctx.IsGet() {
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
					myLog.Log.Debugf("sucsess get: %+v", order.Id)

					WriteJsonOrder(ctx, order)
				}
			}
		} else {
			WriteJson(ctx, my_errors.ErrMethodNotAllowed.Error())
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Debugf("MethodNotAllowed")
		}
	}, "Get")
}

func (hb *HandlersBuilder) UpdateStatus() func(ctx *fasthttp.RequestCtx) {
	myLog.Log.Infof("Start func UpdateStatus")
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if ctx.IsPut() {
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

func (hb *HandlersBuilder) CheckDeliveryStart() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func CheckDeliveryStart: %+v", string(ctx.Request.Body()))
		if ctx.IsPost() {
			delivery_man_id_ := string(ctx.QueryArgs().Peek("delivery_man_id"))
			if delivery_man_id_ == "" {
				myLog.Log.Debugf("equql reqeust: new status")
			} else {
				delivery_man_id, err := strconv.Atoi(delivery_man_id_)
				if err != nil {
					myLog.Log.Errorf("Invalid order id: %+v", err.Error())
				} else {
					result, err := hb.srv.CheckDeliveryStart(delivery_man_id)
					if err != nil {
						if result {
							myLog.Log.Debugf("Start delivery")
						} else {
							myLog.Log.Debugf("Not ready for delivery")
						}
					} else {
						myLog.Log.Errorf("Error CheckDeliveryStart: %+v", err.Error())
					}
				}
			}
		} else {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			myLog.Log.Warnf("message from func CheckDeliveryStart %v", myErrors.ErrMethodNotAllowed.Error())
		}
	}, "CheckDeliveryStart")
}
