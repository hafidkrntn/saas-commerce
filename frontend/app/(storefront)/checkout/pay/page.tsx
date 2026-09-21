"use client"

import { Suspense, useEffect, useMemo, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import Link from "next/link"
import {
  ArrowLeft,
  Banknote,
  CheckCircle2,
  Clock,
  CreditCard,
  ExternalLink,
  Lock,
  Package,
  Smartphone,
  XCircle,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { toast } from "sonner"
import { formatCurrency } from "@/lib/utils"
import { getOrder } from "@/lib/api/orders"
import { getPaymentMethods, type PaymentMethodOption } from "@/lib/api/settings"
import { chargeOrder, simulatePayment, getPaymentByOrder } from "@/lib/api/payments"
import type { Order } from "@/lib/types"

const PAY_TIMEOUT_SECONDS = 15 * 60 // 15 menit

function PayGateway() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const orderId = searchParams.get("order") ?? ""

  const [order, setOrder] = useState<Order | null>(null)
  const [methods, setMethods] = useState<PaymentMethodOption[]>([])
  const [transactionId, setTransactionId] = useState("")
  const [loading, setLoading] = useState(true)
  const [paying, setPaying] = useState(false)
  const [failed, setFailed] = useState(false)
  const [secondsLeft, setSecondsLeft] = useState(PAY_TIMEOUT_SECONDS)

  // Load order + payment state
  useEffect(() => {
    async function load() {
      if (!orderId) {
        setLoading(false)
        return
      }

      const [ord, payment, paymentMethods] = await Promise.all([
        getOrder(orderId),
        getPaymentByOrder(orderId),
        getPaymentMethods().catch(() => []),
      ])

      setMethods(paymentMethods)

      if (!ord) {
        setLoading(false)
        return
      }

      // Sudah dibayar → langsung ke halaman sukses
      if (ord.paymentStatus === "paid") {
        router.replace(`/checkout/success?order=${ord.id}`)
        return
      }

      // Payment record sudah ada → reuse txn id (hindari duplikasi saat refresh)
      if (payment && payment.status === "pending" && payment.provider_transaction_id) {
        setTransactionId(payment.provider_transaction_id)
      } else if (!payment || payment.status !== "paid") {
        try {
          const charge = await chargeOrder(orderId)
          setTransactionId(charge.payment.provider_transaction_id)
        } catch {
          toast.error("Gagal memulai pembayaran, coba lagi")
        }
      }

      setOrder(ord)
      setLoading(false)
    }
    load()
  }, [orderId, router])

  // Countdown
  useEffect(() => {
    if (loading || secondsLeft <= 0) return
    const timer = setInterval(() => setSecondsLeft((s) => Math.max(0, s - 1)), 1000)
    return () => clearInterval(timer)
  }, [loading, secondsLeft])

  const expired = secondsLeft <= 0

  const minutes = Math.floor(secondsLeft / 60)
  const seconds = secondsLeft % 60

  const method = useMemo(
    () => methods.find((m) => m.name === order?.paymentMethod) ?? null,
    [methods, order],
  )

  // Nomor VA mock: deterministik dari order id
  const vaNumber = useMemo(() => {
    if (!order) return ""
    const digits = order.id.replace(/[^0-9a-f]/g, "").slice(-8).toUpperCase()
    return `888 0001 ${digits}`
  }, [order])

  const handlePay = async () => {
    if (!transactionId) return
    setPaying(true)
    try {
      await simulatePayment(transactionId, "paid")
      toast.success("Pembayaran berhasil!")
      router.push(`/checkout/success?order=${order?.id}`)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Pembayaran gagal")
      setFailed(true)
    } finally {
      setPaying(false)
    }
  }

  const handleSimulateFail = async () => {
    if (!transactionId) return
    setPaying(true)
    try {
      await simulatePayment(transactionId, "failed")
      setFailed(true)
    } catch {
      setFailed(true)
    } finally {
      setPaying(false)
    }
  }

  // =========================================================================
  // Loading
  // =========================================================================
  if (loading) {
    return (
      <div className="mx-auto max-w-lg px-4 py-10">
        <Skeleton className="h-8 w-40 mx-auto" />
        <Skeleton className="mt-6 h-64 rounded-2xl" />
      </div>
    )
  }

  // =========================================================================
  // Order tidak ditemukan
  // =========================================================================
  if (!order) {
    return (
      <div className="mx-auto max-w-lg px-4 py-16 text-center">
        <Package className="mx-auto h-12 w-12 text-muted-foreground/40" />
        <h1 className="mt-4 text-xl font-semibold">Order tidak ditemukan</h1>
        <p className="mt-1 text-sm text-muted-foreground">Periksa kembali link pembayaran Anda.</p>
        <Button asChild className="mt-6"><Link href="/">Kembali ke beranda</Link></Button>
      </div>
    )
  }

  // =========================================================================
  // Gagal
  // =========================================================================
  if (failed) {
    return (
      <div className="mx-auto max-w-lg px-4 py-16 text-center">
        <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl bg-destructive/10">
          <XCircle className="h-8 w-8 text-destructive" />
        </div>
        <h1 className="mt-6 text-2xl font-bold tracking-tight">Pembayaran gagal</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Transaksi {order.orderNumber} gagal diproses. Silakan coba lagi.
        </p>
        <Button className="mt-6" onClick={() => router.push(`/checkout/pay?order=${order.id}`)}>
          Coba bayar lagi
        </Button>
      </div>
    )
  }

  const MethodIcon = method?.type === "e_wallet" ? Smartphone : method?.type === "card" ? CreditCard : Banknote

  // =========================================================================
  // Halaman gateway
  // =========================================================================
  return (
    <div className="mx-auto max-w-lg px-4 py-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img src="/logo.png" alt="LAKU" className="h-7 w-auto" />
          <span className="text-sm font-medium text-muted-foreground">Pembayaran aman</span>
        </div>
        <span className="flex items-center gap-1 rounded-full bg-success/10 px-2.5 py-1 text-xs font-medium text-success">
          <Lock className="h-3 w-3" /> Terenkripsi
        </span>
      </div>

      {/* Status + amount */}
      <Card className="mt-6 overflow-hidden">
        <div className="bg-primary px-6 py-4 text-primary-foreground">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs opacity-80">Order</p>
              <p className="text-sm font-semibold font-mono">{order.orderNumber}</p>
            </div>
            <span className={`flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium ${
              expired ? "bg-destructive/20 text-destructive" : "bg-white/15"
            }`}>
              <Clock className="h-3 w-3" />
              {expired ? "Kedaluwarsa" : `${String(minutes).padStart(2, "0")}:${String(seconds).padStart(2, "0")}`}
            </span>
          </div>
        </div>

        <CardContent className="p-6">
          <p className="text-xs text-muted-foreground">Total pembayaran</p>
          <p className="mt-1 text-4xl font-bold tracking-tight tabular-nums">
            {formatCurrency(order.total)}
          </p>

          <div className="mt-6 rounded-xl border bg-muted/30 p-4">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
                <MethodIcon className="h-5 w-5 text-primary" />
              </div>
              <div>
                <p className="text-sm font-medium">{method?.name ?? order.paymentMethod ?? "Pembayaran"}</p>
                <p className="text-xs text-muted-foreground">{method?.description ?? "Transfer bank"}</p>
              </div>
            </div>

            <div className="mt-4 rounded-lg border border-dashed bg-background p-3 text-center">
              <p className="text-xs text-muted-foreground">Nomor Virtual Account</p>
              <p className="mt-1 font-mono text-lg font-semibold tracking-widest">{vaNumber}</p>
            </div>

            <ul className="mt-4 space-y-1.5 text-xs text-muted-foreground">
              <li>1. Transfer ke nomor VA di atas dari m-banking / ATM / e-wallet mana pun.</li>
              <li>2. Pembayaran akan terverifikasi otomatis dalam beberapa menit.</li>
            </ul>
          </div>

          {expired && (
            <p className="mt-4 text-center text-xs text-destructive">
              Sesi pembayaran telah berakhir. Silakan buat order baru.
            </p>
          )}

          <div className="mt-6 space-y-2">
            <Button className="w-full h-11 gap-2 text-base" onClick={handlePay} disabled={paying || expired || !transactionId}>
              <ExternalLink className="h-4 w-4" />
              {paying ? "Memproses..." : "Bayar Sekarang"}
            </Button>
            <div className="flex gap-2">
              <Button variant="outline" className="flex-1 h-10 text-sm" onClick={handleSimulateFail} disabled={paying || expired || !transactionId}>
                Simulasi gagal
              </Button>
              <Button variant="ghost" className="flex-1 h-10 text-sm" onClick={() => router.push("/")}>
                <ArrowLeft className="h-4 w-4" /> Batalkan
              </Button>
            </div>
          </div>

          <p className="mt-4 text-center text-[11px] text-muted-foreground/70">
            Mode demo: klik &quot;Bayar Sekarang&quot; untuk menyimulasikan pembayaran berhasil.
          </p>
        </CardContent>
      </Card>
    </div>
  )
}

export default function PayPage() {
  return (
    <Suspense
      fallback={
        <div className="mx-auto max-w-lg px-4 py-10">
          <Skeleton className="h-8 w-40 mx-auto" />
          <Skeleton className="mt-6 h-64 rounded-2xl" />
        </div>
      }
    >
      <PayGateway />
    </Suspense>
  )
}
