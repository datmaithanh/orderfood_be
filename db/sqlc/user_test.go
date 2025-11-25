package db

import (
	"context"
	"testing"

	"github.com/datmaithanh/orderfood/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	user1 := CreateUserParams{
		Username:     utils.RandomString(6),
		HashPassword: "secret",
		FullName:     utils.RandomString(6),
		Email:        utils.RandomString(6) + "@gmail.com",
	}

	user2, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:     user1.Username,
		HashPassword: user1.HashPassword,
		FullName:     user1.FullName,
		Email:        user1.Email,
	})

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashPassword, user2.HashPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	require.NotZero(t, user2.ID)
	require.NotZero(t, user2.CreatedAt)
}

func TestGetUser(t *testing.T) {
	user1 := CreateUserParams{
		Username:     utils.RandomString(6),
		HashPassword: "secret",
		FullName:     utils.RandomString(6),
		Email:        utils.RandomString(6) + "@gmail.com",
	}

	userCreated, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:     user1.Username,
		HashPassword: user1.HashPassword,
		FullName:     user1.FullName,
		Email:        user1.Email,
	})
	require.NoError(t, err)
	require.NotEmpty(t, userCreated)

	userFetched, err := testQueries.GetUser(context.Background(), userCreated.Username)
	require.NoError(t, err)
	require.NotEmpty(t, userFetched)

	require.Equal(t, userCreated.ID, userFetched.ID)
	require.Equal(t, userCreated.Username, userFetched.Username)
	require.Equal(t, userCreated.HashPassword, userFetched.HashPassword)
	require.Equal(t, userCreated.FullName, userFetched.FullName)
	require.Equal(t, userCreated.Role, userFetched.Role)
	require.Equal(t, userCreated.Email, userFetched.Email)
	require.WithinDuration(t, userCreated.CreatedAt, userFetched.CreatedAt, 0)
}

func TestDeleteUser(t *testing.T) {
	user1 := CreateUserParams{
		Username:     utils.RandomString(6),
		HashPassword: "secret",
		FullName:     utils.RandomString(6),
		Email:        utils.RandomString(6) + "@gmail.com",
	}

	userCreated, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:     user1.Username,
		HashPassword: user1.HashPassword,
		FullName:     user1.FullName,
		Email:        user1.Email,
	})
	require.NoError(t, err)
	require.NotEmpty(t, userCreated)

	err = testQueries.DeleteUser(context.Background(), userCreated.Username)
	require.NoError(t, err)

	userFetched, err := testQueries.GetUser(context.Background(), userCreated.Username)
	require.Error(t, err)
	require.Empty(t, userFetched)
}
