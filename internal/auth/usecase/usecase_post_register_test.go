package usecase

import (
	"context"
	"errors"
	"github.com/bogatyr285/auth-go/internal/auth/entity"
	m "github.com/bogatyr285/auth-go/internal/auth/mocks"
	"github.com/bogatyr285/auth-go/internal/gateway/http/gen"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestPostRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := m.NewMockUserRepository(ctrl)
	mockCryptoPassword := m.NewMockCryptoPassword(ctrl)

	tests := []struct {
		name             string
		inputRequest     gen.PostRegisterRequestObject
		mockHashPassword func()
		mockRegisterUser func()
		expectedResponse gen.PostRegisterResponseObject
		expectedError    error
	}{
		{
			name: "Successful registration",
			inputRequest: gen.PostRegisterRequestObject{
				Body: &gen.PostRegisterJSONRequestBody{
					Username: "testuser",
					Password: "testpassword",
				},
			},
			mockHashPassword: func() {
				mockCryptoPassword.EXPECT().
					HashPassword("testpassword").
					Return([]byte("hashedpassword"), nil)
			},
			mockRegisterUser: func() {
				mockUserRepo.EXPECT().
					RegisterUser(gomock.Any(), entity.UserAccount{
						Username: "testuser",
						Password: "hashedpassword",
					}).
					Return(nil)
			},
			expectedResponse: gen.PostRegister201JSONResponse{
				Username: "testuser",
			},
			expectedError: nil,
		},
		{
			name: "Error in registering user",
			inputRequest: gen.PostRegisterRequestObject{
				Body: &gen.PostRegisterJSONRequestBody{
					Username: "testuser",
					Password: "testpassword",
				},
			},
			mockHashPassword: func() {
				mockCryptoPassword.EXPECT().
					HashPassword("testpassword").
					Return([]byte("hashedpassword"), nil)
			},
			mockRegisterUser: func() {
				mockUserRepo.EXPECT().
					RegisterUser(gomock.Any(), entity.UserAccount{
						Username: "testuser",
						Password: "hashedpassword",
					}).
					Return(errors.New("register error"))
			},
			expectedResponse: gen.PostRegister500JSONResponse{},
			expectedError:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Настройка моков
			tt.mockHashPassword()
			tt.mockRegisterUser()

			// Создаем экземпляр AuthUseCase с моками
			authUseCase := AuthUseCase{
				cp: mockCryptoPassword,
				ur: mockUserRepo,
			}

			// Выполняем тестируемую функцию
			resp, err := authUseCase.PostRegister(context.Background(), tt.inputRequest)

			// Проверяем результат
			assert.Equal(t, tt.expectedResponse, resp)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
