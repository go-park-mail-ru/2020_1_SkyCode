package repository

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserRepository_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	testUser := &models.User{
		ID:        uint64(1),
		Email:     "testemail@m.ru",
		Password:  "1234",
		FirstName: "A",
		LastName:  "B",
		Phone:     "79865433211",
		Avatar:    "default.jpg",
		Role:      "User",
	}

	rows := sqlmock.NewRows([]string{"firstName", "lastName", "email", "password", "avatar", "phone", "role"}).
		AddRow(testUser.FirstName, testUser.LastName, testUser.Email, testUser.Password,
			testUser.Avatar, testUser.Phone, testUser.Role)

	mock.ExpectQuery("SELECT").
		WithArgs(testUser.ID).
		WillReturnRows(rows)

	resultUser := &models.User{
		ID: testUser.ID,
	}

	err = repo.GetById(resultUser)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, testUser, resultUser)
}

func TestUserRepository_GetByPhone(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	testUser := &models.User{
		ID:        uint64(1),
		Email:     "testemail@m.ru",
		Password:  "1234",
		FirstName: "A",
		LastName:  "B",
		Phone:     "79865433211",
		Avatar:    "default.jpg",
		Role:      "User",
	}

	rows := sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password", "avatar", "role"}).
		AddRow(testUser.ID, testUser.FirstName, testUser.LastName, testUser.Email,
			testUser.Password, testUser.Avatar, testUser.Role)

	mock.ExpectQuery("SELECT").
		WithArgs(testUser.Phone).
		WillReturnRows(rows)

	resultUser := &models.User{
		Phone: testUser.Phone,
	}

	err = repo.GetByPhone(resultUser)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, testUser, resultUser)
}

func TestUserRepository_InsertInto(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	testUser := &models.User{
		Email:     "testemail@m.ru",
		Password:  "1234",
		FirstName: "A",
		LastName:  "B",
		Phone:     "79865433211",
		Avatar:    "default.jpg",
		Role:      "User",
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO").
		WithArgs(testUser.FirstName, testUser.LastName, testUser.Email, testUser.Phone,
			testUser.Password, testUser.Avatar, testUser.Role).
		WillReturnRows(rows)

	err = repo.InsertInto(testUser)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, 1, testUser.ID)
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	testUser := &models.User{
		ID:        1,
		Email:     "testemail@m.ru",
		FirstName: "A",
		LastName:  "B",
	}

	rows := sqlmock.NewRows([]string{"firstName", "lastName", "email"}).
		AddRow(testUser.FirstName, testUser.LastName, testUser.Email)

	mock.ExpectQuery("UPDATE users").
		WithArgs(testUser.ID, testUser.FirstName, testUser.LastName, testUser.Email).
		WillReturnRows(rows)

	resultUser := &models.User{
		ID:        testUser.ID,
		Email:     testUser.Email,
		FirstName: testUser.FirstName,
		LastName:  testUser.LastName,
	}

	err = repo.Update(resultUser)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, testUser, resultUser)
}

func TestUserRepository_UpdateAvatar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	testUser := &models.User{
		Email:     "testemail@m.ru",
		FirstName: "A",
		LastName:  "B",
		Phone:     "79865433211",
		Avatar:    "default.jpg",
		Role:      "User",
	}

	rows := sqlmock.NewRows([]string{"firstName", "lastName", "email", "phone", "role"}).
		AddRow(testUser.FirstName, testUser.LastName, testUser.Email,
			testUser.Phone, testUser.Role)

	mock.ExpectQuery("UPDATE users").
		WithArgs(testUser.ID, testUser.Avatar).
		WillReturnRows(rows)

	resultUser := &models.User{
		ID:     testUser.ID,
		Avatar: testUser.Avatar,
	}

	err = repo.UpdateAvatar(resultUser)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, testUser, resultUser)
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	testUser := &models.User{
		Email:     "testemail@m.ru",
		Password:  "1234",
		FirstName: "A",
		LastName:  "B",
		Avatar:    "default.jpg",
		Role:      "User",
	}

	rows := sqlmock.NewRows([]string{"firstName", "lastName", "email", "avatar", "role"}).
		AddRow(testUser.FirstName, testUser.LastName, testUser.Email,
			testUser.Avatar, testUser.Role)

	mock.ExpectQuery("UPDATE users").
		WithArgs(testUser.ID, testUser.Password).
		WillReturnRows(rows)

	resultUser := &models.User{
		ID:       testUser.ID,
		Password: testUser.Password,
	}

	err = repo.UpdatePassword(resultUser)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, testUser, resultUser)
}

func TestUserRepository_UpdatePhone(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	testUser := &models.User{
		Email:     "testemail@m.ru",
		FirstName: "A",
		LastName:  "B",
		Phone:     "79865433211",
		Avatar:    "default.jpg",
		Role:      "User",
	}

	rows := sqlmock.NewRows([]string{"firstName", "lastName", "email", "avatar", "role"}).
		AddRow(testUser.FirstName, testUser.LastName, testUser.Email,
			testUser.Avatar, testUser.Role)

	mock.ExpectQuery("UPDATE users").
		WithArgs(testUser.ID, testUser.Phone).
		WillReturnRows(rows)

	resultUser := &models.User{
		ID:    testUser.ID,
		Phone: testUser.Phone,
	}

	err = repo.UpdatePhone(resultUser)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, testUser, resultUser)
}
