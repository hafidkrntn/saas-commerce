"use client"

import { cn } from "@/lib/utils"
import { Skeleton } from "@/components/ui/skeleton"

interface LoadingSkeletonProps {
  className?: string
  count?: number
}

export function TableSkeleton({ className, count = 5 }: LoadingSkeletonProps) {
  return (
    <div className={cn("space-y-3", className)}>
      <div className="flex items-center gap-4">
        <Skeleton className="h-8 w-8 rounded-lg" />
        <Skeleton className="h-8 flex-1" />
        <Skeleton className="h-8 w-24" />
        <Skeleton className="h-8 w-20" />
        <Skeleton className="h-8 w-16" />
      </div>
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className="flex items-center gap-4">
          <Skeleton className="h-8 w-8 rounded-lg" />
          <Skeleton className="h-8 flex-1" />
          <Skeleton className="h-8 w-24" />
          <Skeleton className="h-8 w-20" />
          <Skeleton className="h-8 w-16" />
        </div>
      ))}
    </div>
  )
}

export function StatCardSkeleton({ className }: { className?: string }) {
  return (
    <div className={cn("rounded-xl border p-6", className)}>
      <div className="flex items-center justify-between">
        <Skeleton className="h-4 w-24" />
        <Skeleton className="h-9 w-9 rounded-xl" />
      </div>
      <Skeleton className="mt-3 h-8 w-28" />
      <Skeleton className="mt-2 h-3 w-20" />
    </div>
  )
}

export function ChartSkeleton({ className }: { className?: string }) {
  return (
    <div className={cn("rounded-xl border p-6", className)}>
      <Skeleton className="h-5 w-32" />
      <Skeleton className="mt-6 h-[300px] w-full rounded-xl" />
    </div>
  )
}
