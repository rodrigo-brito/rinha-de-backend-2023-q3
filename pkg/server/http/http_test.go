package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"rinha/etc/grpc"
	"rinha/testdata/mocks"
)

func TestServer_CreateHandle(t *testing.T) {
	tt := []struct {
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "empty",
			body:       `{}`,
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "invalid type",
			body:       `{"nome": 1}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "nick null",
			body:       `{"apelido" : null,"nome" : "test","nascimento" : "1985-01-23", "stack" : null}`,
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "invalid birth",
			body:       `{"apelido" : "foo","nome" : "Fulano", "nascimento" : "1234", "stack" : ["GO"]}`,
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "valid",
			body:       `{"apelido" : "foo","nome" : "Fulano", "nascimento" : "1985-01-23", "stack" : ["GO"]}`,
			statusCode: http.StatusCreated,
		},
	}

	ctrl := gomock.NewController(t)
	storage := mocks.NewMockStorage(ctrl)
	server := Server{
		storageClient: storage,
	}

	storage.EXPECT().Save(gomock.Any(), gomock.Any()).Return(&grpc.SaveResponse{
		Id: "1",
	}, nil).AnyTimes()

	s := httptest.NewServer(http.HandlerFunc(server.CreateHandle))
	defer s.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.Post(s.URL, "application/json", strings.NewReader(tc.body))
			require.NoError(t, err)
			require.Equal(t, tc.statusCode, resp.StatusCode)
		})
	}
}
