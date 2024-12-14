package service

import (
	myLog "consumer/internal/logger"
	"consumer/internal/models"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

var TokensUser = make(map[string]string)

func (s *Srv) Registration(user models.User) error {
	myLog.Log.Debugf("Registration SRV")
	err := s.db.Registration(user)
	if err != nil {
		myLog.Log.Errorf("Error registration")

	}
	return err

}

func (s *Srv) CheckRegistration(phone string) error {
	myLog.Log.Debugf("CheckRegistration SRV")
	err := s.db.CheckRegistration(phone)
	if err != nil {
		myLog.Log.Errorf("Error Login")
	}
	return err
}

func (s *Srv) GenerateRandomToken(phone string) (string, error) {
	// Создаем массив байтов заданной длины
	bytes := make([]byte, 20)

	// Генерируем случайные байты
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Кодируем байты в строку Base64
	token := base64.URLEncoding.EncodeToString(bytes)
	TokensUser[phone] = token
	return token, nil
}

// Функция для проверки токена
func (s *Srv) ValidateToken(tokenString string) (string, error) {
	// Проверяем, существует ли токен в карте
	for phone, token := range TokensUser {
		if token == tokenString {

			// Если токен действителен, возвращаем номер телефона
			return phone, nil
		}
	}
	return "", fmt.Errorf("токен не найден в карте")
}
