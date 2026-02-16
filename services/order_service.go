package services

import (
	"errors"

	"backend/models"
	"gorm.io/gorm"
)

type OrderItemResponse struct {
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderResponse struct {
	ID          uint                `json:"id"`
	TotalAmount float64             `json:"total_amount"`
	Address     string              `json:"address"`
	Status      string              `json:"status"`
	CreatedAt   interface{}         `json:"created_at"`
	UserName    string              `json:"user_name"`
	Items       []OrderItemResponse `json:"items"`
}
func CreateOrder(db *gorm.DB, userID uint, address string) (*OrderResponse, error) {
	var cartItems []models.CartItem
	if err := db.Preload("Product").
		Where("user_id = ?", userID).
		Find(&cartItems).Error; err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	tx := db.Begin()

	order := models.Order{
		UserID:  userID,
		Address: address,
		Status:  "pending",
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	total := 0.0
	for _, item := range cartItems {
		itemTotal := float64(item.Quantity) * item.Product.Price
		total += itemTotal

		orderItem := models.OrderItem{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			UnitPrice:  item.Product.Price,
			Quantity:   item.Quantity,
			TotalPrice: itemTotal,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	order.TotalAmount = total
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return GetOrderByID(db, order.ID, userID)
}
func GetUserOrders(db *gorm.DB, userID uint) ([]OrderResponse, error) {
	var orders []models.Order
	if err := db.Preload("User").
		Preload("OrderItems.Product").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&orders).Error; err != nil {
		return nil, err
	}

	resp := []OrderResponse{}
	for _, o := range orders {
		items := []OrderItemResponse{}
		for _, oi := range o.OrderItems {
			items = append(items, OrderItemResponse{
				ProductID: oi.ProductID,
				Name:      oi.Product.Name,
				Quantity:  oi.Quantity,
				Price:     oi.TotalPrice,
			})
		}

		resp = append(resp, OrderResponse{
			ID:          o.ID,
			TotalAmount: o.TotalAmount,
			Address:     o.Address,
			Status:      o.Status,
			CreatedAt:   o.CreatedAt,
			UserName:    o.User.FullName,
			Items:       items,
		})
	}
	return resp, nil
}
func GetOrderByID(db *gorm.DB, orderID uint, userID uint) (*OrderResponse, error) {
	var order models.Order
	if err := db.Preload("User").
		Preload("OrderItems.Product").
		First(&order, orderID).Error; err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	items := []OrderItemResponse{}
	for _, oi := range order.OrderItems {
		items = append(items, OrderItemResponse{
			ProductID: oi.ProductID,
			Name:      oi.Product.Name,
			Quantity:  oi.Quantity,
			Price:     oi.TotalPrice,
		})
	}

	return &OrderResponse{
		ID:          order.ID,
		TotalAmount: order.TotalAmount,
		Address:     order.Address,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		UserName:    order.User.FullName,
		Items:       items,
	}, nil
}
func GetAllOrders(db *gorm.DB, status string) ([]OrderResponse, error) {
	var orders []models.Order
	query := db.Preload("User").Preload("OrderItems.Product")

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}

	resp := []OrderResponse{}
	for _, o := range orders {
		items := []OrderItemResponse{}
		for _, oi := range o.OrderItems {
			items = append(items, OrderItemResponse{
				ProductID: oi.ProductID,
				Name:      oi.Product.Name,
				Quantity:  oi.Quantity,
				Price:     oi.TotalPrice,
			})
		}

		resp = append(resp, OrderResponse{
			ID:          o.ID,
			TotalAmount: o.TotalAmount,
			Address:     o.Address,
			Status:      o.Status,
			CreatedAt:   o.CreatedAt,
			UserName:    o.User.FullName,
			Items:       items,
		})
	}
	return resp, nil
}
