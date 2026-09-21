"use client"

import { useEffect, useState } from "react"
import { useParams, useRouter } from "next/navigation"
import { ChevronLeft, Plus, X, Trash2, GripVertical, Package } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Switch } from "@/components/ui/switch"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { toast } from "sonner"
import { readImageFile } from "@/lib/utils"
import { getProduct, getCategoryOptions, updateProduct, type ProductCategory } from "@/lib/api/products"
import type { Product } from "@/lib/types"

export default function EditProductPage() {
  const params = useParams()
  const router = useRouter()
  const [product, setProduct] = useState<Product | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [category, setCategory] = useState("")
  const [categoryOptions, setCategoryOptions] = useState<ProductCategory[]>([])
  const [images, setImages] = useState<string[]>([])

  useEffect(() => {
    async function load() {
      const [p, cats] = await Promise.all([
        getProduct(params.id as string),
        getCategoryOptions(),
      ])
      setProduct(p)
      setCategory(p?.category ?? "")
      setImages(p?.images ?? [])
      setCategoryOptions(cats)
      setLoading(false)
    }
    load()
  }, [params.id])

  const addImage = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    try {
      const dataUrl = await readImageFile(file)
      setImages((prev) => [...prev, dataUrl])
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Gagal membaca gambar")
    }
    e.target.value = ""
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)

    const val = (id: string) => (document.getElementById(id) as HTMLInputElement | null)?.value ?? ""
    const categoryOption = categoryOptions.find((c) => c.name === category)

    const payload = {
      name: val("name"),
      sku: val("sku"),
      description: val("description"),
      price: Number(val("price") || 0),
      compare_at_price: Number(val("compare") || 0),
      cost_price: Number(val("cost") || 0),
      low_stock_threshold: Number(val("lowStock") || 5),
      category_id: categoryOption?.id ?? null,
      images,
    }

    try {
      await updateProduct(params.id as string, payload)
      toast.success("Product updated successfully!")
      router.push(`/admin/products/${params.id}`)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update product")
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="p-6 max-w-4xl">
        <Skeleton className="h-8 w-32" />
        <div className="mt-6 space-y-6">
          <Skeleton className="h-48 rounded-xl" />
          <Skeleton className="h-32 rounded-xl" />
          <Skeleton className="h-32 rounded-xl" />
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

  return (
    <div className="p-6 max-w-4xl">
      <div className="flex items-center gap-4 mb-8">
        <Button variant="ghost" size="icon" onClick={() => router.back()}>
          <ChevronLeft className="h-4 w-4" />
        </Button>
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Edit Product</h1>
          <p className="text-sm text-muted-foreground">{product.name}</p>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-8">
        <Card>
          <CardHeader>
            <CardTitle>General Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2 sm:col-span-2">
                <Label htmlFor="name">Product name</Label>
                <Input id="name" defaultValue={product.name} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="sku">SKU</Label>
                <Input id="sku" defaultValue={product.sku} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="category">Category</Label>
                <Select value={category} onValueChange={setCategory}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {categoryOptions.map((c) => (
                      <SelectItem key={c.id} value={c.name}>{c.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2 sm:col-span-2">
                <Label htmlFor="description">Description</Label>
                <Textarea id="description" defaultValue={product.description} rows={4} />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Pricing</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 sm:grid-cols-3">
              <div className="space-y-2">
                <Label htmlFor="price">Price ($)</Label>
                <Input id="price" type="number" step="0.01" defaultValue={product.price} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="compare">Compare at price ($)</Label>
                <Input id="compare" type="number" step="0.01" defaultValue={product.compareAtPrice || ""} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cost">Cost price ($)</Label>
                <Input id="cost" type="number" step="0.01" defaultValue={product.costPrice || ""} />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Inventory</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 sm:grid-cols-3">
              <div className="space-y-2">
                <Label htmlFor="stock">Stock quantity</Label>
                <Input id="stock" type="number" defaultValue={product.stock} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="lowStock">Low stock threshold</Label>
                <Input id="lowStock" type="number" defaultValue={product.lowStockThreshold} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="incoming">Incoming stock</Label>
                <Input id="incoming" type="number" defaultValue={product.incoming} />
              </div>
            </div>
            <div className="mt-4 flex items-center gap-2">
              <Switch id="track" defaultChecked />
              <Label htmlFor="track">Track inventory</Label>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Media</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-4 gap-3">
              {images.map((img, i) => (
                <div key={i} className="relative aspect-square overflow-hidden rounded-xl border bg-muted/30">
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img src={img} alt="" className="h-full w-full object-cover" />
                  <button
                    type="button"
                    onClick={() => setImages(images.filter((_, idx) => idx !== i))}
                    className="absolute right-1.5 top-1.5 flex h-6 w-6 items-center justify-center rounded-full bg-background/90 shadow-sm hover:bg-destructive hover:text-white"
                  >
                    <X className="h-3 w-3" />
                  </button>
                </div>
              ))}
              <label className="flex aspect-square cursor-pointer items-center justify-center rounded-xl border-2 border-dashed border-border hover:border-primary/50 transition-colors bg-muted/30">
                <div className="text-center">
                  <Plus className="mx-auto h-6 w-6 text-muted-foreground" />
                  <span className="mt-1 block text-xs text-muted-foreground">Add image</span>
                </div>
                <input type="file" accept="image/*" className="hidden" onChange={addImage} />
              </label>
            </div>
          </CardContent>
        </Card>

        <div className="flex items-center gap-3 pb-8">
          <Button type="submit" disabled={saving}>
            {saving ? "Saving..." : "Save changes"}
          </Button>
          <Button type="button" variant="outline" onClick={() => router.back()}>
            Cancel
          </Button>
        </div>
      </form>
    </div>
  )
}
