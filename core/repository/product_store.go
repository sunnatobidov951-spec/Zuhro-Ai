package repository

// Product — структура нашего товара
type Product struct {
	ID    int
	Name  string
	Price float64
}

// GetProductByID — функция, которая будет доставать товары
// молниеносно быстро из памяти нашего "Склада"
func GetProductByID(id int) (*Product, error) {
	// Пока это заглушка, потом мы подключим сюда реальную базу данных
	return &Product{ID: id, Name: "Мощный товар Zuhro", Price: 99.9}, nil
}
