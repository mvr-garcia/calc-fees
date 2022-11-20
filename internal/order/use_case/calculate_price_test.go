package usecase

import (
	"database/sql"
	"testing"

	"github.com/mvr-garcia/calc-fees/internal/order/entity"
	"github.com/mvr-garcia/calc-fees/internal/order/infra/database"
	"github.com/stretchr/testify/suite"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

type CalculateFinalPriceUseCaseTestSuite struct {
	suite.Suite
	OrderRepository entity.OrderRepositoryInterface
	DB              *sql.DB
}

func (suite *CalculateFinalPriceUseCaseTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)

	_, err = db.Exec("CREATE TABLE orders(id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.NoError(err)

	suite.DB = db
	suite.OrderRepository = database.NewOrderRepository(suite.DB)
}

func (suite *CalculateFinalPriceUseCaseTestSuite) TearDownSuite() {
	suite.DB.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CalculateFinalPriceUseCaseTestSuite))
}

func (suite *CalculateFinalPriceUseCaseTestSuite) TestCalculateFinalPrice() {

	// given
	order, err := entity.NewOrder("123", 10, 2)
	suite.NoError(err)

	// when
	useCase := NewCalculateFinalPriceUseCase(suite.OrderRepository)
	orderOutput, err := useCase.Execute(
		OrderInputDTO{
			ID:    order.ID,
			Price: order.Price,
			Tax:   order.Tax,
		},
	)
	suite.NoError(err)

	// then
	var orderResult entity.Order
	result := suite.DB.QueryRow(
		"SELECT id, price, tax, final_price FROM orders WHERE id = ?",
		order.ID,
	)
	err = result.Scan(
		&orderResult.ID,
		&orderResult.Price,
		&orderResult.Tax,
		&orderResult.FinalPrice,
	)

	suite.NoError(err)
	suite.Equal(orderOutput.ID, orderResult.ID)
	suite.Equal(orderOutput.Price, orderResult.Price)
	suite.Equal(orderOutput.Tax, orderResult.Tax)
	suite.Equal(orderOutput.FinalPrice, orderResult.FinalPrice)
}
