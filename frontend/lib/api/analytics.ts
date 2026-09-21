import { http } from "./http"
import type { AnalyticsSummary, RevenueData, SalesData, TopCategory, TopProduct } from "@/lib/types"

interface SummaryDTO {
  revenue: number
  revenue_change: number
  profit: number
  profit_change: number
  orders: number
  orders_change: number
  conversion_rate: number
  conversion_change: number
  customers: number
  customers_change: number
  average_order_value: number
  aov_change: number
}

interface RevenuePointDTO {
  month: string
  revenue: number
  profit: number
  orders: number
  expenses: number
}

interface SalesPointDTO {
  date: string
  sales: number
  orders: number
}

interface TopProductDTO {
  id: string
  name: string
  image: string
  revenue: number
  units_sold: number
  percentage: number
}

interface TopCategoryDTO {
  name: string
  revenue: number
  percentage: number
}

const MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]

function monthLabel(month: string): string {
  const parts = month.split("-")
  if (parts.length !== 2) return month
  const idx = Number(parts[1]) - 1
  return MONTHS[idx] ?? month
}

export async function getAnalyticsSummary(): Promise<AnalyticsSummary> {
  const d = await http.get<SummaryDTO>("/analytics/summary")
  return {
    revenue: d.revenue,
    revenueChange: d.revenue_change,
    profit: d.profit,
    profitChange: d.profit_change,
    orders: d.orders,
    ordersChange: d.orders_change,
    conversionRate: d.conversion_rate,
    conversionChange: d.conversion_change,
    customers: d.customers,
    customersChange: d.customers_change,
    averageOrderValue: d.average_order_value,
    aovChange: d.aov_change,
  }
}

export async function getRevenueData(): Promise<RevenueData[]> {
  const res = await http.get<RevenuePointDTO[]>("/analytics/revenue?months=12")
  return res.map((p) => ({
    month: monthLabel(p.month),
    revenue: p.revenue,
    profit: p.profit,
    orders: p.orders,
    expenses: p.expenses,
  }))
}

export async function getSalesData(): Promise<SalesData[]> {
  const res = await http.get<SalesPointDTO[]>("/analytics/sales?days=30")
  return res.map((p) => ({ date: p.date, sales: p.sales, orders: p.orders }))
}

export async function getTopProducts(): Promise<TopProduct[]> {
  const res = await http.get<TopProductDTO[]>("/analytics/top-products?limit=8")
  return res.map((p) => ({
    id: p.id,
    name: p.name,
    image: p.image,
    revenue: p.revenue,
    unitsSold: p.units_sold,
    percentage: p.percentage,
  }))
}

export async function getTopCategories(): Promise<TopCategory[]> {
  const res = await http.get<TopCategoryDTO[]>("/analytics/top-categories?limit=8")
  return res.map((c) => ({ name: c.name, revenue: c.revenue, percentage: c.percentage }))
}
