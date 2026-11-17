package postgres_test

import (
	"database/sql"
	"testing"

	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/postgres"
	"github.com/alvarolucio2007/bazaargo/internal/domain"
	"github.com/go-openapi/testify/require"
)

func createUserTest(t *testing.T, user *domain.User) {
	r := postgres.NewUserRepo(testDB)
	err := r.Create(t.Context(), user)
	require.NoError(t, err)
}
func cleanupUser(t *testing.T, userID int64) {
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM users where id=$1", userID)
		require.NoError(t, err)
	})
}
func TestUserRepo_Create(t *testing.T) {
	user := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	r := postgres.NewUserRepo(testDB)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *sql.DB
		// Named input parameters for target function.
		user    *domain.User
		wantErr bool
	}{
		{
			name:    "should accept normal users with all default data",
			user:    user,
			wantErr: false,
		},
		{
			name:    "should not accept users without some data",
			user:    &domain.User{Name: "TestWithoutOtherParameters"},
			wantErr: true,
		},
		{
			name:    "should not accept users with same CPF",
			user:    &domain.User{Name: "TestSameCPF", CPF: "00000000000", Email: "test2@gmail.com", PasswordHash: []byte("1234567")},
			wantErr: true,
		},
		{
			name:    "should not accept users with same email",
			user:    &domain.User{Name: "TestSameEmail", CPF: "00000000001", Email: "test@gmail.com", PasswordHash: []byte("123456")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := r.Create(t.Context(), tt.user)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Create() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Create() succeeded unexpectedly")
			}
		})
	}
	cleanupUser(t, user.ID)
}

func TestUserRepo_Read(t *testing.T) {
	user1 := domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, &user1)
	r := postgres.NewUserRepo(testDB)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *sql.DB
		// Named input parameters for target function.
		userID  int64
		want    domain.User
		wantErr bool
	}{
		{
			name:    "should accept normal data and return existing user",
			userID:  user1.ID,
			want:    user1,
			wantErr: false,
		},
		{
			name:    "should return error when user doesn't exist",
			userID:  0,
			want:    domain.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := r.Read(t.Context(), tt.userID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetByID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetByID() succeeded unexpectedly")
			}
			require.EqualValues(t, got.ID, tt.want.ID)
			require.EqualValues(t, got.Email, tt.want.Email)
			require.EqualValues(t, got.CPF, tt.want.CPF)
			require.EqualValues(t, got.CreatedAt, tt.want.CreatedAt)
			require.EqualValues(t, got.UpdatedAt, tt.want.UpdatedAt)
		})
	}
	cleanupUser(t, user1.ID)
}

func TestUserRepo_ReadByEmail(t *testing.T) {
	user1 := domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, &user1)
	r := postgres.NewUserRepo(testDB)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *sql.DB
		// Named input parameters for target function.
		email   string
		want    domain.User
		wantErr bool
	}{
		{
			name:    "should accept normal email and return existing user",
			email:   "test@gmail.com",
			want:    user1,
			wantErr: false,
		},
		{
			name:    "shouldn't return anything if e-mail is not in the DB",
			email:   "notexistent@gmail.com",
			want:    domain.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := r.ReadByEmail(t.Context(), tt.email)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetByID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetByID() succeeded unexpectedly")
			}
			require.EqualValues(t, tt.want.ID, got.ID)
			require.EqualValues(t, tt.want.Email, got.Email)
			require.EqualValues(t, tt.want.PasswordHash, got.PasswordHash)
		})
	}
	cleanupUser(t, user1.ID)
}

func TestUserRepo_Update(t *testing.T) {
	user1 := domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, &user1)
	r := postgres.NewUserRepo(testDB)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *sql.DB
		// Named input parameters for target function.
		user    domain.User
		wantErr bool
	}{
		{
			name:    "shouldn't edit user if all params are empty",
			user:    domain.User{ID: user1.ID, Name: "", CPF: "", Email: "", PasswordHash: []byte("")},
			wantErr: false,
		}, {
			name:    "should edit only name",
			user:    domain.User{ID: user1.ID, Name: "changed", CPF: "", Email: "", PasswordHash: []byte("")},
			wantErr: false,
		}, {
			name:    "should edit only CPF",
			user:    domain.User{ID: user1.ID, Name: "", CPF: "00000000001", Email: "", PasswordHash: []byte("")},
			wantErr: false,
		}, {
			name:    "should edit only email",
			user:    domain.User{ID: user1.ID, Name: "", CPF: "", Email: "changed@gmail.com", PasswordHash: []byte("")},
			wantErr: false,
		}, {
			name:    "should edit only password",
			user:    domain.User{ID: user1.ID, Name: "", CPF: "", Email: "", PasswordHash: []byte("changed")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := r.Update(t.Context(), tt.user)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Update() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Update() succeeded unexpectedly")
			}
			updatedUser, err := r.Read(t.Context(), user1.ID)
			require.NoError(t, err)
			if tt.user.Name != "" {
				require.NotEqual(t, user1.Name, updatedUser.Name)
			}
			if tt.user.CPF != "" {
				require.NotEqual(t, user1.CPF, updatedUser.CPF)
			}
			if tt.user.Email != "" {
				require.NotEqual(t, user1.Email, updatedUser.Email)
			}
			if string(tt.user.PasswordHash) != "" {
				require.NotEqual(t, updatedUser.PasswordHash, tt.user.PasswordHash)
			}
		})
	}
	cleanupUser(t, user1.ID)
}

func TestUserRepo_Delete(t *testing.T) {
	user := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user)
	r := postgres.NewUserRepo(testDB)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *sql.DB
		// Named input parameters for target function.
		userID  int64
		wantErr bool
	}{
		{
			name:    "should delete with valid ID",
			userID:  user.ID,
			wantErr: false,
		}, {
			name:    "shouldn't delete with invalid ID",
			userID:  user.ID + 1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := r.Delete(t.Context(), tt.userID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Delete() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Delete() succeeded unexpectedly")
			}
		})
	}
}

func TestUserRepo_Ping(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db      *sql.DB
		wantErr bool
	}{
		{
			name:    "should ping valid DB",
			db:      testDB,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := postgres.NewUserRepo(tt.db)
			gotErr := u.Ping(t.Context())
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Ping() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Ping() succeeded unexpectedly")
			}
		})
	}
}
