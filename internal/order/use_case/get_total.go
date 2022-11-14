package usecase

import (
	"github.com/mvr-garcia/calc-fees/internal/order/entity"
)

type TotalOutputDTO struct {
	Total int
}

type GetTotalUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewGetTotalUseCase(orderRepository entity.OrderRepositoryInterface) *GetTotalUseCase {
	return &GetTotalUseCase{
		OrderRepository: orderRepository,
	}
}

func (g *GetTotalUseCase) GetTotal() (*TotalOutputDTO, error) {
	total, err := g.OrderRepository.GetTotal()
	if err != nil {
		return nil, err
	}

	return &TotalOutputDTO{
		Total: total,
	}, nil
}
