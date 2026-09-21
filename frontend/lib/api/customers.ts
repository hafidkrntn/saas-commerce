import { http, type Paginated } from "./http"
import type { ActivityEvent, Address, Customer } from "@/lib/types"

interface AddressDTO {
  id: string
  label: string
  line1: string
  line2: string
  city: string
  state: string
  zip: string
  country: string
  is_default: boolean
}

interface CustomerDTO {
  id: string
  name: string
  email: string
  phone: string
  avatar: string
  status: "active" | "inactive" | "blocked"
  tags: string[]
  total_orders: number
  total_spent: number
  lifetime_value: number
  average_order_value: number
  last_order_date: string | null
  addresses: AddressDTO[]
  created_at: string
  updated_at: string
}

interface ActivityDTO {
  id: string
  type: string
  description: string
  timestamp: string
}

function mapAddress(a: AddressDTO): Address {
  return {
    id: a.id,
    label: a.label,
    line1: a.line1,
    line2: a.line2,
    city: a.city,
    state: a.state,
    zip: a.zip,
    country: a.country,
    isDefault: a.is_default,
  }
}

function mapCustomer(c: CustomerDTO): Customer {
  return {
    id: c.id,
    name: c.name,
    email: c.email,
    phone: c.phone,
    avatar: c.avatar || undefined,
    status: c.status,
    totalOrders: c.total_orders,
    totalSpent: c.total_spent,
    lifetimeValue: c.lifetime_value,
    averageOrderValue: c.average_order_value,
    lastOrderDate: c.last_order_date ?? "",
    createdAt: c.created_at,
    addresses: (c.addresses ?? []).map(mapAddress),
    tags: c.tags ?? [],
    recentActivity: [],
  }
}

function mapActivity(a: ActivityDTO): ActivityEvent {
  return {
    id: a.id,
    type: a.type as ActivityEvent["type"],
    description: a.description,
    timestamp: a.timestamp,
  }
}

export async function getCustomers(): Promise<Customer[]> {
  const res = await http.get<Paginated<CustomerDTO>>("/customers?limit=100")
  return res.data.map(mapCustomer)
}

export async function getCustomer(id: string): Promise<Customer | null> {
  try {
    const res = await http.get<CustomerDTO>(`/customers/${id}`)
    const customer = mapCustomer(res)
    customer.recentActivity = await getCustomerActivities(id)
    return customer
  } catch {
    return null
  }
}

export async function getCustomerActivities(id: string, limit = 20): Promise<ActivityEvent[]> {
  const res = await http.get<ActivityDTO[]>(`/customers/${id}/activities?limit=${limit}`)
  return res.map(mapActivity)
}

export interface CustomerPayload {
  name: string
  email?: string
  phone?: string
  avatar?: string
  status?: string
  tags?: string[]
}

export async function createCustomer(payload: CustomerPayload): Promise<Customer> {
  const res = await http.post<CustomerDTO>("/customers", payload)
  return mapCustomer(res)
}

export async function updateCustomer(id: string, payload: Partial<CustomerPayload>): Promise<Customer> {
  const res = await http.put<CustomerDTO>(`/customers/${id}`, payload)
  return mapCustomer(res)
}
