"use client"

import React, { useState, useCallback } from 'react';
import { StreamingGraph, RisingWaveNodeData } from "@/components/streaming-graph";
import { ProgressView } from "./progress-view";

interface DatabaseInsightProps {
  height?: string;
  databaseSchema?: RisingWaveNodeData[];
  result?: { type: 'success' | 'error', message: string, rows?: any[] };
}

export function DatabaseInsight({ height = '30vh', databaseSchema = [], result }: DatabaseInsightProps) {
  const [activeResultTab, setActiveResultTab] = useState<'result' | 'graph' | 'progress'>('result')

  const handleCancelProgress = useCallback((name: string) => {
    // Here you would typically make an API call to cancel the materialized view creation
    console.log(`Cancelling materialized view creation: ${name}`);
  }, []);

  return (
    <div className="flex flex-col h-full overflow-hidden" style={{ height }}>
      <div className="border-b flex">
        <button
          onClick={() => setActiveResultTab('result')}
          className={`px-4 py-2 text-sm font-medium ${activeResultTab === 'result'
            ? 'border-b-2 border-primary text-foreground'
            : 'text-muted-foreground hover:text-foreground'
            }`}
        >
          Result
        </button>
        <button
          onClick={() => setActiveResultTab('graph')}
          className={`px-4 py-2 text-sm font-medium ${activeResultTab === 'graph'
            ? 'border-b-2 border-primary text-foreground'
            : 'text-muted-foreground hover:text-foreground'
            }`}
        >
          Streaming Graph
        </button>
        <button
          onClick={() => setActiveResultTab('progress')}
          className={`px-4 py-2 text-sm font-medium ${activeResultTab === 'progress'
            ? 'border-b-2 border-primary text-foreground'
            : 'text-muted-foreground hover:text-foreground'
            }`}
        >
          Progress
        </button>
      </div>
      <div className="flex-1 min-h-0 bg-muted/30 overflow-hidden">
        {activeResultTab === 'result' && result && (
          <div className="p-4 h-full overflow-auto">
            <div className={`mb-2 text-sm ${result.type === 'error' ? 'text-red-500' : 'text-green-500'}`}>
              {result.message}
            </div>
            {result.rows && (
              <div className="overflow-auto">
                <table className="w-full border-collapse">
                  <thead>
                    <tr>
                      {Object.keys(result.rows[0]).map((key: string) => (
                        <th key={key} className="text-left p-2 border bg-muted font-medium text-sm">
                          {key}
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {result.rows.map((row: Record<string, any>, i: number) => (
                      <tr key={i}>
                        {Object.values(row).map((value: any, j: number) => (
                          <td key={j} className="p-2 border text-sm">
                            {String(value)}
                          </td>
                        ))}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}
        {activeResultTab === 'graph' && (
          <div className="h-full">
            <StreamingGraph
              data={databaseSchema}
              height="100%"
              className="w-full h-full"
            />
          </div>
        )}
        {activeResultTab === 'progress' && (
          <div className="h-full overflow-auto">
            <ProgressView onCancel={handleCancelProgress} />
          </div>
        )}
      </div>
    </div>
  );
}
