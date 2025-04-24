package integration_test_arrange

import (
	"aggregationframework/cmd/provider"
	database "aggregationframework/internal/db"
	model "aggregationframework/internal/model/domain"
	"context"
	"testing"
)

func CreateTestCache(t *testing.T, ctx context.Context) *database.Cache {
	provider := provider.NewProvider("test")
	return provider.ProvideCache(ctx)
}

func AddCachedFollowersToCache(t *testing.T, cache *database.Cache, username, lastFollowerId string, limit int, followers []model.Follower) {
	cache.Client.SetUserFollowers(username, lastFollowerId, limit, followers)
}

func AddCachedFolloweesToCache(t *testing.T, cache *database.Cache, username, lastFolloweeId string, limit int, followees []model.Followee) {
	cache.Client.SetUserFollowees(username, lastFolloweeId, limit, followees)
}
