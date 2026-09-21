package analytics

import (
	"backend-go/pkg/apperror"
	"context"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) ServiceInterface {
	return &Service{db: db}
}

// periodBounds resolves a default 30-day window (to = now, from = to-30d).
func periodBounds(from, to time.Time) (time.Time, time.Time) {
	if to.IsZero() {
		to = time.Now()
	}
	if from.IsZero() {
		from = to.AddDate(0, 0, -30)
	}
	return from, to
}

// Summary returns KPIs for the current period plus % change vs the previous
// equal-length period.
func (s *Service) Summary(ctx context.Context, tenantID string, from, to time.Time) (Summary, error) {
	from, to = periodBounds(from, to)
	window := to.Sub(from)
	prevTo := from.Add(-time.Hour)
	prevFrom := prevTo.Add(-window)

	current, err := s.kpiRow(ctx, tenantID, from, to)
	if err != nil {
		return Summary{}, err
	}
	previous, err := s.kpiRow(ctx, tenantID, prevFrom, prevTo)
	if err != nil {
		return Summary{}, err
	}

	return Summary{
		Revenue:          current.Revenue,
		RevenueChange:    pctChange(previous.Revenue, current.Revenue),
		Profit:           current.Profit,
		ProfitChange:     pctChange(previous.Profit, current.Profit),
		Orders:           current.Orders,
		OrdersChange:     pctChange(float64(previous.Orders), float64(current.Orders)),
		ConversionRate:   current.ConversionRate,
		ConversionChange: current.ConversionRate - previous.ConversionRate,
		Customers:        current.Customers,
		CustomersChange:  pctChange(float64(previous.Customers), float64(current.Customers)),
		AverageOrderValue: current.AOV,
		AOVChange:        pctChange(previous.AOV, current.AOV),
	}, nil
}

type kpi struct {
	Revenue        float64
	PaidSubtotal   float64
	Profit         float64
	Orders         int64
	PaidOrders     int64
	Customers      int64
	ConversionRate float64
	AOV            float64
}

func (s *Service) kpiRow(ctx context.Context, tenantID string, from, to time.Time) (kpi, error) {
	var k kpi

	// Revenue (paid), order counts
	if err := s.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(SUM(CASE WHEN payment_status = 'paid' THEN total_amount ELSE 0 END), 0) AS revenue,
			COALESCE(SUM(CASE WHEN payment_status = 'paid' THEN subtotal ELSE 0 END), 0) AS paid_subtotal,
			COUNT(*) AS orders,
			COUNT(*) FILTER (WHERE payment_status = 'paid') AS paid_orders
		FROM orders
		WHERE tenant_id = ? AND created_at BETWEEN ? AND ?`,
		tenantID, from, to).Scan(&k).Error; err != nil {
		return k, apperror.Internal("gagal menghitung metrik penjualan", err)
	}

	// COGS of paid orders
	var cogs float64
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(COALESCE(p.cost_price, 0) * oi.quantity), 0)
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE oi.tenant_id = ? AND o.payment_status = 'paid' AND o.created_at BETWEEN ? AND ?`,
		tenantID, from, to).Scan(&cogs).Error; err != nil {
		return k, apperror.Internal("gagal menghitung harga pokok penjualan", err)
	}
		k.Profit = k.PaidSubtotal - cogs

	if k.Orders > 0 {
		k.ConversionRate = float64(k.PaidOrders) / float64(k.Orders) * 100
		k.AOV = k.Revenue / float64(k.Orders)
	}

	// New customers
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM customers
		WHERE tenant_id = ? AND created_at BETWEEN ? AND ?`,
		tenantID, from, to).Scan(&k.Customers).Error; err != nil {
		return k, apperror.Internal("gagal menghitung customer baru", err)
	}

	return k, nil
}

// Revenue returns monthly aggregates for the last n months.
func (s *Service) Revenue(ctx context.Context, tenantID string, months int) ([]RevenuePoint, error) {
	if months <= 0 || months > 24 {
		months = 12
	}

	rows, err := s.db.WithContext(ctx).Raw(`
		SELECT
			TO_CHAR(orders.created_at, 'YYYY-MM') AS month,
			COALESCE(SUM(CASE WHEN orders.payment_status = 'paid' THEN orders.total_amount ELSE 0 END), 0) AS revenue,
			COALESCE(SUM(CASE WHEN orders.payment_status = 'paid' THEN orders.subtotal ELSE 0 END), 0)
				- COALESCE(SUM(COALESCE(p.cost_price, 0) * oi.quantity), 0) AS profit,
			COUNT(*) AS orders,
			COALESCE(SUM(p.cost_price * oi.quantity), 0) AS expenses
		FROM orders
		LEFT JOIN order_items oi ON oi.order_id = orders.id
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE orders.tenant_id = ?
			AND orders.created_at >= date_trunc('month', now()) - make_interval(months => ?)
		GROUP BY TO_CHAR(orders.created_at, 'YYYY-MM')
		ORDER BY month ASC`, tenantID, months).Rows()
	if err != nil {
		return nil, apperror.Internal("gagal mengambil data revenue", err)
	}
	defer rows.Close()

	var points []RevenuePoint
	for rows.Next() {
		var p RevenuePoint
		if err := rows.Scan(&p.Month, &p.Revenue, &p.Profit, &p.Orders, &p.Expenses); err != nil {
			return nil, apperror.Internal("gagal membaca data revenue", err)
		}
		points = append(points, p)
	}
	return points, nil
}

// Sales returns daily aggregates for the last n days.
func (s *Service) Sales(ctx context.Context, tenantID string, days int) ([]SalesPoint, error) {
	if days <= 0 || days > 365 {
		days = 30
	}

	rows, err := s.db.WithContext(ctx).Raw(`
		SELECT
			TO_CHAR(created_at, 'YYYY-MM-DD') AS date,
			COALESCE(SUM(CASE WHEN payment_status = 'paid' THEN total_amount ELSE 0 END), 0) AS sales,
			COUNT(*) AS orders
		FROM orders
		WHERE tenant_id = ? AND created_at >= now() - make_interval(days => ?)
		GROUP BY TO_CHAR(created_at, 'YYYY-MM-DD')
		ORDER BY date ASC`, tenantID, days).Rows()
	if err != nil {
		return nil, apperror.Internal("gagal mengambil data penjualan", err)
	}
	defer rows.Close()

	var points []SalesPoint
	for rows.Next() {
		var p SalesPoint
		if err := rows.Scan(&p.Date, &p.Sales, &p.Orders); err != nil {
			return nil, apperror.Internal("gagal membaca data penjualan", err)
		}
		points = append(points, p)
	}
	return points, nil
}

// TopProducts returns the best-selling products in a period.
func (s *Service) TopProducts(ctx context.Context, tenantID string, from, to time.Time, limit int) ([]TopProduct, error) {
	from, to = periodBounds(from, to)
	if limit <= 0 || limit > 100 {
		limit = 5
	}

	rows, err := s.db.WithContext(ctx).Raw(`
		SELECT
			p.id::text, p.name,
			COALESCE((SELECT url FROM product_images pi
				WHERE pi.product_id = p.id ORDER BY is_primary DESC, sort_order ASC LIMIT 1), '') AS image,
			SUM(oi.unit_price * oi.quantity) AS revenue,
			SUM(oi.quantity) AS units_sold
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id AND o.tenant_id = oi.tenant_id
		JOIN products p ON p.id = oi.product_id
		WHERE oi.tenant_id = ? AND o.created_at BETWEEN ? AND ? AND o.payment_status = 'paid'
		GROUP BY p.id, p.name
		ORDER BY revenue DESC
		LIMIT ?`, tenantID, from, to, limit).Rows()
	if err != nil {
		return nil, apperror.Internal("gagal mengambil produk terlaris", err)
	}
	defer rows.Close()

	total := 0.0
	var products []TopProduct
	for rows.Next() {
		var p TopProduct
		if err := rows.Scan(&p.ID, &p.Name, &p.Image, &p.Revenue, &p.UnitsSold); err != nil {
			return nil, apperror.Internal("gagal membaca produk terlaris", err)
		}
		total += p.Revenue
		products = append(products, p)
	}

	if total > 0 {
		for i := range products {
			products[i].Percentage = products[i].Revenue / total * 100
		}
	}
	return products, nil
}

// TopCategories returns revenue per category in a period.
func (s *Service) TopCategories(ctx context.Context, tenantID string, from, to time.Time, limit int) ([]TopCategory, error) {
	from, to = periodBounds(from, to)
	if limit <= 0 || limit > 100 {
		limit = 5
	}

	rows, err := s.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(c.name, 'Tanpa Kategori') AS name,
			SUM(oi.unit_price * oi.quantity) AS revenue
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id AND o.tenant_id = oi.tenant_id
		LEFT JOIN products p ON p.id = oi.product_id
		LEFT JOIN categories c ON c.id = p.category_id
		WHERE oi.tenant_id = ? AND o.created_at BETWEEN ? AND ? AND o.payment_status = 'paid'
		GROUP BY c.name
		ORDER BY revenue DESC
		LIMIT ?`, tenantID, from, to, limit).Rows()
	if err != nil {
		return nil, apperror.Internal("gagal mengambil kategori terlaris", err)
	}
	defer rows.Close()

	total := 0.0
	var categories []TopCategory
	for rows.Next() {
		var c TopCategory
		if err := rows.Scan(&c.Name, &c.Revenue); err != nil {
			return nil, apperror.Internal("gagal membaca kategori terlaris", err)
		}
		total += c.Revenue
		categories = append(categories, c)
	}

	if total > 0 {
		for i := range categories {
			categories[i].Percentage = categories[i].Revenue / total * 100
		}
	}
	return categories, nil
}

func pctChange(previous, current float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return (current - previous) / previous * 100
}
