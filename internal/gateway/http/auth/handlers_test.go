package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bogatyr285/auth-go/internal/auth/entity"
	"github.com/bogatyr285/auth-go/internal/auth/mocks"
	"github.com/bogatyr285/auth-go/internal/auth/usecase"
	"github.com/bogatyr285/auth-go/internal/buildinfo"
	"github.com/bogatyr285/auth-go/internal/gateway/http/gen"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)
	mockCryptoPassword := mock.NewMockCryptoPassword(ctrl)
	mockJWTManager := mock.NewMockJWTManager(ctrl)

	tests := []struct {
		name               string
		requestBody        interface{}
		expectedResponse   gen.PostLoginResponseObject
		setupMocks         func()
		expectedError      error
		mockPostLogin      func(ctx context.Context, request gen.PostLoginRequestObject) (gen.PostLoginResponseObject, error)
		expectedStatusCode int
	}{
		{
			setupMocks: func() {
				mockUserRepo.EXPECT().
					FindUserByEmail(gomock.Any(), "user1").
					Return(entity.UserAccount{ID: 1, Username: "user1", Password: "hashedpassword"}, nil)

				mockUserRepo.EXPECT().
					ExistsTokenByUserID(gomock.Any(), 1).
					Return("refreshToken", nil)

				mockCryptoPassword.EXPECT().
					ComparePasswords("hashedpassword", "hashedpassword").
					Return(true)

				mockJWTManager.EXPECT().
					IssueToken("user1").
					Return("validtoken", nil)
			},
			name: "Success case",
			requestBody: gen.PostLoginJSONRequestBody{
				Username: "user1",
				Password: "hashedpassword",
			},
			expectedResponse: gen.PostLogin200JSONResponse{
				AccessToken:  "token",
				RefreshToken: "refreshToken",
			},
			expectedError: nil,
			mockPostLogin: func(ctx context.Context, request gen.PostLoginRequestObject) (gen.PostLoginResponseObject, error) {
				return gen.PostLogin200JSONResponse{
					AccessToken:  "token",
					RefreshToken: "refreshToken",
				}, nil
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			useCase := usecase.NewUseCase(
				mockUserRepo,
				mockCryptoPassword,
				mockJWTManager,
				buildinfo.New(),
			)

			router := chi.NewRouter()

			// Создаем тестовый сервер
			ts := httptest.NewServer(gen.HandlerFromMux(gen.NewStrictHandler(useCase, nil), router))
			defer ts.Close()

			var bodyBytes []byte
			if tt.requestBody != nil {
				bodyBytes, _ = json.Marshal(tt.requestBody)
			}

			req, err := http.NewRequest(http.MethodPost, ts.URL+"/login", bytes.NewReader(bodyBytes))
			tt.setupMocks()

			client := &http.Client{}
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to do request: %v", err)
			}
			defer resp.Body.Close()

			// Проверяем результат
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}

}
