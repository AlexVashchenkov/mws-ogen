package test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	gen "mws-ogen/gen"
	"mws-ogen/server"
)

func setupClient(t *testing.T) gen.Client {
	t.Helper()

	h := server.NewUserStorage()
	srv, err := gen.NewServer(h)
	require.NoError(t, err)

	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)

	client, err := gen.NewClient(ts.URL, gen.WithClient(ts.Client()))
	require.NoError(t, err)

	return *client
}

func TestAddUser(t *testing.T) {
	// Arrange
	client := setupClient(t)
	ctx := context.Background()
	user := gen.User{
		Name:  "Alexey",
		Email: gen.NewOptString("litmo@litmo.ru"),
	}

	// Act
	res, err := client.AddUser(ctx, &user)

	// Assert
	require.NoError(t, err)
	created, ok := res.(*gen.User)
	require.True(t, ok)
	require.Equal(t, "Alexey", created.Name)
	require.True(t, created.Email.IsSet())
	require.Equal(t, "litmo@litmo.ru", created.Email.Value)
}

func TestListUsers(t *testing.T) {
	// Arrange
	client := setupClient(t)
	ctx := context.Background()
	_, err := client.AddUser(ctx, &gen.User{Name: "Alexey"})
	require.NoError(t, err)
	_, err = client.AddUser(ctx, &gen.User{Name: "Yulia"})
	require.NoError(t, err)

	// Act
	res, err := client.ListUsers(ctx)

	// Assert
	require.NoError(t, err)
	list, ok := res.(*gen.ListUsersOKApplicationJSON)
	require.True(t, ok)
	require.GreaterOrEqual(t, len(*list), 2)
}

func TestGetUser(t *testing.T) {
	// Arrange
	client := setupClient(t)
	ctx := context.Background()
	newUser := gen.User{Name: "Alexey"}
	createdRes, err := client.AddUser(ctx, &newUser)
	require.NoError(t, err)
	created, ok := createdRes.(*gen.User)
	require.True(t, ok)

	// Act
	res, err := client.GetUser(ctx, gen.GetUserParams{UserId: created.ID.Or(0)})

	// Assert
	require.NoError(t, err)
	u, ok := res.(*gen.User)
	require.True(t, ok)
	require.Equal(t, created.Name, u.Name)
}

func TestUpdateUser(t *testing.T) {
	// Arrange
	client := setupClient(t)
	ctx := context.Background()
	newUser := gen.User{Name: "Alexey"}
	createdRes, err := client.AddUser(ctx, &newUser)
	require.NoError(t, err)
	created, ok := createdRes.(*gen.User)
	require.True(t, ok)
	updated := *created
	updated.Name = "Aleksander"

	// Act
	res, err := client.UpdateUser(ctx, &updated, gen.UpdateUserParams{UserId: created.ID.Or(0)})

	// Assert
	require.NoError(t, err)
	u, ok := res.(*gen.User)
	require.True(t, ok)
	require.Equal(t, "Aleksander", u.Name)
}

func TestDeleteUser(t *testing.T) {
	// Arrange
	client := setupClient(t)
	ctx := context.Background()
	u := gen.User{Name: "Alexey"}
	createdRes, err := client.AddUser(ctx, &u)
	require.NoError(t, err)
	created, ok := createdRes.(*gen.User)
	require.True(t, ok)

	// Act
	delRes, err := client.DeleteUser(ctx, gen.DeleteUserParams{UserId: created.ID.Or(0)})
	getAfterDel, err2 := client.GetUser(ctx, gen.GetUserParams{UserId: created.ID.Or(0)})

	// Assert
	require.NoError(t, err)
	_, ok = delRes.(*gen.DeleteUserNoContent)
	require.True(t, ok)

	require.NoError(t, err2)
	_, ok = getAfterDel.(*gen.GetUserNotFound)
	require.True(t, ok)
}
