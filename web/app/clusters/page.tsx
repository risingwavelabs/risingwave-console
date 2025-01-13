"use client"

import { AxiosError } from "axios"
import { useEffect, useState } from "react"
import { ClusterList, Cluster as UICluster } from "@/components/ui/cluster-list"
import { ClusterDialog, ClusterFormData } from "@/components/ui/new-cluster-dialog"
import { ApiError, DefaultService } from "@/api-gen"
import { Cluster as APICluster } from "@/api-gen/models/Cluster"
import toast from "react-hot-toast"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"

function mapAPIClusterToUICluster(apiCluster: APICluster): UICluster {
  return {
    id: apiCluster.ID.toString(),
    name: apiCluster.name,
    status: "running", // You might want to derive this from API data
    host: apiCluster.host,
    sqlPort: apiCluster.sqlPort,
    metaPort: apiCluster.metaPort,
  }
}

export default function ClustersPage() {
  const [clusters, setClusters] = useState<UICluster[]>([])
  const [loading, setLoading] = useState(true)
  const [cascadeDeleteCluster, setCascadeDeleteCluster] = useState<UICluster | null>(null)
  const [cascadeDeleteError, setCascadeDeleteError] = useState<string | null>(null)

  useEffect(() => {
    const fetchClusters = async () => {
      try {
        const data = await DefaultService.listClusters()
        setClusters(data.map(mapAPIClusterToUICluster))
      } catch (error) {
        toast.error("Failed to load clusters")
        console.error("Error loading clusters:", error)
      } finally {
        setLoading(false)
      }
    }

    void fetchClusters()
  }, [])

  const handleCreateCluster = async (data: ClusterFormData) => {
    try {
      const newCluster = await DefaultService.createCluster(data)
      setClusters(prev => [...prev, mapAPIClusterToUICluster(newCluster)])
      toast.success("Cluster created successfully")
    } catch (error) {
      toast.error("Failed to create cluster")
      console.error("Error creating cluster:", error)
    }
  }

  const handleEditCluster = async (cluster: UICluster) => {
    try {
      const updatedCluster = await DefaultService.updateCluster(cluster.id, {
        name: cluster.name,
        host: cluster.host,
        sqlPort: cluster.sqlPort,
        metaPort: cluster.metaPort,
      })
      setClusters(prev => prev.map(c =>
        c.id === cluster.id ? mapAPIClusterToUICluster(updatedCluster) : c
      ))
      toast.success("Cluster updated successfully")
    } catch (error) {
      toast.error("Failed to update cluster")
      console.error("Error updating cluster:", error)
    }
  }

  const handleDeleteCluster = async (cluster: UICluster) => {
    try {
      await DefaultService.deleteCluster(cluster.id)
      setClusters(prev => prev.filter(c => c.id !== cluster.id))
      toast.success("Cluster deleted successfully")
    } catch (error) {
      console.log("cascade delete cluster")
      if (error instanceof ApiError && error.status === 409) {
        const errorMessage = error.message || "This cluster has active database connections."
        setCascadeDeleteCluster(cluster)
        setCascadeDeleteError(errorMessage)
      } else {
        toast.error("Failed to delete cluster")
        console.error("Error deleting cluster:", error)
      }
    }
  }

  const handleCascadeDelete = async () => {
    if (!cascadeDeleteCluster) return

    try {
      await DefaultService.deleteCluster(cascadeDeleteCluster.id, true)
      setClusters(prev => prev.filter(c => c.id !== cascadeDeleteCluster.id))
      toast.success("Cluster and associated connections deleted successfully")
    } catch (error) {
      toast.error("Failed to delete cluster")
      console.error("Error deleting cluster:", error)
    } finally {
      setCascadeDeleteCluster(null)
    }
  }

  if (loading) {
    return (
      <div className="p-8">
        <div className="flex items-center justify-center h-[200px]">
          <div className="text-muted-foreground">Loading clusters...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-semibold">Clusters</h1>
          <p className="text-sm text-muted-foreground">
            Manage your database clusters
          </p>
        </div>
        <ClusterDialog mode="create" onSubmit={handleCreateCluster} />
      </div>
      <ClusterList
        clusters={clusters}
        onEdit={handleEditCluster}
        onDelete={handleDeleteCluster}
      />
      {cascadeDeleteCluster && (
        <Dialog open={!!cascadeDeleteCluster} onOpenChange={() => setCascadeDeleteCluster(null)}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Delete Cluster</DialogTitle>
              <DialogDescription>
                {cascadeDeleteError}
                This cluster has active database connections. Would you like to delete the cluster and all associated connections?
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button variant="outline" onClick={() => setCascadeDeleteCluster(null)}>
                Cancel
              </Button>
              <Button variant="destructive" onClick={handleCascadeDelete}>
                Yes, delete everything
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </div>
  )
}
