"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { Package, ArrowLeft, CreditCard, ChevronDown, ChevronUp, Lock } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Separator } from "@/components/ui/separator"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { toast } from "sonner"
import { formatCurrency } from "@/lib/utils"
import { useCart } from "@/lib/cart-context"
import { createOrder } from "@/lib/api/orders"
import { getPaymentMethods } from "@/lib/api/settings"
import type { PaymentMethod } from "@/lib/types"

const SHIPPING_COST = 25000

export default function CheckoutPage() {
  const router = useRouter()
  const { items, subtotal, clearCart } = useCart()
  const [loading, setLoading] = useState(false)
  const [summaryOpen, setSummaryOpen] = useState(false)
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([])
  const [paymentMethodId, setPaymentMethodId] = useState("")

  useEffect(() => {
    getPaymentMethods()
      .then((methods) => {
        const enabled = methods.filter((m) => m.isEnabled)
        setPaymentMethods(enabled)
        if (enabled.length > 0) setPaymentMethodId(enabled[0].id)
      })
      .catch(() => setPaymentMethods([]))
  }, [])

  const shipping = subtotal >= 500000 ? 0 : SHIPPING_COST
  const total = subtotal + shipping

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (items.length === 0) {
      toast.error("Your cart is empty")
      return
    }
    setLoading(true)

    const val = (id: string) => (document.getElementById(id) as HTMLInputElement | null)?.value ?? ""
    const fullName = `${val("fname")} ${val("lname")}`.trim()

    try {
      const order = await createOrder({
        customer_name: fullName || "Guest",
        customer_email: val("email"),
        items: items.map((i) => ({ product_id: i.product.id, quantity: i.quantity })),
        payment_method_id: paymentMethodId || null,
        shipping_method: "JNE",
        shipping_cost: shipping,
        shipping_address: {
          line1: val("address"),
          city: val("city"),
          state: val("province"),
          zip: val("postal"),
          country: val("country") || "Indonesia",
        },
      })

      clearCart()
      toast.success("Order placed! Lanjutkan ke pembayaran.")
      router.push(`/checkout/pay?order=${order.id}`)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to place order")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="mx-auto max-w-4xl px-4 py-6 sm:px-6 lg:px-8 pb-28 lg:pb-8">
      <button onClick={() => router.push("/cart")} className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors mb-4">
        <ArrowLeft className="h-3.5 w-3.5" /> Back to cart
      </button>

      <h1 className="text-xl font-bold tracking-tight mb-6 lg:text-2xl lg:mb-8">Checkout</h1>

      {items.length === 0 && (
        <div className="rounded-xl border bg-card p-10 text-center">
          <Package className="mx-auto h-10 w-10 text-muted-foreground/40" />
          <p className="mt-3 text-sm text-muted-foreground">Keranjang kosong. Tambahkan produk dulu.</p>
          <Button asChild className="mt-4"><Link href="/products">Browse products</Link></Button>
        </div>
      )}

      {items.length > 0 && (
        <form onSubmit={handleSubmit}>
          {/* Mobile: collapsible order summary */}
          <div className="mb-6 lg:hidden">
            <button
              type="button"
              onClick={() => setSummaryOpen(!summaryOpen)}
              className="flex w-full items-center justify-between rounded-xl border bg-card p-4 text-left shadow-sm"
            >
              <div>
                <p className="text-xs text-muted-foreground">Total order</p>
                <p className="text-lg font-bold">{formatCurrency(total)}</p>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-xs text-muted-foreground">{items.length} items</span>
                {summaryOpen ? <ChevronUp className="h-4 w-4 text-muted-foreground" /> : <ChevronDown className="h-4 w-4 text-muted-foreground" />}
              </div>
            </button>
            {summaryOpen && (
              <div className="mt-2 rounded-xl border bg-card p-4 shadow-sm space-y-2">
                {items.map(({ product, quantity }) => (
                  <div key={product.id} className="flex justify-between text-sm">
                    <span className="text-muted-foreground truncate">{product.name} × {quantity}</span>
                    <span className="tabular-nums shrink-0">{formatCurrency(product.price * quantity)}</span>
                  </div>
                ))}
                <Separator />
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Subtotal</span>
                  <span className="tabular-nums">{formatCurrency(subtotal)}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Shipping</span>
                  <span className="tabular-nums">{shipping === 0 ? "Free" : formatCurrency(shipping)}</span>
                </div>
                <Separator />
                <div className="flex justify-between font-semibold">
                  <span>Total</span>
                  <span className="tabular-nums">{formatCurrency(total)}</span>
                </div>
              </div>
            )}
          </div>

          <div className="grid gap-6 lg:gap-8 lg:grid-cols-5">
            <div className="space-y-6 lg:col-span-3">
              <Card>
                <CardHeader>
                  <CardTitle>Shipping</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="grid gap-3 sm:grid-cols-2">
                    <div className="space-y-1">
                      <Label className="text-xs" htmlFor="fname">First name</Label>
                      <Input id="fname" placeholder="John" className="h-10 lg:h-9" required />
                    </div>
                    <div className="space-y-1">
                      <Label className="text-xs" htmlFor="lname">Last name</Label>
                      <Input id="lname" placeholder="Doe" className="h-10 lg:h-9" required />
                    </div>
                  </div>
                  <div className="space-y-1">
                    <Label className="text-xs" htmlFor="email">Email</Label>
                    <Input id="email" type="email" placeholder="john@example.com" className="h-10 lg:h-9" required />
                  </div>
                  <div className="space-y-1">
                    <Label className="text-xs" htmlFor="address">Address</Label>
                    <Input id="address" placeholder="Jl. Contoh No. 123" className="h-10 lg:h-9" required />
                  </div>
                  <div className="grid gap-3 grid-cols-2 sm:grid-cols-4">
                    <div className="space-y-1">
                      <Label className="text-xs" htmlFor="city">City</Label>
                      <Input id="city" placeholder="Jakarta" className="h-10 lg:h-9" required />
                    </div>
                    <div className="space-y-1">
                      <Label className="text-xs" htmlFor="province">Province</Label>
                      <Input id="province" placeholder="DKI" className="h-10 lg:h-9" required />
                    </div>
                    <div className="space-y-1">
                      <Label className="text-xs" htmlFor="postal">Postal</Label>
                      <Input id="postal" placeholder="12345" className="h-10 lg:h-9" required />
                    </div>
                    <div className="space-y-1">
                      <Label className="text-xs" htmlFor="country">Country</Label>
                      <Input id="country" placeholder="Indonesia" className="h-10 lg:h-9" defaultValue="Indonesia" />
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Payment</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="space-y-1">
                    <Label className="text-xs">Payment method</Label>
                    <Select value={paymentMethodId} onValueChange={setPaymentMethodId}>
                      <SelectTrigger className="h-10 lg:h-9">
                        <SelectValue placeholder="Select payment method" />
                      </SelectTrigger>
                      <SelectContent>
                        {paymentMethods.length === 0 && (
                          <SelectItem value="__none" disabled>No payment methods</SelectItem>
                        )}
                        {paymentMethods.map((m) => (
                          <SelectItem key={m.id} value={m.id}>{m.name} — {m.type.replace("_", " ")}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    Mode demo: setelah order dibuat, Anda akan diarahkan ke simulasi pembayaran.
                  </p>
                </CardContent>
              </Card>
            </div>

            {/* Desktop: sidebar summary */}
            <div className="hidden lg:block lg:col-span-2">
              <Card className="sticky top-24">
                <CardHeader>
                  <CardTitle>Order summary</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  {items.map(({ product, quantity }) => (
                    <div key={product.id} className="flex justify-between text-sm">
                      <span className="text-muted-foreground truncate">{product.name} × {quantity}</span>
                      <span className="tabular-nums shrink-0">{formatCurrency(product.price * quantity)}</span>
                    </div>
                  ))}
                  <Separator />
                  <div className="flex justify-between text-sm">
                    <span className="text-muted-foreground">Subtotal</span>
                    <span className="tabular-nums">{formatCurrency(subtotal)}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-muted-foreground">Shipping</span>
                    <span className="tabular-nums">{shipping === 0 ? "Free" : formatCurrency(shipping)}</span>
                  </div>
                  <Separator />
                  <div className="flex justify-between font-semibold text-base">
                    <span>Total</span>
                    <span className="tabular-nums">{formatCurrency(total)}</span>
                  </div>
                  <Button type="submit" className="w-full h-10 gap-2" disabled={loading}>
                    {loading ? "Processing..." : <><CreditCard className="h-4 w-4" /> Place order</>}
                  </Button>
                  <p className="flex items-center justify-center gap-1 text-xs text-muted-foreground">
                    <Lock className="h-3 w-3" /> Secure checkout
                  </p>
                </CardContent>
              </Card>
            </div>
          </div>

          {/* Mobile: sticky bottom bar */}
          <div className="fixed bottom-0 left-0 right-0 z-40 border-t bg-background p-4 lg:hidden shadow-[0_-4px_12px_rgba(0,0,0,0.05)]">
            <div className="flex items-center justify-between mb-3">
              <div>
                <p className="text-xs text-muted-foreground">Total</p>
                <p className="text-lg font-bold">{formatCurrency(total)}</p>
              </div>
              <p className="text-xs text-muted-foreground flex items-center gap-1">
                <Lock className="h-3 w-3" /> Secure
              </p>
            </div>
            <Button type="submit" className="w-full h-11 gap-2 text-base" disabled={loading}>
              {loading ? "Processing..." : <><CreditCard className="h-4 w-4" /> Place order</>}
            </Button>
          </div>
        </form>
      )}

    </div>
  )
}
