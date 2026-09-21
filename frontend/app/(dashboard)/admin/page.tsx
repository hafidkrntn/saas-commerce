"use client"

import { useEffect, useState } from "react"
import {
  DollarSign,
  ShoppingBag,
  Users,
  Package,
  ArrowRight,
  TrendingUp,
} from "lucide-react"
import {
  AreaChart,
  Area,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip as ChartTooltip,
  ResponsiveContainer,
} from "recharts"
import Link from "next/link"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { PageHeader } from "@/components/shared/page-header"
import { StatCard } from "@/components/shared/stat-card"
import { StatusBadge } from "@/components/shared/status-badge"
import { StatCardSkeleton, ChartSkeleton } from "@/components/shared/loading-skeleton"
import { formatCurrency, formatRelativeTime, getInitials } from "@/lib/utils"
import { getAnalyticsSummary, getRevenueData, getSalesData, getTopProducts } from "@/lib/api/analytics"
import { getRecentOrders } from "@/lib/api/orders"
import type { AnalyticsSummary, RevenueData, SalesData, TopProduct } from "@/lib/types"
import type { Order } from "@/lib/types"

export default function DashboardPage() {
  const [loading, setLoading] = useState(true)
  const [summary, setSummary] = useState<AnalyticsSummary | null>(null)
  const [revenueData, setRevenueData] = useState<RevenueData[]>([])
  const [salesData, setSalesData] = useState<SalesData[]>([])
  const [topProducts, setTopProducts] = useState<TopProduct[]>([])
  const [recentOrders, setRecentOrders] = useState<Order[]>([])

  useEffect(() => {
    async function load() {
      const [s, rd, sd, tp, ro] = await Promise.all([
        getAnalyticsSummary(),
        getRevenueData(),
        getSalesData(),
        getTopProducts(),
        getRecentOrders(5),
      ])
      setSummary(s)
      setRevenueData(rd)
      setSalesData(sd)
      setTopProducts(tp)
      setRecentOrders(ro)
      setLoading(false)
    }
    load()
  }, [])

  if (loading) {
    return (
      <div className="p-6 lg:p-8">
        <PageHeader title="Dashboard" description="Overview of your store" />
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => <StatCardSkeleton key={i} />)}
        </div>
        <div className="mt-6 grid gap-6 lg:grid-cols-2">
          <ChartSkeleton /><ChartSkeleton />
        </div>
      </div>
    )
  }

  return (
    <div className="p-6 lg:p-8">
      <PageHeader title="Dashboard" description="Here&apos;s what&apos;s happening with your store today." />

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard title="Revenue" value={formatCurrency(summary?.revenue || 0)} change={summary?.revenueChange} icon={DollarSign} />
        <StatCard title="Orders" value={String(summary?.orders || 0)} change={summary?.ordersChange} icon={ShoppingBag} />
        <StatCard title="Customers" value={String(summary?.customers || 0)} change={summary?.customersChange} icon={Users} />
        <StatCard title="Avg. order" value={formatCurrency(summary?.averageOrderValue || 0)} change={summary?.aovChange} icon={TrendingUp} trend={summary?.aovChange && summary.aovChange >= 0 ? "up" : "down"} />
      </div>

      <div className="mt-6 grid gap-6 lg:grid-cols-7">
        <Card className="lg:col-span-4">
          <CardHeader className="flex flex-row items-start justify-between">
            <div>
              <CardTitle>Revenue</CardTitle>
              <CardDescription>Monthly revenue vs profit</CardDescription>
            </div>
            <Badge variant="outline" className="text-[11px] font-normal rounded-md gap-1 px-2 py-0.5">
              <TrendingUp className="h-3 w-3 text-success" />
              +12.5%
            </Badge>
          </CardHeader>
          <CardContent>
            <div className="h-[260px]">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={revenueData} margin={{ top: 5, right: 5, left: -20, bottom: 0 }}>
                  <defs>
                    <linearGradient id="revGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#2563eb" stopOpacity={0.1} />
                      <stop offset="95%" stopColor="#2563eb" stopOpacity={0} />
                    </linearGradient>
                    <linearGradient id="profGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#059669" stopOpacity={0.1} />
                      <stop offset="95%" stopColor="#059669" stopOpacity={0} />
                    </linearGradient>
                  </defs>
                  <XAxis dataKey="month" axisLine={false} tickLine={false} tick={{ fontSize: 11, fill: "var(--color-muted-foreground)" }} dy={5} />
                  <YAxis axisLine={false} tickLine={false} tick={{ fontSize: 11, fill: "var(--color-muted-foreground)" }} tickFormatter={(v) => `Rp${(v / 1000).toFixed(0)}k`} />
                  <ChartTooltip
                    cursor={false}
                    contentStyle={{ borderRadius: "8px", border: "1px solid var(--color-border)", background: "var(--color-card)", boxShadow: "0 4px 12px rgba(0,0,0,0.08)", fontSize: "12px", padding: "6px 10px" }}
                    formatter={(value) => [formatCurrency(Number(value ?? 0)), undefined]}
                  />
                  <Area type="monotone" dataKey="revenue" stroke="#2563eb" strokeWidth={2} fill="url(#revGrad)" dot={false} activeDot={{ r: 3, strokeWidth: 2, fill: "#fff" }} />
                  <Area type="monotone" dataKey="profit" stroke="#059669" strokeWidth={2} fill="url(#profGrad)" dot={false} activeDot={{ r: 3, strokeWidth: 2, fill: "#fff" }} />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        <Card className="lg:col-span-3">
          <CardHeader>
            <CardTitle>Daily sales</CardTitle>
            <CardDescription>Last 30 days</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[260px]">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={salesData} margin={{ top: 5, right: 5, left: -20, bottom: 0 }}>
                  <XAxis dataKey="date" axisLine={false} tickLine={false} tick={{ fontSize: 10, fill: "var(--color-muted-foreground)" }} tickFormatter={(v) => new Date(v).getDate().toString()} dy={5} />
                  <YAxis axisLine={false} tickLine={false} tick={{ fontSize: 11, fill: "var(--color-muted-foreground)" }} tickFormatter={(v) => `Rp${(v / 1000).toFixed(0)}k`} />
                  <ChartTooltip
                    cursor={false}
                    contentStyle={{ borderRadius: "8px", border: "1px solid var(--color-border)", background: "var(--color-card)", boxShadow: "0 4px 12px rgba(0,0,0,0.08)", fontSize: "12px", padding: "6px 10px" }}
                    formatter={(value) => [formatCurrency(Number(value ?? 0)), "Sales"]}
                    labelFormatter={(label) => new Date(label).toLocaleDateString("en-US", { month: "short", day: "numeric" })}
                  />
                  <Bar dataKey="sales" fill="#2563eb" radius={[3, 3, 0, 0]} maxBarSize={20} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="mt-6 grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>Recent orders</CardTitle>
            <Button variant="ghost" size="sm" className="text-xs text-muted-foreground h-7 px-2" asChild>
              <Link href="/admin/orders">View all <ArrowRight className="h-3 w-3" /></Link>
            </Button>
          </CardHeader>
          <CardContent className="px-0">
            {recentOrders.map((order) => (
              <div key={order.id} className="flex items-center gap-3 px-5 py-2.5 hover:bg-muted/20 transition-colors">
                <Avatar className="h-7 w-7 shrink-0">
                  <AvatarImage src={order.customer.avatar} />
                  <AvatarFallback className="text-[9px]">{getInitials(order.customer.name)}</AvatarFallback>
                </Avatar>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium truncate">{order.customer.name}</p>
                  <p className="text-xs text-muted-foreground/60 font-mono">{order.orderNumber}</p>
                </div>
                <div className="hidden sm:block">
                  <StatusBadge status={order.status} />
                </div>
                <div className="text-right shrink-0">
                  <p className="text-sm font-medium tabular-nums">{formatCurrency(order.total)}</p>
                  <p className="text-xs text-muted-foreground/60">{formatRelativeTime(order.createdAt)}</p>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Top products</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {topProducts.slice(0, 5).map((product) => (
              <div key={product.id} className="flex items-center gap-3">
                <div className="flex h-[30px] w-[30px] shrink-0 items-center justify-center rounded-lg bg-muted">
                  <Package className="h-3.5 w-3.5 text-muted-foreground/50" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium truncate">{product.name}</p>
                  <p className="text-xs text-muted-foreground/60">{product.unitsSold} sold</p>
                </div>
                <p className="text-sm font-medium tabular-nums">{formatCurrency(product.revenue)}</p>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
