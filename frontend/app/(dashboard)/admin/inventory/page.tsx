"use client"

import { useEffect, useState, useMemo } from "react"
import {
  Search,
  ClipboardList,
  AlertTriangle,
  Package,
  Truck,
  ArrowUpDown,
  AlertCircle,
  CheckCircle2,
} from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { Progress } from "@/components/ui/progress"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Pagination } from "@/components/ui/pagination"
import { PageHeader } from "@/components/shared/page-header"
import { StatusBadge } from "@/components/shared/status-badge"
import { StatCard } from "@/components/shared/stat-card"
import { TableSkeleton, ChartSkeleton } from "@/components/shared/loading-skeleton"
import { formatCurrency, formatDate } from "@/lib/utils"
import { getInventoryItems, getStockMovements, getIncomingShipments } from "@/lib/api/inventory"
import { getProducts } from "@/lib/api/products"
import type { InventoryItem, StockMovement, IncomingShipment } from "@/lib/types"

export default function InventoryPage() {
  const [items, setItems] = useState<InventoryItem[]>([])
  const [movements, setMovements] = useState<StockMovement[]>([])
  const [incoming, setIncoming] = useState<IncomingShipment[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState("")
  const [statusFilter, setStatusFilter] = useState("all")
  const [warehouseFilter, setWarehouseFilter] = useState("all")
  const [page, setPage] = useState(1)
  const perPage = 10

  useEffect(() => {
    async function load() {
      const [i, m, inc] = await Promise.all([
        getInventoryItems(),
        getStockMovements(),
        getIncomingShipments(),
      ])
      setItems(i)
      setMovements(m)
      setIncoming(inc)
      setLoading(false)
    }
    load()
  }, [])

  const filtered = useMemo(() => {
    let result = items
    if (search) {
      const q = search.toLowerCase()
      result = result.filter((i) => i.productName.toLowerCase().includes(q) || i.productSku.toLowerCase().includes(q))
    }
    if (statusFilter !== "all") result = result.filter((i) => i.status === statusFilter)
    if (warehouseFilter !== "all") result = result.filter((i) => i.warehouse === warehouseFilter)
    return result
  }, [items, search, statusFilter, warehouseFilter])

  const totalPages = Math.ceil(filtered.length / perPage)
  const paginated = filtered.slice((page - 1) * perPage, page * perPage)

  const warehouses = useMemo(() => [...new Set(items.map((i) => i.warehouse))], [items])
  const lowStockItems = items.filter((i) => i.status === "low_stock")
  const outOfStockItems = items.filter((i) => i.status === "out_of_stock")

  if (loading) {
    return (
      <div className="p-6">
        <PageHeader title="Inventory" description="Monitor stock levels and movements" />
        <div className="grid gap-4 sm:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => <ChartSkeleton key={i} />)}
        </div>
      </div>
    )
  }

  return (
    <div className="p-6">
      <PageHeader title="Inventory" description="Monitor stock levels and movements" />

      <div className="grid gap-4 sm:grid-cols-3">
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-success/10">
              <Package className="h-5 w-5 text-success" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">In Stock</p>
              <p className="text-xl font-bold">{items.filter((i) => i.status === "in_stock").length}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-warning/10">
              <AlertTriangle className="h-5 w-5 text-warning" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Low Stock</p>
              <p className="text-xl font-bold text-warning">{lowStockItems.length}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-destructive/10">
              <AlertCircle className="h-5 w-5 text-destructive" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Out of Stock</p>
              <p className="text-xl font-bold text-destructive">{outOfStockItems.length}</p>
            </div>
          </CardContent>
        </Card>
      </div>

      <Tabs defaultValue="all" className="mt-6">
        <TabsList>
          <TabsTrigger value="all">All Stock</TabsTrigger>
          <TabsTrigger value="low" className="gap-2">
            Low Stock
            {lowStockItems.length > 0 && <Badge variant="warning" className="h-5 px-1.5 text-[10px]">{lowStockItems.length}</Badge>}
          </TabsTrigger>
          <TabsTrigger value="out" className="gap-2">
            Out of Stock
            {outOfStockItems.length > 0 && <Badge variant="destructive" className="h-5 px-1.5 text-[10px]">{outOfStockItems.length}</Badge>}
          </TabsTrigger>
          <TabsTrigger value="incoming" className="gap-2">
            <Truck className="h-4 w-4" /> Incoming
          </TabsTrigger>
          <TabsTrigger value="movements" className="gap-2">
            Movements
          </TabsTrigger>
        </TabsList>

        <TabsContent value="all" className="mt-4">
          <Card>
            <div className="flex flex-col sm:flex-row sm:items-center gap-3 px-5 pt-5 pb-0">
              <div className="relative flex-1 max-w-xs">
                <Search className="absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground/40" />
                <Input placeholder="Search inventory..." value={search} onChange={(e) => { setSearch(e.target.value); setPage(1) }} className="h-8 pl-9 text-xs" />
              </div>
              <div className="flex items-center gap-2">
                <Select value={statusFilter} onValueChange={(v) => { setStatusFilter(v); setPage(1) }}>
                  <SelectTrigger className="h-8 w-[120px] text-xs">
                    <SelectValue placeholder="Status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All status</SelectItem>
                    <SelectItem value="in_stock">In stock</SelectItem>
                    <SelectItem value="low_stock">Low stock</SelectItem>
                    <SelectItem value="out_of_stock">Out of stock</SelectItem>
                    <SelectItem value="overstocked">Overstocked</SelectItem>
                  </SelectContent>
                </Select>
                <Select value={warehouseFilter} onValueChange={(v) => { setWarehouseFilter(v); setPage(1) }}>
                  <SelectTrigger className="h-8 w-[140px] text-xs">
                    <SelectValue placeholder="Warehouse" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All warehouses</SelectItem>
                    {warehouses.map((w) => (
                      <SelectItem key={w} value={w}>{w}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            <CardContent className="p-0 mt-4">
              <table className="w-full">
                <thead>
                  <tr className="border-y text-left">
                    <th className="px-5 py-3 text-xs font-medium text-muted-foreground/70 uppercase tracking-wider">Product</th>
                    <th className="px-2 py-3 text-xs font-medium text-muted-foreground/70 uppercase tracking-wider hidden md:table-cell">SKU</th>
                    <th className="px-2 py-3 text-xs font-medium text-muted-foreground/70 uppercase tracking-wider text-right">Stock</th>
                    <th className="px-2 py-3 text-xs font-medium text-muted-foreground/70 uppercase tracking-wider text-right hidden sm:table-cell">Reserved</th>
                    <th className="px-2 py-3 text-xs font-medium text-muted-foreground/70 uppercase tracking-wider text-right hidden sm:table-cell">Available</th>
                    <th className="px-2 py-3 text-xs font-medium text-muted-foreground/70 uppercase tracking-wider">Status</th>
                    <th className="px-5 py-3 text-xs font-medium text-muted-foreground/70 uppercase tracking-wider hidden lg:table-cell">Warehouse</th>
                  </tr>
                </thead>
                <tbody>
                  {paginated.map((item) => (
                    <tr key={item.id} className="border-b last:border-0 hover:bg-muted/20 transition-colors">
                      <td className="px-5 py-3">
                        <div className="flex items-center gap-3">
                          <div className="flex h-[34px] w-[34px] shrink-0 items-center justify-center rounded-lg bg-muted">
                            <Package className="h-4 w-4 text-muted-foreground/60" />
                          </div>
                          <span className="text-sm font-medium truncate max-w-[180px]">{item.productName}</span>
                        </div>
                      </td>
                      <td className="px-2 py-3 hidden md:table-cell"><span className="text-xs text-muted-foreground/70 font-mono">{item.productSku}</span></td>
                      <td className="px-2 py-3 text-right">
                        <span className={`text-sm font-medium tabular-nums ${
                          item.status === "out_of_stock" ? "text-destructive" :
                          item.status === "low_stock" ? "text-warning" : ""
                        }`}>
                          {item.currentStock}
                        </span>
                      </td>
                      <td className="px-2 py-3 text-right text-sm text-muted-foreground/70 tabular-nums hidden sm:table-cell">{item.reserved}</td>
                      <td className="px-2 py-3 text-right text-sm font-medium tabular-nums hidden sm:table-cell">{item.available}</td>
                      <td className="px-2 py-3"><StatusBadge status={item.status} /></td>
                      <td className="px-5 py-3 text-sm text-muted-foreground/70 hidden lg:table-cell">{item.warehouse}</td>
                    </tr>
                  ))}
                </tbody>
              </table>

              {totalPages > 1 && (
                <div className="flex items-center justify-between px-5 py-3 border-t">
                  <p className="text-xs text-muted-foreground/60">
                    {(page - 1) * perPage + 1}–{Math.min(page * perPage, filtered.length)} of {filtered.length}
                  </p>
                  <Pagination currentPage={page} totalPages={totalPages} onPageChange={setPage} />
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="low" className="mt-4">
          <Card>
            <CardContent className="p-6">
              {lowStockItems.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-16">
                  <CheckCircle2 className="h-12 w-12 text-success/60" />
                  <h3 className="mt-4 text-base font-semibold">No low stock items</h3>
                  <p className="mt-1 text-sm text-muted-foreground">All products are well stocked</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {lowStockItems.map((item) => (
                    <div key={item.id} className="flex items-center gap-4 p-3 rounded-xl bg-muted/50">
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium">{item.productName}</p>
                        <p className="text-xs text-muted-foreground">{item.productSku}</p>
                      </div>
                      <div className="text-right">
                        <p className={`text-sm font-medium ${item.currentStock === 0 ? "text-destructive" : "text-warning"}`}>
                          {item.currentStock} / {item.lowStockThreshold}
                        </p>
                        {item.incoming > 0 && (
                          <p className="text-xs text-success">+{item.incoming} incoming</p>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="out" className="mt-4">
          <Card>
            <CardContent className="p-6">
              {outOfStockItems.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-16">
                  <CheckCircle2 className="h-12 w-12 text-success/60" />
                  <h3 className="mt-4 text-base font-semibold">No out of stock items</h3>
                  <p className="mt-1 text-sm text-muted-foreground">All products are in stock</p>
                </div>
              ) : (
                <div className="space-y-3">
                  {outOfStockItems.map((item) => (
                    <div key={item.id} className="flex items-center gap-4 p-3 rounded-xl bg-muted/50">
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium">{item.productName}</p>
                        <p className="text-xs text-muted-foreground">{item.productSku}</p>
                      </div>
                      <Badge variant="destructive">Out of stock</Badge>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="incoming" className="mt-4">
          <Card>
            <CardContent className="p-6">
              {incoming.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-16">
                  <Truck className="h-12 w-12 text-muted-foreground/40" />
                  <h3 className="mt-4 text-base font-semibold">No incoming shipments</h3>
                </div>
              ) : (
                <div className="space-y-4">
                  {incoming.map((shipment) => (
                    <div key={shipment.id} className="flex items-center gap-4 p-3 rounded-xl border">
                      <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-muted">
                        <Truck className="h-5 w-5 text-muted-foreground" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium">{shipment.productName}</p>
                        <p className="text-xs text-muted-foreground">{shipment.supplier}</p>
                      </div>
                      <div className="text-right">
                        <p className="text-sm font-medium">+{shipment.quantity}</p>
                        <p className="text-xs text-muted-foreground">{formatDate(shipment.expectedDate, "short")}</p>
                      </div>
                      <StatusBadge status={shipment.status} />
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="movements" className="mt-4">
          <Card>
            <CardContent className="p-6">
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b text-left">
                      <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Product</th>
                      <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Type</th>
                      <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground text-right">Qty</th>
                      <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Reference</th>
                      <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Note</th>
                      <th className="pb-3 text-xs font-medium text-muted-foreground">Date</th>
                    </tr>
                  </thead>
                  <tbody>
                    {movements.slice(0, 20).map((mov) => (
                      <tr key={mov.id} className="border-b last:border-0">
                        <td className="py-3 pr-4 text-sm font-medium">{mov.productName}</td>
                        <td className="py-3 pr-4">
                          <Badge variant={
                            mov.type === "in" ? "success" :
                            mov.type === "out" ? "destructive" :
                            mov.type === "return" ? "info" : "warning"
                          } className="capitalize">
                            {mov.type}
                          </Badge>
                        </td>
                        <td className="py-3 pr-4 text-right">
                          <span className={`text-sm font-medium ${
                            mov.type === "in" || mov.type === "return" ? "text-success" : "text-destructive"
                          }`}>
                            {mov.type === "in" || mov.type === "return" ? "+" : ""}{mov.quantity}
                          </span>
                        </td>
                        <td className="py-3 pr-4 text-sm text-muted-foreground font-mono">{mov.reference}</td>
                        <td className="py-3 pr-4 text-sm text-muted-foreground max-w-[200px] truncate">{mov.note}</td>
                        <td className="py-3 text-sm text-muted-foreground">{formatDate(mov.date, "short")}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
