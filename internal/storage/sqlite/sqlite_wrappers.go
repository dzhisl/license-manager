package sqlite

import (
	_ "github.com/mattn/go-sqlite3"
)

func (s *Storage) GetLicenseById(userId string) (*UserLicense, error) {
	return s.getLicense(`SELECT id, license, UserId, createdAt, updatedAt, expiresAt, hwid, status FROM UserLicense WHERE UserId = ?`, userId, "UserId")
}

func (s *Storage) GetLicenseByLicense(license string) (*UserLicense, error) {
	return s.getLicense(`SELECT id, license, UserId, createdAt, updatedAt, expiresAt, hwid, status FROM UserLicense WHERE license = ?`, license, "License")
}

func (s *Storage) BindHwidToLicenseByLicense(license, hwid string) error {
	return s.updateHwid(license, hwid, "bind_hwid")
}

func (s *Storage) UnbindHwidFromLicense(license string) error {
	return s.updateHwid(license, "", "unbind_hwid")
}

func (s *Storage) FreezeLicenseById(userId string) error {
	return s.updateLicenseStatus(userId, "frozen")
}

func (s *Storage) UnfreezeLicenseById(userId string) error {
	return s.updateLicenseStatus(userId, "active")
}
