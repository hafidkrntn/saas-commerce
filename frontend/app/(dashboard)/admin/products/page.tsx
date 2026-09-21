"use client"

import { useEffect, useState, useMemo } from "react"
import { useRouter } from "next/navigation"
import {
  Plus,
  Search,
  MoreHorizontal,
  Edit,
  Trash2,
  Copy,
  Eye,
  Package,
  X,
} from "lucide-react"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Checkbox } from "@/components/ui/checkbox"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
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
import { formatCurrency } from "@/lib/utils"
import { getProducts, getCategories, deleteProduct } from "@/lib/api/products"
import type { Product } from "@/lib/types"
import { toast } from "sonner"

export default function ProductsPage() {
  const router = useRouter()
  const [products, setProducts] = useState<Product[]>([])
  const [categories, setCategories] = useState<string[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState("")
  const [categoryFilter, setCategoryFilter] = useState("all")
  const [stockFilter, setStockFilter] = useState("all")
  const [page, setPage] = useState(1)
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set())
  const perPage = 12

  useEffect(() => {
    async function load() {
      const [p, c] = await Promise.all([getProducts(), getCategories()])
      setProducts(p)
      setCategories(c)
      setLoading(false)
    }
    load()
  }, [])

  const handleDelete = async (id: string) => {
    try {
      await deleteProduct(id)
      setProducts((prev) => prev.filter((p) => p.id !== id))
      setSelectedIds((prev) => {
        const next = new Set(prev)
        next.delete(id)
        return next
      })
      toast.success("Product deleted")
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete product")
    }
  }

  const filtered = useMemo(() => {
    let result = products
    if (search) {
      const q = search.toLowerCase()
      result = result.filter((p) => p.name.toLowerCase().includes(q) || p.sku.toLowerCase().includes(q))
    }
    if (categoryFilter !== "all") result = result.filter((p) => p.category === categoryFilter)
    if (stockFilter === "in_stock") result = result.filter((p) => p.stock > 0)
    else if (stockFilter === "low") result = result.filter((p) => p.stock > 0 && p.stock <= p.lowStockThreshold)
    else if (stockFilter === "out") result = result.filter((p) => p.stock === 0)
    return result
  }, [products, search, categoryFilter, stockFilter])

  const totalPages = Math.ceil(filtered.length / perPage)
  const paginated = filtered.slice((page - 1) * perPage, page * perPage)

  const toggleAll = () => {
    setSelectedIds(selectedIds.size === paginated.length ? new Set() : new Set(paginated.map((p) => p.id)))
  }
  const toggleOne = (id: string) => {
    const next = new Set(selectedIds)
    next.has(id) ? next.delete(id) : next.add(id)
    setSelectedIds(next)
  }

  const hasFilters = search || categoryFilter !== "all" || stockFilter !== "all"
  const clearFilters = () => { setSearch(""); setCategoryFilter("all"); setStockFilter("all"); setPage(1) }

  if (loading) {
    return (
      <div className="p-6 lg:p-8">
        <PageHeader title="Products" description="Manage your product catalog">
          <Button disabled><Plus className="h-4 w-4" /> Add product</Button>
        </PageHeader>
        <Card><CardContent className="p-0"><TableSkeleton count={8} /></CardContent></Card>
      </div>
    )
  }

  return (
    <div className="p-6 lg:p-8">
      <PageHeader title="Products" description={`${products.length} products`}>
        <Button onClick={() => router.push("/admin/products/new")}>
          <Plus className="h-4 w-4" /> Add product
        </Button>
      </PageHeader>

      <Card>
        <div className="flex flex-col sm:flex-row sm:items-center gap-3 px-5 pt-5 pb-0">
          <div className="relative flex-1 max-w-xs">
            <Search className="absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground/40" />
            <Input
              placeholder="Search products..."
              value={search}
              onChange={(e) => { setSearch(e.target.value); setPage(1) }}
              className="h-8 pl-9 text-xs"
            />
          </div>
          <div className="flex items-center gap-2">
            <Select value={categoryFilter} onValueChange={(v) => { setCategoryFilter(v); setPage(1) }}>
              <SelectTrigger className="h-8 w-[130px] text-xs">
                <SelectValue placeholder="Category" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All categories</SelectItem>
                {categories.map((c) => <SelectItem key={c} value={c}>{c}</SelectItem>)}
              </SelectContent>
            </Select>
            <Select value={stockFilter} onValueChange={(v) => { setStockFilter(v); setPage(1) }}>
              <SelectTrigger className="h-8 w-[110px] text-xs">
                <SelectValue placeholder="Stock" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All stock</SelectItem>
                <SelectItem value="in_stock">In stock</SelectItem>
                <SelectItem value="low">Low stock</SelectItem>
                <SelectItem value="out">Out of stock</SelectItem>
              </SelectContent>
            </Select>
            {hasFilters && (
              <Button variant="ghost" size="icon-sm" onClick={clearFilters}>
                <X className="h-3.5 w-3.5" />
              </Button>
            )}
          </div>
          {selectedIds.size > 0 && (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm" className="h-8 text-xs">{selectedIds.size} selected</Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-36">
                <DropdownMenuItem className="text-destructive" onClick={() => setSelectedIds(new Set())}>
                  <Trash2 className="mr-2 h-3.5 w-3.5" /> Delete
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          )}
        </div>

        <CardContent className="p-0 mt-4">
          <table className="w-full">
            <thead>
              <tr className="border-y text-left">
                <th className="w-10 px-5 py-2.5">
                  <Checkbox checked={paginated.length > 0 && selectedIds.size === paginated.length} onCheckedChange={toggleAll} />
                </th>
                <th className="px-2 py-2.5 text-[11px] font-medium text-muted-foreground/60 uppercase tracking-wider">Product</th>
                <th className="px-2 py-2.5 text-[11px] font-medium text-muted-foreground/60 uppercase tracking-wider hidden md:table-cell">SKU</th>
                <th className="px-2 py-2.5 text-[11px] font-medium text-muted-foreground/60 uppercase tracking-wider hidden lg:table-cell">Category</th>
                <th className="px-2 py-2.5 text-[11px] font-medium text-muted-foreground/60 uppercase tracking-wider text-right">Price</th>
                <th className="px-2 py-2.5 text-[11px] font-medium text-muted-foreground/60 uppercase tracking-wider text-right">Stock</th>
                <th className="px-2 py-2.5 text-[11px] font-medium text-muted-foreground/60 uppercase tracking-wider hidden sm:table-cell">Status</th>
                <th className="w-10 px-5 py-2.5" />
              </tr>
            </thead>
            <tbody>
              {paginated.map((product) => (
                <tr
                  key={product.id}
                  className="border-b last:border-0 hover:bg-muted/20 transition-colors cursor-pointer"
                  onClick={() => router.push(`/admin/products/${product.id}`)}
                >
                  <td className="px-5 py-2.5" onClick={(e) => e.stopPropagation()}>
                    <Checkbox checked={selectedIds.has(product.id)} onCheckedChange={() => toggleOne(product.id)} />
                  </td>
                  <td className="px-2 py-2.5">
                    <div className="flex items-center gap-3">
                      <div className="flex h-[32px] w-[32px] shrink-0 items-center justify-center rounded-lg bg-muted overflow-hidden">
                        {product.images[0] ? (
                          // eslint-disable-next-line @next/next/no-img-element
                          <img src={product.images[0]} alt={product.name} className="h-full w-full object-cover" />
                        ) : (
                          <Package className="h-3.5 w-3.5 text-muted-foreground/50" />
                        )}
                      </div>
                      <div className="min-w-0">
                        <p className="text-sm font-medium truncate max-w-[200px]">{product.name}</p>
                        <div className="flex items-center gap-1.5 mt-0.5">
                          <span className="text-[11px] text-amber-400">{"★".repeat(Math.round(product.rating))}{"☆".repeat(5 - Math.round(product.rating))}</span>
                          <span className="text-[11px] text-muted-foreground/50">({product.reviewCount})</span>
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="px-2 py-2.5 hidden md:table-cell">
                    <span className="text-xs text-muted-foreground/60 font-mono">{product.sku}</span>
                  </td>
                  <td className="px-2 py-2.5 hidden lg:table-cell">
                    <Badge variant="secondary" className="font-normal text-[11px]">{product.category}</Badge>
                  </td>
                  <td className="px-2 py-2.5 text-right">
                    <span className="text-sm font-medium tabular-nums">{formatCurrency(product.price)}</span>
                    {product.compareAtPrice && (
                      <span className="ml-1 text-[11px] text-muted-foreground/40 line-through tabular-nums">{formatCurrency(product.compareAtPrice)}</span>
                    )}
                  </td>
                  <td className="px-2 py-2.5 text-right">
                    <span className={`text-sm font-medium tabular-nums ${product.stock === 0 ? "text-destructive" : product.stock <= product.lowStockThreshold ? "text-warning" : ""}`}>
                      {product.stock}
                    </span>
                    {product.incoming > 0 && (
                      <span className="ml-1 text-[11px] text-success tabular-nums">+{product.incoming}</span>
                    )}
                  </td>
                  <td className="px-2 py-2.5 hidden sm:table-cell">
                    <StatusBadge status={product.status} />
                  </td>
                  <td className="px-5 py-2.5 text-right" onClick={(e) => e.stopPropagation()}>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon-sm" className="h-7 w-7">
                          <MoreHorizontal className="h-3.5 w-3.5" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end" className="w-36">
                        <DropdownMenuItem onClick={() => router.push(`/admin/products/${product.id}`)}>
                          <Eye className="mr-2 h-3.5 w-3.5" /> View
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => router.push(`/admin/products/${product.id}/edit`)}>
                          <Edit className="mr-2 h-3.5 w-3.5" /> Edit
                        </DropdownMenuItem>
                        <DropdownMenuItem>
                          <Copy className="mr-2 h-3.5 w-3.5" /> Duplicate
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem className="text-destructive" onClick={() => handleDelete(product.id)}>
                          <Trash2 className="mr-2 h-3.5 w-3.5" /> Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {paginated.length === 0 && (
            <div className="flex flex-col items-center justify-center py-16 px-6">
              <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-muted">
                <Package className="h-6 w-6 text-muted-foreground/40" />
              </div>
              <h3 className="mt-4 text-base font-semibold">No products found</h3>
              <p className="mt-1 text-sm text-muted-foreground/60 text-center max-w-sm">
                {hasFilters ? "Try adjusting your search or filters." : "Add your first product to get started."}
              </p>
              {hasFilters ? (
                <Button variant="outline" size="sm" className="mt-5" onClick={clearFilters}>Clear filters</Button>
              ) : (
                <Button size="sm" className="mt-5" onClick={() => router.push("/admin/products/new")}>
                  <Plus className="h-4 w-4" /> Add product
                </Button>
              )}
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
    </div>
  )
}
