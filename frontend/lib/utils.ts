import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount)
}

export function formatNumber(num: number): string {
  return new Intl.NumberFormat("en-US").format(num)
}

export function formatDate(date: string | Date, format: "short" | "medium" | "long" = "medium"): string {
  const d = typeof date === "string" ? new Date(date) : date
  const options: Record<string, Intl.DateTimeFormatOptions> = {
    short: { month: "numeric", day: "numeric", year: "2-digit" },
    medium: { month: "short", day: "numeric", year: "numeric" },
    long: { month: "long", day: "numeric", year: "numeric" },
  }
  return d.toLocaleDateString("en-US", options[format])
}

export function formatRelativeTime(date: string | Date): string {
  const d = typeof date === "string" ? new Date(date) : date
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (seconds < 60) return "just now"
  if (minutes < 60) return `${minutes}m ago`
  if (hours < 24) return `${hours}h ago`
  if (days < 7) return `${days}d ago`
  return formatDate(d, "short")
}

export function getInitials(name: string): string {
  return name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2)
}

export function slugify(str: string): string {
  return str
    .toLowerCase()
    .replace(/[^\w\s-]/g, "")
    .replace(/\s+/g, "-")
    .replace(/-+/g, "-")
    .trim()
}

// readImageFile reads an image file and downscales it to maxEdge pixels,
// returning a compact JPEG data URL (keeps the stored payload small).
export function readImageFile(file: File, maxEdge = 800): Promise<string> {
  return new Promise((resolve, reject) => {
    if (!file.type.startsWith("image/")) {
      reject(new Error("File harus berupa gambar"))
      return
    }
    if (file.size > 5 * 1024 * 1024) {
      reject(new Error("Ukuran gambar maksimal 5 MB"))
      return
    }

    const reader = new FileReader()
    reader.onerror = () => reject(new Error("Gagal membaca file"))
    reader.onload = () => {
      const img = new Image()
      img.onerror = () => reject(new Error("Gambar tidak valid"))
      img.onload = () => {
        const scale = Math.min(1, maxEdge / Math.max(img.width, img.height))
        const canvas = document.createElement("canvas")
        canvas.width = Math.max(1, Math.round(img.width * scale))
        canvas.height = Math.max(1, Math.round(img.height * scale))
        const ctx = canvas.getContext("2d")
        if (!ctx) {
          reject(new Error("Canvas tidak tersedia"))
          return
        }
        ctx.drawImage(img, 0, 0, canvas.width, canvas.height)
        resolve(canvas.toDataURL("image/jpeg", 0.82))
      }
      img.src = reader.result as string
    }
    reader.readAsDataURL(file)
  })
}
