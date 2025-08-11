package daemons

import (
	"sync"
	"time"

	"github.com/frenswifbenefits/myfren/internal/metrics"

	"go.uber.org/zap"

	"github.com/frenswifbenefits/myfren/internal/entity"
	"github.com/frenswifbenefits/myfren/internal/repository"
)

type CustomerPool struct {
	logger        *zap.Logger
	mu            sync.RWMutex
	tokenIndex    map[string]*entity.Customer
	idIndex       map[int64]*entity.Customer
	usernameIndex map[string]*entity.Customer

	repo *repository.Repository
}

func NewCustomerPool(logger *zap.Logger, repo *repository.Repository) *CustomerPool {
	logger = logger.With(zap.String("component", "customer_pool"))
	return &CustomerPool{
		logger:        logger,
		idIndex:       make(map[int64]*entity.Customer),
		tokenIndex:    make(map[string]*entity.Customer),
		usernameIndex: make(map[string]*entity.Customer),
		repo:          repo,
	}
}

func (cp *CustomerPool) GetByToken(token string) (*entity.Customer, bool) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	customer, ok := cp.tokenIndex[token]
	return customer, ok
}

func (cp *CustomerPool) GetByUsername(username string) (*entity.Customer, bool) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	customer, ok := cp.usernameIndex[username]
	return customer, ok
}

func (cp *CustomerPool) GetByID(id int64) (*entity.Customer, bool) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	customer, ok := cp.idIndex[id]
	return customer, ok
}

func (cp *CustomerPool) GetAll() []entity.Customer {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	customers := make([]entity.Customer, 0, len(cp.idIndex))
	for _, customer := range cp.idIndex {
		customers = append(customers, *customer)
	}
	return customers
}

func (cp *CustomerPool) Invalidate() error {
	customers, err := cp.repo.GetCustomers()
	if err != nil {
		return err
	}

	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.idIndex = make(map[int64]*entity.Customer)
	cp.usernameIndex = make(map[string]*entity.Customer)
	for _, customer := range customers {
		cs := customer
		cp.idIndex[customer.Id] = &cs
		cp.usernameIndex[customer.Username] = &cs
	}

	for token, customer := range cp.tokenIndex {
		username := customer.Username
		if _, ok := cp.usernameIndex[username]; !ok {
			cp.logger.Warn("delete token from customer pool", zap.String("token", token))
			delete(cp.tokenIndex, token)
		}
	}

	metrics.CustomerCount.Set(float64(len(cp.idIndex)))
	return nil
}

func (cp *CustomerPool) RunInvalidate(period time.Duration) error {
	err := cp.Invalidate()
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(period)
		defer ticker.Stop()

		for range ticker.C {
			err := cp.Invalidate()
			if err != nil {
				cp.logger.Error("failed to invalidate", zap.Error(err))
			}
		}
	}()

	return nil
}

func (cp *CustomerPool) AttachToken(token string, customer *entity.Customer) {
	cp.mu.Lock()
	cp.tokenIndex[token] = cp.usernameIndex[customer.Username]
	cp.mu.Unlock()
	cp.logger.Info("added token to pool", zap.String("token", token), zap.String("username", customer.Username))
}
