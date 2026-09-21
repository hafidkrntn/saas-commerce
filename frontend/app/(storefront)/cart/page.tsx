"use client"

import Link from "next/link"
import { Package, Minus, Plus, Trash2, ArrowRight, ShoppingCart as CartIcon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { formatCurrency } from "@/lib/utils"
import { useCart } from "@/lib/cart-context"

export default function CartPage() {
  const { items, updateQuantity, removeItem, subtotal } = useCart()

  const shipping = subtotal >= 500000 ? 0 : 25000
  const total = subtotal + shipping

  if (items.length === 0) {
    return (
      <div className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8 text-center">
        <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-muted mx-auto">
          <CartIcon className="h-8 w-8 text-muted-foreground/40" />
        </div>
        <h1 className="mt-6 text-2xl font-bold tracking-tight">Your cart is empty</h1>
        <p className="mt-2 text-sm text-muted-foreground">Looks like you haven&apos;t added anything yet.</p>
        <Button asChild className="mt-6">
          <Link href="/products">Browse products <ArrowRight className="h-4 w-4" /></Link>
        </Button>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8 pb-28 lg:pb-8">
      <h1 className="text-xl font-bold tracking-tight mb-6 lg:text-2xl lg:mb-8">
        Shopping cart ({items.length} {items.length === 1 ? "item" : "items"})
      </h1>

      <div className="grid gap-6 lg:gap-8 lg:grid-cols-3">
        <div className="lg:col-span-2 space-y-3 lg:space-y-4">
          {items.map(({ product, quantity }) => (
            <Card key={product.id}>
              <CardContent className="p-3 lg:p-4 flex items-center gap-3 lg:gap-4">
                <div className="flex h-16 w-16 lg:h-20 lg:w-20 shrink-0 items-center justify-center rounded-lg bg-muted overflow-hidden">
                  {product.images[0] ? (
                    // eslint-disable-next-line @next/next/no-img-element
                    <img src={product.images[0]} alt={product.name} className="h-full w-full object-cover" />
                  ) : (
                    <Package className="h-6 w-6 lg:h-8 lg:w-8 text-muted-foreground/30" />
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <Link href={`/products/${product.id}`} className="text-sm font-medium hover:text-primary transition-colors truncate block">
                    {product.name}
                  </Link>
                  <p className="text-xs text-muted-foreground font-mono mt-0.5">{product.sku}</p>
                  <p className="text-sm font-semibold mt-1">{formatCurrency(product.price)}</p>
                </div>
                <div className="flex items-center rounded-lg border shrink-0">
                  <button onClick={() => updateQuantity(product.id, quantity - 1)} className="flex h-7 w-7 lg:h-8 lg:w-8 items-center justify-center text-muted-foreground hover:text-foreground transition-colors">
                    <Minus className="h-3 w-3" />
                  </button>
                  <span className="flex h-7 w-8 lg:h-8 lg:w-10 items-center justify-center text-sm font-medium tabular-nums border-x">
                    {quantity}
                  </span>
                  <button onClick={() => updateQuantity(product.id, quantity + 1)} className="flex h-7 w-7 lg:h-8 lg:w-8 items-center justify-center text-muted-foreground hover:text-foreground transition-colors">
                    <Plus className="h-3 w-3" />
                  </button>
                </div>
                <div className="hidden sm:block text-sm font-semibold tabular-nums w-20 text-right shrink-0">
                  {formatCurrency(product.price * quantity)}
                </div>
                <button onClick={() => removeItem(product.id)} className="flex h-7 w-7 lg:h-8 lg:w-8 items-center justify-center rounded-lg text-muted-foreground hover:text-destructive hover:bg-destructive/10 transition-colors shrink-0">
                  <Trash2 className="h-3.5 w-3.5" />
                </button>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Desktop summary sidebar */}
        <div className="hidden lg:block">
          <Card className="sticky top-24">
            <CardContent className="p-5 space-y-4">
              <h3 className="text-sm font-semibold">Order summary</h3>
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Subtotal</span>
                  <span className="tabular-nums">{formatCurrency(subtotal)}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Shipping</span>
                  <span className="tabular-nums">{shipping === 0 ? "Free" : formatCurrency(shipping)}</span>
                </div>
                {shipping === 0 && <p className="text-xs text-success">Gratis ongkir minimal belanja Rp 500.000</p>}
              </div>
              <Separator />
              <div className="flex justify-between font-semibold">
                <span>Total</span>
                <span className="tabular-nums">{formatCurrency(total)}</span>
              </div>
              <Button asChild className="w-full h-10 gap-2">
                <Link href="/checkout">Checkout <ArrowRight className="h-4 w-4" /></Link>
              </Button>
              <Button asChild variant="outline" className="w-full h-9 text-xs">
                <Link href="/products">Continue shopping</Link>
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Mobile sticky bottom bar */}
      <div className="fixed bottom-0 left-0 right-0 z-40 border-t bg-background p-4 lg:hidden shadow-[0_-4px_12px_rgba(0,0,0,0.05)]">
        <div className="flex items-center justify-between mb-3">
          <div>
            <p className="text-xs text-muted-foreground">Total</p>
            <p className="text-lg font-bold">{formatCurrency(total)}</p>
          </div>
          <p className="text-xs text-muted-foreground">{shipping === 0 ? "Free shipping" : `+${formatCurrency(shipping)} shipping`}</p>
        </div>
        <Button asChild className="w-full h-11 gap-2 text-base">
          <Link href="/checkout">Checkout <ArrowRight className="h-4 w-4" /></Link>
        </Button>
      </div>
    </div>
  )
}
