"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { useParams } from "next/navigation"
import { Package, Star } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { Pagination } from "@/components/ui/pagination"
import { formatCurrency } from "@/lib/utils"
import { getProducts, getCategories } from "@/lib/api/products"
import type { Product } from "@/lib/types"

export default function CategoryPage() {
  const params = useParams()
  const slug = (params.slug as string).replace(/-/g, " ")
  const categoryName = slug.charAt(0).toUpperCase() + slug.slice(1)
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const perPage = 12

  useEffect(() => {
    async function load() {
      const all = await getProducts()
      const filtered = all.filter(
        (p) => p.category.toLowerCase() === categoryName.toLowerCase() && p.status === "active"
      )
      setProducts(filtered)
      setLoading(false)
    }
    load()
  }, [categoryName])

  const totalPages = Math.ceil(products.length / perPage)
  const paginated = products.slice((page - 1) * perPage, page * perPage)

  if (loading) {
    return (
      <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="h-8 w-48 rounded-lg bg-muted animate-pulse" />
        <div className="mt-8 grid grid-cols-2 gap-3 sm:gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
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
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold tracking-tight">{categoryName}</h1>
        <p className="mt-1 text-sm text-muted-foreground">{products.length} products</p>
      </div>

      {paginated.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20">
          <Package className="h-14 w-14 text-muted-foreground/30" />
          <h3 className="mt-4 text-lg font-semibold">No products in this category</h3>
          <Button asChild variant="outline" className="mt-5">
            <Link href="/products">Browse all products</Link>
          </Button>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-2 gap-3 sm:gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
            {paginated.map((product) => (
              <Link key={product.id} href={`/products/${product.id}`} className="group">
                <Card className="overflow-hidden transition-shadow hover:shadow-md h-full">
                  <div className="aspect-square bg-muted flex items-center justify-center">
                    <Package className="h-12 w-12 text-muted-foreground/20" />
                  </div>
                  <CardContent className="p-4">
                    <h3 className="text-sm font-medium group-hover:text-primary transition-colors truncate">{product.name}</h3>
                    <div className="mt-1 flex items-center gap-0.5">
                      {Array.from({ length: 5 }).map((_, i) => (
                        <Star key={i} className={`h-3 w-3 ${i < Math.round(product.rating) ? "fill-amber-400 text-amber-400" : "text-muted-foreground/20"}`} />
                      ))}
                      <span className="ml-1 text-xs text-muted-foreground">({product.reviewCount})</span>
                    </div>
                    <p className="mt-2 text-sm font-semibold">{formatCurrency(product.price)}</p>
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>

          {totalPages > 1 && (
            <div className="mt-10 flex items-center justify-center">
              <Pagination currentPage={page} totalPages={totalPages} onPageChange={setPage} />
            </div>
          )}
        </>
      )}
    </div>
  )
}
