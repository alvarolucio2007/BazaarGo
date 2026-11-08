package app

import (
	"crypto/rand"
	"database/sql"
	"log/slog"

	api "github.com/alvarolucio2007/bazaargo/internal/adapters/inbound/http"
	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/argon2"
	iMongo "github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/mongo"
	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/paseto"
	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/postgres"
	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/valkey"
	"github.com/alvarolucio2007/bazaargo/internal/service"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	Handler *api.ServerHandler
}

func New(db *sql.DB, mColl *mongo.Collection, rdb *redis.Client, logger *slog.Logger, symmetricKey []byte) (*App, error) {
	pasetoAuth, err := paseto.NewPasetoAuthenticator(symmetricKey)
	if err != nil {
		return nil, err
	}
	hasher := argon2.NewArgon2Hasher(64*1024, 3, 2)
	userRepo := postgres.NewUserRepo(db)
	userService := service.NewUserService(userRepo, hasher, pasetoAuth)

	productCache := valkey.NewProductCache(rdb)
	recentlyViewedCache := valkey.NewRecentlyViewedRepo(rdb)

	productRepo := iMongo.NewProductRepo(mColl, logger)
	productService := service.NewProductService(productRepo, productCache, recentlyViewedCache, logger)

	handler := api.NewServerHandler(userService, productService, pasetoAuth, logger)
	return &App{handler}, nil
}
func GenerateSymmetricKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}
