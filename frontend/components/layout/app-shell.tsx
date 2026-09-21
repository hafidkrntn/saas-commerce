"use client"

import { Sidebar } from "@/components/layout/sidebar"
import { Navbar } from "@/components/layout/navbar"
import { usePathname } from "next/navigation"

const authPaths = ["/login"]

export function AppShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const isAuthPage = authPaths.includes(pathname)

  if (isAuthPage) {
    return <>{children}</>
  }

  return (
    <div className="flex min-h-screen bg-muted/30">
      <Sidebar />
      <div className="flex flex-1 flex-col pl-[240px] transition-all duration-200">
        <Navbar />
        <main className="flex-1">{children}</main>
      </div>
    </div>
  )
}
