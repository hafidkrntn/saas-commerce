"use client"

import { useEffect, useState, useMemo } from "react"
import { useSearchParams } from "next/navigation"
import Link from "next/link"
import { Package, Star, Search as SearchIcon } from "lucide-react"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { formatCurrency } from "@/lib/utils"
import { getProducts } from "@/lib/api/products"
import type { Product } from "@/lib/types"

export default function SearchPage() {
  const searchParams = useSearchParams()
  const initialQuery = searchParams.get("q") || ""
  const [products, setProducts] = useState<Product[]>([])
  const [query, setQuery] = useState(initialQuery)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      const p = await getProducts()
      setProducts(p.filter((p) => p.status === "active"))
      setLoading(false)
    }
    load()
  }, [])

  const results = useMemo(() => {
    if (!query) return []
    const q = query.toLowerCase()
    return products.filter(
      (p) =>
        p.name.toLowerCase().includes(q) ||
        p.description.toLowerCase().includes(q) ||
        p.category.toLowerCase().includes(q) ||
        p.tags.some((t) => t.toLowerCase().includes(q))
    )
  }, [products, query])

  return (
    <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <div className="max-w-xl mx-auto mb-8">
        <div className="relative">
          <SearchIcon className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground/50" />
          <Input
            placeholder="Search products..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="h-10 pl-9 text-sm"
            autoFocus
          />
        </div>
      </div>

      {loading ? (
        <div className="grid grid-cols-2 gap-3 sm:gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="rounded-xl border">
              <div className="aspect-square bg-muted animate-pulse" />
              <div className="p-4 space-y-2">
                <div className="h-3 w-16 rounded bg-muted animate-pulse" />
                <div className="h-4 w-3/4 rounded bg-muted animate-pulse" />
              </div>
            </div>
          ))}
        </div>
      ) : query && results.length === 0 ? (
        <div className="text-center py-20">
          <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-muted mx-auto">
            <SearchIcon className="h-7 w-7 text-muted-foreground/40" />
          </div>
          <h2 className="mt-4 text-lg font-semibold">No results for &quot;{query}&quot;</h2>
          <p className="mt-1 text-sm text-muted-foreground">Try different keywords or browse our categories.</p>
          <Button asChild variant="outline" className="mt-5">
            <Link href="/products">Browse all products</Link>
          </Button>
        </div>
      ) : query ? (
        <>
          <p className="text-sm text-muted-foreground mb-6">{results.length} results for &quot;{query}&quot;</p>
          <div className="grid grid-cols-2 gap-3 sm:gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
            {results.map((product) => (
              <Link key={product.id} href={`/products/${product.id}`} className="group">
                <Card className="overflow-hidden transition-shadow hover:shadow-md h-full">
                  <div className="aspect-square bg-muted flex items-center justify-center">
                    <Package className="h-12 w-12 text-muted-foreground/20" />
                  </div>
                  <CardContent className="p-4">
                    <p className="text-xs text-muted-foreground mb-1">{product.category}</p>
                    <h3 className="text-sm font-medium group-hover:text-primary transition-colors truncate">{product.name}</h3>
                    <div className="mt-1 flex items-center gap-0.5">
                      {Array.from({ length: 5 }).map((_, i) => (
                        <Star key={i} className={`h-3 w-3 ${i < Math.round(product.rating) ? "fill-amber-400 text-amber-400" : "text-muted-foreground/20"}`} />
                      ))}
                    </div>
                    <p className="mt-2 text-sm font-semibold">{formatCurrency(product.price)}</p>
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>
        </>
      ) : (
        <div className="text-center py-20">
          <SearchIcon className="h-10 w-10 text-muted-foreground/30 mx-auto" />
          <h2 className="mt-4 text-lg font-semibold">Search our store</h2>
          <p className="mt-1 text-sm text-muted-foreground">Type above to find what you&apos;re looking for.</p>
        </div>
      )}
    </div>
  )
}
