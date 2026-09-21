"use client"

import { useEffect, useState } from "react"
import {
  Building2,
  Store,
  Percent,
  Truck,
  CreditCard,
  Bell,
  Users,
  Shield,
  Palette,
  Save,
  Plus,
  Trash2,
} from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import { Separator } from "@/components/ui/separator"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { PageHeader } from "@/components/shared/page-header"
import { StatusBadge } from "@/components/shared/status-badge"
import { Skeleton } from "@/components/ui/skeleton"
import { getInitials, formatCurrency } from "@/lib/utils"
import { toast } from "sonner"
import {
  getCompanySettings,
  getStoreSettings,
  getTaxRates,
  getShippingZones,
  getPaymentMethods,
  getNotificationSettings,
  getUsers,
  getRoles,
  updateCompanySettings,
  updateStoreSettings,
  createTaxRate,
  deleteTaxRate,
  createShippingZone,
  deleteShippingZone,
  updatePaymentMethod,
  updateNotificationSetting,
} from "@/lib/api/settings"
import type {
  CompanySettings, StoreSettings, TaxRate, ShippingZone,
  PaymentMethod, NotificationSetting, User, Role,
} from "@/lib/types"

const tabs = [
  { id: "company", label: "Company", icon: Building2 },
  { id: "store", label: "Store", icon: Store },
  { id: "tax", label: "Tax", icon: Percent },
  { id: "shipping", label: "Shipping", icon: Truck },
  { id: "payment", label: "Payment", icon: CreditCard },
  { id: "notifications", label: "Notifications", icon: Bell },
  { id: "users", label: "Users", icon: Users },
  { id: "roles", label: "Roles", icon: Shield },
  { id: "appearance", label: "Appearance", icon: Palette },
]

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState("company")
  const [loading, setLoading] = useState(true)
  const [company, setCompany] = useState<CompanySettings | null>(null)
  const [store, setStore] = useState<StoreSettings | null>(null)
  const [taxRates, setTaxRates] = useState<TaxRate[]>([])
  const [shippingZones, setShippingZones] = useState<ShippingZone[]>([])
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([])
  const [notifications, setNotifications] = useState<NotificationSetting[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [roles, setRoles] = useState<Role[]>([])

  useEffect(() => {
    async function load() {
      const [c, s, t, sz, pm, n, u, r] = await Promise.all([
        getCompanySettings(), getStoreSettings(), getTaxRates(), getShippingZones(),
        getPaymentMethods(), getNotificationSettings(), getUsers(), getRoles(),
      ])
      setCompany(c); setStore(s); setTaxRates(t); setShippingZones(sz)
      setPaymentMethods(pm); setNotifications(n); setUsers(u); setRoles(r)
      setLoading(false)
    }
    load()
  }, [])

  if (loading) {
    return (
      <div className="p-6">
        <PageHeader title="Settings" description="Manage your store settings" />
        <div className="flex gap-8">
          <div className="w-48 space-y-2">
            {Array.from({ length: 6 }).map((_, i) => <Skeleton key={i} className="h-10 w-full rounded-xl" />)}
          </div>
          <div className="flex-1">
            <Skeleton className="h-64 rounded-xl" />
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="p-6">
      <PageHeader title="Settings" description="Manage your store settings" />

      <div className="flex flex-col gap-8 lg:flex-row">
        <nav className="flex shrink-0 gap-1 overflow-x-auto lg:w-48 lg:flex-col">
          {tabs.map((tab) => {
            const Icon = tab.icon
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 rounded-xl px-3 py-2.5 text-sm font-medium whitespace-nowrap transition-colors ${
                  activeTab === tab.id
                    ? "bg-primary/10 text-primary"
                    : "text-muted-foreground hover:bg-muted hover:text-foreground"
                }`}
              >
                <Icon className="h-4 w-4" />
                {tab.label}
              </button>
            )
          })}
        </nav>

        <div className="flex-1 min-w-0 space-y-6">
          {activeTab === "company" && company && (
            <CompanySettingsForm
              data={company}
              onSaved={(c) => setCompany(c)}
            />
          )}
          {activeTab === "store" && store && (
            <StoreSettingsForm
              data={store}
              onSaved={(s) => setStore(s)}
            />
          )}
          {activeTab === "tax" && (
            <TaxSettingsForm
              data={taxRates}
              onCreated={(t) => setTaxRates((prev) => [...prev, t])}
              onDeleted={(id) => setTaxRates((prev) => prev.filter((t) => t.id !== id))}
            />
          )}
          {activeTab === "shipping" && (
            <ShippingSettingsForm
              data={shippingZones}
              onCreated={(z) => setShippingZones((prev) => [...prev, z])}
              onDeleted={(id) => setShippingZones((prev) => prev.filter((z) => z.id !== id))}
            />
          )}
          {activeTab === "payment" && (
            <PaymentSettingsForm
              data={paymentMethods}
              onChanged={(m) => setPaymentMethods((prev) => prev.map((p) => (p.id === m.id ? m : p)))}
            />
          )}
          {activeTab === "notifications" && (
            <NotificationSettingsForm
              data={notifications}
              onChanged={(n) => setNotifications((prev) => prev.map((x) => (x.id === n.id ? n : x)))}
            />
          )}
          {activeTab === "users" && <UsersSettingsForm data={users} />}
          {activeTab === "roles" && <RolesSettingsForm data={roles} />}
          {activeTab === "appearance" && <AppearanceSettings />}
        </div>
      </div>
    </div>
  )
}

// =============================================================================
// Company
// =============================================================================

function CompanySettingsForm({
  data,
  onSaved,
}: {
  data: CompanySettings
  onSaved: (c: CompanySettings) => void
}) {
  const [saving, setSaving] = useState(false)

  const handleSave = async () => {
    const val = (id: string) => (document.getElementById(id) as HTMLInputElement | null)?.value ?? ""
    setSaving(true)
    try {
      const updated = await updateCompanySettings({
        name: val("c_name"),
        legal_name: val("c_legal"),
        email: val("c_email"),
        phone: val("c_phone"),
        website: val("c_website"),
        tax_id: val("c_taxid"),
        address: {
          line1: val("c_line1"),
          line2: val("c_line2"),
          city: val("c_city"),
          state: val("c_state"),
          zip: val("c_zip"),
          country: "Indonesia",
        },
      })
      onSaved(updated)
      toast.success("Company settings saved")
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save")
    } finally {
      setSaving(false)
    }
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Company Information</CardTitle>
        <Button size="sm" onClick={handleSave} disabled={saving}><Save className="h-3.5 w-3.5" /> Save</Button>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="c_name">Company name</Label>
            <Input id="c_name" defaultValue={data.name} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="c_legal">Legal name</Label>
            <Input id="c_legal" defaultValue={data.legalName} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="c_email">Email</Label>
            <Input id="c_email" defaultValue={data.email} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="c_phone">Phone</Label>
            <Input id="c_phone" defaultValue={data.phone} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="c_website">Website</Label>
            <Input id="c_website" defaultValue={data.website} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="c_taxid">Tax ID</Label>
            <Input id="c_taxid" defaultValue={data.taxId} />
          </div>
        </div>
        <Separator />
        <div className="space-y-2">
          <Label>Address</Label>
          <Input id="c_line1" defaultValue={data.address.line1} placeholder="Address line 1" />
          <Input id="c_line2" defaultValue={data.address.line2} placeholder="Address line 2" />
          <div className="grid gap-4 sm:grid-cols-3">
            <Input id="c_city" defaultValue={data.address.city} placeholder="City" />
            <Input id="c_state" defaultValue={data.address.state} placeholder="State" />
            <Input id="c_zip" defaultValue={data.address.zip} placeholder="ZIP code" />
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

// =============================================================================
// Store
// =============================================================================

function StoreSettingsForm({
  data,
  onSaved,
}: {
  data: StoreSettings
  onSaved: (s: StoreSettings) => void
}) {
  const [saving, setSaving] = useState(false)

  const handleSave = async () => {
    const val = (id: string) => (document.getElementById(id) as HTMLInputElement | null)?.value ?? ""
    setSaving(true)
    try {
      const updated = await updateStoreSettings({
        name: val("s_name"),
        currency: val("s_currency"),
        timezone: val("s_timezone"),
        order_prefix: val("s_prefix"),
        low_stock_threshold: Number(val("s_lowstock") || 5),
      })
      onSaved(updated)
      toast.success("Store settings saved")
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save")
    } finally {
      setSaving(false)
    }
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Store Settings</CardTitle>
        <Button size="sm" onClick={handleSave} disabled={saving}><Save className="h-3.5 w-3.5" /> Save</Button>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="s_name">Store name</Label>
            <Input id="s_name" defaultValue={data.name} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="s_currency">Currency</Label>
            <Input id="s_currency" defaultValue={data.currency} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="s_timezone">Timezone</Label>
            <Input id="s_timezone" defaultValue={data.timezone} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="s_prefix">Order prefix</Label>
            <Input id="s_prefix" defaultValue={data.orderPrefix} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="s_lowstock">Low stock threshold</Label>
            <Input id="s_lowstock" type="number" defaultValue={data.lowStockThreshold} />
          </div>
        </div>
        <Separator />
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div><Label>Enable reviews</Label><p className="text-xs text-muted-foreground">Allow customers to review products</p></div>
            <Switch defaultChecked={data.enableReviews} />
          </div>
          <div className="flex items-center justify-between">
            <div><Label>Enable wishlist</Label><p className="text-xs text-muted-foreground">Allow customers to save products to wishlist</p></div>
            <Switch defaultChecked={data.enableWishlist} />
          </div>
          <div className="flex items-center justify-between">
            <div><Label>Enable gift cards</Label><p className="text-xs text-muted-foreground">Allow customers to purchase gift cards</p></div>
            <Switch defaultChecked={data.enableGiftCards} />
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

// =============================================================================
// Tax rates
// =============================================================================

function TaxSettingsForm({
  data,
  onCreated,
  onDeleted,
}: {
  data: TaxRate[]
  onCreated: (t: TaxRate) => void
  onDeleted: (id: string) => void
}) {
  const [adding, setAdding] = useState(false)
  const [name, setName] = useState("")
  const [rate, setRate] = useState("")

  const handleAdd = async () => {
    if (!name || !rate) {
      toast.error("Name and rate are required")
      return
    }
    try {
      const created = await createTaxRate({ name, rate: Number(rate), type: "vat", applies_to: "all" })
      onCreated(created)
      setName(""); setRate(""); setAdding(false)
      toast.success("Tax rate added")
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to add tax rate")
    }
  }

  const handleDelete = async (id: string) => {
    try {
      await deleteTaxRate(id)
      onDeleted(id)
      toast.success("Tax rate deleted")
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete")
    }
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Tax Rates</CardTitle>
        <Button variant="outline" size="sm" onClick={() => setAdding(!adding)}><Plus className="h-4 w-4" /> Add rate</Button>
      </CardHeader>
      <CardContent>
        {adding && (
          <div className="mb-4 flex items-end gap-2 rounded-xl border p-3">
            <div className="flex-1 space-y-1">
              <Label className="text-xs">Name</Label>
              <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="e.g. PPN 11%" />
            </div>
            <div className="w-28 space-y-1">
              <Label className="text-xs">Rate (%)</Label>
              <Input type="number" value={rate} onChange={(e) => setRate(e.target.value)} placeholder="11" />
            </div>
            <Button size="sm" onClick={handleAdd}>Add</Button>
          </div>
        )}
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b text-left">
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Name</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Rate</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Region</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Type</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Status</th>
                <th className="pb-3 text-xs font-medium text-muted-foreground text-right">Actions</th>
              </tr>
            </thead>
            <tbody>
              {data.map((tax) => (
                <tr key={tax.id} className="border-b last:border-0">
                  <td className="py-3 pr-4 text-sm font-medium">{tax.name}</td>
                  <td className="py-3 pr-4 text-sm">{tax.rate}%</td>
                  <td className="py-3 pr-4 text-sm text-muted-foreground">{tax.region}</td>
                  <td className="py-3 pr-4">
                    <Badge variant="secondary" className="font-normal uppercase text-xs">{tax.type}</Badge>
                  </td>
                  <td className="py-3"><StatusBadge status={tax.isActive ? "active" : "inactive"} /></td>
                  <td className="py-3 text-right">
                    <Button variant="ghost" size="icon-sm" onClick={() => handleDelete(tax.id)}>
                      <Trash2 className="h-3.5 w-3.5 text-destructive" />
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  )
}

// =============================================================================
// Shipping zones
// =============================================================================

function ShippingSettingsForm({
  data,
  onCreated,
  onDeleted,
}: {
  data: ShippingZone[]
  onCreated: (z: ShippingZone) => void
  onDeleted: (id: string) => void
}) {
  const [adding, setAdding] = useState(false)
  const [name, setName] = useState("")
  const [rate, setRate] = useState("")

  const handleAdd = async () => {
    if (!name) {
      toast.error("Zone name is required")
      return
    }
    try {
      const created = await createShippingZone({
        name,
        regions: [],
        method: "flat_rate",
        rate: Number(rate || 0),
        estimated_days: "2-4 hari",
      })
      onCreated(created)
      setName(""); setRate(""); setAdding(false)
      toast.success("Shipping zone added")
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to add zone")
    }
  }

  const handleDelete = async (id: string) => {
    try {
      await deleteShippingZone(id)
      onDeleted(id)
      toast.success("Shipping zone deleted")
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete")
    }
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Shipping Zones</CardTitle>
        <Button variant="outline" size="sm" onClick={() => setAdding(!adding)}><Plus className="h-4 w-4" /> Add zone</Button>
      </CardHeader>
      <CardContent>
        {adding && (
          <div className="mb-4 flex items-end gap-2 rounded-xl border p-3">
            <div className="flex-1 space-y-1">
              <Label className="text-xs">Zone name</Label>
              <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="e.g. Jawa Barat" />
            </div>
            <div className="w-28 space-y-1">
              <Label className="text-xs">Rate (IDR)</Label>
              <Input type="number" value={rate} onChange={(e) => setRate(e.target.value)} placeholder="20000" />
            </div>
            <Button size="sm" onClick={handleAdd}>Add</Button>
          </div>
        )}
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b text-left">
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Zone</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Method</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground text-right">Rate</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Est. Days</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Status</th>
                <th className="pb-3 text-xs font-medium text-muted-foreground text-right">Actions</th>
              </tr>
            </thead>
            <tbody>
              {data.map((zone) => (
                <tr key={zone.id} className="border-b last:border-0">
                  <td className="py-3 pr-4 text-sm font-medium">{zone.name}</td>
                  <td className="py-3 pr-4 text-sm capitalize text-muted-foreground">{zone.method.replace("_", " ")}</td>
                  <td className="py-3 pr-4 text-right text-sm">
                    {zone.method === "free" ? "Free" : zone.method === "calculated" ? "Calculated" : formatCurrency(zone.rate)}
                  </td>
                  <td className="py-3 pr-4 text-sm text-muted-foreground">{zone.estimatedDays}</td>
                  <td className="py-3"><StatusBadge status={zone.isActive ? "active" : "inactive"} /></td>
                  <td className="py-3 text-right">
                    <Button variant="ghost" size="icon-sm" onClick={() => handleDelete(zone.id)}>
                      <Trash2 className="h-3.5 w-3.5 text-destructive" />
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  )
}

// =============================================================================
// Payment methods
// =============================================================================

function PaymentSettingsForm({
  data,
  onChanged,
}: {
  data: PaymentMethod[]
  onChanged: (m: PaymentMethod) => void
}) {
  const toggle = async (pm: PaymentMethod, enabled: boolean) => {
    try {
      const updated = await updatePaymentMethod(pm.id, { is_enabled: enabled })
      onChanged(updated)
      toast.success(`${pm.name} ${enabled ? "enabled" : "disabled"}`)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update")
    }
  }

  return (
    <Card>
      <CardHeader><CardTitle>Payment Methods</CardTitle></CardHeader>
      <CardContent className="space-y-4">
        {data.map((pm) => (
          <div key={pm.id} className="flex items-center justify-between p-3 rounded-xl border">
            <div>
              <p className="text-sm font-medium">{pm.name}</p>
              <p className="text-xs text-muted-foreground">{pm.description}</p>
            </div>
            <Switch checked={pm.isEnabled} onCheckedChange={(v) => toggle(pm, v)} />
          </div>
        ))}
      </CardContent>
    </Card>
  )
}

// =============================================================================
// Notifications
// =============================================================================

function NotificationSettingsForm({
  data,
  onChanged,
}: {
  data: NotificationSetting[]
  onChanged: (n: NotificationSetting) => void
}) {
  const toggle = async (n: NotificationSetting, field: "email" | "sms" | "inApp", value: boolean) => {
    try {
      const updated = await updateNotificationSetting(n.id, {
        email: field === "email" ? value : n.email,
        sms: field === "sms" ? value : n.sms,
        in_app: field === "inApp" ? value : n.inApp,
      })
      onChanged(updated)
      toast.success("Notification updated")
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update")
    }
  }

  return (
    <Card>
      <CardHeader><CardTitle>Notifications</CardTitle></CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b text-left">
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Event</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Email</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">SMS</th>
                <th className="pb-3 text-xs font-medium text-muted-foreground">In-App</th>
              </tr>
            </thead>
            <tbody>
              {data.map((n) => (
                <tr key={n.id} className="border-b last:border-0">
                  <td className="py-3 pr-4">
                    <p className="text-sm font-medium">{n.event}</p>
                    <p className="text-xs text-muted-foreground">{n.description}</p>
                  </td>
                  <td className="py-3 pr-4"><Switch checked={n.email} onCheckedChange={(v) => toggle(n, "email", v)} /></td>
                  <td className="py-3 pr-4"><Switch checked={n.sms} onCheckedChange={(v) => toggle(n, "sms", v)} /></td>
                  <td className="py-3"><Switch checked={n.inApp} onCheckedChange={(v) => toggle(n, "inApp", v)} /></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  )
}

// =============================================================================
// Users & roles (read-only views)
// =============================================================================

function UsersSettingsForm({ data }: { data: User[] }) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Team Members</CardTitle>
        <Button variant="outline" size="sm"><Plus className="h-4 w-4" /> Invite user</Button>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b text-left">
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">User</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Role</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Status</th>
                <th className="pb-3 text-xs font-medium text-muted-foreground">Last Login</th>
              </tr>
            </thead>
            <tbody>
              {data.map((user) => (
                <tr key={user.id} className="border-b last:border-0">
                  <td className="py-3 pr-4">
                    <div className="flex items-center gap-3">
                      <Avatar className="h-8 w-8">
                        <AvatarImage src={user.avatar} />
                        <AvatarFallback className="text-xs">{getInitials(user.name)}</AvatarFallback>
                      </Avatar>
                      <div>
                        <p className="text-sm font-medium">{user.name}</p>
                        <p className="text-xs text-muted-foreground">{user.email}</p>
                      </div>
                    </div>
                  </td>
                  <td className="py-3 pr-4 text-sm text-muted-foreground">{user.role}</td>
                  <td className="py-3 pr-4"><StatusBadge status={user.status} /></td>
                  <td className="py-3 text-sm text-muted-foreground">{user.lastLogin ? new Date(user.lastLogin).toLocaleDateString() : "—"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  )
}

function RolesSettingsForm({ data }: { data: Role[] }) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Roles & Permissions</CardTitle>
        <Button variant="outline" size="sm"><Plus className="h-4 w-4" /> Add role</Button>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b text-left">
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Role</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground">Description</th>
                <th className="pb-3 pr-4 text-xs font-medium text-muted-foreground text-right">Permissions</th>
                <th className="pb-3 text-xs font-medium text-muted-foreground text-right">Users</th>
              </tr>
            </thead>
            <tbody>
              {data.map((role) => (
                <tr key={role.id} className="border-b last:border-0">
                  <td className="py-3 pr-4 text-sm font-medium">{role.name}</td>
                  <td className="py-3 pr-4 text-sm text-muted-foreground">{role.description}</td>
                  <td className="py-3 pr-4 text-right text-sm">{role.permissions}</td>
                  <td className="py-3 text-right text-sm">{role.users}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  )
}

function AppearanceSettings() {
  return (
    <Card>
      <CardHeader><CardTitle>Appearance</CardTitle></CardHeader>
      <CardContent className="space-y-6">
        <div>
          <Label className="text-base">Accent Color</Label>
          <p className="text-sm text-muted-foreground mb-3">Choose your brand color</p>
          <div className="flex gap-2">
            {["#2563eb", "#059669", "#d97706", "#dc2626", "#7c3aed", "#0891b2"].map((color) => (
              <button
                key={color}
                className="h-10 w-10 rounded-xl border-2 border-border hover:scale-110 transition-transform"
                style={{ backgroundColor: color }}
              />
            ))}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
