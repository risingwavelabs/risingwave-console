"use client"

import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { useEffect, useState } from "react"
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
import { DefaultService } from "@/api-gen"
import toast from "react-hot-toast"

interface ClusterData {
  id: string
  name: string
  status: "running" | "stopped" | "error"
  host: string
  sqlPort: number
  metaPort: number
  httpPort: number
  nodes: number
  snapshots: Array<{
    id: number
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

interface ClusterPageProps {
  params: {
    id: number
  }
}

export default function ClusterPage({ params }: ClusterPageProps) {
  const router = useRouter()
  const clusterId = params.id
  const [clusterData, setClusterData] = useState<ClusterData | null>(null)
  const [loading, setLoading] = useState(true)
  const [interval, setInterval] = useState("")
  const [expiration, setExpiration] = useState("")
  const [noExpiration, setNoExpiration] = useState(false)
  const [currentPage, setCurrentPage] = useState(1)
  const [snapshotPage, setSnapshotPage] = useState(1)
  const [deleteSnapshotId, setDeleteSnapshotId] = useState<number | null>(null)
  const [dateRange, setDateRange] = useState<DateRange | undefined>()
  const [autoBackupEnabled, setAutoBackupEnabled] = useState(false)
  const [autoBackupInterval, setAutoBackupInterval] = useState("24h")
  const [autoBackupKeepCount, setAutoBackupKeepCount] = useState(7)
  const [isCreatingSnapshot, setIsCreatingSnapshot] = useState(false)
  const [risectlCommand, setRisectlCommand] = useState("")
  const [risectlResult, setRisectlResult] = useState<{ exitCode: number; result: string; error: string } | null>(null)
  const [isRunningCommand, setIsRunningCommand] = useState(false)

  useEffect(() => {
    const fetchClusterData = async () => {
      try {
        const data = await DefaultService.getCluster(clusterId)
        const snapshots = await DefaultService.listClusterSnapshots(clusterId)
        // Transform API data to match our UI needs
        const transformedData: ClusterData = {
          id: data.ID.toString(),
          name: data.name,
          status: "running", // You might want to derive this from API data
          host: data.host,
          sqlPort: data.sqlPort,
          metaPort: data.metaPort,
          httpPort: data.httpPort,
          nodes: 1, // Set default or get from API if available
          snapshots: snapshots.map(s => ({
            id: s.ID,
            name: s.name,
            created_at: s.createdAt
          })),
          autoBackup: {
            enabled: false,
            interval: "24h",
            keepCount: 7
          },
          diagnostics: {
            interval: "1h",
            expiration: "7d",
            noExpiration: false,
            history: []
          }
        }
        setClusterData(transformedData)
        // Initialize state with the fetched data
        setInterval(transformedData.diagnostics.interval)
        setExpiration(transformedData.diagnostics.expiration)
        setNoExpiration(transformedData.diagnostics.noExpiration)
        setAutoBackupEnabled(transformedData.autoBackup.enabled)
        setAutoBackupInterval(transformedData.autoBackup.interval)
        setAutoBackupKeepCount(transformedData.autoBackup.keepCount)
      } catch (error) {
        console.error("Error fetching cluster:", error)
        toast.error("Failed to load cluster details")
        router.push('/clusters')
      } finally {
        setLoading(false)
      }
    }

    void fetchClusterData()
  }, [clusterId, router])

  if (loading) {
    return (
      <div className="p-8">
        <div className="flex items-center justify-center h-[200px]">
          <div className="text-muted-foreground">Loading cluster details...</div>
        </div>
      </div>
    )
  }

  if (!clusterData) {
    return null
  }

  const ITEMS_PER_PAGE = 5
  const paginatedSnapshots = clusterData?.snapshots.slice(
    (snapshotPage - 1) * ITEMS_PER_PAGE,
    snapshotPage * ITEMS_PER_PAGE
  ) || []
  const totalSnapshotPages = Math.ceil((clusterData?.snapshots.length || 0) / ITEMS_PER_PAGE)

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

  const totalPages = Math.ceil(filteredItems.length / ITEMS_PER_PAGE)
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE
  const currentItems = filteredItems.slice(startIndex, startIndex + ITEMS_PER_PAGE)

  const handleDeleteSnapshot = async (id: number) => {
    try {
      await DefaultService.deleteClusterSnapshot(clusterId, id)
      setClusterData(prev => prev ? {
        ...prev,
        snapshots: prev.snapshots.filter(s => s.id !== id)
      } : null)
      toast.success("Snapshot deleted successfully")
    } catch (error) {
      console.error("Error deleting snapshot:", error)
      toast.error("Failed to delete snapshot")
    }
    setDeleteSnapshotId(null)
  }

  const parseCommandArgs = (command: string): string[] => {
    const args: string[] = []
    let currentArg = ''
    let insideQuotes = false
    
    // Helper to add the current argument to args array
    const pushArg = () => {
      const trimmed = currentArg.trim()
      if (trimmed) args.push(trimmed)
      currentArg = ''
    }

    for (let i = 0; i < command.length; i++) {
      const char = command[i]
      
      if (char === '"') {
        if (insideQuotes) {
          // End of quoted section
          pushArg()
          insideQuotes = false
        } else {
          // Start of quoted section
          if (currentArg.trim()) pushArg() // Push any existing arg
          insideQuotes = true
        }
        continue
      }
      
      if (!insideQuotes && (char === ' ' || char === '\n')) {
        pushArg()
        continue
      }
      
      currentArg += char
    }
    
    // Add any remaining argument
    if (currentArg) pushArg()
    
    return args
  }

  const runRisectl = async () => {
    if (!risectlCommand.trim()) {
      toast.error("Please enter a command")
      return
    }

    setIsRunningCommand(true)
    try {
      const args = parseCommandArgs(risectlCommand)
      const result = await DefaultService.runRisectlCommand(clusterId, {
        args
      })
      setRisectlResult({
        exitCode: result.exitCode,
        result: result.result,
        error: result.err
      })
      if (result.exitCode === 0) {
        toast.success("Command executed successfully")
      } else {
        toast.error("Command failed")
      }
    } catch (error) {
      console.error("Error running risectl command:", error)
      toast.error("Failed to run command")
      setRisectlResult({
        exitCode: -1,
        result: "",
        error: "Failed to execute command"
      })
    } finally {
      setIsRunningCommand(false)
    }
  }

  const createSnapshot = async () => {
    setIsCreatingSnapshot(true)
    try {
      const snapshot = await DefaultService.createClusterSnapshot(clusterId, {
        name: `snapshot-${new Date().toISOString().split('T')[0]}`
      })
      // Transform the API snapshot to match our local type
      const transformedSnapshot = {
        id: snapshot.ID,
        name: snapshot.name,
        created_at: snapshot.createdAt
      }
      setClusterData(prev => prev ? {
        ...prev,
        snapshots: [...prev.snapshots, transformedSnapshot]
      } : null)
      toast.success("Snapshot created successfully")
    } catch (error) {
      console.error("Error creating snapshot:", error)
      toast.error("Failed to create snapshot")
    } finally {
      setIsCreatingSnapshot(false)
    }
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
        <div>
          <h3 className="text-lg font-semibold">Cluster Information</h3>
        </div>
        <div className="inline-block">
          <Card className="shadow-none">
            <CardContent className="p-6">
              <div className="inline-flex flex-wrap items-center gap-6">
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
                <div>
                  <p className="text-sm text-muted-foreground">Host</p>
                  <p className="text-sm font-medium">{clusterData.host}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">SQL Port</p>
                  <p className="text-sm font-medium">{clusterData.sqlPort}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Meta Node Port</p>
                  <p className="text-sm font-medium">{clusterData.metaPort}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">HTTP Port</p>
                  <p className="text-sm font-medium">{clusterData.httpPort}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Risectl Command Section */}
      <div className="space-y-4">
        <div className="space-y-1">
          <h3 className="text-lg font-semibold">Risectl Command</h3>
          <p className="text-sm text-muted-foreground">
            Run risectl commands directly on the cluster
          </p>
        </div>

        <div className="max-w-4xl space-y-4 border rounded-lg p-4">
          <div className="flex gap-2">
            <textarea
              value={risectlCommand}
              onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setRisectlCommand(e.target.value)}
              placeholder={'Enter risectl command arguments\nAguments can be separated by spaces and also line breaks\nEnter help to check available commands'}
              className="flex-1 min-h-[80px] px-3 py-2 border rounded-md text-sm font-mono resize-y"
              onKeyDown={(e: React.KeyboardEvent<HTMLTextAreaElement>) => {
                if (e.key === 'Enter' && e.metaKey) {
                  e.preventDefault()
                  void runRisectl()
                }
              }}
            />
            <div className="flex flex-col justify-start">
              <Button 
                onClick={() => void runRisectl()} 
                disabled={isRunningCommand}
                className="whitespace-nowrap select-none"
                size="sm"
              >
                {isRunningCommand ? (
                  <>
                    <svg className="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Running...
                  </>
                ) : (
                  "Run Command"
                )}
              </Button>
              {/* <p className="text-xs text-muted-foreground mt-2 text-center">
                {navigator.platform.toLowerCase().includes('mac') ? '⌘' : 'Ctrl'}+Enter
              </p> */}
            </div>
          </div>

          {risectlResult && (
            <Collapsible defaultOpen>
              <CollapsibleTrigger className="flex items-center gap-2 w-full hover:bg-accent/50 p-2 rounded-md transition-colors">
                <ChevronDown className="h-4 w-4" />
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium">Result</span>
                  <span className={`text-sm ${risectlResult.exitCode === 0 ? 'text-green-600' : 'text-red-600'}`}>
                    (Exit Code: {risectlResult.exitCode})
                  </span>
                </div>
              </CollapsibleTrigger>
              <CollapsibleContent className="pt-2">
                <div className="space-y-2">
                  {risectlResult.result && (
                    <div className="space-y-1">
                      <span className="text-sm font-medium">Output:</span>
                      <div className="relative max-h-[400px] overflow-auto rounded-md">
                        <pre className="p-3 bg-muted text-sm whitespace-pre overflow-auto">
                          {risectlResult.result}
                        </pre>
                      </div>
                    </div>
                  )}
                  
                  {risectlResult.error && (
                    <div className="space-y-1">
                      <span className="text-sm font-medium text-red-600">Error:</span>
                      <div className="relative max-h-[400px] overflow-auto rounded-md">
                        <pre className="p-3 bg-red-50 text-red-600 text-sm whitespace-pre overflow-auto">
                          {risectlResult.error}
                        </pre>
                      </div>
                    </div>
                  )}
                </div>
              </CollapsibleContent>
            </Collapsible>
          )}
        </div>
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
          <Button className="select-none" size="sm" onClick={createSnapshot} disabled={isCreatingSnapshot}>
            {isCreatingSnapshot ? (
              <>
                <svg className="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Creating...
              </>
            ) : (
              "Create Snapshot"
            )}
          </Button>
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

        <div className="space-y-4 max-w-4xl">
          {clusterData.snapshots.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground border-2 border-dashed rounded-lg">
              No snapshots available. Create a snapshot to backup your cluster metadata.
            </div>
          ) : (
            <>
              {paginatedSnapshots.map(snapshot => (
                <div key={snapshot.id} className="flex items-center justify-between p-4 border rounded-lg bg-card">
                  <div>
                    <p className="text-sm font-medium">{snapshot.name}</p>
                    <p className="text-sm text-muted-foreground">{snapshot.created_at}</p>
                  </div>
                  <div className="flex gap-2">
                    {/* <Button variant="outline" size="sm">Restore</Button> */}
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

              <Pagination>
                <PaginationContent>
                  <PaginationItem>
                    <PaginationPrevious
                      href="#"
                      onClick={(e) => {
                        e.preventDefault();
                        if (snapshotPage > 1) setSnapshotPage(p => p - 1);
                      }}
                    />
                  </PaginationItem>
                  {[...Array(totalSnapshotPages)].map((_, i) => (
                    <PaginationItem key={i + 1}>
                      <PaginationLink
                        href="#"
                        onClick={(e) => {
                          e.preventDefault();
                          setSnapshotPage(i + 1);
                        }}
                        isActive={snapshotPage === i + 1}
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
                        if (snapshotPage < totalSnapshotPages) setSnapshotPage(p => p + 1);
                      }}
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </>
          )}
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

          <div className="space-y-4 max-w-4xl">
            <div className="flex items-center justify-between">
              <h4 className="text-sm font-medium">Collection History</h4>
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    size="sm"
                    className="justify-start text-left font-normal"
                    disabled={filteredItems.length === 0}
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

            {filteredItems.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground border-2 border-dashed rounded-lg">
                No diagnostic data available. Data will appear here once collection begins.
              </div>
            ) : (
              <>
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
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  )
} 
