"use client"

import { useSearchParams } from "next/navigation";
import ClusterPage from "@/components/clusterDetails";
import { Suspense } from "react";

function ClusterDetails() {
  const searchParams = useSearchParams();
  const id = searchParams.get("id");
  if (id) {
    return <ClusterPage params={{ id }} />;
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
