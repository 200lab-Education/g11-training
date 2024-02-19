package usecase

import (
	"context"
	"errors"
	"my-app/module/user/domain"
	"testing"
)

type mockHasher struct{}

func (mockHasher) RandomStr(length int) (string, error) {
	return "abcd", nil
}

func (mockHasher) HashPassword(salt, password string) (string, error) {
	return "asjhdjashdjashd", nil
}

type mockUserRepo struct{}

func (mockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if email == "existed@gmail.com" {
		return &domain.User{}, nil
	}

	if email == "error@gmail.com" {
		return nil, errors.New("cannot get record")
	}

	return &domain.User{}, nil
}

func (mockUserRepo) Create(ctx context.Context, data *domain.User) error {
	return nil
}

func TestUseCase_Register(t *testing.T) {
	uc := NewUseCase(mockUserRepo{}, mockHasher{})

	type testData struct {
		Input    EmailPasswordRegistrationDTO
		Expected error
	}

	table := []testData{
		{
			Input: EmailPasswordRegistrationDTO{
				FirstName: "Viet",
				LastName:  "Tran",
				Email:     "existed@gmail.com",
				Password:  "123456",
			},
			Expected: domain.ErrEmailHasExisted,
		},
		{
			Input: EmailPasswordRegistrationDTO{
				FirstName: "Viet",
				LastName:  "Tran",
				Email:     "error@gmail.com",
				Password:  "123456",
			},
			Expected: errors.New("cannot get record"),
		},
	}

	for i := range table {
		actualError := uc.Register(context.Background(), table[i].Input)

		if actualError.Error() != table[i].Expected.Error() {
			t.Errorf("Register failed. Expected %s, but actual is %s", table[i].Expected.Error(), actualError.Error())
		}
	}
}
