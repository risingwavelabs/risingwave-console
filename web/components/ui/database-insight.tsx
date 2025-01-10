"use client"

import React from 'react';
import { StreamingGraph, RisingWaveNodeData } from "@/components/streaming-graph";
import { ProgressView } from "./progress-view";
import { ExecutionHistory } from "./execution-history";
import { QueryResult } from "./query-result";

interface DatabaseInsightProps {
  height?: string;
  databaseSchema?: RisingWaveNodeData[];
  result?: { type: 'success' | 'error', message: string, rows?: any[] };
  selectedDatabaseId?: string | null;
  onCancelProgress?: (ddlId: string) => void;
  executionHistory?: Array<{
    query: string;
    timestamp: string;
    status: 'success' | 'error';
    message: string;
    rowsAffected?: number;
  }>;
  activeTab: 'result' | 'graph' | 'progress' | 'history';
  onTabChange: (tab: 'result' | 'graph' | 'progress' | 'history') => void;
  isQueryLoading?: boolean;
}

export function DatabaseInsight({
  height = '30vh',
  databaseSchema = [],
  result,
  selectedDatabaseId,
  onCancelProgress,
  executionHistory = [],
  activeTab,
  onTabChange,
  isQueryLoading = false
}: DatabaseInsightProps) {
  return (
    <div className="flex flex-col h-full overflow-hidden" style={{ height }}>
      <div className="border-b flex">
        <button
          onClick={() => onTabChange('result')}
          className={`px-4 py-2 text-sm font-medium ${activeTab === 'result'
            ? 'border-b-2 border-primary text-foreground'
            : 'text-muted-foreground hover:text-foreground'
            }`}
        >
          Result
        </button>
        <button
          onClick={() => onTabChange('history')}
          className={`px-4 py-2 text-sm font-medium ${activeTab === 'history'
            ? 'border-b-2 border-primary text-foreground'
            : 'text-muted-foreground hover:text-foreground'
            }`}
        >
          History
        </button>
        <button
          onClick={() => onTabChange('graph')}
          className={`px-4 py-2 text-sm font-medium ${activeTab === 'graph'
            ? 'border-b-2 border-primary text-foreground'
            : 'text-muted-foreground hover:text-foreground'
            }`}
        >
          Streaming Graph
        </button>
        <button
          onClick={() => onTabChange('progress')}
          className={`px-4 py-2 text-sm font-medium ${activeTab === 'progress'
            ? 'border-b-2 border-primary text-foreground'
            : 'text-muted-foreground hover:text-foreground'
            }`}
        >
          Progress
        </button>
      </div>
      <div className="flex-1 min-h-0 bg-muted/30 overflow-hidden">
        {activeTab === 'result' && result && (
          <QueryResult result={result} isLoading={isQueryLoading} />
        )}
        {activeTab === 'graph' && (
          <div className="h-full">
            <StreamingGraph
              data={databaseSchema}
              height="100%"
              className="w-full h-full"
            />
          </div>
        )}
        {activeTab === 'progress' && (
          <ProgressView
            databaseId={selectedDatabaseId}
            onCancel={onCancelProgress}
          />
        )}
        {activeTab === 'history' && (
          <ExecutionHistory history={executionHistory} />
        )}
      </div>
    </div>
  );
}
