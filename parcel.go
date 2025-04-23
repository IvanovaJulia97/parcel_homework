package main

import (
	"database/sql"
	"errors"
)

type SQLParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return &SQLParcelStore{db: db}
}

func (s *SQLParcelStore) Add(p Parcel) (int, error) {
	// добавляем строку в таблицу parcel
	// Данные Parcel:
	// Number    int
	// Client    int
	// Status    string
	// Address   string
	// CreatedAt string
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)",
		p.Client, p.Status, p.Address, p.CreatedAt)

	if err != nil {
		return 0, err
	}

	number, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return int(number), nil
}

func (s *SQLParcelStore) Get(number int) (Parcel, error) {
	// реализуем чтение строки по заданному number

	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = ?", number)

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s *SQLParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуем чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// заполняем срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		result := Parcel{}

		err := rows.Scan(&result.Number, &result.Client, &result.Status, &result.Address, &result.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, result)
	}

	return res, nil
}

func (s *SQLParcelStore) SetStatus(number int, status string) error {
	// обновляем статус в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLParcelStore) SetAddress(number int, address string) error {
	// обновляем адрес в таблице parcel
	// менять адрес можно только если статус зарегистрирован
	var statusRegistr string

	err := s.db.QueryRow("SELECT status FROM parcel WHERE number = ?", number).Scan(&statusRegistr)
	if err != nil {
		return err
	}

	if statusRegistr != ParcelStatusRegistered {
		return errors.New("у посылки незарегистрирован статус, нельзя обновить адрес")
	}

	_, err = s.db.Exec("UPDATE parcel SET address = ? WHERE number = ?", address, number)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLParcelStore) Delete(number int) error {
	// удаляем строку из таблицы parcel
	// удалять строку можно только если статус зарегистрирован
	var statusRegistr string

	err := s.db.QueryRow("SELECT status FROM parcel WHERE number = ?", number).Scan(&statusRegistr)
	if err != nil {
		return err
	}

	// if statusRegistr != ParcelStatusRegistered {
	//  return errors.New("у посылки незарегистрирован статус, нельзя удалить строку")
	// }

	_, err = s.db.Exec("DELETE FROM parcel WHERE number = ?", number)
	if err != nil {
		return err
	}

	return nil
}
