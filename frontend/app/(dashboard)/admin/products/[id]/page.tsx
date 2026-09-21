"use client"

import { useEffect, useState } from "react"
import { useParams, useRouter } from "next/navigation"
import {
  Package,
  ChevronLeft,
  Edit,
  Star,
  ShoppingCart,
  BarChart3,
  MessageSquare,
  ClipboardList,
  Trash2,
} from "lucide-react"
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip as ChartTooltip,
  ResponsiveContainer,
} from "recharts"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { StatusBadge } from "@/components/shared/status-badge"
import { Skeleton } from "@/components/ui/skeleton"
import { formatCurrency, getInitials, formatDate } from "@/lib/utils"
import { getProduct, getProductReviews, getRelatedProducts } from "@/lib/api/products"
import type { Product, Review } from "@/lib/types"

export default function ProductDetailPage() {
  const params = useParams()
  const router = useRouter()
  const [product, setProduct] = useState<Product | null>(null)
  const [reviews, setReviews] = useState<Review[]>([])
  const [related, setRelated] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedImage, setSelectedImage] = useState(0)

  useEffect(() => {
    async function load() {
      const p = await getProduct(params.id as string)
      if (!p) { setLoading(false); return }
      const [r, rel] = await Promise.all([
        getProductReviews(p.id),
        getRelatedProducts(p.id),
      ])
      setProduct(p)
      setReviews(r)
      setRelated(rel)
      setLoading(false)
    }
    load()
  }, [params.id])

  if (loading) {
    return (
      <div className="p-6">
        <Skeleton className="h-8 w-24" />
        <div className="mt-6 grid gap-6 lg:grid-cols-2">
          <Skeleton className="h-[400px] rounded-xl" />
          <div className="space-y-4">
            <Skeleton className="h-8 w-3/4" />
            <Skeleton className="h-6 w-1/4" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-2/3" />
          </div>
        </div>
      </div>
    )
  }

  if (!product) {
    return (
      <div className="flex flex-col items-center justify-center py-24">
        <Package className="h-16 w-16 text-muted-foreground/40" />
        <h2 className="mt-4 text-xl font-semibold">Product not found</h2>
        <Button variant="outline" className="mt-4" onClick={() => router.push("/admin/products")}>
          Back to products
        </Button>
      </div>
    )
  }

  const salesHistory = Array.from({ length: 12 }, (_, i) => ({
    month: new Date(2024, i).toLocaleDateString("en-US", { month: "short" }),
    sales: ((i * 13 + product.stock) % 80) + 20,
    revenue: ((i * 29 + product.stock) % 40) * product.price * 0.25 + product.price * 5,
  }))

  const inventoryHistory = Array.from({ length: 12 }, (_, i) => ({
    month: new Date(2024, i).toLocaleDateString("en-US", { month: "short" }),
    stock: Math.max(0, product.stock + Math.round(Math.sin(i * 0.8) * 30)),
  }))

  return (
    <div className="p-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => router.push("/admin/products")}>
          <ChevronLeft className="h-4 w-4" />
        </Button>
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">{product.name}</h1>
          <p className="text-sm text-muted-foreground">{product.sku}</p>
        </div>
        <div className="ml-auto flex items-center gap-2">
          <StatusBadge status={product.status} />
          <Button variant="outline" size="sm" onClick={() => router.push(`/admin/products/${product.id}/edit`)}>
            <Edit className="h-4 w-4" /> Edit
          </Button>
        </div>
      </div>

      <div className="mt-6 grid gap-6 lg:grid-cols-5">
        <div className="lg:col-span-2 space-y-4">
          <Card className="overflow-hidden">
            <div className="relative aspect-square bg-muted flex items-center justify-center">
              {product.images[selectedImage] ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  src={product.images[selectedImage]}
                  alt={product.name}
                  className="h-full w-full object-cover"
                />
              ) : (
                <Package className="h-24 w-24 text-muted-foreground/30" />
              )}
            </div>
          </Card>
          {product.images.length > 1 && (
            <div className="flex gap-2">
              {product.images.map((img, i) => (
                <button
                  key={i}
                  onClick={() => setSelectedImage(i)}
                  className={`h-16 w-16 rounded-lg border-2 overflow-hidden ${
                    i === selectedImage ? "border-primary" : "border-border"
                  }`}
                >
                  <div className="h-full w-full bg-muted flex items-center justify-center">
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img src={img} alt="" className="h-full w-full object-cover" />
                  </div>
                </button>
              ))}
            </div>
          )}
        </div>

        <div className="lg:col-span-3 space-y-6">
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
            <Card>
              <CardContent className="p-4 text-center">
                <p className="text-xs text-muted-foreground">Price</p>
                <p className="mt-1 text-lg font-bold">{formatCurrency(product.price)}</p>
                {product.compareAtPrice && (
                  <p className="text-xs text-muted-foreground line-through">{formatCurrency(product.compareAtPrice)}</p>
                )}
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 text-center">
                <p className="text-xs text-muted-foreground">Stock</p>
                <p className={`mt-1 text-lg font-bold ${
                  product.stock === 0 ? "text-destructive" :
                  product.stock <= product.lowStockThreshold ? "text-warning" : ""
                }`}>
                  {product.stock}
                </p>
                {product.incoming > 0 && (
                  <p className="text-xs text-muted-foreground">+{product.incoming} incoming</p>
                )}
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 text-center">
                <p className="text-xs text-muted-foreground">Rating</p>
                <p className="mt-1 text-lg font-bold flex items-center justify-center gap-1">
                  {product.rating}
                  <Star className="h-4 w-4 fill-amber-400 text-amber-400" />
                </p>
                <p className="text-xs text-muted-foreground">{product.reviewCount} reviews</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 text-center">
                <p className="text-xs text-muted-foreground">Category</p>
                <Badge variant="secondary" className="mt-1">{product.category}</Badge>
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Description</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm leading-relaxed text-muted-foreground">{product.description}</p>
              {product.tags.length > 0 && (
                <div className="mt-4 flex flex-wrap gap-2">
                  {product.tags.map((tag) => (
                    <Badge key={tag} variant="secondary" className="font-normal">
                      {tag}
                    </Badge>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>

          <Tabs defaultValue="sales">
            <TabsList>
              <TabsTrigger value="sales" className="gap-2">
                <BarChart3 className="h-4 w-4" /> Sales History
              </TabsTrigger>
              <TabsTrigger value="inventory" className="gap-2">
                <ClipboardList className="h-4 w-4" /> Inventory History
              </TabsTrigger>
              <TabsTrigger value="reviews" className="gap-2">
                <MessageSquare className="h-4 w-4" /> Reviews ({reviews.length})
              </TabsTrigger>
              <TabsTrigger value="related" className="gap-2">
                <ShoppingCart className="h-4 w-4" /> Related
              </TabsTrigger>
            </TabsList>

            <TabsContent value="sales">
              <Card>
                <CardContent className="p-6">
                  <div className="h-[250px]">
                    <ResponsiveContainer width="100%" height="100%">
                      <AreaChart data={salesHistory}>
                        <defs>
                          <linearGradient id="salesGrad" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="5%" stopColor="#2563eb" stopOpacity={0.15} />
                            <stop offset="95%" stopColor="#2563eb" stopOpacity={0} />
                          </linearGradient>
                        </defs>
                        <CartesianGrid strokeDasharray="3 3" stroke="var(--color-border)" vertical={false} />
                        <XAxis dataKey="month" tick={{ fontSize: 12 }} tickLine={false} axisLine={false} />
                        <YAxis tick={{ fontSize: 12 }} tickLine={false} axisLine={false} />
                        <ChartTooltip
                          contentStyle={{
                            borderRadius: "12px",
                            border: "1px solid var(--color-border)",
                            background: "var(--color-card)",
                          }}
                        />
                        <Area type="monotone" dataKey="sales" stroke="#2563eb" strokeWidth={2} fill="url(#salesGrad)" />
                      </AreaChart>
                    </ResponsiveContainer>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="inventory">
              <Card>
                <CardContent className="p-6">
                  <div className="h-[250px]">
                    <ResponsiveContainer width="100%" height="100%">
                      <AreaChart data={inventoryHistory}>
                        <defs>
                          <linearGradient id="invGrad" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="5%" stopColor="#059669" stopOpacity={0.15} />
                            <stop offset="95%" stopColor="#059669" stopOpacity={0} />
                          </linearGradient>
                        </defs>
                        <CartesianGrid strokeDasharray="3 3" stroke="var(--color-border)" vertical={false} />
                        <XAxis dataKey="month" tick={{ fontSize: 12 }} tickLine={false} axisLine={false} />
                        <YAxis tick={{ fontSize: 12 }} tickLine={false} axisLine={false} />
                        <ChartTooltip
                          contentStyle={{
                            borderRadius: "12px",
                            border: "1px solid var(--color-border)",
                            background: "var(--color-card)",
                          }}
                        />
                        <Area type="monotone" dataKey="stock" stroke="#059669" strokeWidth={2} fill="url(#invGrad)" />
                      </AreaChart>
                    </ResponsiveContainer>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="reviews">
              <Card>
                <CardContent className="p-6 space-y-4">
                  {reviews.length === 0 && (
                    <p className="text-sm text-muted-foreground text-center py-8">No reviews yet</p>
                  )}
                  {reviews.map((review) => (
                    <div key={review.id} className="border-b last:border-0 pb-4 last:pb-0">
                      <div className="flex items-center gap-3">
                        <Avatar className="h-8 w-8">
                          <AvatarFallback className="text-xs">{getInitials(review.customerName)}</AvatarFallback>
                        </Avatar>
                        <div>
                          <p className="text-sm font-medium">{review.customerName}</p>
                          <div className="flex items-center gap-1">
                            {Array.from({ length: 5 }).map((_, i) => (
                              <Star
                                key={i}
                                className={`h-3 w-3 ${i < review.rating ? "fill-amber-400 text-amber-400" : "text-muted-foreground/30"}`}
                              />
                            ))}
                            {review.isVerified && (
                              <Badge variant="success" className="ml-2 text-[10px] px-1.5 py-0">Verified</Badge>
                            )}
                          </div>
                        </div>
                        <span className="ml-auto text-xs text-muted-foreground">{formatDate(review.createdAt)}</span>
                      </div>
                      <p className="mt-2 text-sm font-medium">{review.title}</p>
                      <p className="mt-1 text-sm text-muted-foreground">{review.content}</p>
                    </div>
                  ))}
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="related">
              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
                {related.length === 0 && (
                  <p className="text-sm text-muted-foreground col-span-full text-center py-8">No related products</p>
                )}
                {related.map((rp) => (
                  <Card
                    key={rp.id}
                    className="cursor-pointer hover:shadow-md transition-shadow"
                    onClick={() => router.push(`/admin/products/${rp.id}`)}
                  >
                    <CardContent className="p-4">
                      <div className="flex h-20 items-center justify-center rounded-lg bg-muted mb-3">
                        <Package className="h-8 w-8 text-muted-foreground" />
                      </div>
                      <p className="text-sm font-medium truncate">{rp.name}</p>
                      <p className="text-sm text-muted-foreground">{formatCurrency(rp.price)}</p>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  )
}
