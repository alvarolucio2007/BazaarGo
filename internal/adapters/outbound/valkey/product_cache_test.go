package valkey_test

import (
	"testing"

	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/valkey"
	"github.com/alvarolucio2007/bazaargo/internal/domain"
	"github.com/go-openapi/testify/require"
	"github.com/redis/go-redis/v9"
)

func TestProductCache_Create(t *testing.T) {
	product := domain.Product{ID: "1", Specs: map[string]any{}}
	cleanupTest(t)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		item    domain.Product
		wantErr bool
	}{
		{
			name:    "should cache valid user",
			rdb:     testRDB,
			item:    product,
			wantErr: false,
		},
		{
			name:    "should not cache invalid product",
			rdb:     testRDB,
			item:    domain.Product{},
			wantErr: true,
		},
		{
			name:    "should throw err if cache is down",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			item:    domain.Product{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewProductCache(tt.rdb)
			gotErr := c.Create(t.Context(), tt.item)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Create() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Create() succeeded unexpectedly")
			}
			cachedProduct, err := c.Read(t.Context(), product.ID)
			require.NoError(t, err)
			require.Equal(t, product, cachedProduct)
		})
	}
}

func TestProductCache_Read(t *testing.T) {
	product := domain.Product{ID: "1", Specs: map[string]any{}}
	productCache := valkey.NewProductCache(testRDB)
	err := productCache.Create(t.Context(), product)
	require.NoError(t, err)
	cleanupTest(t)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		itemID  string
		want    domain.Product
		wantErr bool
	}{
		{
			name:    "should read valid product",
			itemID:  product.ID,
			rdb:     testRDB,
			want:    product,
			wantErr: false,
		}, {
			name:    "should throw err if user doesnt exist",
			itemID:  "",
			rdb:     testRDB,
			want:    domain.Product{},
			wantErr: true,
		}, {
			name:    "should throw err if cache is down",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			itemID:  "",
			want:    domain.Product{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := valkey.NewProductCache(tt.rdb)
			got, gotErr := p.Read(t.Context(), tt.itemID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Read() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Read() succeeded unexpectedly")
			}
			require.Equal(t, product, got)
		})
	}
}

func TestProductCache_Delete(t *testing.T) {
	product := domain.Product{ID: "1", Specs: map[string]any{}}
	productCache := valkey.NewProductCache(testRDB)
	err := productCache.Create(t.Context(), product)
	require.NoError(t, err)
	cleanupTest(t)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		itemID  string
		wantErr bool
	}{
		{
			name:    "should delete valid product",
			rdb:     testRDB,
			itemID:  product.ID,
			wantErr: false,
		}, {
			name:    "should throw err if cache is down",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			itemID:  "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := valkey.NewProductCache(tt.rdb)
			gotErr := p.Delete(t.Context(), tt.itemID)
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
