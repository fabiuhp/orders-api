package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fabiuhp/orders-api/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("erro ao serializar ordem: %w", err)
	}

	key := orderIDKey(order.OrderID)
	res := r.Client.SetNX(ctx, key, data, 0)
	if err := res.Err(); err != nil {
		return fmt.Errorf("erro ao inserir ordem: %w", err)
	}

	return nil
}

func (r *RedisRepo) GetByID(ctx context.Context, id uint64) (model.Order, error) {
	key := orderIDKey(id)
	res := r.Client.Get(ctx, key)
	if res.Err() != nil {
		return model.Order{}, fmt.Errorf("erro ao buscar ordem: %w", res.Err())
	}

	var order model.Order
	err := json.Unmarshal([]byte(res.Val()), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("erro ao deserializar ordem: %w", err)
	}

	return order, nil
}

func (r *RedisRepo) DeleteById(ctx context.Context, id uint64) error {
	key := orderIDKey(id)
	res := r.Client.Del(ctx, key)
	if res.Err() != nil {
		return fmt.Errorf("erro ao deletar ordem: %w", res.Err())
	}

	return nil
}

func (r *RedisRepo) FindAll(ctx context.Context) ([]model.Order, error) {
	keys := r.Client.Keys(ctx, "order:*")
	if keys.Err() != nil {
		return nil, fmt.Errorf("erro ao buscar todas as ordens: %w", keys.Err())
	}

	orders := make([]model.Order, 0)
	for _, key := range keys.Val() {
		res := r.Client.Get(ctx, key)
		if res.Err() != nil {
			return nil, fmt.Errorf("erro ao buscar ordem: %w", res.Err())
		}

		var order model.Order
		err := json.Unmarshal([]byte(res.Val()), &order)
		if err != nil {
			return nil, fmt.Errorf("erro ao deserializar ordem: %w", err)
		}

		orders = append(orders, order)
	}

	return orders, nil
}
