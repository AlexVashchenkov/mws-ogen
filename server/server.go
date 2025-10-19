package server

import (
	"context"
	api "mws-ogen/gen"
	gen "mws-ogen/gen"
	"net/http"
	"sync"
)

type UserStorage struct {
	mu    sync.Mutex
	users map[int64]gen.User
}

func NewUserStorage() *UserStorage {
	return &UserStorage{users: make(map[int64]gen.User)}
}

func (s *UserStorage) GetUser(ctx context.Context, params gen.GetUserParams) (gen.GetUserRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[params.UserId]
	if !ok {
		return &gen.GetUserNotFound{Message: api.NewOptString("User not found")}, nil
	}
	u := user
	return &u, nil
}

func (s *UserStorage) ListUsers(ctx context.Context) (gen.ListUsersRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var users []gen.User
	for _, u := range s.users {
		users = append(users, u)
	}
	ok := gen.ListUsersOKApplicationJSON(users)
	return &ok, nil
}

func (s *UserStorage) AddUser(ctx context.Context, req *gen.User) (gen.AddUserRes, error) {
	if req == nil || req.Name == "" {
		return &api.AddUserBadRequest{Message: api.NewOptString("Name is required")}, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	id := req.ID.Or(0)
	if id == 0 {
		// простая авто-выдача id
		var max int64
		for k := range s.users {
			if k > max {
				max = k
			}
		}
		id = max + 1
		req.SetID(api.NewOptInt64(id))
	}

	u := api.User{ID: req.ID, Name: req.Name, Email: req.Email}
	s.users[id] = u

	return &u, nil
}

func (s *UserStorage) UpdateUser(ctx context.Context, req *gen.User, params gen.UpdateUserParams) (gen.UpdateUserRes, error) {
	if req == nil || req.Name == "" {
		return &api.UpdateUserBadRequest{Message: api.NewOptString("Name is required")}, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.users[params.UserId]
	if !ok {
		return &api.UpdateUserNotFound{Message: api.NewOptString("User not found")}, nil
	}
	u := api.User{
		ID:    api.NewOptInt64(params.UserId),
		Name:  req.Name,
		Email: req.Email,
	}
	s.users[params.UserId] = u
	return &u, nil
}

func (s *UserStorage) DeleteUser(ctx context.Context, params gen.DeleteUserParams) (gen.DeleteUserRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[params.UserId]; !ok {
		return &api.DeleteUserNotFound{Message: api.NewOptString("user not found")}, nil
	}
	delete(s.users, params.UserId)
	return &gen.DeleteUserNoContent{}, nil
}

func (s *UserStorage) NewError(ctx context.Context, err error) *api.ErrRespStatusCode {
	return &api.ErrRespStatusCode{StatusCode: http.StatusInternalServerError, Response: gen.ErrResp{Message: gen.NewOptString("User not found")}}
}
