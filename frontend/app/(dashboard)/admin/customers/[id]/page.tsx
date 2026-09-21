"use client"

import { useEffect, useState } from "react"
import { useParams, useRouter } from "next/navigation"
import { ChevronLeft, Mail, Phone, MapPin, ShoppingBag, DollarSign, Clock, Tag, Star } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Separator } from "@/components/ui/separator"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { StatusBadge } from "@/components/shared/status-badge"
import { Skeleton } from "@/components/ui/skeleton"
import { TableSkeleton } from "@/components/shared/loading-skeleton"
import { formatCurrency, formatDate, formatRelativeTime, getInitials } from "@/lib/utils"
import { getCustomer } from "@/lib/api/customers"
import { getOrders } from "@/lib/api/orders"
import type { Customer, Order } from "@/lib/types"

export default function CustomerDetailPage() {
  const params = useParams()
  const router = useRouter()
  const [customer, setCustomer] = useState<Customer | null>(null)
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      const c = await getCustomer(params.id as string)
      if (!c) { setLoading(false); return }
      const allOrders = await getOrders()
      const customerOrders = allOrders.filter((o) => o.customer.id === c.id)
      setCustomer(c)
      setOrders(customerOrders)
      setLoading(false)
    }
    load()
  }, [params.id])

  if (loading) {
    return (
      <div className="p-6">
        <Skeleton className="h-8 w-24" />
        <div className="mt-6 grid gap-6 lg:grid-cols-3">
          <Skeleton className="h-64 rounded-xl lg:col-span-1" />
          <Skeleton className="h-64 rounded-xl lg:col-span-2" />
        </div>
      </div>
    )
  }

  if (!customer) {
    return (
      <div className="flex flex-col items-center justify-center py-24">
        <h2 className="text-xl font-semibold">Customer not found</h2>
        <Button variant="outline" className="mt-4" onClick={() => router.push("/admin/customers")}>
          Back to customers
        </Button>
      </div>
    )
  }

  const defaultAddress = customer.addresses.find((a) => a.isDefault) || customer.addresses[0]

  return (
    <div className="p-6">
      <div className="flex items-center gap-4 mb-8">
        <Button variant="ghost" size="icon" onClick={() => router.push("/admin/customers")}>
          <ChevronLeft className="h-4 w-4" />
        </Button>
        <h1 className="text-2xl font-semibold tracking-tight">Customer Profile</h1>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <div className="space-y-6">
          <Card>
            <CardContent className="p-6 text-center">
              <Avatar className="mx-auto h-20 w-20">
                <AvatarImage src={customer.avatar} />
                <AvatarFallback className="text-lg">{getInitials(customer.name)}</AvatarFallback>
              </Avatar>
              <h2 className="mt-4 text-lg font-semibold">{customer.name}</h2>
              <p className="text-sm text-muted-foreground">{customer.email}</p>
              <div className="mt-3">
                <StatusBadge status={customer.status} />
              </div>
              {customer.tags.length > 0 && (
                <div className="mt-3 flex flex-wrap justify-center gap-1">
                  {customer.tags.map((tag) => (
                    <Badge key={tag} variant="secondary" className="font-normal text-xs">{tag}</Badge>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="text-sm">Contact</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex items-center gap-2 text-sm">
                <Mail className="h-4 w-4 text-muted-foreground" />
                <span>{customer.email}</span>
              </div>
              <div className="flex items-center gap-2 text-sm">
                <Phone className="h-4 w-4 text-muted-foreground" />
                <span>{customer.phone}</span>
              </div>
            </CardContent>
          </Card>

          {defaultAddress && (
            <Card>
              <CardHeader>
                <CardTitle className="text-sm flex items-center gap-2">
                  <MapPin className="h-4 w-4 text-muted-foreground" />
                  Default Address
                </CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground space-y-0.5">
                <p>{defaultAddress.line1}</p>
                {defaultAddress.line2 && <p>{defaultAddress.line2}</p>}
                <p>{defaultAddress.city}, {defaultAddress.state} {defaultAddress.zip}</p>
                <p>{defaultAddress.country}</p>
              </CardContent>
            </Card>
          )}
        </div>

        <div className="lg:col-span-2 space-y-6">
          <div className="grid grid-cols-3 gap-4">
            <Card>
              <CardContent className="p-4 text-center">
                <p className="text-xs text-muted-foreground">Total Orders</p>
                <p className="mt-1 text-2xl font-bold">{customer.totalOrders}</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 text-center">
                <p className="text-xs text-muted-foreground">Lifetime Value</p>
                <p className="mt-1 text-2xl font-bold">{formatCurrency(customer.lifetimeValue)}</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 text-center">
                <p className="text-xs text-muted-foreground">Avg Order Value</p>
                <p className="mt-1 text-2xl font-bold">{formatCurrency(customer.averageOrderValue)}</p>
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Recent Activity</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-0">
                {customer.recentActivity.slice(0, 5).map((event) => (
                  <div key={event.id} className="flex items-start gap-3 py-3 border-b last:border-0">
                    <div className="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-muted">
                      <Clock className="h-4 w-4 text-muted-foreground" />
                    </div>
                    <div className="flex-1">
                      <p className="text-sm">{event.description}</p>
                      <p className="text-xs text-muted-foreground">{formatRelativeTime(event.timestamp)}</p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Purchase History ({orders.length} orders)</CardTitle>
            </CardHeader>
            <CardContent>
              {orders.length === 0 ? (
                <p className="text-sm text-muted-foreground text-center py-8">No orders yet</p>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b text-left">
                        <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Order</th>
                        <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Status</th>
                        <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground text-right">Total</th>
                        <th className="pb-3 text-xs font-medium text-muted-foreground">Date</th>
                      </tr>
                    </thead>
                    <tbody>
                      {orders.slice(0, 10).map((order) => (
                        <tr key={order.id} className="border-b last:border-0">
                          <td className="py-3 pr-4 text-sm font-medium font-mono">{order.orderNumber}</td>
                          <td className="py-3 pr-4"><StatusBadge status={order.status} /></td>
                          <td className="py-3 pr-4 text-right text-sm font-medium">{formatCurrency(order.total)}</td>
                          <td className="py-3 text-sm text-muted-foreground">{formatDate(order.createdAt, "short")}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
