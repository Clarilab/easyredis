package easyredis_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Clarilab/easyredis"
	"github.com/stretchr/testify/assert"
)

func Test_SetGet(t *testing.T) {
	redisClient, err := easyredis.ConnectToRedis("localhost", "6379", "0", "guest")
	assert.NoError(t, err)

	service := easyredis.New(redisClient)

	ctx := context.TODO()

	t.Run("Test Set", func(t *testing.T) {
		err = service.Set(ctx, time.Hour, "1234", 4)
		assert.NoError(t, err)
	})

	t.Run("Test Get", func(t *testing.T) {
		val, err := service.Get(ctx, "1234")
		if assert.NoError(t, err) {
			var result int
			err = json.Unmarshal(val, &result)
			assert.NoError(t, err)
		}
	})

}
