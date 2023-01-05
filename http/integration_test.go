//go:build integration
// +build integration

package http_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	domain "github.com/yugovtr/bingo/domain/game"
	"github.com/yugovtr/bingo/domain/repository"
	"github.com/yugovtr/bingo/http"
	"github.com/yugovtr/bingo/http/routes"
	"github.com/yugovtr/bingo/infra/db"
)

func TestServer_BingoWithRethinkDB(t *testing.T) {
	ctx := context.Background()
	address := StartDB(t, ctx)

	cli := db.Connect(context.TODO(), address)
	require.NotNil(t, cli)

	err := db.Migrate(cli.Session)
	require.NoError(t, err)

	rethinkDB := repository.NewRethinkDB(cli.Session)
	game := domain.NewGameWithCaller(rethinkDB, func() int { return 1 })
	routes := routes.NewBingo(routes.New(), game)

	client := AssertServer(t, http.ServerConfig{Routes: routes})

	AcceptanceBingo(t, client)
}

func StartDB(t *testing.T, ctx context.Context) string {
	t.Helper()

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image:        "rethinkdb",
			ExposedPorts: []string{"28015/tcp", "29015/tcp"},
			WaitingFor:   wait.ForListeningPort("28015/tcp"),
		},
		Started: true,
	})

	require.NoError(t, err)
	t.Cleanup(func() {
		err := container.Terminate(ctx)
		assert.NoError(t, err)
	})

	host, _ := container.Host(ctx)
	p, _ := container.MappedPort(ctx, "28015/tcp")

	return fmt.Sprintf("%s:%d", host, p.Int())
}
