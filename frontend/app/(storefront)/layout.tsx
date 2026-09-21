"use client"

import { usePathname } from "next/navigation"
import { StoreNavbar } from "@/components/storefront/navbar"
import { StoreFooter } from "@/components/storefront/footer"

export default function StorefrontLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const pathname = usePathname()
  const isCheckout = pathname === "/checkout"

  return (
    <div className="flex min-h-screen flex-col bg-background">
      <StoreNavbar />
      <main className="flex-1">{children}</main>
      {!isCheckout && <StoreFooter />}
    </div>
  )
}
