package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/fabiuhp/orders-api/model"
	"github.com/redis/go-redis/v9"
)

const orderKeyPrefix = "order:"

type RedisRepo struct {
	Client *redis.Client
}

func orderIDKey(id uint64) string {
	return fmt.Sprintf("%s%d", orderKeyPrefix, id)
}

// Insert adds a new order to the Redis database.
func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to serialize order: %w", err)
	}

	key := orderIDKey(order.OrderID)
	res := r.Client.SetNX(ctx, key, data, 0)
	if err := res.Err(); err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	log.Printf("Order %d inserted successfully", order.OrderID)
	return nil
}

// GetByID retrieves an order by its ID from the Redis database.
func (r *RedisRepo) GetByID(ctx context.Context, id uint64) (model.Order, error) {
	key := orderIDKey(id)
	res := r.Client.Get(ctx, key)
	if res.Err() != nil {
		return model.Order{}, fmt.Errorf("failed to fetch order: %w", res.Err())
	}

	var order model.Order
	err := json.Unmarshal([]byte(res.Val()), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to deserialize order: %w", err)
	}

	log.Printf("Order %d retrieved successfully", id)
	return order, nil
}

// DeleteById removes an order by its ID from the Redis database.
func (r *RedisRepo) DeleteById(ctx context.Context, id uint64) error {
	key := orderIDKey(id)
	res := r.Client.Del(ctx, key)
	if res.Err() != nil {
		return fmt.Errorf("failed to delete order: %w", res.Err())
	}

	log.Printf("Order %d deleted successfully", id)
	return nil
}

// FindAll retrieves all orders from the Redis database.
func (r *RedisRepo) FindAll(ctx context.Context) ([]model.Order, error) {
	keys := r.Client.Keys(ctx, orderKeyPrefix+"*")
	if keys.Err() != nil {
		return nil, fmt.Errorf("failed to fetch all orders: %w", keys.Err())
	}

	orders := make([]model.Order, 0)
	pipe := r.Client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys.Val()))

	for i, key := range keys.Val() {
		cmds[i] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute pipeline: %w", err)
	}

	for _, cmd := range cmds {
		var order model.Order
		err := json.Unmarshal([]byte(cmd.Val()), &order)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize order: %w", err)
		}
		orders = append(orders, order)
	}

	log.Printf("All orders retrieved successfully")
	return orders, nil
}
