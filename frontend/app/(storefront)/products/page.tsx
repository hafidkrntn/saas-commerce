"use client"

import { useEffect, useState, useMemo } from "react"
import Link from "next/link"
import { useSearchParams } from "next/navigation"
import { Package, Star, SlidersHorizontal, ArrowUpDown, X } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Pagination } from "@/components/ui/pagination"
import { formatCurrency } from "@/lib/utils"
import { getProducts, getCategories } from "@/lib/api/products"
import type { Product } from "@/lib/types"

export default function CatalogPage() {
  const searchParams = useSearchParams()
  const [products, setProducts] = useState<Product[]>([])
  const [categories, setCategories] = useState<string[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState(searchParams.get("q") || "")
  const [category, setCategory] = useState(searchParams.get("category") || "all")
  const [sort, setSort] = useState("featured")
  const [page, setPage] = useState(1)
  const perPage = 12

  useEffect(() => {
    async function load() {
      const [p, c] = await Promise.all([getProducts(), getCategories()])
      setProducts(p.filter((p) => p.status === "active"))
      setCategories(c)
      setLoading(false)
    }
    load()
  }, [])

  const filtered = useMemo(() => {
    let result = [...products]
    if (search) {
      const q = search.toLowerCase()
      result = result.filter((p) => p.name.toLowerCase().includes(q) || p.description.toLowerCase().includes(q))
    }
    if (category !== "all") result = result.filter((p) => p.category === category)

    switch (sort) {
      case "price-asc": result.sort((a, b) => a.price - b.price); break
      case "price-desc": result.sort((a, b) => b.price - a.price); break
      case "rating": result.sort((a, b) => b.rating - a.rating); break
      case "newest": result.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()); break
      default: break
    }

    return result
  }, [products, search, category, sort])

  const totalPages = Math.ceil(filtered.length / perPage)
  const paginated = filtered.slice((page - 1) * perPage, page * perPage)

  const hasFilters = search || category !== "all"
  const clearFilters = () => { setSearch(""); setCategory("all"); setPage(1) }

  if (loading) {
    return (
      <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="h-8 w-48 rounded-lg bg-muted animate-pulse" />
        <div className="mt-8 grid grid-cols-2 gap-3 sm:gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
          {Array.from({ length: 8 }).map((_, i) => (
            <div key={i} className="rounded-xl border">
              <div className="aspect-square bg-muted animate-pulse" />
              <div className="p-4 space-y-2">
                <div className="h-3 w-16 rounded bg-muted animate-pulse" />
                <div className="h-4 w-3/4 rounded bg-muted animate-pulse" />
                <div className="h-3 w-20 rounded bg-muted animate-pulse" />
              </div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">All Products</h1>
          <p className="mt-1 text-sm text-muted-foreground">{filtered.length} products</p>
        </div>
        <div className="flex items-center gap-2">
          <div className="relative">
            <input
              placeholder="Search..."
              value={search}
              onChange={(e) => { setSearch(e.target.value); setPage(1) }}
              className="h-9 w-40 rounded-lg border bg-background px-3 text-xs ring-offset-background placeholder:text-muted-foreground/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            />
          </div>
          <Select value={category} onValueChange={(v) => { setCategory(v); setPage(1) }}>
            <SelectTrigger className="h-9 w-[130px] text-xs">
              <SelectValue placeholder="Category" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All categories</SelectItem>
              {categories.map((c) => <SelectItem key={c} value={c}>{c}</SelectItem>)}
            </SelectContent>
          </Select>
          <Select value={sort} onValueChange={(v) => setSort(v)}>
            <SelectTrigger className="h-9 w-[120px] text-xs">
              <SelectValue placeholder="Sort" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="featured">Featured</SelectItem>
              <SelectItem value="price-asc">Price: Low</SelectItem>
              <SelectItem value="price-desc">Price: High</SelectItem>
              <SelectItem value="rating">Rating</SelectItem>
              <SelectItem value="newest">Newest</SelectItem>
            </SelectContent>
          </Select>
          {hasFilters && (
            <Button variant="ghost" size="icon-sm" onClick={clearFilters}>
              <X className="h-3.5 w-3.5" />
            </Button>
          )}
        </div>
      </div>

      {paginated.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20">
          <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-muted">
            <Package className="h-7 w-7 text-muted-foreground/40" />
          </div>
          <h3 className="mt-4 text-base font-semibold">No products found</h3>
          <p className="mt-1 text-sm text-muted-foreground">Try adjusting your search or filters.</p>
          {hasFilters && (
            <Button variant="outline" size="sm" className="mt-5" onClick={clearFilters}>
              Clear filters
            </Button>
          )}
        </div>
      ) : (
        <>
          <div className="grid grid-cols-2 gap-3 sm:gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
            {paginated.map((product) => (
              <Link key={product.id} href={`/products/${product.id}`} className="group">
                <Card className="overflow-hidden transition-shadow hover:shadow-md h-full">
                  <div className="relative aspect-square bg-muted flex items-center justify-center">
                    {product.images[0] ? (
                      // eslint-disable-next-line @next/next/no-img-element
                      <img src={product.images[0]} alt={product.name} className="h-full w-full object-cover" />
                    ) : (
                      <Package className="h-16 w-16 text-muted-foreground/20" />
                    )}
                    {product.compareAtPrice && (
                      <Badge className="absolute left-3 top-3 text-xs px-2 py-0.5">
                        -{Math.round((1 - product.price / product.compareAtPrice) * 100)}%
                      </Badge>
                    )}
                    {product.stock === 0 && (
                      <div className="absolute inset-0 bg-background/60 flex items-center justify-center">
                        <Badge variant="outline" className="text-sm">Out of stock</Badge>
                      </div>
                    )}
                  </div>
                  <CardContent className="p-4">
                    <p className="text-xs text-muted-foreground mb-1">{product.category}</p>
                    <h3 className="text-sm font-medium group-hover:text-primary transition-colors truncate">
                      {product.name}
                    </h3>
                    <div className="mt-1 flex items-center gap-0.5">
                      {Array.from({ length: 5 }).map((_, i) => (
                        <Star
                          key={i}
                          className={`h-3 w-3 ${i < Math.round(product.rating) ? "fill-amber-400 text-amber-400" : "text-muted-foreground/20"}`}
                        />
                      ))}
                      <span className="ml-1 text-xs text-muted-foreground">({product.reviewCount})</span>
                    </div>
                    <div className="mt-2 flex items-center gap-2">
                      <span className="text-base font-semibold">{formatCurrency(product.price)}</span>
                      {product.compareAtPrice && (
                        <span className="text-sm text-muted-foreground line-through">{formatCurrency(product.compareAtPrice)}</span>
                      )}
                    </div>
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>

          {totalPages > 1 && (
            <div className="mt-10 flex items-center justify-between">
              <p className="text-xs text-muted-foreground/60">
                {(page - 1) * perPage + 1}–{Math.min(page * perPage, filtered.length)} of {filtered.length}
              </p>
              <Pagination currentPage={page} totalPages={totalPages} onPageChange={setPage} />
            </div>
          )}
        </>
      )}
    </div>
  )
}
