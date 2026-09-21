"use client"

import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"

const statusMap: Record<string, "default" | "secondary" | "destructive" | "outline" | "success" | "warning" | "info" | "neutral"> = {
  active: "success",
  inactive: "neutral",
  draft: "neutral",
  archived: "neutral",
  pending: "warning",
  confirmed: "info",
  processing: "info",
  shipped: "info",
  delivered: "success",
  cancelled: "destructive",
  refunded: "destructive",
  paid: "success",
  failed: "destructive",
  in_stock: "success",
  low_stock: "warning",
  out_of_stock: "destructive",
  overstocked: "neutral",
  scheduled: "neutral",
  in_transit: "info",
  delayed: "warning",
  invited: "neutral",
  blocked: "destructive",
}

interface StatusBadgeProps {
  status: string
  className?: string
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const variant = statusMap[status.toLowerCase()] || "neutral"

  return (
    <Badge variant={variant} className={cn("font-normal capitalize rounded-md text-[11px] px-2 py-0.5", className)}>
      {status.replace(/_/g, " ")}
    </Badge>
  )
}
