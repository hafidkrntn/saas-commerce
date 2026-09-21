import Link from "next/link"
import { Package } from "lucide-react"
import { Button } from "@/components/ui/button"

export default function NotFound() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background px-6">
      <div className="flex h-20 w-20 items-center justify-center rounded-2xl bg-muted">
        <Package className="h-10 w-10 text-muted-foreground" />
      </div>
      <h1 className="mt-6 text-4xl font-bold tracking-tight">404</h1>
      <p className="mt-2 text-lg text-muted-foreground">Page not found</p>
      <p className="mt-1 text-sm text-muted-foreground max-w-sm text-center">
        The page you&apos;re looking for doesn&apos;t exist or has been moved.
      </p>
      <Button asChild className="mt-8">
        <Link href="/">Go to Dashboard</Link>
      </Button>
    </div>
  )
}
