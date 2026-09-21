"use client"

import { useEffect, useState } from "react"
import { useParams, useRouter } from "next/navigation"
import Link from "next/link"
import { Package, Star, ChevronLeft, ShoppingCart, Minus, Plus, Truck, Shield, RefreshCw } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Skeleton } from "@/components/ui/skeleton"
import { Separator } from "@/components/ui/separator"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Label } from "@/components/ui/label"
import { toast } from "sonner"
import { formatCurrency, formatDate } from "@/lib/utils"
import { useCart } from "@/lib/cart-context"
import { getProduct, getProductReviews, getRelatedProducts, createReview } from "@/lib/api/products"
import type { Product, Review } from "@/lib/types"

export default function ProductDetailPage() {
  const params = useParams()
  const router = useRouter()
  const { addItem } = useCart()
  const [product, setProduct] = useState<Product | null>(null)
  const [reviews, setReviews] = useState<Review[]>([])
  const [related, setRelated] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [quantity, setQuantity] = useState(1)
  const [reviewName, setReviewName] = useState("")
  const [reviewRating, setReviewRating] = useState(5)
  const [reviewTitle, setReviewTitle] = useState("")
  const [reviewContent, setReviewContent] = useState("")
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    async function load() {
      const p = await getProduct(params.id as string)
      if (!p) { setLoading(false); return }
      const [r, rel] = await Promise.all([getProductReviews(p.id), getRelatedProducts(p.id)])
      setProduct(p)
      setReviews(r)
      setRelated(rel)
      setLoading(false)
    }
    load()
  }, [params.id])

  const addToCart = () => {
    if (!product) return
    addItem(product, quantity)
    toast.success(`Added ${quantity} × ${product.name} to cart`)
  }

  const submitReview = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!product) return
    setSubmitting(true)
    try {
      await createReview({
        product_id: product.id,
        customer_name: reviewName || "Anonymous",
        rating: reviewRating,
        title: reviewTitle,
        content: reviewContent,
      })
      toast.success("Review submitted! It will appear after moderation.")
      setReviewTitle(""); setReviewContent(""); setReviewRating(5)
      const refreshed = await getProductReviews(product.id)
      setReviews(refreshed)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to submit review")
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <Skeleton className="h-6 w-24" />
        <div className="mt-6 grid gap-8 lg:grid-cols-2">
          <Skeleton className="aspect-square rounded-xl" />
          <div className="space-y-4">
            <Skeleton className="h-4 w-20" />
            <Skeleton className="h-8 w-3/4" />
            <Skeleton className="h-6 w-24" />
            <Skeleton className="h-20 w-full" />
          </div>
        </div>
      </div>
    )
  }

  if (!product) {
    return (
      <div className="flex flex-col items-center justify-center py-24">
        <Package className="h-16 w-16 text-muted-foreground/30" />
        <h2 className="mt-4 text-xl font-semibold">Product not found</h2>
        <Button variant="outline" className="mt-4" onClick={() => router.push("/products")}>
          Browse products
        </Button>
      </div>
    )
  }

  const discount = product.compareAtPrice
    ? Math.round((1 - product.price / product.compareAtPrice) * 100)
    : 0

  return (
    <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <Link href="/products" className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors mb-6">
        <ChevronLeft className="h-3.5 w-3.5" /> Back to products
      </Link>

      <div className="grid gap-8 lg:grid-cols-2">
        <div className="space-y-4">
          <div className="relative aspect-square rounded-xl bg-muted flex items-center justify-center overflow-hidden">
            {product.images[0] ? (
              // eslint-disable-next-line @next/next/no-img-element
              <img src={product.images[0]} alt={product.name} className="h-full w-full object-cover" />
            ) : (
              <Package className="h-32 w-32 text-muted-foreground/15" />
            )}
            {discount > 0 && (
              <Badge className="absolute left-4 top-4 text-sm px-3 py-1">
                -{discount}%
              </Badge>
            )}
            {product.stock === 0 && (
              <div className="absolute inset-0 bg-background/60 flex items-center justify-center backdrop-blur-sm">
                <Badge variant="outline" className="text-lg px-6 py-2">Out of stock</Badge>
              </div>
            )}
          </div>
        </div>

        <div className="flex flex-col gap-6">
          <div>
            <p className="text-sm text-muted-foreground mb-1">{product.category}</p>
            <h1 className="text-2xl font-bold tracking-tight lg:text-3xl">{product.name}</h1>
            <div className="mt-2 flex items-center gap-2">
              <div className="flex items-center gap-0.5">
                {Array.from({ length: 5 }).map((_, i) => (
                  <Star key={i} className={`h-4 w-4 ${i < Math.round(product.rating) ? "fill-amber-400 text-amber-400" : "text-muted-foreground/20"}`} />
                ))}
              </div>
              <span className="text-sm text-muted-foreground">{product.rating} ({product.reviewCount} reviews)</span>
            </div>
          </div>

          <div className="flex items-baseline gap-3">
            <span className="text-3xl font-bold">{formatCurrency(product.price)}</span>
            {product.compareAtPrice && (
              <span className="text-lg text-muted-foreground line-through">{formatCurrency(product.compareAtPrice)}</span>
            )}
          </div>

          <p className="text-sm text-muted-foreground leading-relaxed">{product.description}</p>

          {product.tags.length > 0 && (
            <div className="flex flex-wrap gap-1.5">
              {product.tags.map((tag) => (
                <Badge key={tag} variant="secondary" className="text-xs font-normal">{tag}</Badge>
              ))}
            </div>
          )}

          <Separator />

          <div className="flex items-center gap-4">
            <div className="flex items-center rounded-lg border">
              <button onClick={() => setQuantity(Math.max(1, quantity - 1))} className="flex h-10 w-10 items-center justify-center text-muted-foreground hover:text-foreground transition-colors">
                <Minus className="h-3.5 w-3.5" />
              </button>
              <span className="flex h-10 w-12 items-center justify-center text-sm font-medium tabular-nums border-x">
                {quantity}
              </span>
              <button onClick={() => setQuantity(quantity + 1)} className="flex h-10 w-10 items-center justify-center text-muted-foreground hover:text-foreground transition-colors">
                <Plus className="h-3.5 w-3.5" />
              </button>
            </div>
            <Button className="flex-1 h-10 gap-2" onClick={addToCart} disabled={product.stock === 0}>
              <ShoppingCart className="h-4 w-4" />
              {product.stock === 0 ? "Out of stock" : "Add to cart"}
            </Button>
          </div>

          {product.stock > 0 && product.stock <= product.lowStockThreshold && (
            <p className="text-xs text-warning">Only {product.stock} left in stock</p>
          )}

          <div className="grid grid-cols-3 gap-3 rounded-xl border bg-muted/30 p-4">
            {[
              { icon: Truck, text: "Gratis ongkir min. Rp 500rb" },
              { icon: Shield, text: "Secure checkout" },
              { icon: RefreshCw, text: "30-day returns" },
            ].map((item) => (
              <div key={item.text} className="flex items-center gap-2">
                <item.icon className="h-4 w-4 text-primary shrink-0" />
                <span className="text-xs text-muted-foreground">{item.text}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      <Tabs defaultValue="reviews" className="mt-12">
        <TabsList>
          <TabsTrigger value="reviews" className="gap-2">
            Reviews ({reviews.length})
          </TabsTrigger>
          <TabsTrigger value="details" className="gap-2">
            Details
          </TabsTrigger>
        </TabsList>

        <TabsContent value="reviews" className="mt-6">
          <Card className="mb-6">
            <CardContent className="p-6">
              <h3 className="text-base font-semibold mb-4">Write a review</h3>
              <form onSubmit={submitReview} className="space-y-4">
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-1.5">
                    <Label className="text-xs font-medium">Name</Label>
                    <Input value={reviewName} onChange={(e) => setReviewName(e.target.value)} placeholder="Your name" />
                  </div>
                  <div className="space-y-1.5">
                    <Label className="text-xs font-medium">Rating</Label>
                    <div className="flex items-center gap-1">
                      {[1, 2, 3, 4, 5].map((r) => (
                        <button key={r} type="button" onClick={() => setReviewRating(r)}>
                          <Star className={`h-5 w-5 ${r <= reviewRating ? "fill-amber-400 text-amber-400" : "text-muted-foreground/20"}`} />
                        </button>
                      ))}
                    </div>
                  </div>
                </div>
                <div className="space-y-1.5">
                  <Label className="text-xs font-medium">Title</Label>
                  <Input value={reviewTitle} onChange={(e) => setReviewTitle(e.target.value)} placeholder="Review title" />
                </div>
                <div className="space-y-1.5">
                  <Label className="text-xs font-medium">Review</Label>
                  <Textarea value={reviewContent} onChange={(e) => setReviewContent(e.target.value)} rows={3} placeholder="Share your experience..." />
                </div>
                <Button type="submit" size="sm" disabled={submitting}>
                  {submitting ? "Submitting..." : "Submit review"}
                </Button>
              </form>
            </CardContent>
          </Card>

          {reviews.length === 0 ? (
            <p className="text-sm text-muted-foreground text-center py-12">No reviews yet.</p>
          ) : (
            <div className="grid gap-4 sm:grid-cols-2">
              {reviews.map((review) => (
                <Card key={review.id}>
                  <CardContent className="p-4">
                    <div className="flex items-center gap-3">
                      <div className="flex h-9 w-9 items-center justify-center rounded-full bg-muted text-xs font-medium">
                        {review.customerName.charAt(0)}
                      </div>
                      <div>
                        <p className="text-sm font-medium">{review.customerName}</p>
                        <div className="flex items-center gap-1">
                          {Array.from({ length: 5 }).map((_, i) => (
                            <Star key={i} className={`h-3 w-3 ${i < review.rating ? "fill-amber-400 text-amber-400" : "text-muted-foreground/20"}`} />
                          ))}
                          {review.isVerified && (
                            <Badge variant="success" className="ml-1 text-[10px] px-1.5">Verified</Badge>
                          )}
                        </div>
                      </div>
                      <span className="ml-auto text-xs text-muted-foreground">{formatDate(review.createdAt)}</span>
                    </div>
                    <p className="mt-2 text-sm font-medium">{review.title}</p>
                    <p className="mt-1 text-sm text-muted-foreground">{review.content}</p>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="details" className="mt-6">
          <Card>
            <CardContent className="p-6">
              <dl className="grid gap-4 sm:grid-cols-2">
                <div>
                  <dt className="text-xs text-muted-foreground">SKU</dt>
                  <dd className="text-sm font-medium font-mono">{product.sku}</dd>
                </div>
                <div>
                  <dt className="text-xs text-muted-foreground">Category</dt>
                  <dd className="text-sm font-medium">{product.category}</dd>
                </div>
                <div>
                  <dt className="text-xs text-muted-foreground">Stock</dt>
                  <dd className={`text-sm font-medium ${product.stock === 0 ? "text-destructive" : ""}`}>{product.stock} units</dd>
                </div>
                <div>
                  <dt className="text-xs text-muted-foreground">Rating</dt>
                  <dd className="text-sm font-medium">{product.rating} / 5</dd>
                </div>
              </dl>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {related.length > 0 && (
        <section className="mt-12">
          <h2 className="text-xl font-bold tracking-tight mb-6">You might also like</h2>
          <div className="grid grid-cols-2 gap-3 sm:gap-4 sm:grid-cols-3 lg:grid-cols-4">
            {related.slice(0, 4).map((rp) => (
              <Link key={rp.id} href={`/products/${rp.id}`} className="group">
                <Card className="overflow-hidden transition-shadow hover:shadow-md">
                  <div className="aspect-square bg-muted flex items-center justify-center">
                    <Package className="h-12 w-12 text-muted-foreground/20" />
                  </div>
                  <CardContent className="p-4">
                    <p className="text-xs text-muted-foreground mb-1">{rp.category}</p>
                    <h3 className="text-sm font-medium group-hover:text-primary transition-colors truncate">{rp.name}</h3>
                    <p className="mt-1 text-sm font-semibold">{formatCurrency(rp.price)}</p>
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>
        </section>
      )}
    </div>
  )
}
