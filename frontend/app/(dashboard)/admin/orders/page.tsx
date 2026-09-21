"use client"

import { useEffect, useState, useMemo } from "react"
import {
  Search,
  ShoppingBag,
  X,
  ChevronLeft,
  Package,
  Truck,
  CheckCircle2,
  Clock,
  AlertCircle,
  CreditCard,
  MapPin,
  MessageSquare,
  ArrowRight,
} from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Separator } from "@/components/ui/separator"
import { ScrollArea } from "@/components/ui/scroll-area"
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
import { toast } from "sonner"
import { formatCurrency, formatDate, formatRelativeTime, getInitials } from "@/lib/utils"
import { getOrders, getOrder, updateOrderStatus } from "@/lib/api/orders"
import type { Order } from "@/lib/types"

export default function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState("")
  const [statusFilter, setStatusFilter] = useState("all")
  const [page, setPage] = useState(1)
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null)
  const perPage = 12

  useEffect(() => {
    async function load() {
      const o = await getOrders()
      setOrders(o)
      setLoading(false)
    }
    load()
  }, [])

  const filtered = useMemo(() => {
    let result = orders
    if (search) {
      const q = search.toLowerCase()
      result = result.filter(
        (o) => o.orderNumber.toLowerCase().includes(q) || o.customer.name.toLowerCase().includes(q)
      )
    }
    if (statusFilter !== "all") result = result.filter((o) => o.status === statusFilter)
    return result
  }, [orders, search, statusFilter])

  const totalPages = Math.ceil(filtered.length / perPage)
  const paginated = filtered.slice((page - 1) * perPage, page * perPage)

  return (
    <div className="p-6 lg:p-8">
      <PageHeader title="Orders" description={`${orders.length} orders total`} />

      <Card>
        <div className="flex flex-col sm:flex-row sm:items-center gap-3 px-5 pt-5 pb-0">
          <div className="relative flex-1 max-w-xs">
            <Search className="absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground/40" />
            <Input
              placeholder="Search orders..."
              value={search}
              onChange={(e) => { setSearch(e.target.value); setPage(1) }}
              className="h-8 pl-9 text-xs"
            />
          </div>
          <Select value={statusFilter} onValueChange={(v) => { setStatusFilter(v); setPage(1) }}>
            <SelectTrigger className="h-8 w-[130px] text-xs">
              <SelectValue placeholder="Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All status</SelectItem>
              <SelectItem value="pending">Pending</SelectItem>
              <SelectItem value="confirmed">Confirmed</SelectItem>
              <SelectItem value="processing">Processing</SelectItem>
              <SelectItem value="shipped">Shipped</SelectItem>
              <SelectItem value="delivered">Delivered</SelectItem>
              <SelectItem value="cancelled">Cancelled</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <CardContent className="p-0 mt-4">
          <table className="w-full">
            <thead>
              <tr className="border-y text-left">
                <th className="px-5 py-2.5 text-[11px] font-medium text-muted-foreground/60 uppercase tracking-wider">Order</th>
                <th className="px-2 py-2.5 text-[11px] font-medium text-muted-foreground/60 uppercase tracking-wider">Customer</th>
                <th className="px-2 py-2.5 text-[11px] font-medium text-muted-foreground/60 uppercase tracking-wider hidden sm:table-cell">Status</th>
                <th className="px-2 py-2.5 text-[11px] font-medium text-muted-foreground/60 uppercase tracking-wider hidden md:table-cell">Payment</th>
                <th className="px-2 py-2.5 text-[11px] font-medium text-muted-foreground/60 uppercase tracking-wider text-right">Total</th>
                <th className="px-2 py-2.5 text-[11px] font-medium text-muted-foreground/60 uppercase tracking-wider hidden md:table-cell">Date</th>
                <th className="px-5 py-2.5 text-right text-[11px] font-medium text-muted-foreground/60 uppercase tracking-wider">Items</th>
              </tr>
            </thead>
            <tbody>
              {paginated.map((order) => (
                <tr
                  key={order.id}
                  className="border-b last:border-0 hover:bg-muted/20 transition-colors cursor-pointer"
                  onClick={async () => {
                    const full = await getOrder(order.id)
                    setSelectedOrder(full ?? order)
                  }}
                >
                  <td className="px-5 py-2.5">
                    <span className="text-sm font-medium font-mono">{order.orderNumber}</span>
                  </td>
                  <td className="px-2 py-2.5">
                    <div className="flex items-center gap-2.5">
                      <Avatar className="h-7 w-7 shrink-0">
                        <AvatarImage src={order.customer.avatar} />
                        <AvatarFallback className="text-[9px]">{getInitials(order.customer.name)}</AvatarFallback>
                      </Avatar>
                      <span className="text-sm truncate max-w-[140px]">{order.customer.name}</span>
                    </div>
                  </td>
                  <td className="px-2 py-2.5 hidden sm:table-cell">
                    <StatusBadge status={order.status} />
                  </td>
                  <td className="px-2 py-2.5 hidden md:table-cell">
                    <div className="flex items-center gap-1.5">
                      <StatusBadge status={order.paymentStatus} />
                      <span className="text-[11px] text-muted-foreground/50">{order.paymentMethod}</span>
                    </div>
                  </td>
                  <td className="px-2 py-2.5 text-right">
                    <span className="text-sm font-medium tabular-nums">{formatCurrency(order.total)}</span>
                  </td>
                  <td className="px-2 py-2.5 hidden md:table-cell">
                    <span className="text-sm text-muted-foreground/60">{formatRelativeTime(order.createdAt)}</span>
                  </td>
                  <td className="px-5 py-2.5 text-right text-sm text-muted-foreground/60">{order.items.length}</td>
                </tr>
              ))}
            </tbody>
          </table>

          {paginated.length === 0 && (
            <div className="flex flex-col items-center justify-center py-16 px-6">
              <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-muted">
                <ShoppingBag className="h-6 w-6 text-muted-foreground/40" />
              </div>
              <h3 className="mt-4 text-base font-semibold">No orders found</h3>
              <p className="mt-1 text-sm text-muted-foreground/60">Try adjusting your search or filters.</p>
            </div>
          )}
        </CardContent>

        {totalPages > 1 && paginated.length > 0 && (
          <div className="flex items-center justify-between px-5 py-3 border-t">
            <p className="text-xs text-muted-foreground/60">
              {(page - 1) * perPage + 1}–{Math.min(page * perPage, filtered.length)} of {filtered.length}
            </p>
            <Pagination currentPage={page} totalPages={totalPages} onPageChange={setPage} />
          </div>
        )}
      </Card>

      {selectedOrder && (
        <OrderDrawer
          order={selectedOrder}
          onClose={() => setSelectedOrder(null)}
          onStatusChange={async (id, status) => {
            try {
              const updated = await updateOrderStatus(id, status)
              setOrders((prev) => prev.map((o) => (o.id === id ? updated : o)))
              setSelectedOrder(updated)
              toast.success(`Order ${status}`)
            } catch (err) {
              toast.error(err instanceof Error ? err.message : "Failed to update order")
            }
          }}
        />
      )}
    </div>
  )
}

const timelineIcons: Record<string, React.ElementType> = {
  created: Clock, confirmed: CheckCircle2, processing: Package,
  shipped: Truck, delivered: CheckCircle2, cancelled: X,
  refunded: AlertCircle, payment: CreditCard, note: MessageSquare,
}

function OrderDrawer({
  order,
  onClose,
  onStatusChange,
}: {
  order: Order
  onClose: () => void
  onStatusChange: (id: string, status: string) => Promise<void>
}) {
  const [acting, setActing] = useState(false)

  const nextAction = (() => {
    switch (order.status) {
      case "pending": return { status: "confirmed", label: "Confirm order" }
      case "confirmed": return { status: "processing", label: "Start processing" }
      case "processing": return { status: "shipped", label: "Ship order" }
      case "shipped": return { status: "delivered", label: "Mark delivered" }
      default: return null
    }
  })()

  const handleAdvance = async () => {
    if (!nextAction) return
    setActing(true)
    await onStatusChange(order.id, nextAction.status)
    setActing(false)
  }

  const handleCancel = async () => {
    setActing(true)
    await onStatusChange(order.id, "cancelled")
    setActing(false)
  }
  return (
    <div className="fixed inset-0 z-50 flex justify-end">
      <div className="fixed inset-0 bg-black/40" onClick={onClose} />
      <div className="relative w-full max-w-lg bg-background border-l shadow-2xl animate-in slide-in-from-right">
        <div className="flex items-center justify-between px-6 h-[65px] border-b">
          <div className="flex items-center gap-3">
            <Button variant="ghost" size="icon-sm" onClick={onClose}>
              <ChevronLeft className="h-4 w-4" />
            </Button>
            <div>
              <h2 className="text-sm font-semibold">{order.orderNumber}</h2>
              <p className="text-xs text-muted-foreground/70">{formatDate(order.createdAt, "long")}</p>
            </div>
          </div>
          <StatusBadge status={order.status} />
        </div>

        <ScrollArea className="h-[calc(100vh-65px)]">
          <div className="p-6 space-y-6">
            <div>
              <h3 className="text-xs font-semibold text-muted-foreground/70 uppercase tracking-wider mb-3">Timeline</h3>
              <div className="space-y-0">
                {order.timeline.map((event, i) => {
                  const Icon = timelineIcons[event.type] || Clock
                  return (
                    <div key={event.id} className="flex gap-3">
                      <div className="flex flex-col items-center">
                        <div className={`flex h-6 w-6 items-center justify-center rounded-full ${
                          ["cancelled", "refunded"].includes(event.type)
                            ? "bg-destructive/10 text-destructive"
                            : event.type === "delivered"
                            ? "bg-success/10 text-success"
                            : "bg-muted text-muted-foreground/60"
                        }`}>
                          <Icon className="h-3 w-3" />
                        </div>
                        {i < order.timeline.length - 1 && <div className="w-px flex-1 bg-border" />}
                      </div>
                      <div className="pb-4 flex-1">
                        <p className="text-sm">{event.title}</p>
                        {event.description && <p className="text-xs text-muted-foreground/70 mt-0.5">{event.description}</p>}
                        <p className="text-xs text-muted-foreground/50 mt-0.5">{formatRelativeTime(event.timestamp)}</p>
                      </div>
                    </div>
                  )
                })}
              </div>
            </div>

            <Separator />

            <div>
              <h3 className="text-xs font-semibold text-muted-foreground/70 uppercase tracking-wider mb-3">Items ({order.items.length})</h3>
              <div className="space-y-2.5">
                {order.items.map((item) => (
                  <div key={item.id} className="flex items-center gap-3 p-2.5 rounded-lg bg-muted/50">
                    <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-background">
                      <Package className="h-4 w-4 text-muted-foreground/60" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium truncate">{item.productName}</p>
                      <p className="text-xs text-muted-foreground/70">{item.productSku} × {item.quantity}</p>
                    </div>
                    <p className="text-sm font-medium tabular-nums">{formatCurrency(item.subtotal)}</p>
                  </div>
                ))}
              </div>
            </div>

            <Separator />

            <div className="space-y-1.5">
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground/70">Subtotal</span>
                <span className="tabular-nums">{formatCurrency(order.subtotal)}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground/70">Shipping</span>
                <span className="tabular-nums">{order.shippingCost === 0 ? "Free" : formatCurrency(order.shippingCost)}</span>
              </div>
              {order.discount > 0 && (
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground/70">Discount</span>
                  <span className="text-success tabular-nums">-{formatCurrency(order.discount)}</span>
                </div>
              )}
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground/70">Tax</span>
                <span className="tabular-nums">{formatCurrency(order.tax)}</span>
              </div>
              <Separator />
              <div className="flex justify-between font-semibold text-base">
                <span>Total</span>
                <span className="tabular-nums">{formatCurrency(order.total)}</span>
              </div>
            </div>

            <Separator />

            <div className="space-y-3">
              <div>
                <h3 className="text-xs font-semibold text-muted-foreground/70 uppercase tracking-wider mb-2">Customer</h3>
                <div className="flex items-center gap-3">
                  <Avatar className="h-9 w-9">
                    <AvatarImage src={order.customer.avatar} />
                    <AvatarFallback className="text-xs">{getInitials(order.customer.name)}</AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="text-sm font-medium">{order.customer.name}</p>
                    <p className="text-xs text-muted-foreground/70">{order.customer.email}</p>
                  </div>
                </div>
              </div>

              <div>
                <h3 className="text-xs font-semibold text-muted-foreground/70 uppercase tracking-wider mb-2 flex items-center gap-1.5">
                  <MapPin className="h-3.5 w-3.5" />
                  Shipping
                </h3>
                <div className="text-sm text-muted-foreground/80 space-y-0.5 bg-muted/50 rounded-lg p-3">
                  <p>{order.shippingAddress.line1}</p>
                  {order.shippingAddress.line2 && <p>{order.shippingAddress.line2}</p>}
                  <p>{order.shippingAddress.city}, {order.shippingAddress.state} {order.shippingAddress.zip}</p>
                  <p>{order.shippingAddress.country}</p>
                </div>
              </div>
            </div>

            <div className="flex gap-2 pb-6">
              {nextAction ? (
                <Button className="flex-1 text-sm h-9" onClick={handleAdvance} disabled={acting}>
                  <Truck className="h-4 w-4" /> {nextAction.label}
                </Button>
              ) : (
                <Button className="flex-1 text-sm h-9" disabled>
                  {order.status === "delivered" ? "Completed" : order.status}
                </Button>
              )}
              {!["cancelled", "refunded", "delivered"].includes(order.status) && (
                <Button variant="outline" className="flex-1 text-sm h-9" onClick={handleCancel} disabled={acting}>
                  <X className="h-4 w-4" /> Cancel
                </Button>
              )}
            </div>
          </div>
        </ScrollArea>
      </div>
    </div>
  )
}
