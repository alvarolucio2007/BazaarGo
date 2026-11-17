package valkey_test

import (
	"testing"

	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/valkey"
	"github.com/go-openapi/testify/require"
	"github.com/redis/go-redis/v9"
)

func TestRecentlyViewedRepo_Add(t *testing.T) {
	validProductID := []string{"product"}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		clientID  string
		productID string
		wantErr   bool
	}{
		{
			name:      "should cache valid productID and clientID",
			rdb:       testRDB,
			clientID:  "client",
			productID: "product",
			wantErr:   false,
		}, {
			name:      "should not cache invalid productID",
			rdb:       testRDB,
			clientID:  "client",
			productID: "",
			wantErr:   true,
		}, {
			name:      "should not cache invalid productID",
			rdb:       testRDB,
			clientID:  "",
			productID: "product",
			wantErr:   true,
		}, {
			name:      "should throw err if cache is down",
			rdb:       redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			clientID:  "",
			productID: "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := valkey.NewRecentlyViewedRepo(tt.rdb)
			gotErr := r.Add(t.Context(), tt.clientID, tt.productID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Add() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Add() succeeded unexpectedly")
			}
			cachedID, err := r.Get(t.Context(), tt.clientID)
			require.NoError(t, err)
			require.Equal(t, validProductID, cachedID)
		})
	}
}

func TestRecentlyViewedRepo_Get(t *testing.T) {
	clientID := "client"
	productIDs := []string{"product1", "product2", "product3", "product4", "product5"}
	r := valkey.NewRecentlyViewedRepo(testRDB)
	for _, productID := range productIDs {
		err := r.Add(t.Context(), clientID, productID)
		require.NoError(t, err)
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		clientID string
		want     []string
		wantErr  bool
	}{
		{
			name:     "should fetch all 5 productIDs from a valid client",
			rdb:      testRDB,
			clientID: clientID,
			want:     productIDs,
			wantErr:  false,
		}, {
			name:     "should throw err and not accept empty user",
			rdb:      testRDB,
			clientID: "",
			want:     []string{},
			wantErr:  true,
		}, {
			name:     "shouldn't return anything if user isn't cached",
			rdb:      testRDB,
			clientID: "not a cached client",
			want:     []string{},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := valkey.NewRecentlyViewedRepo(tt.rdb)
			got, gotErr := r.Get(t.Context(), tt.clientID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Get() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Get() succeeded unexpectedly")
			}
			for i, product := range got {
				inverse := 4 - i
				require.Equal(t, productIDs[inverse], product)
			}
		})
	}
}
