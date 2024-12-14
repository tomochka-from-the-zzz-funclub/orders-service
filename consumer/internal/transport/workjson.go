package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	myErrors "consumer/internal/errors"
	myLog "consumer/internal/logger"
	"consumer/internal/models"

	"github.com/valyala/fasthttp"
)

func ParseJsonOrder(ctx *fasthttp.RequestCtx) (models.Order, error) {
	var order models.Order
	err := json.NewDecoder(bytes.NewReader(ctx.Request.Body())).Decode(&order)
	if err != nil {
		return models.Order{}, myErrors.ErrParseJSON
	}

	_, err = time.Parse("2006/01/02", order.DateCreated)
	if err != nil {
		return order, nil
	}

	return models.Order{}, myErrors.ErrParseJSON
}

func ParseJsonDEliveryMan(ctx *fasthttp.RequestCtx) (models.DeliveryMan, error) {
	var del models.DeliveryMan
	err := json.NewDecoder(bytes.NewReader(ctx.Request.Body())).Decode(&del)
	if err != nil {
		return models.DeliveryMan{}, myErrors.ErrParseJSON
	}

	return del, nil
}

func ParseJsonUser(ctx *fasthttp.RequestCtx) (models.User, error) {
	var user models.User
	err := json.NewDecoder(bytes.NewReader(ctx.Request.Body())).Decode(&user)
	if err != nil {
		return models.User{}, myErrors.ErrParseJSON
	}

	return user, nil
}

func ParseJsonLogin(ctx *fasthttp.RequestCtx) (string, error) {
	var user struct {
		Phone string `json:"phone"` // Измените на заглавную букву, чтобы поле было экспортируемым
	}

	// Декодируем JSON из тела запроса
	err := json.NewDecoder(bytes.NewReader(ctx.Request.Body())).Decode(&user)
	if err != nil {
		return "", myErrors.ErrParseJSON // Возвращаем ошибку в случае неудачи
	}

	// Выводим номер телефона для отладки
	fmt.Println("Parsed phone: ", user.Phone)

	return user.Phone, nil // Возвращаем номер телефона
}

func WriteJson(ctx *fasthttp.RequestCtx, s string) error {
	ctx.SetContentType("application/json")
	ctx.Response.BodyWriter()
	err := json.NewEncoder((*ctx).Response.BodyWriter()).Encode(s)
	if err != nil {
		return err
	}
	return nil
}

func WriteJsonOrder(ctx *fasthttp.RequestCtx, order models.Order) error {
	ctx.SetContentType("application/json")
	ctx.Response.BodyWriter()
	err := json.NewEncoder((*ctx).Response.BodyWriter()).Encode(order)
	if err != nil {
		return err
	}
	return nil
}

func WriteJsonStatus(ctx *fasthttp.RequestCtx, stat string, create_time string) error {
	ctx.SetContentType("application/json")
	var status struct {
		Status     string `json:"status"`
		Updated_at string `json:"updated_at"`
	}
	status.Status = stat
	status.Updated_at = create_time
	ctx.Response.BodyWriter()
	err := json.NewEncoder((*ctx).Response.BodyWriter()).Encode(status)
	if err != nil {
		return err
	}
	return nil
}

func WriteJsonToken(ctx *fasthttp.RequestCtx, token string) error {
	ctx.SetContentType("application/json")
	var tok struct {
		Token string `json:"token"` // Изменено с маленькой буквы на большую
	}
	tok.Token = token // Устанавливаем значение токена
	// Убедимся, что мы используем правильный метод для записи в тело ответа
	err := json.NewEncoder(ctx).Encode(tok) // Используем ctx напрямую
	if err != nil {
		myLog.Log.Errorf("Error writing JSON: %v", err) // Добавил вывод ошибки
		return err
	}
	return nil
}
