package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, number)
	// get

	parcels, err := store.Get(number)
	require.NoError(t, err)

	parcel.Number = number

	require.Equal(t, parcel, parcels)

	// require.Equal(t, parcel.Client, parcels.Client)
	// require.Equal(t, parcel.Status, parcels.Status)
	// require.Equal(t, parcel.Address, parcels.Address)
	// require.Equal(t, parcel.CreatedAt, parcels.CreatedAt)

	// delete
	err = store.Delete(number)
	require.NoError(t, err)

	_, err = store.Get(number)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД,
	// убедитесь в отсутствии ошибки и наличии идентификатора
	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, number)

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь,
	// что адрес обновился
	upd, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, newAddress, upd.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД,
	// убедитесь в отсутствии ошибки и наличии идентификатора
	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, number)

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	statusNew := ParcelStatusSent
	err = store.SetStatus(number, statusNew)
	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь,
	// что статус обновился
	newParcel, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, statusNew, newParcel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client
	// add
	for i := 0; i < len(parcels); i++ {
		// добавьте новую посылку в БД,
		// убедитесь в отсутствии ошибки и наличии идентификатора
		number, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotZero(t, number)
		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = number

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[number] = parcels[i]
	}

	// get by client
	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что
	// количество полученных посылок совпадает
	// с количеством добавленных
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels))

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
		expct, ok := parcelMap[parcel.Number]
		require.True(t, ok, "посылка под номером %d", parcel.Number)
		require.Equal(t, expct, parcel)

		// require.Equal(t, expct.Client, parcel.Client)
		// require.Equal(t, expct.Status, parcel.Status)
		// require.Equal(t, expct.Address, parcel.Address)
		// require.Equal(t, expct.CreatedAt, parcel.CreatedAt)

	}
}
