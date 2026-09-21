import { getProducts } from "./products"
import { http, type Paginated } from "./http"
import type { IncomingShipment, InventoryItem, StockMovement } from "@/lib/types"

interface MovementDTO {
  id: string
  product_id: string
  type: string
  quantity: number
  reference: string
  note: string
  user: string
  date: string
}

interface ShipmentDTO {
  id: string
  product_id: string
  quantity: number
  expected_date: string | null
  supplier: string
  status: string
  notes: string
  created_at: string
}

interface WarehouseDTO {
  id: string
  name: string
  address: string
  is_default: boolean
}

function inventoryStatus(stock: number, threshold: number): InventoryItem["status"] {
  if (stock <= 0) return "out_of_stock"
  if (stock <= threshold) return "low_stock"
  if (stock > threshold * 4) return "overstocked"
  return "in_stock"
}

export async function getInventoryItems(): Promise<InventoryItem[]> {
  const [products, warehouses] = await Promise.all([
    getProducts(),
    http.get<WarehouseDTO[]>("/inventory/warehouses"),
  ])
  const defaultWarehouse = warehouses.find((w) => w.is_default) ?? warehouses[0]

  return products.map((p) => ({
    id: p.id,
    productId: p.id,
    productName: p.name,
    productSku: p.sku,
    productImage: p.images[0] ?? "",
    currentStock: p.stock,
    reserved: p.reserved,
    available: Math.max(p.stock - p.reserved, 0),
    incoming: p.incoming,
    lowStockThreshold: p.lowStockThreshold,
    status: inventoryStatus(p.stock, p.lowStockThreshold),
    warehouse: defaultWarehouse?.name ?? "Gudang Utama",
    lastRestocked: "",
  }))
}

export async function getStockMovements(): Promise<StockMovement[]> {
  const [res, products] = await Promise.all([
    http.get<Paginated<MovementDTO>>("/inventory/stock-movements?limit=100"),
    getProducts(),
  ])
  const productMap = new Map(products.map((p) => [p.id, p]))

  return res.data.map((m) => {
    const product = productMap.get(m.product_id)
    return {
      id: m.id,
      productId: m.product_id,
      productName: product?.name ?? "",
      productSku: product?.sku ?? "",
      type: m.type as StockMovement["type"],
      quantity: m.quantity,
      reference: m.reference,
      note: m.note,
      date: m.date,
      user: m.user,
    }
  })
}

export async function getIncomingShipments(): Promise<IncomingShipment[]> {
  const [res, products] = await Promise.all([
    http.get<Paginated<ShipmentDTO>>("/inventory/shipments?limit=100"),
    getProducts(),
  ])
  const productMap = new Map(products.map((p) => [p.id, p]))

  return res.data.map((s) => {
    const product = productMap.get(s.product_id)
    return {
      id: s.id,
      productId: s.product_id,
      productName: product?.name ?? "",
      productSku: product?.sku ?? "",
      quantity: s.quantity,
      expectedDate: s.expected_date ?? "",
      supplier: s.supplier,
      status: s.status as IncomingShipment["status"],
      notes: s.notes || undefined,
    }
  })
}
