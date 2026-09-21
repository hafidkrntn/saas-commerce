import { http } from "./http"

export interface ChargeResult {
  payment: {
    id: string
    order_id: string
    provider_transaction_id: string
    status: string
    amount: number
    currency: string
  }
  redirect_url: string
  status: string
}

export async function chargeOrder(orderId: string): Promise<ChargeResult> {
  return http.post<ChargeResult>("/payments/charge", { order_id: orderId })
}

export async function getPaymentByOrder(orderId: string): Promise<ChargeResult["payment"] | null> {
  try {
    return await http.get<ChargeResult["payment"]>(`/payments/order/${orderId}`)
  } catch {
    return null
  }
}

// simulatePayment marks a mock charge as paid/failed via the dev webhook.
export async function simulatePayment(
  transactionId: string,
  status: "paid" | "failed",
): Promise<ChargeResult["payment"]> {
  return http.post<ChargeResult["payment"]>("/payments/webhook", {
    transaction_id: transactionId,
    status,
  })
}
