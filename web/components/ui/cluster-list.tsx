"use client"

import { useEffect, useState, useRef } from "react"
import { DndProvider, useDrag, useDrop } from "react-dnd"
import { HTML5Backend } from "react-dnd-html5-backend"
import { GripVertical, LayoutGrid, List, Trash2 } from "lucide-react"
import { useRouter } from "next/navigation"
import { Button } from "./button"
import { Card, CardContent } from "./card"
import { ClusterDialog, ClusterFormData } from "./new-cluster-dialog"
import { ConfirmationPopup } from "./confirmation-popup"
import { DefaultService } from "@/api-gen"
import { MetricsStore } from "@/api-gen/models/MetricsStore"

export interface Cluster {
  id: number
  name: string
  status: "running" | "stopped" | "error"
  host: string
  sqlPort: number
  metaPort: number
  httpPort: number
  version: string
  metricsStoreID?: number | null
}

type ViewMode = "grid" | "list"

interface ClusterListProps {
  clusters: Cluster[]
  onEdit?: (cluster: Cluster) => void
  onDelete?: (cluster: Cluster) => void
}

const STORAGE_KEY = "cluster-order"
const VIEW_MODE_KEY = "cluster-view-mode"
const ItemTypes = {
  CLUSTER: "cluster"
}

interface DraggableClusterItemProps {
  cluster: Cluster
  index: number
  moveCluster: (dragIndex: number, hoverIndex: number) => void
  viewMode: ViewMode
  onEdit?: (cluster: Cluster) => void
  onDelete?: (cluster: Cluster) => void
}

interface DragItem {
  index: number
}

function DraggableClusterItem({ cluster, index, moveCluster, viewMode, onEdit, onDelete }: DraggableClusterItemProps) {
  const ref = useRef<HTMLDivElement>(null)
  const router = useRouter()
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false)
  const [metricsStoreName, setMetricsStoreName] = useState<string | null>(null)
  const [loadingStore, setLoadingStore] = useState(false)
  const [storeError, setStoreError] = useState(false)

  useEffect(() => {
    const fetchMetricsStoreName = async () => {
      if (cluster.metricsStoreID) {
        setLoadingStore(true)
        setStoreError(false)
        try {
          const store = await DefaultService.getMetricsStore(cluster.metricsStoreID)
          setMetricsStoreName(store.name)
        } catch (error) {
          console.error("Error fetching metrics store:", error)
          setStoreError(true)
        } finally {
          setLoadingStore(false)
        }
      } else {
        setMetricsStoreName(null)
      }
    }
    
    void fetchMetricsStoreName()
  }, [cluster.metricsStoreID])

  const [{ handlerId }, drop] = useDrop<DragItem, void, { handlerId: string | symbol | null }>({
    accept: ItemTypes.CLUSTER,
    collect(monitor) {
      return {
        handlerId: monitor.getHandlerId(),
      }
    },
    hover(item: DragItem, monitor) {
      if (!ref.current) {
        return
      }
      const dragIndex = item.index
      const hoverIndex = index

      if (dragIndex === hoverIndex) {
        return
      }

      const hoverBoundingRect = ref.current?.getBoundingClientRect()
      const hoverMiddleY = (hoverBoundingRect.bottom - hoverBoundingRect.top) / 2
      const clientOffset = monitor.getClientOffset()
      const hoverClientY = clientOffset!.y - hoverBoundingRect.top

      if (dragIndex < hoverIndex && hoverClientY < hoverMiddleY) {
        return
      }
      if (dragIndex > hoverIndex && hoverClientY > hoverMiddleY) {
        return
      }

      moveCluster(dragIndex, hoverIndex)
      item.index = hoverIndex
    },
  })

  const [{ isDragging }, drag] = useDrag<DragItem, void, { isDragging: boolean }>({
    type: ItemTypes.CLUSTER,
    item: () => ({ index }),
    collect: (monitor) => ({
      isDragging: monitor.isDragging(),
    }),
  })

  drag(drop(ref))

  return (
    <div
      ref={ref}
      style={{
        opacity: isDragging ? 0.5 : 1,
        transform: isDragging ? 'scale(1.02)' : 'scale(1)',
        transition: 'transform 0.2s ease, opacity 0.2s ease'
      }}
      className={`relative ${viewMode === 'grid' ? '' : ''}`}
      data-handler-id={handlerId}
    >
      <Card className={`
        border shadow-sm hover:shadow-md 
        transition-shadow duration-200 ease-in-out
        ${isDragging ? 'ring-2 ring-primary ring-opacity-50' : ''}
      `}>
        <CardContent className="p-4">
          <div className="absolute left-2 top-1/2 -translate-y-1/2 cursor-move opacity-30 hover:opacity-100 transition-opacity">
            <GripVertical className="h-5 w-5" />
          </div>
          <div className="flex-1 ml-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-4">
                <div>
                  <h3 className="font-semibold text-lg">{cluster.name}</h3>
                  <div className="flex items-center gap-2 mt-1">
                    <div className={`w-2 h-2 rounded-full ${cluster.status === "running" ? "bg-green-500" :
                      cluster.status === "stopped" ? "bg-yellow-500" :
                        "bg-red-500"
                      }`} />
                    <span className="text-sm text-muted-foreground capitalize">
                      {cluster.status}
                    </span>
                  </div>
                </div>
                <div className="h-10 w-px bg-border" /> {/* Vertical divider */}
                <div className="flex flex-col">
                  <span className="text-sm text-muted-foreground">Version</span>
                  <span className="text-sm font-medium">{cluster.version}</span>
                </div>
              </div>
              <div className="flex gap-2">
                <ClusterDialog
                  mode="edit"
                  defaultValues={{
                    name: cluster.name,
                    host: cluster.host,
                    sqlPort: cluster.sqlPort,
                    metaPort: cluster.metaPort,
                    httpPort: cluster.httpPort,
                    version: cluster.version,
                    metricsStoreID: cluster.metricsStoreID
                  }}
                  trigger={<Button variant="outline" size="sm">Edit</Button>}
                  onSubmit={(data: ClusterFormData) => onEdit?.({
                    ...data,
                    id: cluster.id,
                    status: cluster.status
                  })}
                />
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => router.push(`/clusters/details?id=${cluster.id}`)}
                >
                  Manage
                </Button>
                <div className="relative">
                  <Button
                    variant="outline"
                    size="sm"
                    className="text-red-500 hover:text-red-600"
                    onClick={() => setDeleteConfirmOpen(true)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                  {deleteConfirmOpen && (
                    <ConfirmationPopup
                      message="Delete this cluster?"
                      onConfirm={() => {
                        onDelete?.(cluster)
                        setDeleteConfirmOpen(false)
                      }}
                      onCancel={() => setDeleteConfirmOpen(false)}
                    />
                  )}
                </div>
              </div>
            </div>
            
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2 p-3 bg-muted/30 rounded-lg">
                <div className="flex items-center gap-2">
                  <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path strokeLinecap="round" strokeLinejoin="round" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                    </svg>
                  </div>
                  <div>
                    <div className="text-sm text-muted-foreground">Host</div>
                    <div className="font-medium">{cluster.host}</div>
                  </div>
                </div>
                {cluster.metricsStoreID && (
                  <div className="flex items-center gap-2">
                    <div className="w-8 h-8 rounded-full bg-orange-500/10 flex items-center justify-center">
                      <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-orange-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M16 8v8m-4-5v5m-4-2v2m-2 4h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                      </svg>
                    </div>
                    <div>
                      <div className="text-sm text-muted-foreground">Metrics Store</div>
                      {loadingStore ? (
                        <div className="flex items-center gap-1">
                          <svg className="animate-spin h-3 w-3 text-muted-foreground" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                          </svg>
                          <span className="text-xs">Loading...</span>
                        </div>
                      ) : storeError ? (
                        <div className="text-xs text-red-500">Error loading store</div>
                      ) : (
                        <div className="font-medium">{metricsStoreName || "Unknown"}</div>
                      )}
                    </div>
                  </div>
                )}
              </div>

              <div className="space-y-2 p-3 bg-muted/30 rounded-lg">
                <div className="flex items-center gap-2">
                  <div className="w-8 h-8 rounded-full bg-green-500/10 flex items-center justify-center">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-green-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path strokeLinecap="round" strokeLinejoin="round" d="M4 7v10c0 2 1 3 3 3h10c2 0 3-1 3-3V7c0-2-1-3-3-3H7C5 4 4 5 4 7z" />
                      <path strokeLinecap="round" strokeLinejoin="round" d="M8 12h8" />
                    </svg>
                  </div>
                  <div>
                    <div className="text-sm text-muted-foreground">Ports</div>
                    <div className="grid grid-cols-3 gap-2 font-medium">
                      <div>
                        <span className="text-xs text-muted-foreground">SQL</span>
                        <div>{cluster.sqlPort}</div>
                      </div>
                      <div>
                        <span className="text-xs text-muted-foreground">Meta</span>
                        <div>{cluster.metaPort}</div>
                      </div>
                      <div>
                        <span className="text-xs text-muted-foreground">HTTP</span>
                        <div>{cluster.httpPort}</div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

export function ClusterList({ clusters: initialClusters, onEdit, onDelete }: ClusterListProps) {
  const [clusters, setClusters] = useState<Cluster[]>(initialClusters)
  const [viewMode, setViewMode] = useState<ViewMode>("list")

  // Load saved order and view mode after mount
  useEffect(() => {
    const savedOrder = localStorage.getItem(STORAGE_KEY)
    if (savedOrder) {
      const orderIds = JSON.parse(savedOrder) as number[]
      const orderedClusters = [...initialClusters]
      orderedClusters.sort((a, b) => {
        const aIndex = orderIds.indexOf(a.id)
        const bIndex = orderIds.indexOf(b.id)
        if (aIndex === -1) return 1
        if (bIndex === -1) return -1
        return aIndex - bIndex
      })
      setClusters(orderedClusters)
    }

    const savedViewMode = localStorage.getItem(VIEW_MODE_KEY) as ViewMode
    if (savedViewMode) {
      setViewMode(savedViewMode)
    }
  }, [initialClusters])

  // Save order when it changes
  useEffect(() => {
    const orderIds = clusters.map(c => c.id)
    localStorage.setItem(STORAGE_KEY, JSON.stringify(orderIds))
  }, [clusters])

  // Save view mode when it changes
  useEffect(() => {
    localStorage.setItem(VIEW_MODE_KEY, viewMode)
  }, [viewMode])

  const moveCluster = (dragIndex: number, hoverIndex: number) => {
    const dragCluster = clusters[dragIndex]
    setClusters(prevClusters => {
      const newClusters = [...prevClusters]
      newClusters.splice(dragIndex, 1)
      newClusters.splice(hoverIndex, 0, dragCluster)
      return newClusters
    })
  }

  const toggleViewMode = () => {
    setViewMode(prev => prev === "list" ? "grid" : "list")
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h2 className="text-lg font-semibold">All Clusters</h2>
          <p className="text-sm text-muted-foreground">
            Manage and monitor your database clusters
          </p>
        </div>
        <Button variant="outline" size="icon" onClick={toggleViewMode}>
          {viewMode === "list" ? <LayoutGrid size={20} /> : <List size={20} />}
        </Button>
      </div>

      {clusters.length === 0 ? (
        <div className="text-center py-8 text-muted-foreground border-2 border-dashed rounded-lg">
          No clusters found. Create a new cluster to get started.
        </div>
      ) : (
        <DndProvider backend={HTML5Backend}>
          <div className={`grid ${viewMode === 'grid' ? 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3' : 'grid-cols-1'} gap-4`}>
            {clusters.map((cluster, index) => (
              <DraggableClusterItem
                key={`${cluster.id}-${index}`}
                index={index}
                cluster={cluster}
                moveCluster={moveCluster}
                viewMode={viewMode}
                onEdit={onEdit}
                onDelete={onDelete}
              />
            ))}
          </div>
        </DndProvider>
      )}
    </div>
  )
} 