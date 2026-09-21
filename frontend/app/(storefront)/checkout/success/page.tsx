"use client"

import { Suspense, useEffect, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import Link from "next/link"
import { CheckCircle2, Package } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { Separator } from "@/components/ui/separator"
import { formatCurrency, formatDate } from "@/lib/utils"
import { getOrder } from "@/lib/api/orders"
import type { Order } from "@/lib/types"

function SuccessView() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const orderId = searchParams.get("order") ?? ""
  const [order, setOrder] = useState<Order | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      if (!orderId) {
        setLoading(false)
        return
      }
      const ord = await getOrder(orderId)
      setOrder(ord)
      setLoading(false)
    }
    load()
  }, [orderId])

  if (loading) {
    return (
      <div className="mx-auto max-w-lg px-4 py-16">
        <Skeleton className="h-10 w-40 mx-auto" />
        <Skeleton className="mt-6 h-64 rounded-2xl" />
      </div>
    )
  }

  if (!order) {
    return (
      <div className="mx-auto max-w-lg px-4 py-16 text-center">
        <Package className="mx-auto h-12 w-12 text-muted-foreground/40" />
        <h1 className="mt-4 text-xl font-semibold">Order tidak ditemukan</h1>
        <Button asChild className="mt-6"><Link href="/">Kembali ke beranda</Link></Button>
      </div>
    )
  }

  const address = order.shippingAddress

  return (
    <div className="mx-auto max-w-lg px-4 py-12">
      <div className="text-center">
        <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl bg-success/10">
          <CheckCircle2 className="h-8 w-8 text-success" />
        </div>
        <h1 className="mt-6 text-2xl font-bold tracking-tight">Pembayaran Berhasil</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Terima kasih! Pesanan Anda sedang diproses.
        </p>
      </div>

      <Card className="mt-8">
        <CardContent className="p-6 space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs text-muted-foreground">Nomor Order</p>
              <p className="text-base font-semibold font-mono">{order.orderNumber}</p>
            </div>
            <div className="text-right">
              <p className="text-xs text-muted-foreground">Total</p>
              <p className="text-base font-semibold tabular-nums">{formatCurrency(order.total)}</p>
            </div>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-muted-foreground">Metode pembayaran</span>
            <span className="font-medium">{order.paymentMethod || "—"}</span>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-muted-foreground">Tanggal</span>
            <span className="font-medium">{formatDate(order.createdAt, "long")}</span>
          </div>

          <Separator />

          <div>
            <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
              Item ({order.items.length})
            </p>
            <div className="space-y-2">
              {order.items.map((item) => (
                <div key={item.id} className="flex justify-between text-sm">
                  <span className="text-muted-foreground truncate">
                    {item.productName} × {item.quantity}
                  </span>
                  <span className="tabular-nums shrink-0">{formatCurrency(item.subtotal)}</span>
                </div>
              ))}
            </div>
          </div>

          {address.line1 && (
            <>
              <Separator />
              <div>
                <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
                  Pengiriman ke
                </p>
                <p className="text-sm text-muted-foreground">
                  {order.customer.name}
                  <br />
                  {address.line1}
                  {address.line2 ? <><br />{address.line2}</> : null}
                  <br />
                  {[address.city, address.state, address.zip].filter(Boolean).join(", ")} — {address.country}
                </p>
              </div>
            </>
          )}
        </CardContent>
      </Card>

      <div className="mt-6 flex gap-2">
        <Button asChild className="flex-1 h-10">
          <Link href="/products">Lanjut Belanja</Link>
        </Button>
        <Button asChild variant="outline" className="flex-1 h-10">
          <Link href="/">Beranda</Link>
        </Button>
      </div>
    </div>
  )
}

export default function SuccessPage() {
  return (
    <Suspense
      fallback={
        <div className="mx-auto max-w-lg px-4 py-16">
          <Skeleton className="h-10 w-40 mx-auto" />
          <Skeleton className="mt-6 h-64 rounded-2xl" />
        </div>
      }
    >
      <SuccessView />
    </Suspense>
  )
}
