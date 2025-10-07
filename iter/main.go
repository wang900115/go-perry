package main

// Eager 立即執行
// Lazy 延遲執行

import (
	"cmp"
	"fmt"
	"iter"
	"slices"
)

type Order struct {
	OrderID      string
	CustomerName string
	Amount       float64
	Status       string
}

func filter[V any](it iter.Seq[V], keep func(V) bool) iter.Seq[V] {
	seq := func(yield func(V) bool) {
		for v := range it {
			if keep(v) {
				if !yield(v) {
					break
				}
			}
		}
	}
	return seq
}

func display(it iter.Seq[Order]) {
	for order := range it {
		fmt.Printf("Order ID: %s, Customer: %s, Amount: %.2f, Status: %s\n", order.OrderID, order.CustomerName, order.Amount, order.Status)
	}
}

func main() {
	orders := []Order{
		{OrderID: "1", CustomerName: "Alice", Amount: 75.50, Status: "delivered"},
		{OrderID: "2", CustomerName: "Bob", Amount: 45.90, Status: "delivered"},
		{OrderID: "3", CustomerName: "Candy", Amount: 12.00, Status: "pending"},
		{OrderID: "4", CustomerName: "Dan", Amount: 30.99, Status: "canceled"},
		{OrderID: "5", CustomerName: "Eve", Amount: 71.22, Status: "delivered"},
	}

	// Filter orders wuth amount > 50
	highValueOrders := filter(slices.Values(orders), func(order Order) bool {
		return order.Amount > 50
	})
	display(highValueOrders)

	// Filter orders with status "delivered"
	deliveredOrders := filter(slices.Values(orders), func(order Order) bool {
		return order.Status == "delivered"
	})
	display(deliveredOrders)

	// Explore Collect() function
	seq := func(yield func(Order) bool) {
		for _, o := range orders {
			if o.Amount < 20 {
				if !yield(o) {
					return
				}
			}
		}
	}

	filteredOrders := slices.Collect(seq)
	for _, o := range filteredOrders {
		fmt.Printf("%s: %.2f\n", o.OrderID, o.Amount)
	}

	// Sort orders by amount
	sortFunc := func(a, b Order) int {
		return cmp.Compare(b.Amount, a.Amount)
	}
	sortedOrders := slices.SortedFunc(slices.Values(orders), sortFunc)
	for _, o := range sortedOrders {
		fmt.Printf("OrderID: %s, Customer: %s, Amount: %.2f, Status: %s\n", o.OrderID, o.CustomerName, o.Amount, o.Status)
	}

	// Chunk 3 order
	for c := range slices.Chunk(orders, 3) {
		fmt.Println(c)
	}

}
