"use client"

import { useEffect, useState } from "react"
import {
  DollarSign,
  TrendingUp,
  ShoppingBag,
  Percent,
  BarChart3,
} from "lucide-react"
import {
  AreaChart,
  Area,
  BarChart,
  Bar,
  PieChart as RechartPie,
  Pie,
  Cell,
  XAxis,
  YAxis,
  Tooltip as ChartTooltip,
  ResponsiveContainer,
} from "recharts"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { PageHeader } from "@/components/shared/page-header"
import { StatCard } from "@/components/shared/stat-card"
import { StatCardSkeleton, ChartSkeleton } from "@/components/shared/loading-skeleton"
import { formatCurrency } from "@/lib/utils"
import { getAnalyticsSummary, getRevenueData, getTopProducts, getTopCategories } from "@/lib/api/analytics"
import type { AnalyticsSummary, RevenueData, TopProduct, TopCategory } from "@/lib/types"

const COLORS = ["#2563eb", "#059669", "#d97706", "#dc2626", "#7c3aed", "#0891b2", "#db2777", "#ca8a04"]

export default function AnalyticsPage() {
  const [loading, setLoading] = useState(true)
  const [summary, setSummary] = useState<AnalyticsSummary | null>(null)
  const [revenueData, setRevenueData] = useState<RevenueData[]>([])
  const [topProducts, setTopProducts] = useState<TopProduct[]>([])
  const [topCategories, setTopCategories] = useState<TopCategory[]>([])
  const [period, setPeriod] = useState("year")

  useEffect(() => {
    async function load() {
      const [s, rd, tp, tc] = await Promise.all([
        getAnalyticsSummary(),
        getRevenueData(),
        getTopProducts(),
        getTopCategories(),
      ])
      setSummary(s)
      setRevenueData(rd)
      setTopProducts(tp)
      setTopCategories(tc)
      setLoading(false)
    }
    load()
  }, [])

  if (loading) {
    return (
      <div className="p-6 lg:p-8">
        <PageHeader title="Analytics" description="Performance metrics and insights" />
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
      <PageHeader title="Analytics" description="Performance metrics and insights">
        <Tabs value={period} onValueChange={setPeriod}>
          <TabsList className="h-8">
            <TabsTrigger value="7d" className="text-xs px-3">7d</TabsTrigger>
            <TabsTrigger value="30d" className="text-xs px-3">30d</TabsTrigger>
            <TabsTrigger value="year" className="text-xs px-3">1y</TabsTrigger>
          </TabsList>
        </Tabs>
      </PageHeader>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard title="Revenue" value={formatCurrency(summary?.revenue || 0)} change={summary?.revenueChange} icon={DollarSign} />
        <StatCard title="Profit" value={formatCurrency(summary?.profit || 0)} change={summary?.profitChange} icon={TrendingUp} />
        <StatCard title="Orders" value={String(summary?.orders || 0)} change={summary?.ordersChange} icon={ShoppingBag} />
        <StatCard title="Conversion" value={`${summary?.conversionRate || 0}%`} change={summary?.conversionChange} icon={Percent} />
      </div>

      <div className="mt-6 grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>Revenue vs Profit</CardTitle>
              <CardDescription>Monthly comparison</CardDescription>
            </div>
            <Badge variant="outline" className="text-xs font-normal rounded-md gap-1">
              <BarChart3 className="h-3 w-3" />
              12 months
            </Badge>
          </CardHeader>
          <CardContent>
            <div className="h-[280px]">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={revenueData} margin={{ top: 5, right: 5, left: -20, bottom: 0 }}>
                  <defs>
                    <linearGradient id="rev" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#2563eb" stopOpacity={0.1} />
                      <stop offset="95%" stopColor="#2563eb" stopOpacity={0} />
                    </linearGradient>
                    <linearGradient id="prof" x1="0" y1="0" x2="0" y2="1">
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
                  <Area type="monotone" dataKey="revenue" stroke="#2563eb" strokeWidth={2} fill="url(#rev)" dot={false} activeDot={{ r: 3, strokeWidth: 2, fill: "#fff" }} />
                  <Area type="monotone" dataKey="profit" stroke="#059669" strokeWidth={2} fill="url(#prof)" dot={false} activeDot={{ r: 4, stroke: "#059669", strokeWidth: 2, fill: "#fff" }} />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Orders</CardTitle>
            <CardDescription>Monthly volume</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[280px]">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={revenueData} margin={{ top: 5, right: 5, left: -20, bottom: 0 }}>
                  <XAxis dataKey="month" axisLine={false} tickLine={false} tick={{ fontSize: 11, fill: "var(--color-muted-foreground)" }} dy={5} />
                  <YAxis axisLine={false} tickLine={false} tick={{ fontSize: 11, fill: "var(--color-muted-foreground)" }} />
                  <ChartTooltip
                    cursor={false}
                    contentStyle={{ borderRadius: "8px", border: "1px solid var(--color-border)", background: "var(--color-card)", boxShadow: "0 4px 12px rgba(0,0,0,0.08)", fontSize: "12px", padding: "6px 10px" }}
                  />
                  <Bar dataKey="orders" fill="#2563eb" radius={[3, 3, 0, 0]} maxBarSize={24} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="mt-6 grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Top products</CardTitle>
            <CardDescription>By revenue</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {topProducts.slice(0, 6).map((product, i) => (
              <div key={product.id} className="space-y-1.5">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-sm font-medium truncate">{product.name}</span>
                  <span className="text-sm text-muted-foreground tabular-nums">{formatCurrency(product.revenue)}</span>
                </div>
                <div className="flex items-center gap-3">
                  <div className="flex-1 h-1.5 rounded-full bg-muted overflow-hidden">
                    <div className="h-full rounded-full bg-primary" style={{ width: `${Math.max(2, product.percentage)}%` }} />
                  </div>
                  <span className="text-xs text-muted-foreground/60 w-8 text-right tabular-nums">{product.percentage}%</span>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Revenue by category</CardTitle>
            <CardDescription>Distribution</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[280px]">
              <ResponsiveContainer width="100%" height="100%">
                <RechartPie>
                  <Pie data={topCategories} cx="50%" cy="50%" innerRadius={55} outerRadius={95} paddingAngle={2} dataKey="revenue" nameKey="name">
                    {topCategories.map((_, i) => <Cell key={i} fill={COLORS[i % COLORS.length]} />)}
                  </Pie>
                  <ChartTooltip
                    contentStyle={{ borderRadius: "8px", border: "1px solid var(--color-border)", background: "var(--color-card)", boxShadow: "0 4px 12px rgba(0,0,0,0.08)", fontSize: "12px", padding: "6px 10px" }}
                    formatter={(value) => [formatCurrency(Number(value ?? 0)), "Revenue"]}
                  />
                </RechartPie>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
