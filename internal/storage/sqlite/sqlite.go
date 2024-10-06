package sqlite

import (
	"database/sql"
	"fmt"
	"time"
)

// DeleteLicenseById deletes a license by UserId
func (s *Storage) DeleteLicenseById(userId string) error {
	const op = "storage.sqlite.DeleteLicenseById"

	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`DELETE FROM UserLicense WHERE UserId = ?`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get affected rows: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: no license found for Id: %s", op, userId)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: transaction commit failed: %w", op, err)
	}

	// Log transaction
	return s.LogTransaction(fmt.Sprintf("action=delete_license user_id=%s", userId))
}

// Get all licenses
func (s *Storage) GetAllLicenses() ([]UserLicense, error) {
	const op = "storage.sqlite.GetAllLicenses"

	rows, err := s.db.Query(`SELECT id, license, UserId, createdAt, updatedAt, expiresAt, hwid, status FROM UserLicense`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var licenses []UserLicense

	for rows.Next() {
		var license UserLicense
		var hwid sql.NullString

		if err := rows.Scan(&license.ID, &license.License, &license.UserId, &license.CreatedAt, &license.UpdatedAt, &license.ExpiresAt, &hwid, &license.Status); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if hwid.Valid {
			license.HWID = &hwid.String
		} else {
			license.HWID = nil
		}

		licenses = append(licenses, license)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return licenses, nil
}

func (s *Storage) AddLicense(license, UserId, status string, hwid *string, expiresAt time.Time) (int64, error) {
	const op = "storage.sqlite.AddLicense"

	stmt, err := s.db.Prepare(`
INSERT INTO UserLicense (license, UserId, createdAt, updatedAt, expiresAt, hwid, status)
VALUES (?, ?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now()
	hwidValue := sql.NullString{Valid: false}
	if hwid != nil {
		hwidValue = sql.NullString{String: *hwid, Valid: true}
	}

	res, err := stmt.Exec(license, UserId, now, now, expiresAt, hwidValue, status)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	err = s.LogTransaction(fmt.Sprintf(
		"action=add_license license=%s user_id=%s",
		license, UserId,
	))
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// RenewLicenseById renews the license by extending its expiration date
func (s *Storage) RenewLicenseById(userId string, days int) (time.Time, error) {
	const op = "storage.sqlite.RenewLicenseById"

	expirationTime := time.Now().AddDate(0, 0, days)

	stmt, err := s.db.Prepare(`UPDATE UserLicense SET expiresAt = ?, updatedAt = ? WHERE UserId = ?`)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(expirationTime, time.Now(), userId)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return time.Time{}, fmt.Errorf("%s: no license found for UserId: %s", op, userId)
	}

	return expirationTime, s.LogTransaction(fmt.Sprintf("action=renew_license user_id=%s days=%d", userId, days))
}

// Common method to retrieve a license
func (s *Storage) getLicense(query, param, paramName string) (*UserLicense, error) {
	const op = "storage.sqlite.getLicense"

	row := s.db.QueryRow(query, param)

	var license UserLicense
	var hwid sql.NullString

	err := row.Scan(&license.ID, &license.License, &license.UserId, &license.CreatedAt, &license.UpdatedAt, &license.ExpiresAt, &hwid, &license.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: no license found for %s: %s", op, paramName, param)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if hwid.Valid {
		license.HWID = &hwid.String
	} else {
		license.HWID = nil
	}

	return &license, nil
}

// Update HWID helper
func (s *Storage) updateHwid(license, hwid, action string) error {
	const op = "storage.sqlite.updateHwid"

	stmt, err := s.db.Prepare(`UPDATE UserLicense SET hwid = ?, updatedAt = ? WHERE license = ?`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(hwid, time.Now(), license)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("%s: no license found for License: %s", op, license)
	}

	return s.LogTransaction(fmt.Sprintf("action=%s hwid=%s license=%s", action, hwid, license))
}

// Freeze/Unfreeze license helper
func (s *Storage) updateLicenseStatus(userId, status string) error {
	stmt, err := s.db.Prepare(`UPDATE UserLicense SET status = ?, updatedAt = ? WHERE UserId = ?`)
	if err != nil {
		return fmt.Errorf("storage.sqlite.updateLicenseStatus: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(status, time.Now(), userId)
	if err != nil {
		return fmt.Errorf("storage.sqlite.updateLicenseStatus: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("storage.sqlite.updateLicenseStatus: no license found for UserId: %s", userId)
	}

	return s.LogTransaction(fmt.Sprintf("action=%s_license user_id=%s", status, userId))
}

// LogTransaction logs transaction actions
func (s *Storage) LogTransaction(description string) error {
	_, err := s.db.Exec(`INSERT INTO TransactionLogs (description) VALUES (?)`, description)
	return err
}
