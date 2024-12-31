"use client"

import { useEffect, useState, useRef } from "react"
import { DndProvider, useDrag, useDrop } from "react-dnd"
import { HTML5Backend } from "react-dnd-html5-backend"
import { GripVertical, LayoutGrid, List, Database } from "lucide-react"
import { useRouter } from "next/navigation"
import { Button } from "./button"
import { Card, CardContent } from "./card"
import { ClusterDialog, ClusterFormData } from "./new-cluster-dialog"

export interface Cluster {
  id: string
  name: string
  status: "running" | "stopped" | "error"
  host: string
  sqlPort: number
  metaNodePort: number
  user: string
  password?: string
  database: string
}

type ViewMode = "grid" | "list"

interface ClusterListProps {
  clusters: Cluster[]
  onEdit?: (cluster: Cluster) => void
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
}

function DraggableClusterItem({ cluster, index, moveCluster, viewMode, onEdit }: DraggableClusterItemProps) {
  const router = useRouter()
  const ref = useRef<HTMLDivElement>(null)
  const handleRef = useRef<HTMLDivElement>(null)

  const [{ handlerId }, drop] = useDrop({
    accept: ItemTypes.CLUSTER,
    collect(monitor) {
      return {
        handlerId: monitor.getHandlerId(),
      }
    },
    hover(item: any, monitor) {
      if (!ref.current) {
        return
      }
      const dragIndex = item.index
      const hoverIndex = index

      if (dragIndex === hoverIndex) {
        return
      }

      moveCluster(dragIndex, hoverIndex)
      item.index = hoverIndex
    },
  })

  const [{ isDragging }, drag] = useDrag({
    type: ItemTypes.CLUSTER,
    item: () => ({ id: cluster.id, index }),
    collect: (monitor) => ({
      isDragging: monitor.isDragging(),
    }),
  })

  const opacity = isDragging ? 0.4 : 1
  drop(ref)
  drag(handleRef)

  return (
    <div ref={ref} style={{ opacity }} data-handler-id={handlerId}>
      <Card className="shadow-none hover:bg-accent/50 transition-colors">
        <CardContent className="flex items-center gap-4 p-4">
          <div 
            ref={handleRef}
            className="cursor-grab active:cursor-grabbing p-1 -m-1 text-muted-foreground hover:text-foreground transition-colors"
          >
            <GripVertical size={20} />
          </div>
          <div className="flex-1">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="font-semibold">{cluster.name}</h3>
                <div className="flex items-center gap-2 mt-1">
                  <div className={`w-2 h-2 rounded-full ${
                    cluster.status === "running" ? "bg-green-500" :
                    cluster.status === "stopped" ? "bg-yellow-500" :
                    "bg-red-500"
                  }`} />
                  <span className="text-sm text-muted-foreground capitalize">
                    {cluster.status}
                  </span>
                </div>
              </div>
              <div className="flex gap-2">
                <ClusterDialog
                  mode="edit"
                  defaultValues={cluster}
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
                  onClick={() => router.push(`/clusters/${cluster.id}`)}
                >
                  Manage
                </Button>
              </div>
            </div>
            <div className="grid grid-cols-2 gap-x-4 gap-y-1 mt-3 text-sm text-muted-foreground">
              <div className="flex items-center gap-1">
                <Database className="h-3 w-3" />
                <span>{cluster.database}</span>
              </div>
              <div>Host: {cluster.host}</div>
              <div>SQL Port: {cluster.sqlPort}</div>
              <div>Meta Port: {cluster.metaNodePort}</div>
              <div>User: {cluster.user}</div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

export function ClusterList({ clusters: initialClusters, onEdit }: ClusterListProps) {
  const [clusters, setClusters] = useState<Cluster[]>(() => {
    if (typeof window === "undefined") return initialClusters

    const savedOrder = localStorage.getItem(STORAGE_KEY)
    if (!savedOrder) return initialClusters

    const orderIds = JSON.parse(savedOrder) as string[]
    const orderedClusters = [...initialClusters]
    orderedClusters.sort((a, b) => {
      const aIndex = orderIds.indexOf(a.id)
      const bIndex = orderIds.indexOf(b.id)
      if (aIndex === -1) return 1
      if (bIndex === -1) return -1
      return aIndex - bIndex
    })
    return orderedClusters
  })

  const [viewMode, setViewMode] = useState<ViewMode>(() => {
    if (typeof window === "undefined") return "list"
    return (localStorage.getItem(VIEW_MODE_KEY) as ViewMode) || "list"
  })

  useEffect(() => {
    const orderIds = clusters.map(c => c.id)
    localStorage.setItem(STORAGE_KEY, JSON.stringify(orderIds))
  }, [clusters])

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
        {/* <div className="flex gap-2">
          <Button variant="outline" size="icon" onClick={toggleViewMode}>
            {viewMode === "list" ? <LayoutGrid size={20} /> : <List size={20} />}
          </Button>
          <Button>New Cluster</Button>
        </div> */}
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
                key={cluster.id}
                index={index}
                cluster={cluster}
                moveCluster={moveCluster}
                viewMode={viewMode}
                onEdit={onEdit}
              />
            ))}
          </div>
        </DndProvider>
      )}
    </div>
  )
} 