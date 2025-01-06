"use client"

import { useEffect, useState, useRef } from "react"
import { DndProvider, useDrag, useDrop } from "react-dnd"
import { HTML5Backend } from "react-dnd-html5-backend"
import { GripVertical, LayoutGrid, List } from "lucide-react"
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

interface DragItem {
  index: number
}

function DraggableClusterItem({ cluster, index, moveCluster, viewMode, onEdit }: DraggableClusterItemProps) {
  const ref = useRef<HTMLDivElement>(null)
  const router = useRouter()

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
            <div className="flex items-center justify-between">
              <div>
                <h3 className="font-semibold">{cluster.name}</h3>
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
              <div className="flex gap-2">
                <ClusterDialog
                  mode="edit"
                  defaultValues={{
                    name: cluster.name,
                    host: cluster.host,
                    sqlPort: cluster.sqlPort,
                    metaNodePort: cluster.metaNodePort
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
                  onClick={() => router.push(`/clusters/${cluster.id}`)}
                >
                  Manage
                </Button>
              </div>
            </div>
            <div className="grid grid-cols-2 gap-x-4 gap-y-1 mt-3 text-sm text-muted-foreground">
              <div>Host: <span className="text-foreground">{cluster.host}</span></div>
              <div>SQL Port: <span className="text-foreground">{cluster.sqlPort}</span> Meta Port: <span className="text-foreground">{cluster.metaNodePort}</span></div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

export function ClusterList({ clusters: initialClusters, onEdit }: ClusterListProps) {
  const [clusters, setClusters] = useState<Cluster[]>(initialClusters)
  const [viewMode, setViewMode] = useState<ViewMode>("list")

  // Load saved order and view mode after mount
  useEffect(() => {
    const savedOrder = localStorage.getItem(STORAGE_KEY)
    if (savedOrder) {
      const orderIds = JSON.parse(savedOrder) as string[]
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
              />
            ))}
          </div>
        </DndProvider>
      )}
    </div>
  )
} 