package database

import (
	"database/sql"
	"testing"

	"github.com/mvr-garcia/calc-fees/internal/order/entity"
	"github.com/stretchr/testify/suite"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	DB *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)

	_, err = db.Exec("CREATE TABLE orders(id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.NoError(err)

	suite.DB = db
}

func (suite *OrderRepositoryTestSuite) TearDownSuite() {
	suite.DB.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestSaveOrder() {

	// given
	order, err := entity.NewOrder("123", 10, 2)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())

	// when
	repo := NewOrderRepository(suite.DB)
	err = repo.Save(order)
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
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}
