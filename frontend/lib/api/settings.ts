import { http } from "./http"
import type {
  CompanySettings,
  NotificationSetting,
  PaymentMethod,
  Role,
  ShippingZone,
  StoreSettings,
  TaxRate,
  User,
} from "@/lib/types"

// =============================================================================
// Store & Company
// =============================================================================

interface StoreDTO {
  name: string
  description: string
  currency: string
  timezone: string
  order_prefix: string
  low_stock_threshold: number
  enable_reviews: boolean
  enable_wishlist: boolean
  enable_gift_cards: boolean
  default_weight_unit: string
  default_dimension_unit: string
}

export async function getStoreSettings(): Promise<StoreSettings> {
  const d = await http.get<StoreDTO>("/settings/store")
  return {
    name: d.name,
    description: d.description,
    currency: d.currency,
    timezone: d.timezone,
    orderPrefix: d.order_prefix,
    lowStockThreshold: d.low_stock_threshold,
    enableReviews: d.enable_reviews,
    enableWishlist: d.enable_wishlist,
    enableGiftCards: d.enable_gift_cards,
    defaultWeightUnit: d.default_weight_unit,
    defaultDimensionUnit: d.default_dimension_unit,
  }
}

export async function updateStoreSettings(payload: Partial<Record<string, unknown>>): Promise<StoreSettings> {
  const d = await http.put<StoreDTO>("/settings/store", payload)
  return (await getStoreSettingsFromDTO(d))
}

async function getStoreSettingsFromDTO(d: StoreDTO): Promise<StoreSettings> {
  return {
    name: d.name,
    description: d.description,
    currency: d.currency,
    timezone: d.timezone,
    orderPrefix: d.order_prefix,
    lowStockThreshold: d.low_stock_threshold,
    enableReviews: d.enable_reviews,
    enableWishlist: d.enable_wishlist,
    enableGiftCards: d.enable_gift_cards,
    defaultWeightUnit: d.default_weight_unit,
    defaultDimensionUnit: d.default_dimension_unit,
  }
}

interface CompanyDTO {
  name: string
  legal_name: string
  logo: string
  email: string
  phone: string
  website: string
  address: Record<string, string>
  tax_id: string
  currency: string
  timezone: string
  date_format: string
}

export async function getCompanySettings(): Promise<CompanySettings> {
  const d = await http.get<CompanyDTO>("/settings/company")
  const addr = d.address ?? {}
  return {
    name: d.name,
    legalName: d.legal_name,
    logo: d.logo,
    email: d.email,
    phone: d.phone,
    website: d.website,
    address: {
      id: "",
      label: (addr.label as string) ?? "office",
      line1: (addr.line1 as string) ?? "",
      line2: (addr.line2 as string) ?? "",
      city: (addr.city as string) ?? "",
      state: (addr.state as string) ?? "",
      zip: (addr.zip as string) ?? "",
      country: (addr.country as string) ?? "Indonesia",
      isDefault: true,
    },
    taxId: d.tax_id,
    currency: d.currency,
    timezone: d.timezone,
    dateFormat: d.date_format,
  }
}

export async function updateCompanySettings(
  payload: Partial<Record<string, unknown>>,
): Promise<CompanySettings> {
  await http.put("/settings/company", payload)
  return getCompanySettings()
}

// =============================================================================
// Tax rates
// =============================================================================

interface TaxRateDTO {
  id: string
  name: string
  rate: number
  region: string
  type: string
  applies_to: string
  is_active: boolean
}

function mapTaxRate(t: TaxRateDTO): TaxRate {
  return {
    id: t.id,
    name: t.name,
    rate: t.rate,
    region: t.region,
    type: t.type as TaxRate["type"],
    appliesTo: t.applies_to as TaxRate["appliesTo"],
    isActive: t.is_active,
  }
}

export async function getTaxRates(): Promise<TaxRate[]> {
  const res = await http.get<TaxRateDTO[]>("/settings/tax-rates")
  return res.map(mapTaxRate)
}

export async function createTaxRate(payload: Record<string, unknown>): Promise<TaxRate> {
  const res = await http.post<TaxRateDTO>("/settings/tax-rates", payload)
  return mapTaxRate(res)
}

export async function updateTaxRate(id: string, payload: Record<string, unknown>): Promise<TaxRate> {
  const res = await http.put<TaxRateDTO>(`/settings/tax-rates/${id}`, payload)
  return mapTaxRate(res)
}

export async function deleteTaxRate(id: string): Promise<void> {
  await http.delete(`/settings/tax-rates/${id}`)
}

// =============================================================================
// Shipping zones
// =============================================================================

interface ShippingZoneDTO {
  id: string
  name: string
  regions: string[]
  method: string
  rate: number
  free_above: number | null
  estimated_days: string
  is_active: boolean
}

function mapShippingZone(z: ShippingZoneDTO): ShippingZone {
  return {
    id: z.id,
    name: z.name,
    regions: z.regions ?? [],
    method: z.method as ShippingZone["method"],
    rate: z.rate,
    freeAbove: z.free_above ?? undefined,
    estimatedDays: z.estimated_days,
    isActive: z.is_active,
  }
}

export async function getShippingZones(): Promise<ShippingZone[]> {
  const res = await http.get<ShippingZoneDTO[]>("/settings/shipping-zones")
  return res.map(mapShippingZone)
}

export async function createShippingZone(payload: Record<string, unknown>): Promise<ShippingZone> {
  const res = await http.post<ShippingZoneDTO>("/settings/shipping-zones", payload)
  return mapShippingZone(res)
}

export async function updateShippingZone(id: string, payload: Record<string, unknown>): Promise<ShippingZone> {
  const res = await http.put<ShippingZoneDTO>(`/settings/shipping-zones/${id}`, payload)
  return mapShippingZone(res)
}

export async function deleteShippingZone(id: string): Promise<void> {
  await http.delete(`/settings/shipping-zones/${id}`)
}

// =============================================================================
// Payment methods (payment module)
// =============================================================================

interface PaymentMethodDTO {
  id: string
  name: string
  type: string
  is_enabled: boolean
  description: string
  config?: Record<string, string>
}

export interface PaymentMethodOption extends PaymentMethod {
  config?: Record<string, string>
}

function mapPaymentMethod(m: PaymentMethodDTO): PaymentMethodOption {
  return {
    id: m.id,
    name: m.name,
    type: m.type as PaymentMethod["type"],
    isEnabled: m.is_enabled,
    description: m.description,
    config: m.config,
  }
}

export async function getPaymentMethods(): Promise<PaymentMethodOption[]> {
  const res = await http.get<PaymentMethodDTO[]>("/payments/methods")
  return res.map(mapPaymentMethod)
}

export async function updatePaymentMethod(id: string, payload: Record<string, unknown>): Promise<PaymentMethod> {
  const res = await http.put<PaymentMethodDTO>(`/payments/methods/${id}`, payload)
  return mapPaymentMethod(res)
}

export async function createPaymentMethod(payload: Record<string, unknown>): Promise<PaymentMethod> {
  const res = await http.post<PaymentMethodDTO>("/payments/methods", payload)
  return mapPaymentMethod(res)
}

// =============================================================================
// Notifications
// =============================================================================

interface NotificationDTO {
  id: string
  event: string
  email: boolean
  sms: boolean
  in_app: boolean
  description: string
}

function mapNotification(n: NotificationDTO): NotificationSetting {
  return {
    id: n.id,
    event: n.event,
    email: n.email,
    sms: n.sms,
    inApp: n.in_app,
    description: n.description,
  }
}

export async function getNotificationSettings(): Promise<NotificationSetting[]> {
  const res = await http.get<NotificationDTO[]>("/settings/notifications")
  return res.map(mapNotification)
}

export async function updateNotificationSetting(
  id: string,
  payload: Record<string, unknown>,
): Promise<NotificationSetting> {
  const res = await http.put<NotificationDTO>(`/settings/notifications/${id}`, payload)
  return mapNotification(res)
}

// =============================================================================
// Users & roles
// =============================================================================

interface UserDTO {
  id: string
  name: string
  email: string
  avatar: string
  role: string
  status: string
  last_login: string | null
  created_at: string
}

interface RoleDTO {
  id: string
  name: string
  description: string
  permissions: number
  users: number
  is_system: boolean
}

function mapUser(u: UserDTO): User {
  return {
    id: u.id,
    name: u.name,
    email: u.email,
    avatar: u.avatar || undefined,
    role: u.role,
    status: u.status as User["status"],
    lastLogin: u.last_login ?? undefined,
    createdAt: u.created_at,
  }
}

export async function getUsers(): Promise<User[]> {
  const res = await http.get<{ data: UserDTO[] }>("/users?limit=100")
  return res.data.map(mapUser)
}

export async function createUser(payload: {
  name: string
  email: string
  password: string
  role_id: string
}): Promise<void> {
  await http.post("/users", payload)
}

export async function updateUser(id: string, payload: Record<string, unknown>): Promise<void> {
  await http.put(`/users/${id}`, payload)
}

export async function getRoles(): Promise<Role[]> {
  const res = await http.get<RoleDTO[]>("/roles")
  return res.map((r) => ({
    id: r.id,
    name: r.name,
    description: r.description,
    permissions: r.permissions,
    users: r.users,
  }))
}

export async function createRole(payload: Record<string, unknown>): Promise<void> {
  await http.post("/roles", payload)
}

export async function updateRole(id: string, payload: Record<string, unknown>): Promise<void> {
  await http.put(`/roles/${id}`, payload)
}
