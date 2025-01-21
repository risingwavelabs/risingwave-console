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
import { Card, CardContent } from "@/components/ui/card"
import { DefaultService } from "@/api-gen"
import { DiagnosticData } from "@/api-gen/models/DiagnosticData"
import toast from "react-hot-toast"

const width = "w-full"

const enableAutoDiagnostic = false

const enableAutoBackup = false

interface BaseCluster {
  name: string
  host: string
  sqlPort: number
  metaPort: number
  httpPort: number
  autoBackup?: {
    enabled: boolean
    interval: string
    keepCount: number
  }
  diagnostics?: {
    enabled: boolean
    interval: string
    expiration: string
    noExpiration: boolean
  }
}

// API response type
interface Cluster extends BaseCluster {
  ID: number
}

// UI data type
interface ClusterData extends BaseCluster {
  id: string
  status: "running" | "stopped" | "error"
  nodes: number
  snapshots: Array<{
    id: number
    name: string
    created_at: string
  }>
  autoBackup: Required<NonNullable<BaseCluster['autoBackup']>>
  diagnostics: Required<NonNullable<BaseCluster['diagnostics']>> & {
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
  const [snapshotPage, setSnapshotPage] = useState(1)
  const [deleteSnapshotId, setDeleteSnapshotId] = useState<number | null>(null)
  const [autoBackupEnabled, setAutoBackupEnabled] = useState(false)
  const [autoBackupInterval, setAutoBackupInterval] = useState("*/30 * * * *")
  const [autoBackupKeepCount, setAutoBackupKeepCount] = useState(7)
  const [isCreatingSnapshot, setIsCreatingSnapshot] = useState(false)
  const [isUpdatingBackupConfig, setIsUpdatingBackupConfig] = useState(false)
  const [autoDiagnosticInterval, setAutoDiagnosticInterval] = useState("*/30 * * * *")
  const [autoDiagnosticExpiration, setAutoDiagnosticExpiration] = useState("7d")
  const [autoDiagnosticNoExpiration, setAutoDiagnosticNoExpiration] = useState(false)
  const [autoDiagnosticEnabled, setAutoDiagnosticEnabled] = useState(false)
  const [isUpdatingDiagnosticConfig, setIsUpdatingDiagnosticConfig] = useState(false)
  const [risectlCommand, setRisectlCommand] = useState("")
  const [risectlResult, setRisectlResult] = useState<{ exitCode: number; stdout: string; stderr: string; err: string } | null>(null)
  const [isRunningCommand, setIsRunningCommand] = useState(false)
  const [isDiagnosing, setIsDiagnosing] = useState(false)
  const [isResultOpen, setIsResultOpen] = useState(true)
  const [diagnostics, setDiagnostics] = useState<Array<DiagnosticData>>([])
  const [diagnosticContent, setDiagnosticContent] = useState<Record<number, string>>({})
  const [isDiagnosticContentLoading, setIsDiagnosticContentLoading] = useState<Record<number, boolean>>({})
  const [diagnosticPage, setDiagnosticPage] = useState(1)
  const [totalDiagnostics, setTotalDiagnostics] = useState(0)
  const DIAGNOSTICS_PER_PAGE = 5

  const updateAutoBackupConfig = async (
    enabled?: boolean,
    interval?: string,
    keepCount?: number
  ) => {
    setIsUpdatingBackupConfig(true)
    try {
      const config = {
        enabled: enabled ?? autoBackupEnabled,
        cronExpression: interval ?? autoBackupInterval,
        keepLast: keepCount ?? autoBackupKeepCount
      }
      await DefaultService.updateClusterAutoBackupConfig(clusterId, config)
      toast.success("Auto backup configuration updated")
    } catch (error) {
      console.error("Error updating auto backup config:", error)
      toast.error("Failed to update auto backup configuration")
      // Revert the state changes on error
      setAutoBackupEnabled(clusterData?.autoBackup.enabled ?? false)
      setAutoBackupInterval(clusterData?.autoBackup.interval ?? "*/30 * * * *")
      setAutoBackupKeepCount(clusterData?.autoBackup.keepCount ?? 7)
    } finally {
      setIsUpdatingBackupConfig(false)
    }
  }

  const updateAutoDiagnosticConfig = async (
    enabled?: boolean,
    interval?: string,
    expiration?: string,
    noExpiration?: boolean
  ) => {
    setIsUpdatingDiagnosticConfig(true)
    try {
      const noExp = noExpiration ?? autoDiagnosticNoExpiration
      await DefaultService.updateClusterAutoDiagnosticConfig(clusterId, {
        enabled: enabled ?? autoDiagnosticEnabled,
        cronExpression: interval ?? autoDiagnosticInterval,
        retentionDuration: noExp ? undefined : (expiration ?? autoDiagnosticExpiration)
      })
      toast.success("Auto diagnostic configuration updated")
    } catch (error) {
      console.error("Error updating auto diagnostic config:", error)
      toast.error("Failed to update auto diagnostic configuration")
      // Revert the state changes on error
      setAutoDiagnosticEnabled(clusterData?.diagnostics.enabled ?? false)
      setAutoDiagnosticInterval(clusterData?.diagnostics.interval ?? "*/30 * * * *")
      setAutoDiagnosticExpiration(clusterData?.diagnostics.expiration ?? "7d")
      setAutoDiagnosticNoExpiration(clusterData?.diagnostics.noExpiration ?? false)
    } finally {
      setIsUpdatingDiagnosticConfig(false)
    }
  }

  useEffect(() => {
    const fetchClusterData = async () => {
      try {
        const [data, snapshots, diagnosticsData] = await Promise.all([
          DefaultService.getCluster(clusterId),
          DefaultService.listClusterSnapshots(clusterId),
          DefaultService.listClusterDiagnostics(clusterId)
        ])

        const cluster = data as Cluster // Type assertion to match our interface

        // Transform API data to match our UI needs
        const transformedData: ClusterData = {
          id: cluster.ID.toString(),
          name: cluster.name,
          status: "running", // You might want to derive this from API data
          host: cluster.host,
          sqlPort: cluster.sqlPort,
          metaPort: cluster.metaPort,
          httpPort: cluster.httpPort,
          nodes: 1, // Set default or get from API if available
          snapshots: snapshots.map(s => ({
            id: s.ID,
            name: s.name,
            created_at: s.createdAt
          })),
          autoBackup: {
            enabled: cluster.autoBackup?.enabled ?? false,
            interval: cluster.autoBackup?.interval ?? "*/30 * * * *",
            keepCount: cluster.autoBackup?.keepCount ?? 7
          },
          diagnostics: {
            enabled: cluster.diagnostics?.enabled ?? false,
            interval: cluster.diagnostics?.interval ?? "*/30 * * * *",
            expiration: cluster.diagnostics?.expiration ?? "7d",
            noExpiration: cluster.diagnostics?.noExpiration ?? false,
            history: []
          }
        }
        setClusterData(transformedData)

        // Set diagnostic data
        console.log('Setting diagnostic data:', diagnosticsData)
        setDiagnostics(diagnosticsData)
        setTotalDiagnostics(diagnosticsData.length)

        // Initialize state with the fetched data
        setAutoBackupEnabled(transformedData.autoBackup.enabled)
        setAutoBackupInterval(transformedData.autoBackup.interval)
        setAutoBackupKeepCount(transformedData.autoBackup.keepCount)
        setAutoDiagnosticEnabled(transformedData.diagnostics.enabled)
        setAutoDiagnosticInterval(transformedData.diagnostics.interval)
        setAutoDiagnosticExpiration(transformedData.diagnostics.expiration)
        setAutoDiagnosticNoExpiration(transformedData.diagnostics.noExpiration)
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
    setIsResultOpen(true)
    try {
      const args = parseCommandArgs(risectlCommand)
      const result = await DefaultService.runRisectlCommand(clusterId, {
        args
      })
      setRisectlResult({
        exitCode: result.exitCode,
        stdout: result.stdout,
        stderr: result.stderr,
        err: result.err
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
        stdout: "",
        stderr: "Failed to execute command",
        err: "Failed to execute command"
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

  const fetchDiagnosticContent = async (id: number) => {
    if (diagnosticContent[id] || isDiagnosticContentLoading[id]) return

    setIsDiagnosticContentLoading(prev => ({ ...prev, [id]: true }))
    try {
      const data = await DefaultService.getClusterDiagnostic(clusterId, id)
      setDiagnosticContent(prev => ({ ...prev, [id]: data.content }))
    } catch (error) {
      console.error("Error fetching diagnostic content:", error)
      toast.error("Failed to load diagnostic content")
    } finally {
      setIsDiagnosticContentLoading(prev => ({ ...prev, [id]: false }))
    }
  }

  const runDiagnostic = async () => {
    setIsDiagnosing(true)
    try {
      await DefaultService.createClusterDiagnostic(clusterId, {
        ID: 0, // Server will assign the actual ID
        createdAt: new Date().toISOString(),
        content: "" // Server will collect the diagnostic data
      })
      // Refresh the diagnostics list
      const diagnosticsData = await DefaultService.listClusterDiagnostics(clusterId)
      setDiagnostics(diagnosticsData)
      setTotalDiagnostics(diagnosticsData.length)
      toast.success("Diagnostic data collection initiated")
    } catch (error) {
      console.error("Error running diagnostic:", error)
      toast.error("Failed to run diagnostic")
    } finally {
      setIsDiagnosing(false)
    }
  }

  // Calculate paginated diagnostics
  const paginatedDiagnostics = diagnostics.slice(
    (diagnosticPage - 1) * DIAGNOSTICS_PER_PAGE,
    diagnosticPage * DIAGNOSTICS_PER_PAGE
  )

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

        <div className={`${width} space-y-4 border rounded-lg p-4`}>
          <div className="flex flex-col justify-start gap-2">
            <div className="flex flex-col justify-start">
              <Button
                onClick={() => void runRisectl()}
                disabled={isRunningCommand}
                className="whitespace-nowrap select-none w-fit"
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
            </div>
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

          </div>

          {risectlResult && (
            <Collapsible
              open={isResultOpen}
              onOpenChange={setIsResultOpen}
            >
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
                  {risectlResult.err && (
                    <div className="space-y-1">
                      <span className="text-sm font-medium text-red-600">Error:</span>
                      <div className="relative max-h-[400px] overflow-auto rounded-md">
                        <pre className="p-3 bg-red-50 text-red-600 text-sm whitespace-pre overflow-auto">
                          {risectlResult.err}
                        </pre>
                      </div>
                    </div>
                  )}

                  {risectlResult.stdout && (
                    <div className="space-y-1">
                      <span className="text-sm font-medium">stdout:</span>
                      <div className="relative max-h-[400px] overflow-auto rounded-md">
                        <pre className="p-3 bg-muted text-sm whitespace-pre overflow-auto">
                          {risectlResult.stdout}
                        </pre>
                      </div>
                    </div>
                  )}

                  {risectlResult.stderr && (
                    <div className="space-y-1">
                      <span className="text-sm font-medium text-red-600">stderr:</span>
                      <div className="relative max-h-[400px] overflow-auto rounded-md">
                        <pre className="p-3 bg-muted text-sm whitespace-pre overflow-auto">
                          {risectlResult.stderr}
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
        <div className={`flex items-center justify-between ${width}`}>
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

        <div className={`${width} space-y-4 border rounded-lg p-4 ${enableAutoBackup ? '' : 'hidden'}`}>
          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label className="text-sm font-medium">Auto Backup</Label>
              <p className="text-sm text-muted-foreground">Automatically create snapshots at regular intervals</p>
            </div>
            <Switch
              checked={autoBackupEnabled}
              onCheckedChange={(enabled) => {
                setAutoBackupEnabled(enabled)
                void updateAutoBackupConfig(enabled)
              }}
              disabled={isUpdatingBackupConfig}
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
                onValueChange={(interval) => {
                  setAutoBackupInterval(interval)
                  void updateAutoBackupConfig(undefined, interval)
                }}
                disabled={!autoBackupEnabled || isUpdatingBackupConfig}
              >
                <SelectTrigger className="w-[160px]">
                  <SelectValue placeholder="Select interval" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="*/5 * * * *">Every 5 minutes</SelectItem>
                  <SelectItem value="*/15 * * * *">Every 15 minutes</SelectItem>
                  <SelectItem value="*/30 * * * *">Every 30 minutes</SelectItem>
                  <SelectItem value="*/60 * * * *">Every 1 hour</SelectItem>
                  <SelectItem value="*/120 * * * *">Every 2 hours</SelectItem>
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
                onValueChange={(value) => {
                  const keepCount = parseInt(value)
                  setAutoBackupKeepCount(keepCount)
                  void updateAutoBackupConfig(undefined, undefined, keepCount)
                }}
                disabled={!autoBackupEnabled || isUpdatingBackupConfig}
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

        <div className={`${width} space-y-4`}>
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
        <div className={`flex items-center justify-between ${width}`}>
          <div className="space-y-1">
            <h3 className="text-lg font-semibold">Diagnostic Information</h3>
            <p className="text-sm text-muted-foreground">
              Configure automatic collection of diagnostic data and system metrics
            </p>
          </div>
          <Button
            onClick={() => void runDiagnostic()}
            disabled={isDiagnosing}
            className="select-none"
            size="sm"
          >
            {isDiagnosing ? (
              <>
                <svg className="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Collecting...
              </>
            ) : (
              "Collect Diagnostic"
            )}
          </Button>
        </div>

        <div className="space-y-6">
          {/* Configuration Card */}
          <div className={`${width} space-y-4 border rounded-lg p-4 ${enableAutoDiagnostic ? '' : 'hidden'}`}>
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label className="text-sm font-medium">Auto Diagnostic</Label>
                <p className="text-sm text-muted-foreground">Automatically collect diagnostic data</p>
              </div>
              <Switch
                checked={autoDiagnosticEnabled}
                onCheckedChange={(enabled) => {
                  setAutoDiagnosticEnabled(enabled)
                  void updateAutoDiagnosticConfig(enabled)
                }}
                disabled={isUpdatingDiagnosticConfig}
              />
            </div>
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label className="text-sm font-medium">Collection Interval</Label>
                <p className="text-sm text-muted-foreground">How often to collect data</p>
              </div>
              <Select 
                value={autoDiagnosticInterval} 
                onValueChange={(interval) => {
                  setAutoDiagnosticInterval(interval)
                  void updateAutoDiagnosticConfig(undefined, interval)
                }}
                disabled={!autoDiagnosticEnabled || isUpdatingDiagnosticConfig}
              >
                <SelectTrigger className="w-[180px]">
                  <SelectValue placeholder="Select interval" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="*/5 * * * *">Every 5 minutes</SelectItem>
                  <SelectItem value="*/15 * * * *">Every 15 minutes</SelectItem>
                  <SelectItem value="*/30 * * * *">Every 30 minutes</SelectItem>
                  <SelectItem value="0 * * * *">Every 1 hour</SelectItem>
                  <SelectItem value="0 */2 * * *">Every 2 hours</SelectItem>
                  <SelectItem value="0 */6 * * *">Every 6 hours</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label className="text-sm font-medium">Data Retention</Label>
                <p className="text-sm text-muted-foreground">How long to keep data</p>
              </div>
              <div className="flex flex-col gap-2">
                <Select
                  value={autoDiagnosticExpiration}
                  onValueChange={(expiration) => {
                    setAutoDiagnosticExpiration(expiration)
                    void updateAutoDiagnosticConfig(undefined, undefined, expiration)
                  }}
                  disabled={autoDiagnosticNoExpiration || !autoDiagnosticEnabled || isUpdatingDiagnosticConfig}
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
                <div className="flex items-center justify-end gap-1.5">
                  <Label className="text-sm text-muted-foreground">No expiration</Label>
                  <Switch
                    checked={autoDiagnosticNoExpiration}
                    onCheckedChange={(noExpiration) => {
                      setAutoDiagnosticNoExpiration(noExpiration)
                      void updateAutoDiagnosticConfig(undefined, undefined, undefined, noExpiration)
                    }}
                    disabled={!autoDiagnosticEnabled || isUpdatingDiagnosticConfig}
                  />
                </div>
              </div>

            </div>
          </div>

          {/* Diagnostics List */}
          <div className={`${width} space-y-4`}>
            {!diagnostics || diagnostics.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground border-2 border-dashed rounded-lg">
                No diagnostic data available. Click &quot;Collect Diagnostic&quot; to start collecting data.
              </div>
            ) : (
              <>
                {paginatedDiagnostics.map((diagnostic) => (
                  <Collapsible
                    key={diagnostic.ID}
                    onOpenChange={(isOpen: boolean) => {
                      if (isOpen && !diagnosticContent[diagnostic.ID]) {
                        void fetchDiagnosticContent(diagnostic.ID)
                      }
                      return undefined
                    }}
                  >
                    <CollapsibleTrigger className="flex items-center justify-between w-full p-4 hover:bg-accent/50 transition-colors border rounded-lg">
                      <div className="flex items-center gap-2">
                        <ChevronDown className="h-4 w-4" />
                        <span className="text-sm font-medium">
                          Diagnostic {diagnostic.ID}
                        </span>
                        <span className="text-sm text-muted-foreground">
                          {new Date(diagnostic.createdAt).toLocaleString()}
                        </span>
                      </div>
                    </CollapsibleTrigger>
                    <CollapsibleContent className="pt-2">
                      <div className="p-4 border rounded-lg mt-2">
                        {isDiagnosticContentLoading[diagnostic.ID] ? (
                          <div className="flex items-center justify-center py-4">
                            <svg className="animate-spin h-5 w-5 text-muted-foreground" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                          </div>
                        ) : (
                          <div className="space-y-2">
                            <div className="flex justify-end">
                              <Button
                                variant="ghost"
                                size="sm"
                                className="h-8 w-8 p-0"
                                onClick={() => {
                                  void navigator.clipboard.writeText(diagnosticContent[diagnostic.ID] || '')
                                  toast.success('Content copied to clipboard')
                                }}
                              >
                                <svg
                                  width="15"
                                  height="15"
                                  viewBox="0 0 15 15"
                                  fill="none"
                                  xmlns="http://www.w3.org/2000/svg"
                                  className="h-4 w-4"
                                >
                                  <path
                                    d="M1 9.50006C1 10.3285 1.67157 11.0001 2.5 11.0001H4L4 10.0001H2.5C2.22386 10.0001 2 9.7762 2 9.50006L2 2.50006C2 2.22392 2.22386 2.00006 2.5 2.00006L9.5 2.00006C9.77614 2.00006 10 2.22392 10 2.50006V4.00002H5.5C4.67157 4.00002 4 4.67159 4 5.50002V12.5C4 13.3284 4.67157 14 5.5 14H12.5C13.3284 14 14 13.3284 14 12.5V5.50002C14 4.67159 13.3284 4.00002 12.5 4.00002H11V2.50006C11 1.67163 10.3284 1.00006 9.5 1.00006H2.5C1.67157 1.00006 1 1.67163 1 2.50006V9.50006ZM5 5.50002C5 5.22388 5.22386 5.00002 5.5 5.00002H12.5C12.7761 5.00002 13 5.22388 13 5.50002V12.5C13 12.7762 12.7761 13 12.5 13H5.5C5.22386 13 5 12.7762 5 12.5V5.50002Z"
                                    fill="currentColor"
                                    fillRule="evenodd"
                                    clipRule="evenodd"
                                  />
                                </svg>
                                <span className="sr-only">Copy content</span>
                              </Button>
                            </div>
                            <div className="relative max-h-[400px] overflow-auto">
                              <pre className="whitespace-pre text-sm min-w-max p-4">
                                {diagnosticContent[diagnostic.ID] || 'No content available'}
                              </pre>
                            </div>
                          </div>
                        )}
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
                          if (diagnosticPage > 1) setDiagnosticPage(p => p - 1);
                        }}
                      />
                    </PaginationItem>
                    {[...Array(Math.ceil(totalDiagnostics / DIAGNOSTICS_PER_PAGE))].map((_, i) => (
                      <PaginationItem key={i + 1}>
                        <PaginationLink
                          href="#"
                          onClick={(e) => {
                            e.preventDefault();
                            setDiagnosticPage(i + 1);
                          }}
                          isActive={diagnosticPage === i + 1}
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
                          if (diagnosticPage < Math.ceil(totalDiagnostics / DIAGNOSTICS_PER_PAGE)) {
                            setDiagnosticPage(p => p + 1);
                          }
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
