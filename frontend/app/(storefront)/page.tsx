"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { ArrowRight, Star, Package, Truck, Shield, RefreshCw } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { formatCurrency } from "@/lib/utils"
import { getProducts } from "@/lib/api/products"
import type { Product } from "@/lib/types"

export default function StorefrontHome() {
  const [featured, setFeatured] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      const products = await getProducts()
      setFeatured(products.filter((p) => p.status === "active").slice(0, 8))
      setLoading(false)
    }
    load()
  }, [])

  return (
    <div>
      <section className="relative overflow-hidden border-b bg-gradient-to-b from-background to-muted/20">
        <div className="mx-auto max-w-7xl px-4 py-20 sm:px-6 sm:py-28 lg:px-8">
          <div className="mx-auto max-w-2xl text-center">
            <Badge variant="secondary" className="mb-6 text-xs font-normal px-3 py-1">
              New season collection available
            </Badge>
            <h1 className="text-4xl font-bold tracking-tight sm:text-5xl lg:text-6xl">
              Premium products for{" "}
              <span className="text-primary">modern living</span>
            </h1>
            <p className="mt-6 text-lg leading-relaxed text-muted-foreground max-w-lg mx-auto">
              Discover our curated collection of high-quality products. Gratis ongkir minimal belanja Rp 500.000.
            </p>
            <div className="mt-8 flex items-center justify-center gap-4">
              <Button asChild size="lg" className="h-12 px-8 text-base">
                <Link href="/products">
                  Shop now <ArrowRight className="h-4 w-4" />
                </Link>
              </Button>
              <Button asChild variant="outline" size="lg" className="h-12 px-8 text-base">
                <Link href="/products?sale=true">View sale</Link>
              </Button>
            </div>
          </div>
        </div>
      </section>

      <section className="border-b bg-muted/20">
        <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
          <div className="grid grid-cols-2 gap-6 sm:grid-cols-4">
            {[
              { icon: Truck, title: "Free shipping", desc: "Min. belanja Rp 500rb" },
              { icon: Shield, title: "Secure checkout", desc: "SSL encrypted" },
              { icon: RefreshCw, title: "Easy returns", desc: "30-day returns" },
              { icon: Package, title: "Track delivery", desc: "Real-time tracking" },
            ].map((item) => (
              <div key={item.title} className="flex items-center gap-3">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
                  <item.icon className="h-5 w-5 text-primary" />
                </div>
                <div>
                  <p className="text-sm font-medium">{item.title}</p>
                  <p className="text-xs text-muted-foreground">{item.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="mx-auto max-w-7xl px-4 py-16 sm:px-6 lg:px-8">
        <div className="flex items-end justify-between mb-8">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Featured products</h2>
            <p className="mt-1 text-sm text-muted-foreground">Our most popular items this month</p>
          </div>
          <Button asChild variant="ghost" size="sm" className="text-sm gap-1">
            <Link href="/products">View all <ArrowRight className="h-3.5 w-3.5" /></Link>
          </Button>
        </div>

        <div className="grid grid-cols-2 gap-3 sm:gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
          {featured.map((product) => (
            <Link key={product.id} href={`/products/${product.id}`} className="group">
              <Card className="overflow-hidden transition-shadow hover:shadow-md">
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
      </section>

      <section className="bg-muted/30 border-y">
        <div className="mx-auto max-w-7xl px-4 py-16 sm:px-6 lg:px-8">
          <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
            {[
              { name: "Electronics", count: "24 products", color: "bg-blue-500/10 text-blue-600" },
              { name: "Clothing", count: "18 products", color: "bg-emerald-500/10 text-emerald-600" },
              { name: "Home & Garden", count: "15 products", color: "bg-amber-500/10 text-amber-600" },
            ].map((cat) => (
              <Link
                key={cat.name}
                href={`/categories/${cat.name.toLowerCase().replace(/ & /g, "-")}`}
                className="group rounded-xl border bg-background p-6 transition-shadow hover:shadow-md"
              >
                <div className={`inline-flex h-12 w-12 items-center justify-center rounded-xl ${cat.color} mb-4`}>
                  <Package className="h-6 w-6" />
                </div>
                <h3 className="text-lg font-semibold group-hover:text-primary transition-colors">{cat.name}</h3>
                <p className="mt-1 text-sm text-muted-foreground">{cat.count}</p>
              </Link>
            ))}
          </div>
        </div>
      </section>

      <section className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8 text-center">
        <h2 className="text-2xl font-bold tracking-tight">Stay in the loop</h2>
        <p className="mt-2 text-sm text-muted-foreground max-w-sm mx-auto">
          Get notified about new products, exclusive sales, and more.
        </p>
        <div className="mt-6 mx-auto flex max-w-md items-center gap-3">
          <input
            type="email"
            placeholder="Enter your email"
            className="flex h-10 w-full rounded-lg border bg-background px-4 text-sm ring-offset-background placeholder:text-muted-foreground/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          />
          <Button className="shrink-0">Subscribe</Button>
        </div>
      </section>
    </div>
  )
}
