"use client"

import { useEffect, useState, useMemo } from "react"
import { useRouter } from "next/navigation"
import { Search, Users, Mail, Phone, ShoppingBag, DollarSign } from "lucide-react"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
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
import { TableSkeleton } from "@/components/shared/loading-skeleton"
import { formatCurrency, formatDate, getInitials } from "@/lib/utils"
import { getCustomers } from "@/lib/api/customers"
import type { Customer } from "@/lib/types"

export default function CustomersPage() {
  const router = useRouter()
  const [customers, setCustomers] = useState<Customer[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState("")
  const [statusFilter, setStatusFilter] = useState("all")
  const [page, setPage] = useState(1)
  const perPage = 10

  useEffect(() => {
    async function load() {
      const c = await getCustomers()
      setCustomers(c)
      setLoading(false)
    }
    load()
  }, [])

  const filtered = useMemo(() => {
    let result = customers
    if (search) {
      const q = search.toLowerCase()
      result = result.filter((c) => c.name.toLowerCase().includes(q) || c.email.toLowerCase().includes(q))
    }
    if (statusFilter !== "all") result = result.filter((c) => c.status === statusFilter)
    return result
  }, [customers, search, statusFilter])

  const totalPages = Math.ceil(filtered.length / perPage)
  const paginated = filtered.slice((page - 1) * perPage, page * perPage)

  if (loading) {
    return (
      <div className="p-6">
        <PageHeader title="Customers" description="View your customer base" />
        <Card><CardContent className="p-6"><TableSkeleton count={10} /></CardContent></Card>
      </div>
    )
  }

  return (
    <div className="p-6">
      <PageHeader title="Customers" description={`${customers.length} registered customers`} />

      <Card>
        <div className="flex flex-col sm:flex-row sm:items-center gap-3 px-5 pt-5 pb-0">
          <div className="relative flex-1 max-w-xs">
            <Search className="absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground/40" />
            <Input
              placeholder="Search customers..."
              value={search}
              onChange={(e) => { setSearch(e.target.value); setPage(1) }}
              className="h-8 pl-9 text-xs"
            />
          </div>
          <Select value={statusFilter} onValueChange={(v) => { setStatusFilter(v); setPage(1) }}>
            <SelectTrigger className="h-8 w-[120px] text-xs">
              <SelectValue placeholder="Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All status</SelectItem>
              <SelectItem value="active">Active</SelectItem>
              <SelectItem value="inactive">Inactive</SelectItem>
              <SelectItem value="blocked">Blocked</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <CardContent className="p-0 mt-4">
          <table className="w-full">
            <thead>
              <tr className="border-y text-left">
                <th className="px-5 py-3 text-xs font-medium text-muted-foreground/70 uppercase tracking-wider">Customer</th>
                <th className="px-2 py-3 text-xs font-medium text-muted-foreground/70 uppercase tracking-wider">Status</th>
                <th className="px-2 py-3 text-xs font-medium text-muted-foreground/70 uppercase tracking-wider text-right">Orders</th>
                <th className="px-2 py-3 text-xs font-medium text-muted-foreground/70 uppercase tracking-wider text-right">Spent</th>
                <th className="px-2 py-3 text-xs font-medium text-muted-foreground/70 uppercase tracking-wider text-right hidden md:table-cell">AOV</th>
                <th className="px-5 py-3 text-xs font-medium text-muted-foreground/70 uppercase tracking-wider">Last order</th>
              </tr>
            </thead>
            <tbody>
              {paginated.map((customer) => (
                <tr
                  key={customer.id}
                  className="border-b last:border-0 hover:bg-muted/20 transition-colors cursor-pointer"
                  onClick={() => router.push(`/admin/customers/${customer.id}`)}
                >
                  <td className="px-5 py-3">
                    <div className="flex items-center gap-3">
                      <Avatar className="h-8 w-8 shrink-0">
                        <AvatarImage src={customer.avatar} />
                        <AvatarFallback className="text-[10px]">{getInitials(customer.name)}</AvatarFallback>
                      </Avatar>
                      <div>
                        <p className="text-sm font-medium">{customer.name}</p>
                        <p className="text-xs text-muted-foreground/70">{customer.email}</p>
                      </div>
                    </div>
                  </td>
                  <td className="px-2 py-3"><StatusBadge status={customer.status} /></td>
                  <td className="px-2 py-3 text-right text-sm tabular-nums">{customer.totalOrders}</td>
                  <td className="px-2 py-3 text-right text-sm font-medium tabular-nums">{formatCurrency(customer.totalSpent)}</td>
                  <td className="px-2 py-3 text-right text-sm text-muted-foreground/70 tabular-nums hidden md:table-cell">{formatCurrency(customer.averageOrderValue)}</td>
                  <td className="px-5 py-3 text-sm text-muted-foreground/70">{formatDate(customer.lastOrderDate, "short")}</td>
                </tr>
              ))}
            </tbody>
          </table>

          {paginated.length === 0 && (
            <div className="flex flex-col items-center justify-center py-16 px-6">
              <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-muted">
                <Users className="h-7 w-7 text-muted-foreground/40" />
              </div>
              <h3 className="mt-4 text-base font-semibold">No customers found</h3>
              <p className="mt-1 text-sm text-muted-foreground/70">Try adjusting your search.</p>
            </div>
          )}

          {totalPages > 1 && paginated.length > 0 && (
            <div className="flex items-center justify-between px-5 py-3 border-t">
              <p className="text-xs text-muted-foreground/60">
                {(page - 1) * perPage + 1}–{Math.min(page * perPage, filtered.length)} of {filtered.length}
              </p>
              <Pagination currentPage={page} totalPages={totalPages} onPageChange={setPage} />
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
