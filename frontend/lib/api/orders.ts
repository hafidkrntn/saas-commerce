import { http, type Paginated } from "./http"
import type { Address, Order, OrderItem, OrderStatus, PaymentStatus, TimelineEvent } from "@/lib/types"

interface OrderItemDTO {
  id: string
  product_id: string
  product_name: string
  product_sku: string
  product_image: string
  quantity: number
  unit_price: number
  subtotal: number
}

interface TimelineDTO {
  id: string
  type: string
  title: string
  description: string
  user: string
  timestamp: string
}

interface OrderDTO {
  id: string
  order_number: string
  customer: { id: string; name: string; email: string }
  status: OrderStatus
  payment_status: PaymentStatus
  payment_method_id: string | null
  payment_method: string
  shipping_method: string
  shipping_address: Record<string, string> | null
  subtotal: number
  shipping_cost: number
  tax: number
  discount: number
  total: number
  notes: string
  items: OrderItemDTO[]
  timeline: TimelineDTO[]
  created_at: string
  updated_at: string
}

function mapAddress(raw: Record<string, string> | null): Address {
  const a = raw ?? {}
  return {
    id: "",
    label: (a.label as string) ?? "home",
    line1: (a.line1 as string) ?? "",
    line2: (a.line2 as string) ?? "",
    city: (a.city as string) ?? "",
    state: (a.state as string) ?? "",
    zip: (a.zip as string) ?? "",
    country: (a.country as string) ?? "Indonesia",
    isDefault: false,
  }
}

function mapTimeline(t: TimelineDTO): TimelineEvent {
  return {
    id: t.id,
    type: t.type as TimelineEvent["type"],
    title: t.title,
    description: t.description || undefined,
    timestamp: t.timestamp,
    user: t.user || undefined,
  }
}

function mapOrder(o: OrderDTO): Order {
  return {
    id: o.id,
    orderNumber: o.order_number,
    customer: { id: o.customer.id, name: o.customer.name, email: o.customer.email },
    status: o.status,
    paymentStatus: o.payment_status,
    paymentMethod: o.payment_method,
    shippingMethod: o.shipping_method,
    shippingAddress: mapAddress(o.shipping_address),
    items: o.items.map(
      (i): OrderItem => ({
        id: i.id,
        productId: i.product_id,
        productName: i.product_name,
        productSku: i.product_sku,
        productImage: i.product_image,
        quantity: i.quantity,
        unitPrice: i.unit_price,
        subtotal: i.subtotal,
      }),
    ),
    subtotal: o.subtotal,
    shippingCost: o.shipping_cost,
    tax: o.tax,
    discount: o.discount,
    total: o.total,
    notes: o.notes || undefined,
    timeline: (o.timeline ?? []).map(mapTimeline),
    createdAt: o.created_at,
    updatedAt: o.updated_at,
  }
}

export async function getOrders(): Promise<Order[]> {
  const res = await http.get<Paginated<OrderDTO>>("/orders?limit=100")
  return res.data.map(mapOrder)
}

export async function getOrder(id: string): Promise<Order | null> {
  try {
    const res = await http.get<OrderDTO>(`/orders/${id}`)
    return mapOrder(res)
  } catch {
    return null
  }
}

export async function getRecentOrders(limit = 5): Promise<Order[]> {
  const orders = await getOrders()
  return orders.slice(0, limit)
}

export interface OrderPayload {
  customer_id?: string | null
  customer_name: string
  customer_email?: string
  items: { product_id: string; quantity: number }[]
  payment_method_id?: string | null
  shipping_method?: string
  shipping_address?: Record<string, unknown>
  shipping_cost?: number
  discount?: number
  notes?: string
}

export async function createOrder(payload: OrderPayload): Promise<Order> {
  const res = await http.post<OrderDTO>("/orders", payload)
  return mapOrder(res)
}

export async function updateOrderStatus(id: string, status: string): Promise<Order> {
  const res = await http.patch<OrderDTO>(`/orders/${id}/status`, { status })
  return mapOrder(res)
}
