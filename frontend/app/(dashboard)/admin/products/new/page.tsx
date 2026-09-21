"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { ChevronLeft, Plus, X, Trash2, GripVertical } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Switch } from "@/components/ui/switch"
import { Badge } from "@/components/ui/badge"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { getCategoryOptions, createProduct, type ProductCategory } from "@/lib/api/products"
import { readImageFile } from "@/lib/utils"
import { toast } from "sonner"

const allTags = [
  "bestseller", "new", "sale", "limited", "eco-friendly",
  "premium", "organic", "handmade", "vintage", "smart",
]

export default function NewProductPage() {
  const router = useRouter()
  const [saving, setSaving] = useState(false)
  const [images, setImages] = useState<string[]>([])
  const [tags, setTags] = useState<string[]>([])
  const [tagInput, setTagInput] = useState("")
  const [category, setCategory] = useState("")
  const [categoryOptions, setCategoryOptions] = useState<ProductCategory[]>([])
  const [variants, setVariants] = useState<{ name: string; price: string; stock: string }[]>([])

  useEffect(() => {
    getCategoryOptions().then(setCategoryOptions).catch(() => setCategoryOptions([]))
  }, [])

  const addTag = (tag: string) => {
    if (tag && !tags.includes(tag)) {
      setTags([...tags, tag])
    }
    setTagInput("")
  }

  const removeTag = (tag: string) => {
    setTags(tags.filter((t) => t !== tag))
  }

  const addVariant = () => {
    setVariants([...variants, { name: "", price: "", stock: "" }])
  }

  const updateVariant = (i: number, field: string, value: string) => {
    const next = [...variants]
    next[i] = { ...next[i], [field]: value }
    setVariants(next)
  }

  const removeVariant = (i: number) => {
    setVariants(variants.filter((_, idx) => idx !== i))
  }

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
      sku: val("sku").trim() || `SKU-${Date.now()}`,
      name: val("name"),
      description: val("description"),
      price: Number(val("price") || 0),
      compare_at_price: Number(val("compare") || 0),
      cost_price: Number(val("cost") || 0),
      stock: Number(val("stock") || 0),
      low_stock_threshold: Number(val("lowStock") || 5),
      status: "active" as const,
      category_id: categoryOption?.id ?? null,
      tags,
      images,
    }

    try {
      await createProduct(payload)
      toast.success("Product created successfully!")
      router.push("/admin/products")
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create product")
      setSaving(false)
    }
  }

  return (
    <div className="p-6 max-w-4xl">
      <div className="flex items-center gap-4 mb-8">
        <Button variant="ghost" size="icon" onClick={() => router.back()}>
          <ChevronLeft className="h-4 w-4" />
        </Button>
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Add Product</h1>
          <p className="text-sm text-muted-foreground">Create a new product in your catalog</p>
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
                <Input id="name" placeholder="Enter product name" required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="sku">SKU</Label>
                <Input id="sku" placeholder="e.g. PRD-00001" required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="category">Category</Label>
                <Select value={category} onValueChange={setCategory}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select category" />
                  </SelectTrigger>
                  <SelectContent>
                    {categoryOptions.length === 0 && (
                      <SelectItem value="__none" disabled>No categories yet</SelectItem>
                    )}
                    {categoryOptions.map((c) => (
                      <SelectItem key={c.id} value={c.name}>{c.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2 sm:col-span-2">
                <Label htmlFor="description">Description</Label>
                <Textarea id="description" placeholder="Enter product description" rows={4} />
              </div>
            </div>

            <div className="space-y-2">
              <Label>Tags</Label>
              <div className="flex flex-wrap gap-2 mb-2">
                {tags.map((tag) => (
                  <Badge key={tag} variant="secondary" className="gap-1">
                    {tag}
                    <button onClick={() => removeTag(tag)} className="hover:text-destructive">
                      <X className="h-3 w-3" />
                    </button>
                  </Badge>
                ))}
              </div>
              <div className="flex gap-2">
                <Input
                  placeholder="Add a tag..."
                  value={tagInput}
                  onChange={(e) => setTagInput(e.target.value)}
                  onKeyDown={(e) => { if (e.key === "Enter") { e.preventDefault(); addTag(tagInput) } }}
                />
                <Button type="button" variant="outline" size="icon" onClick={() => addTag(tagInput)}>
                  <Plus className="h-4 w-4" />
                </Button>
              </div>
              <div className="flex flex-wrap gap-1 mt-2">
                {allTags.filter((t) => !tags.includes(t)).map((t) => (
                  <button
                    key={t}
                    type="button"
                    onClick={() => addTag(t)}
                    className="text-xs text-muted-foreground hover:text-primary px-2 py-0.5 rounded-md bg-muted"
                  >
                    + {t}
                  </button>
                ))}
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
                <Input id="price" type="number" step="0.01" placeholder="0.00" required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="compare">Compare at price ($)</Label>
                <Input id="compare" type="number" step="0.01" placeholder="0.00" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cost">Cost price ($)</Label>
                <Input id="cost" type="number" step="0.01" placeholder="0.00" />
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
                <Input id="stock" type="number" placeholder="0" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="lowStock">Low stock threshold</Label>
                <Input id="lowStock" type="number" placeholder="15" defaultValue="15" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="incoming">Incoming stock</Label>
                <Input id="incoming" type="number" placeholder="0" />
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

        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>Variants</CardTitle>
            <Button type="button" variant="outline" size="sm" onClick={addVariant}>
              <Plus className="h-4 w-4" /> Add variant
            </Button>
          </CardHeader>
          <CardContent>
            {variants.length === 0 && (
              <p className="text-sm text-muted-foreground text-center py-6">
                No variants added yet. Click &quot;Add variant&quot; to create size, color, or other variations.
              </p>
            )}
            {variants.map((v, i) => (
              <div key={i} className="flex items-center gap-3 mb-3">
                <GripVertical className="h-4 w-4 text-muted-foreground shrink-0" />
                <Input
                  placeholder="Variant name (e.g. Large, Red)"
                  value={v.name}
                  onChange={(e) => updateVariant(i, "name", e.target.value)}
                  className="flex-1"
                />
                <Input
                  type="number"
                  placeholder="Price"
                  value={v.price}
                  onChange={(e) => updateVariant(i, "price", e.target.value)}
                  className="w-24"
                />
                <Input
                  type="number"
                  placeholder="Stock"
                  value={v.stock}
                  onChange={(e) => updateVariant(i, "stock", e.target.value)}
                  className="w-20"
                />
                <Button type="button" variant="ghost" size="icon" onClick={() => removeVariant(i)}>
                  <Trash2 className="h-4 w-4 text-destructive" />
                </Button>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>SEO</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="slug">URL slug</Label>
              <Input id="slug" placeholder="product-name" />
              <p className="text-xs text-muted-foreground">Auto-generated if left empty</p>
            </div>
            <div className="space-y-2">
              <Label htmlFor="metaTitle">Meta title</Label>
              <Input id="metaTitle" placeholder="SEO title" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="metaDesc">Meta description</Label>
              <Textarea id="metaDesc" placeholder="SEO description" rows={2} />
            </div>
          </CardContent>
        </Card>

        <div className="flex items-center gap-3 pb-8">
          <Button type="submit" disabled={saving}>
            {saving ? "Saving..." : "Create product"}
          </Button>
          <Button type="button" variant="outline" onClick={() => router.back()}>
            Cancel
          </Button>
        </div>
      </form>
    </div>
  )
}
