"use client"

import { useEffect, useState } from "react"
import { ClusterList, Cluster as UICluster } from "@/components/ui/cluster-list"
import { ClusterDialog, ClusterFormData } from "@/components/ui/new-cluster-dialog"
import { DefaultService } from "@/api-gen"
import { Cluster as APICluster } from "@/api-gen/models/Cluster"
import toast from "react-hot-toast"

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
      toast.error("Failed to delete cluster")
      console.error("Error deleting cluster:", error)
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
    </div>
  )
}
