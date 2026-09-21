"use client"

import Link from "next/link"
import Image from "next/image"
import { usePathname } from "next/navigation"
import { cn } from "@/lib/utils"
import {
  LayoutDashboard,
  ShoppingBag,
  Package,
  Users,
  BarChart3,
  ClipboardList,
  Settings,
  ChevronLeft,
  ChevronRight,
} from "lucide-react"
import { useState } from "react"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"

const navItems = [
  { label: "Dashboard", href: "/admin", icon: LayoutDashboard },
  { label: "Orders", href: "/admin/orders", icon: ShoppingBag },
  { label: "Products", href: "/admin/products", icon: Package },
  { label: "Customers", href: "/admin/customers", icon: Users },
  { label: "Analytics", href: "/admin/analytics", icon: BarChart3 },
  { label: "Inventory", href: "/admin/inventory", icon: ClipboardList },
  { label: "Settings", href: "/admin/settings", icon: Settings },
]

export function Sidebar() {
  const pathname = usePathname()
  const [collapsed, setCollapsed] = useState(false)

  const isActive = (href: string) => {
    if (href === "/admin") return pathname === "/admin"
    return pathname.startsWith(href)
  }

  return (
    <TooltipProvider delayDuration={0}>
      <aside
        className={cn(
          "fixed left-0 top-0 z-30 flex h-full flex-col bg-sidebar transition-all duration-200",
          collapsed ? "w-[60px]" : "w-[240px]"
        )}
      >
        <div className={cn(
          "flex h-[57px] items-center border-b border-sidebar-border",
          collapsed ? "justify-center" : "px-4"
        )}>
          <Link href="/admin" className="flex items-center gap-2.5">
            {collapsed ? (
              <Image
                src="/logo.png"
                alt="LAKU"
                width={60}
                height={40}
                className="h-7 w-auto"
                priority
              />
            ) : (
              <Image
                src="/logo.png"
                alt="LAKU"
                width={120}
                height={80}
                className="h-8 w-auto"
                priority
              />
            )}
          </Link>
        </div>

        <nav className="flex-1 space-y-0.5 px-2 py-4">
          {navItems.map((item) => {
            const active = isActive(item.href)
            const Icon = item.icon

            if (collapsed) {
              return (
                <Tooltip key={item.href}>
                  <TooltipTrigger asChild>
                    <Link
                      href={item.href}
                      className={cn(
                        "relative flex h-[34px] w-[34px] items-center justify-center rounded-lg mx-auto transition-colors",
                        active
                          ? "bg-primary/10 text-primary"
                          : "text-sidebar-foreground/40 hover:text-sidebar-foreground/70 hover:bg-sidebar-accent"
                      )}
                    >
                      <Icon className="h-4 w-4" />
                    </Link>
                  </TooltipTrigger>
                  <TooltipContent side="right" className="text-xs py-1 px-2">
                    {item.label}
                  </TooltipContent>
                </Tooltip>
              )
            }

            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  "relative flex items-center gap-2.5 rounded-lg px-2.5 h-[34px] text-sm font-medium transition-all",
                  active
                    ? "bg-primary/10 text-primary"
                    : "text-sidebar-foreground/50 hover:text-sidebar-foreground/80 hover:bg-sidebar-accent"
                )}
              >
                {active && (
                  <span className="absolute left-0 top-1/2 -translate-y-1/2 h-4 w-[3px] rounded-r-full bg-primary" />
                )}
                <Icon className="h-4 w-4 shrink-0" />
                <span className="truncate">{item.label}</span>
              </Link>
            )
          })}
        </nav>

        <div className={cn(
          "border-t border-sidebar-border py-1.5",
          collapsed && "flex justify-center"
        )}>
          <button
            onClick={() => setCollapsed(!collapsed)}
            className={cn(
              "flex items-center gap-2 rounded-lg text-xs text-sidebar-foreground/30 hover:text-sidebar-foreground/60 hover:bg-sidebar-accent transition-colors",
              collapsed ? "justify-center h-[30px] w-[34px] mx-auto" : "px-3 h-[30px] w-full"
            )}
          >
            <ChevronLeft className={cn("h-3.5 w-3.5 transition-transform duration-200", collapsed && "rotate-180")} />
            {!collapsed && <span>Collapse</span>}
          </button>
        </div>
      </aside>
    </TooltipProvider>
  )
}
