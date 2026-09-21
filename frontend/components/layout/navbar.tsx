"use client"

import { useState, useRef, useEffect } from "react"
import { useRouter } from "next/navigation"
import {
  Search,
  Bell,
  Command,
  LogOut,
  User,
  Settings,
  HelpCircle,
  LayoutDashboard,
  ShoppingBag,
  Package,
  Users,
  BarChart3,
  ClipboardList,
  ArrowUpRight,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { cn } from "@/lib/utils"
import { useAuth } from "@/lib/api/auth-context"
import { toast } from "sonner"

const quickLinks = [
  { label: "Dashboard", href: "/admin", icon: LayoutDashboard },
  { label: "Orders", href: "/admin/orders", icon: ShoppingBag },
  { label: "Products", href: "/admin/products", icon: Package },
  { label: "Customers", href: "/admin/customers", icon: Users },
  { label: "Analytics", href: "/admin/analytics", icon: BarChart3 },
  { label: "Inventory", href: "/admin/inventory", icon: ClipboardList },
]

export function Navbar() {
  const router = useRouter()
  const { user, logout } = useAuth()
  const [searchOpen, setSearchOpen] = useState(false)
  const [query, setQuery] = useState("")
  const inputRef = useRef<HTMLInputElement>(null)

  const handleLogout = async () => {
    await logout()
    toast.success("Signed out")
    router.push("/login")
  }

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault()
        inputRef.current?.focus()
      }
    }
    document.addEventListener("keydown", handleKeyDown)
    return () => document.removeEventListener("keydown", handleKeyDown)
  }, [])

  const filteredLinks = query
    ? quickLinks.filter((l) => l.label.toLowerCase().includes(query.toLowerCase()))
    : quickLinks

  const handleSelect = (href: string) => {
    setQuery("")
    setSearchOpen(false)
    router.push(href)
  }

  return (
    <header className="sticky top-0 z-20 flex h-[57px] items-center gap-3 border-b bg-background/80 backdrop-blur-lg px-6">
      <div className="relative flex-1 max-w-sm">
        <Search className="absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground/40" />
        <Input
          ref={inputRef}
          placeholder="Search..."
          value={query}
          onChange={(e) => { setQuery(e.target.value); setSearchOpen(true) }}
          onFocus={() => setSearchOpen(true)}
          onBlur={() => setTimeout(() => setSearchOpen(false), 200)}
          className="h-8 pl-9 pr-10 text-xs bg-muted/50 border-none rounded-lg placeholder:text-muted-foreground/70 focus-visible:bg-background focus-visible:ring-1 focus-visible:ring-ring"
        />
        <div className="absolute right-2.5 top-1/2 -translate-y-1/2 hidden md:flex items-center gap-0.5">
          <kbd className="pointer-events-none inline-flex h-4 select-none items-center gap-0.5 rounded border bg-muted px-1.5 font-mono text-[9px] font-medium text-muted-foreground/40">
            <Command className="h-2 w-2" />K
          </kbd>
        </div>

        {searchOpen && query && (
          <div className="absolute top-full mt-1.5 w-full rounded-lg border bg-background shadow-lg py-1 z-50">
            {filteredLinks.length === 0 ? (
              <p className="px-3 py-2 text-xs text-muted-foreground/60">No results for &quot;{query}&quot;</p>
            ) : (
              filteredLinks.map((link) => (
                <button
                  key={link.href}
                  onMouseDown={() => handleSelect(link.href)}
                  className="flex w-full items-center gap-2.5 px-3 py-2 text-sm hover:bg-accent text-left"
                >
                  <link.icon className="h-3.5 w-3.5 text-muted-foreground/60" />
                  <span>{link.label}</span>
                  <ArrowUpRight className="h-3 w-3 ml-auto text-muted-foreground/30" />
                </button>
              ))
            )}
          </div>
        )}
      </div>

      <div className="flex items-center gap-1">
        <button className="relative flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground/50 hover:text-foreground hover:bg-accent transition-colors">
          <Bell className="h-3.5 w-3.5" />
          <span className="absolute right-2 top-2 h-[5px] w-[5px] rounded-full bg-primary ring-2 ring-background" />
        </button>

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button className="flex h-8 w-8 items-center justify-center rounded-lg ml-1 hover:bg-accent transition-colors">
              <Avatar className="h-6 w-6">
                <AvatarImage src="https://api.dicebear.com/7.x/avataaars/svg?seed=admin" />
                <AvatarFallback className="text-[9px]">AD</AvatarFallback>
              </Avatar>
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56 mt-1">
            <DropdownMenuLabel>
              <div className="flex flex-col">
                <span className="text-sm font-medium">{user?.name ?? "Admin User"}</span>
                <span className="text-xs font-normal text-muted-foreground">{user?.email ?? "admin@laku.id"}</span>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => router.push("/admin/settings")}>
              <User className="mr-2 h-3.5 w-3.5" />
              Profile
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => router.push("/admin/settings")}>
              <Settings className="mr-2 h-3.5 w-3.5" />
              Settings
            </DropdownMenuItem>
            <DropdownMenuItem>
              <HelpCircle className="mr-2 h-3.5 w-3.5" />
              Help
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="text-destructive focus:text-destructive" onClick={handleLogout}>
              <LogOut className="mr-2 h-3.5 w-3.5" />
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  )
}
