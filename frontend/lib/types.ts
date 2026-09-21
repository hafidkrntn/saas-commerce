export interface Product {
  id: string
  name: string
  sku: string
  description: string
  price: number
  compareAtPrice?: number
  costPrice?: number
  stock: number
  reserved: number
  incoming: number
  lowStockThreshold: number
  category: string
  categoryId: string
  tags: string[]
  images: string[]
  status: "active" | "draft" | "archived"
  isActive: boolean
  rating: number
  reviewCount: number
  createdAt: string
  updatedAt: string
}

export interface Order {
  id: string
  orderNumber: string
  customer: CustomerBrief
  status: OrderStatus
  paymentStatus: PaymentStatus
  paymentMethod: string
  shippingMethod: string
  shippingAddress: Address
  items: OrderItem[]
  subtotal: number
  shippingCost: number
  tax: number
  discount: number
  total: number
  notes?: string
  timeline: TimelineEvent[]
  createdAt: string
  updatedAt: string
}

export type OrderStatus = "pending" | "confirmed" | "processing" | "shipped" | "delivered" | "cancelled" | "refunded"
export type PaymentStatus = "pending" | "paid" | "failed" | "refunded" | "partially_refunded"

export interface OrderItem {
  id: string
  productId: string
  productName: string
  productSku: string
  productImage: string
  quantity: number
  unitPrice: number
  subtotal: number
}

export interface CustomerBrief {
  id: string
  name: string
  email: string
  avatar?: string
}

export interface Customer {
  id: string
  name: string
  email: string
  phone: string
  avatar?: string
  status: "active" | "inactive" | "blocked"
  totalOrders: number
  totalSpent: number
  lifetimeValue: number
  averageOrderValue: number
  lastOrderDate: string
  createdAt: string
  addresses: Address[]
  tags: string[]
  recentActivity: ActivityEvent[]
}

export interface Address {
  id: string
  label: string
  line1: string
  line2?: string
  city: string
  state: string
  zip: string
  country: string
  isDefault: boolean
}

export interface TimelineEvent {
  id: string
  type: "created" | "confirmed" | "shipped" | "delivered" | "cancelled" | "payment" | "note" | "refunded"
  title: string
  description?: string
  timestamp: string
  user?: string
}

export interface ActivityEvent {
  id: string
  type: "order" | "login" | "review" | "support" | "payment"
  description: string
  timestamp: string
}

export interface RevenueData {
  month: string
  revenue: number
  profit: number
  orders: number
  expenses: number
}

export interface SalesData {
  date: string
  sales: number
  orders: number
}

export interface AnalyticsSummary {
  revenue: number
  revenueChange: number
  profit: number
  profitChange: number
  orders: number
  ordersChange: number
  conversionRate: number
  conversionChange: number
  customers: number
  customersChange: number
  averageOrderValue: number
  aovChange: number
}

export interface TopProduct {
  id: string
  name: string
  image: string
  revenue: number
  unitsSold: number
  percentage: number
}

export interface TopCategory {
  name: string
  revenue: number
  percentage: number
}

export interface InventoryItem {
  id: string
  productId: string
  productName: string
  productSku: string
  productImage: string
  currentStock: number
  reserved: number
  available: number
  incoming: number
  incomingDate?: string
  lowStockThreshold: number
  status: "in_stock" | "low_stock" | "out_of_stock" | "overstocked"
  warehouse: string
  lastRestocked: string
}

export interface StockMovement {
  id: string
  productId: string
  productName: string
  productSku: string
  type: "in" | "out" | "adjustment" | "return"
  quantity: number
  reference: string
  note: string
  date: string
  user: string
}

export interface IncomingShipment {
  id: string
  productId: string
  productName: string
  productSku: string
  quantity: number
  expectedDate: string
  supplier: string
  status: "scheduled" | "in_transit" | "delivered" | "delayed"
  notes?: string
}

export interface CompanySettings {
  name: string
  legalName: string
  logo: string
  email: string
  phone: string
  website: string
  address: Address
  taxId: string
  currency: string
  timezone: string
  dateFormat: string
}

export interface StoreSettings {
  name: string
  description: string
  currency: string
  timezone: string
  orderPrefix: string
  lowStockThreshold: number
  enableReviews: boolean
  enableWishlist: boolean
  enableGiftCards: boolean
  defaultWeightUnit: string
  defaultDimensionUnit: string
}

export interface TaxRate {
  id: string
  name: string
  rate: number
  region: string
  type: "vat" | "sales_tax" | "gst"
  isActive: boolean
  appliesTo: "all" | "digital" | "physical" | "services"
}

export interface ShippingZone {
  id: string
  name: string
  regions: string[]
  method: "flat_rate" | "free" | "calculated"
  rate: number
  freeAbove?: number
  estimatedDays: string
  isActive: boolean
}

export interface PaymentMethod {
  id: string
  name: string
  type: "card" | "bank_transfer" | "cod" | "e_wallet"
  isEnabled: boolean
  description: string
}

export interface NotificationSetting {
  id: string
  event: string
  email: boolean
  sms: boolean
  inApp: boolean
  description: string
}

export interface User {
  id: string
  name: string
  email: string
  avatar?: string
  role: string
  status: "active" | "inactive" | "invited"
  lastLogin?: string
  createdAt: string
}

export interface Role {
  id: string
  name: string
  description: string
  permissions: number
  users: number
}

export interface Review {
  id: string
  productId: string
  customerName: string
  customerAvatar?: string
  rating: number
  title: string
  content: string
  isVerified: boolean
  createdAt: string
}
