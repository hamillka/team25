package integration_test

import (
	"net/http"
	"testing"

	net "github.com/hamillka/team25/backend/internal/integration_test"
	"github.com/stretchr/testify/require"
)

func TestUserRegister(t *testing.T) {
	client := http.Client{}

	tests := []struct { //nolint:govet // no other ways
		respBody       interface{}
		headers        [][2]string
		Name           string
		Method         string
		route          string
		body           string
		expectedStatus int
		parseResp      bool
	}{
		{
			Name:   "user register success",
			Method: http.MethodPost,
			route:  "/auth/register",
			body: "{\"fio\": \"testfio1\",\"phoneNumber\": \"testnumber1\", " +
				"\"email\": \"testemail1\",\"insurance\": \"insurance1\", \"login\": \"login1231\", " +
				"\"password\": \"pass1231\", \"role\": 1}",
			headers:        [][2]string{{"Content-Type", "application/json"}},
			expectedStatus: http.StatusCreated,
			parseResp:      false,
			respBody:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req, err := net.CreateRequest(tt.Method, tt.route, tt.body, tt.headers)
			require.NoError(t, err)

			net.SendRequest(t, &client, req, tt.expectedStatus, tt.parseResp, tt.respBody)
		})
	}
}
