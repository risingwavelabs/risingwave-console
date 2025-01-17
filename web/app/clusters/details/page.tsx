"use client"

import { useSearchParams } from "next/navigation";
import ClusterPage from "@/components/clusterDetails";
import { Suspense } from "react";

function ClusterDetails() {
  const searchParams = useSearchParams();
  const id = searchParams.get("id");
  if (id) {
    return <div className="w-full">
      <ClusterPage params={{ id: parseInt(id) }} />
    </div>;
  }
  return <div>No cluster id provided</div>;
}

export default function ClusterDetailsPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <ClusterDetails />
    </Suspense>
  );
}
