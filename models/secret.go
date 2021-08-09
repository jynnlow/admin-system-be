package models

import "gorm.io/gorm"

type Secret struct {
	gorm.Model
	Secret string `json:"secret"`
	Type   string `json:"type" gorm:"unique"`
}

type SecretOperations interface {
	GetAdminLoginKey() (string, error)
	GetTokenKey() (string, error)
}

type SecretOperationsImpl struct {
	DbConn *gorm.DB
}

func (s *SecretOperationsImpl) GetAdminLoginKey() (string, error) {
	secret := &Secret{}
	err := s.DbConn.Where("type = ?", "admin-login").First(secret).Error
	if err != nil {
		return "", err
	}
	return secret.Secret, nil
}

func (s *SecretOperationsImpl) GetTokenKey() (string, error) {
	secret := &Secret{}
	err := s.DbConn.Where("type = ?", "jwt-token-key").First(secret).Error
	if err != nil {
		return "", err
	}
	return secret.Secret, nil
}

//func (s *SecretOperationsImpl) Insert(sct, typ string) error {
//	secret := &Secret{
//		Secret: sct,
//		Type:   typ,
//	}
//
//	err := s.DbConn.Select("secret", "type").Create(secret).Error
//	if err != nil {
//		return err
//	}
//	return nil
//}
