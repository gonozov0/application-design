package orders

import (
	"sync"

	"application-design/internal/domain"
	"application-design/internal/domain/orders"
)

type InMemoryRepo struct {
	orders map[domain.OrderID]orders.Order
	mu     sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		orders: make(map[domain.OrderID]orders.Order),
	}
}

func (r *InMemoryRepo) SaveOrder(order orders.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.ID()] = order
	return nil
}

func (r *InMemoryRepo) GetOrder(id domain.OrderID) (*orders.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return nil, orders.ErrOrderNotFound
	}
	return &order, nil
}
