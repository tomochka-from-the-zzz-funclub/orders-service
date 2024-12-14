package service

import (
	myLog "consumer/internal/logger"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

var TokensAdmin = make(map[int]string)

// func (s *Srv) Registration(user models.User) error {
// 	myLog.Log.Debugf("Registration SRV")
// 	err := s.db.Registration(user)
// 	if err != nil {
// 		myLog.Log.Errorf("Error registration")

// 	}
// 	return err

// }

func (s *Srv) CheckAdmin(id_admin int) error {
	myLog.Log.Debugf("CheckRegistration SRV")
	err := s.db.CheckAdmin(id_admin)
	if err != nil {
		myLog.Log.Errorf("Error Login")
	}
	return err
}

func (s *Srv) GenerateAdminToken(id int) (string, error) {
	// Создаем массив байтов заданной длины
	bytes := make([]byte, 20)

	// Генерируем случайные байты
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Кодируем байты в строку Base64
	token := base64.URLEncoding.EncodeToString(bytes)
	TokensAdmin[id] = token
	return token, nil
}

// Функция для проверки токена
func (s *Srv) ValidateTokenAdmin(tokenString string) (int, error) {
	// Проверяем, существует ли токен в карте
	for id, token := range TokensAdmin {
		if token == tokenString {

			// Если токен действителен, возвращаем номер телефона
			return id, nil
		}
	}
	return 0, fmt.Errorf("токен не найден в карте")
}
