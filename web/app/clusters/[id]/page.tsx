"use client"

import { useParams } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { useState } from "react"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Switch } from "@/components/ui/switch"
import { Label } from "@/components/ui/label"
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible"
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination"
import { ChevronDown } from "lucide-react"
import { ConfirmationPopup } from "@/components/ui/confirmation-popup"
import { Calendar } from "@/components/ui/calendar"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { CalendarIcon } from "lucide-react"
import { format } from "date-fns"
import { DateRange } from "react-day-picker"
import { Card, CardContent } from "@/components/ui/card"

interface ClusterData {
  id: string
  name: string
  status: "running" | "stopped" | "error"
  host: string
  sqlPort: number
  metaNodePort: number
  user: string
  password?: string
  database: string
  nodes: number
  snapshots: Array<{
    id: string
    name: string
    created_at: string
  }>
  autoBackup: {
    enabled: boolean
    interval: string
    keepCount: number
  }
  diagnostics: {
    interval: string
    expiration: string
    noExpiration: boolean
    history: Array<{
      id: number
      timestamp: string
      data: string
    }>
  }
}

// Mock data for demonstration
const initialClusterData: ClusterData = {
  id: "cluster-1",
  name: "Production DB",
  status: "running",
  host: "localhost",
  sqlPort: 8080,
  metaNodePort: 9090,
  user: "admin",
  password: "supersecret",
  database: "myapp_production",
  nodes: 5,
  snapshots: [
    { id: "snap-1", name: "Daily Backup", created_at: "2024-01-20 10:00" },
    { id: "snap-2", name: "Pre-deployment", created_at: "2024-01-19 15:30" },
  ],
  autoBackup: {
    enabled: false,
    interval: "24h",
    keepCount: 7
  },
  diagnostics: {
    interval: "1h",
    expiration: "7d",
    noExpiration: false,
    history: [
      { id: 1, timestamp: "2024-01-20 15:00:00", data: "CPU Usage: 78%\nMemory Usage: 4.2GB\nDisk I/O: 250MB/s\nNetwork: 180Mbps\n\nActive Connections: 1250\nQuery Response Time: 45ms\nCache Hit Rate: 92%\n\nReplication Lag: 0.8s\nWrite Operations: 12k/s\nRead Operations: 45k/s" },
      { id: 2, timestamp: "2024-01-20 14:00:00", data: "CPU Usage: 65%\nMemory Usage: 3.8GB\nDisk I/O: 180MB/s\nNetwork: 150Mbps\n\nActive Connections: 980\nQuery Response Time: 38ms\nCache Hit Rate: 94%\n\nReplication Lag: 0.5s\nWrite Operations: 8k/s\nRead Operations: 40k/s" },
    ]
  }
}

export default function ClusterPage() {
  const params = useParams()
  const clusterId = params.id as string
  const [clusterData, setClusterData] = useState<ClusterData>(initialClusterData)
  const [interval, setInterval] = useState(clusterData.diagnostics.interval)
  const [expiration, setExpiration] = useState(clusterData.diagnostics.expiration)
  const [noExpiration, setNoExpiration] = useState(clusterData.diagnostics.noExpiration)
  const [currentPage, setCurrentPage] = useState(1)
  const [deleteSnapshotId, setDeleteSnapshotId] = useState<string | null>(null)
  const itemsPerPage = 5
  const totalPages = Math.ceil(clusterData.diagnostics.history.length / itemsPerPage)
  const startIndex = (currentPage - 1) * itemsPerPage
  const [dateRange, setDateRange] = useState<DateRange | undefined>()
  const [autoBackupEnabled, setAutoBackupEnabled] = useState(clusterData.autoBackup.enabled)
  const [autoBackupInterval, setAutoBackupInterval] = useState(clusterData.autoBackup.interval)
  const [autoBackupKeepCount, setAutoBackupKeepCount] = useState(clusterData.autoBackup.keepCount)

  const filteredItems = clusterData.diagnostics.history
    .filter(item => {
      if (!dateRange?.from && !dateRange?.to) return true
      const itemDate = new Date(item.timestamp)
      if (dateRange.from && dateRange.to) {
        return itemDate >= dateRange.from && itemDate <= dateRange.to
      }
      if (dateRange.from) {
        return itemDate >= dateRange.from
      }
      if (dateRange.to) {
        return itemDate <= dateRange.to
      }
      return true
    })

  const currentItems = filteredItems
    .slice(startIndex, startIndex + itemsPerPage)

  const handleDeleteSnapshot = (id: string) => {
    setClusterData(prev => ({
      ...prev,
      snapshots: prev.snapshots.filter(s => s.id !== id)
    }))
    setDeleteSnapshotId(null)
  }

  return (
    <div className="p-8 space-y-8">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h2 className="text-lg font-semibold">Cluster Details</h2>
          <p className="text-sm text-muted-foreground">
            Manage and monitor your cluster
          </p>
        </div>
      </div>

      {/* Cluster Details Section */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold">Cluster Information</h3>
          <Button variant="outline" size="sm">
            Edit
          </Button>
        </div>
        <Card className="shadow-none">
          <CardContent className="p-6">
            <div className="flex items-center gap-6">
              <div>
                <p className="text-sm text-muted-foreground">Name</p>
                <p className="text-sm font-medium">{clusterData.name}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Status</p>
                <div className="flex items-center gap-2">
                  <div className={`w-2 h-2 rounded-full ${clusterData.status === "running" ? "bg-green-500" :
                    clusterData.status === "stopped" ? "bg-yellow-500" :
                      "bg-red-500"
                    }`} />
                  <p className="text-sm font-medium capitalize">{clusterData.status}</p>
                </div>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Number of Nodes</p>
                <p className="text-sm font-medium">{clusterData.nodes} nodes</p>
              </div>
            </div>
          </CardContent>
        </Card>
        
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold">Connection Information</h3>
          <Button variant="outline" size="sm">
            Edit
          </Button>
        </div>
        <Card className="shadow-none">
          <CardContent className="p-6">
            <div className="grid grid-cols-3 gap-6">
              <div>
                <p className="text-sm text-muted-foreground">Host</p>
                <p className="text-sm font-medium">{clusterData.host}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Database</p>
                <p className="text-sm font-medium">{clusterData.database}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">User</p>
                <p className="text-sm font-medium">{clusterData.user}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">SQL Port</p>
                <p className="text-sm font-medium">{clusterData.sqlPort}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Meta Node Port</p>
                <p className="text-sm font-medium">{clusterData.metaNodePort}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Snapshots Section */}
      <div className="space-y-4">
        <div className="flex items-center justify-between max-w-4xl">
          <div className="space-y-1">
            <h3 className="text-lg font-semibold">Metadata Snapshot</h3>
            <p className="text-sm text-muted-foreground">
              Backup and restore cluster metadata. Keep snapshots minimal as excessive snapshots may affect performance.
            </p>
          </div>
          <Button size="sm">Create Snapshot</Button>
        </div>

        <div className="max-w-4xl space-y-4 border rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label className="text-sm font-medium">Auto Backup</Label>
              <p className="text-sm text-muted-foreground">Automatically create snapshots at regular intervals</p>
            </div>
            <Switch
              checked={autoBackupEnabled}
              onCheckedChange={setAutoBackupEnabled}
            />
          </div>

          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label className="text-sm font-medium">Backup Interval</Label>
                <p className="text-sm text-muted-foreground">How often to create snapshots</p>
              </div>
              <Select 
                value={autoBackupInterval} 
                onValueChange={setAutoBackupInterval}
                disabled={!autoBackupEnabled}
              >
                <SelectTrigger className="w-[160px]">
                  <SelectValue placeholder="Select interval" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="30m">Every 30 minutes</SelectItem>
                  <SelectItem value="6h">Every 6 hours</SelectItem>
                  <SelectItem value="12h">Every 12 hours</SelectItem>
                  <SelectItem value="24h">Every 24 hours</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label className="text-sm font-medium">Keep Last</Label>
                <p className="text-sm text-muted-foreground">Number of automatic snapshots to retain</p>
              </div>
              <Select 
                value={autoBackupKeepCount.toString()} 
                onValueChange={(value) => setAutoBackupKeepCount(parseInt(value))}
                disabled={!autoBackupEnabled}
              >
                <SelectTrigger className="w-[160px]">
                  <SelectValue placeholder="Select count" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="3">3 snapshots</SelectItem>
                  <SelectItem value="5">5 snapshots</SelectItem>
                  <SelectItem value="7">7 snapshots</SelectItem>
                  <SelectItem value="14">14 snapshots</SelectItem>
                  <SelectItem value="30">30 snapshots</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </div>

        <div className="space-y-4">
          {clusterData.snapshots.map(snapshot => (
            <div key={snapshot.id} className="flex items-center justify-between p-4 border rounded-lg bg-card">
              <div>
                <p className="text-sm font-medium">{snapshot.name}</p>
                <p className="text-sm text-muted-foreground">{snapshot.created_at}</p>
              </div>
              <div className="flex gap-2">
                <Button variant="outline" size="sm">Restore</Button>
                <div className="relative">
                  <Button
                    variant="outline"
                    size="sm"
                    className="text-red-500 hover:text-red-600"
                    onClick={() => setDeleteSnapshotId(snapshot.id)}
                  >
                    Delete
                  </Button>
                  {deleteSnapshotId === snapshot.id && (
                    <ConfirmationPopup
                      message="Delete this snapshot?"
                      onConfirm={() => handleDeleteSnapshot(snapshot.id)}
                      onCancel={() => setDeleteSnapshotId(null)}
                    />
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Diagnostics Section */}
      <div className="space-y-4">
        <div className="space-y-1">
          <h3 className="text-lg font-semibold">Diagnostic Information</h3>
          <p className="text-sm text-muted-foreground">
            Configure automatic collection of diagnostic data and system metrics
          </p>
        </div>

        <div className="space-y-6">
          <div className="max-w-4xl space-y-4 border rounded-lg p-4">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label className="text-sm font-medium">Collection Interval</Label>
                <p className="text-sm text-muted-foreground">How often to collect data</p>
              </div>
              <Select value={interval} onValueChange={setInterval}>
                <SelectTrigger className="w-[140px]">
                  <SelectValue placeholder="Select interval" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="15m">15 minutes</SelectItem>
                  <SelectItem value="30m">30 minutes</SelectItem>
                  <SelectItem value="1h">1 hour</SelectItem>
                  <SelectItem value="6h">6 hours</SelectItem>
                  <SelectItem value="12h">12 hours</SelectItem>
                  <SelectItem value="24h">24 hours</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label className="text-sm font-medium">Data Retention</Label>
                <p className="text-sm text-muted-foreground">How long to keep data</p>
              </div>
              <Select
                value={expiration}
                onValueChange={setExpiration}
                disabled={noExpiration}
              >
                <SelectTrigger className="w-[140px]">
                  <SelectValue placeholder="Select retention" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="1d">1 day</SelectItem>
                  <SelectItem value="7d">7 days</SelectItem>
                  <SelectItem value="14d">14 days</SelectItem>
                  <SelectItem value="30d">30 days</SelectItem>
                  <SelectItem value="90d">90 days</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="flex items-center justify-end gap-1.5">
              <Label className="text-sm text-muted-foreground">No expiration</Label>
              <Switch
                checked={noExpiration}
                onCheckedChange={setNoExpiration}
              />
            </div>
          </div>

          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h4 className="text-sm font-medium">Collection History</h4>
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    size="sm"
                    className="justify-start text-left font-normal"
                  >
                    <CalendarIcon className="mr-2 h-4 w-4" />
                    {dateRange?.from ? (
                      dateRange.to ? (
                        <>
                          {format(dateRange.from, "LLL dd, y")} -{" "}
                          {format(dateRange.to, "LLL dd, y")}
                        </>
                      ) : (
                        format(dateRange.from, "LLL dd, y")
                      )
                    ) : (
                      "Pick a date range"
                    )}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="end">
                  <Calendar
                    initialFocus
                    mode="range"
                    defaultMonth={dateRange?.from}
                    selected={dateRange}
                    onSelect={(range) => {
                      setDateRange({ from: range?.from, to: range?.to });
                      setCurrentPage(1); // Reset to first page when filter changes
                    }}
                    numberOfMonths={1}
                  />
                </PopoverContent>
              </Popover>
            </div>
            {currentItems.map((item) => (
              <Collapsible key={item.id} className="border rounded-lg">
                <CollapsibleTrigger className="flex items-center justify-between w-full p-4 hover:bg-accent/50 transition-colors">
                  <span className="text-sm font-medium">{item.timestamp}</span>
                  <ChevronDown className="h-4 w-4" />
                </CollapsibleTrigger>
                <CollapsibleContent>
                  <div className="p-4 pt-0 font-mono text-sm bg-muted/50">
                    <pre className="whitespace-pre-wrap">{item.data}</pre>
                  </div>
                </CollapsibleContent>
              </Collapsible>
            ))}

            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                      if (currentPage > 1) setCurrentPage(p => p - 1);
                    }}
                  />
                </PaginationItem>
                {[...Array(totalPages)].map((_, i) => (
                  <PaginationItem key={i + 1}>
                    <PaginationLink
                      href="#"
                      onClick={(e) => {
                        e.preventDefault();
                        setCurrentPage(i + 1);
                      }}
                      isActive={currentPage === i + 1}
                    >
                      {i + 1}
                    </PaginationLink>
                  </PaginationItem>
                ))}
                <PaginationItem>
                  <PaginationNext
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                      if (currentPage < totalPages) setCurrentPage(p => p + 1);
                    }}
                  />
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          </div>
        </div>
      </div>
    </div>
  )
} 
