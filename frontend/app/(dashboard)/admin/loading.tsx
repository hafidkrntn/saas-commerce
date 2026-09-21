export default function DashboardLoading() {
  return (
    <div className="p-6">
      <div className="h-8 w-48 rounded-lg bg-muted animate-pulse" />
      <div className="mt-1 h-4 w-72 rounded-lg bg-muted animate-pulse" />
      <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="rounded-xl border p-6">
            <div className="h-4 w-24 rounded bg-muted animate-pulse" />
            <div className="mt-3 h-8 w-28 rounded bg-muted animate-pulse" />
            <div className="mt-2 h-3 w-20 rounded bg-muted animate-pulse" />
          </div>
        ))}
      </div>
    </div>
  )
}
