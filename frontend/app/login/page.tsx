"use client"

import { useState } from "react"
import Image from "next/image"
import { useRouter } from "next/navigation"
import { Eye, EyeOff, ArrowRight } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Checkbox } from "@/components/ui/checkbox"
import { useAuth } from "@/lib/api/auth-context"
import { toast } from "sonner"

export default function LoginPage() {
  const router = useRouter()
  const { login } = useAuth()
  const [email, setEmail] = useState("admin@laku.id")
  const [password, setPassword] = useState("admin123")
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await login(email, password)
      toast.success("Signed in successfully")
      router.push("/admin")
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Login failed")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen bg-background">
      <div className="flex flex-1 items-center justify-center px-6 lg:flex-[1.2]">
        <div className="w-full max-w-sm">
          <div className="mb-10 flex items-center gap-2.5">
            <Image src="/logo.png" alt="LAKU" width={120} height={80} className="h-9 w-auto" priority />
          </div>

          <div className="mb-8">
            <h1 className="text-[22px] font-semibold tracking-tight text-black">Sign in</h1>
            <p className="mt-1 text-sm text-muted-foreground">
              Enter your credentials to access your account
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-5">
            <div className="space-y-1.5">
              <Label htmlFor="email" className="text-xs font-medium">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="name@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                autoComplete="email"
                className="h-10"
              />
            </div>

            <div className="space-y-1.5">
              <div className="flex items-center justify-between">
                <Label htmlFor="password" className="text-xs font-medium">Password</Label>
                <button
                  type="button"
                  onClick={() => toast.info("Password reset coming soon!")}
                  className="text-xs text-muted-foreground hover:text-primary transition-colors"
                >
                  Forgot?
                </button>
              </div>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  placeholder="Enter your password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  autoComplete="current-password"
                  className="h-10 pr-10"
                />
                <button
                  type="button"
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
            </div>

            <div className="flex items-center gap-2">
              <Checkbox id="remember" defaultChecked />
              <Label htmlFor="remember" className="text-sm font-normal text-muted-foreground">
                Remember me
              </Label>
            </div>

            <Button type="submit" className="w-full h-10 text-sm" disabled={loading}>
              {loading ? (
                <svg className="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                </svg>
              ) : (
                <span className="flex items-center gap-2">
                  Sign in <ArrowRight className="h-4 w-4" />
                </span>
              )}
            </Button>
          </form>

          <p className="mt-8 text-center text-xs text-muted-foreground">
            Don&apos;t have an account?{" "}
            <button
              type="button"
              onClick={() => toast.info("Free trial signup is coming soon!")}
              className="font-medium text-primary hover:underline"
            >
              Start free trial
            </button>
          </p>
        </div>
      </div>

      <div className="hidden flex-1 items-center justify-center bg-muted/40 lg:flex">
        <div className="relative max-w-md px-8 text-center">
          <div className="relative">
            <div className="mb-10 flex justify-center">
              <Image src="/logo.png" alt="LAKU" width={160} height={106} className="h-20 w-auto" />
            </div>
            <h2 className="text-[28px] font-bold leading-tight tracking-tight text-foreground">
              All-in-one commerce<br />platform
            </h2>
            <p className="mt-4 text-[15px] leading-relaxed text-muted-foreground max-w-sm mx-auto">
              Manage products, orders, customers, and analytics in one place.
            </p>
            <div className="mt-12 grid grid-cols-3 gap-4 text-center">
              {[
                { value: "50K+", label: "Products" },
                { value: "10K+", label: "Orders" },
                { value: "99.9%", label: "Uptime" },
              ].map((stat) => (
                <div key={stat.label} className="rounded-xl bg-background border p-4 shadow-sm">
                  <div className="text-xl font-bold text-foreground">{stat.value}</div>
                  <div className="mt-0.5 text-xs text-muted-foreground">{stat.label}</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
