"use client"

import { useState } from "react"
import { ClusterList } from "@/components/ui/cluster-list"

interface Cluster {
  id: string
  name: string
  status: "running" | "stopped" | "error"
  nodes: number
}

const sampleClusters: Cluster[] = [
  {
    id: "cluster-1",
    name: "Production DB",
    status: "running",
    nodes: 5
  },
  {
    id: "cluster-2",
    name: "Staging Environment",
    status: "running",
    nodes: 3
  },
  {
    id: "cluster-3",
    name: "Development DB",
    status: "stopped",
    nodes: 2
  },
  {
    id: "cluster-4",
    name: "Analytics Cluster",
    status: "running",
    nodes: 4
  },
  {
    id: "cluster-5",
    name: "Testing Environment",
    status: "error",
    nodes: 1
  },
  {
    id: "cluster-6",
    name: "Backup DB",
    status: "stopped",
    nodes: 2
  },
  {
    id: "cluster-7",
    name: "Data Warehouse",
    status: "running",
    nodes: 6
  },
  {
    id: "cluster-8",
    name: "Reporting DB",
    status: "running",
    nodes: 3
  },
  {
    id: "cluster-9",
    name: "QA Environment",
    status: "stopped",
    nodes: 2
  },
  {
    id: "cluster-10",
    name: "Archive DB",
    status: "running",
    nodes: 4
  },
  {
    id: "cluster-11",
    name: "ML Training Cluster",
    status: "error",
    nodes: 8
  },
  {
    id: "cluster-12",
    name: "Cache Cluster",
    status: "running",
    nodes: 3
  }
]

export default function ClustersPage() {
  const [clusters, setClusters] = useState<Cluster[]>(sampleClusters)

  return (
    <div className="p-8">
      <ClusterList clusters={clusters} />
    </div>
  )
}
