"use client"

import { cn } from "@/lib/utils"
import { Card, CardContent } from "@/components/ui/card"
import { ArrowUpRight, ArrowDownRight, LucideIcon } from "lucide-react"

interface StatCardProps {
  title: string
  value: string
  change?: number
  changeLabel?: string
  icon?: LucideIcon
  className?: string
  trend?: "up" | "down"
}

export function StatCard({
  title,
  value,
  change,
  changeLabel,
  icon: Icon,
  className,
  trend = "up",
}: StatCardProps) {
  return (
    <Card className={cn("transition-shadow hover:shadow-md", className)}>
      <CardContent className="p-5">
        <div className="flex items-center justify-between">
          <span className="text-xs font-medium text-muted-foreground tracking-wide uppercase">
            {title}
          </span>
          {Icon && (
            <div className="flex h-[30px] w-[30px] items-center justify-center rounded-lg bg-primary/10">
              <Icon className="h-[14px] w-[14px] text-primary" />
            </div>
          )}
        </div>
        <div className="mt-2 flex items-baseline gap-2">
          <span className="text-[26px] font-semibold tracking-tight">{value}</span>
          {change !== undefined && (
            <span
              className={cn(
                "inline-flex items-center gap-0.5 text-xs font-medium",
                trend === "up" ? "text-success" : "text-destructive"
              )}
            >
              {trend === "up" ? (
                <ArrowUpRight className="h-3.5 w-3.5" />
              ) : (
                <ArrowDownRight className="h-3.5 w-3.5" />
              )}
              {Math.abs(change)}%
            </span>
          )}
        </div>
        {changeLabel && (
          <p className="mt-0.5 text-xs text-muted-foreground/60">{changeLabel}</p>
        )}
      </CardContent>
    </Card>
  )
}
