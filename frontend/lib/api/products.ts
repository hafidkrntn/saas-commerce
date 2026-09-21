import { http, type Paginated } from "./http"
import type { Product, Review } from "@/lib/types"

interface CategoryDTO {
  id: string
  name: string
}

interface ProductImageDTO {
  id: string
  url: string
  alt: string
  is_primary: boolean
}

interface ProductDTO {
  id: string
  sku: string
  name: string
  description: string
  price: number
  compare_at_price: number
  cost_price: number
  stock: number
  reserved: number
  incoming: number
  low_stock_threshold: number
  status: "active" | "draft" | "archived"
  category_id: string | null
  category_name: string
  tags: string[]
  is_active: boolean
  rating: number
  review_count: number
  images: ProductImageDTO[]
  created_at: string
  updated_at: string
}

function mapProduct(p: ProductDTO): Product {
  return {
    id: p.id,
    name: p.name,
    sku: p.sku,
    description: p.description,
    price: p.price,
    compareAtPrice: p.compare_at_price || undefined,
    costPrice: p.cost_price || undefined,
    stock: p.stock,
    reserved: p.reserved,
    incoming: p.incoming,
    lowStockThreshold: p.low_stock_threshold,
    category: p.category_name || "",
    categoryId: p.category_id ?? "",
    tags: p.tags ?? [],
    images: (p.images ?? []).map((i) => i.url),
    status: p.status,
    isActive: p.is_active,
    rating: p.rating,
    reviewCount: p.review_count,
    createdAt: p.created_at,
    updatedAt: p.updated_at,
  }
}

export interface ProductCategory {
  id: string
  name: string
}

export interface ProductPayload {
  sku: string
  name: string
  description?: string
  price: number
  compare_at_price?: number
  cost_price?: number
  stock?: number
  low_stock_threshold?: number
  status?: "active" | "draft" | "archived"
  category_id?: string | null
  tags?: string[]
  images?: string[]
}

export async function getProducts(): Promise<Product[]> {
  const res = await http.get<Paginated<ProductDTO>>("/products?limit=100")
  return res.data.map(mapProduct)
}

export async function getProduct(id: string): Promise<Product | null> {
  try {
    const res = await http.get<ProductDTO>(`/products/${id}`)
    return mapProduct(res)
  } catch {
    return null
  }
}

export async function createProduct(payload: ProductPayload): Promise<Product> {
  const res = await http.post<ProductDTO>("/products", payload)
  return mapProduct(res)
}

export async function updateProduct(id: string, payload: Partial<ProductPayload>): Promise<Product> {
  const res = await http.put<ProductDTO>(`/products/${id}`, payload)
  return mapProduct(res)
}

export async function deleteProduct(id: string): Promise<void> {
  await http.delete(`/products/${id}`)
}

export async function getCategories(): Promise<string[]> {
  const res = await http.get<Paginated<CategoryDTO>>("/categories?limit=100")
  return res.data.map((c) => c.name)
}

export async function getCategoryOptions(): Promise<ProductCategory[]> {
  const res = await http.get<Paginated<CategoryDTO>>("/categories?limit=100")
  return res.data.map((c) => ({ id: c.id, name: c.name }))
}

interface ReviewDTO {
  id: string
  product_id: string
  customer_id: string | null
  customer_name: string
  customer_avatar: string
  rating: number
  title: string
  content: string
  is_verified: boolean
  is_published: boolean
  created_at: string
}

function mapReview(r: ReviewDTO): Review {
  return {
    id: r.id,
    productId: r.product_id,
    customerName: r.customer_name,
    customerAvatar: r.customer_avatar || undefined,
    rating: r.rating,
    title: r.title,
    content: r.content,
    isVerified: r.is_verified,
    createdAt: r.created_at,
  }
}

export async function getProductReviews(productId: string): Promise<Review[]> {
  const res = await http.get<Paginated<ReviewDTO>>(`/reviews?product_id=${productId}&published=true&limit=100`)
  return res.data.map(mapReview)
}

export async function getRelatedProducts(productId: string): Promise<Product[]> {
  const products = await getProducts()
  const current = products.find((p) => p.id === productId)
  if (!current) return []
  return products.filter((p) => p.id !== productId && p.categoryId === current.categoryId).slice(0, 4)
}

export async function createReview(payload: {
  product_id: string
  customer_name: string
  rating: number
  title?: string
  content?: string
}): Promise<void> {
  await http.post("/reviews", payload)
}
