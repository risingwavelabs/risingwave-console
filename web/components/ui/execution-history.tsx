"use client"

import React, { useState, useEffect } from 'react';
import { ChevronDown, ChevronRight } from 'lucide-react';

interface ExecutionHistoryProps {
  history?: Array<{
    query: string;
    timestamp: string;
    status: 'success' | 'error';
    message: string;
    rowsAffected?: number;
    latencyMs?: number;
    rowCount?: number;
  }>;
}

export function ExecutionHistory({ history = [] }: ExecutionHistoryProps) {
  const [expandedItems, setExpandedItems] = useState<Set<number>>(new Set([0]));

  // Update expanded items when history changes to expand the first item
  useEffect(() => {
    if (history.length > 0) {
      setExpandedItems(new Set([0]));
    }
  }, [history]);

  const toggleExpand = (index: number) => {
    setExpandedItems(prev => {
      const newSet = new Set(prev);
      if (newSet.has(index)) {
        newSet.delete(index);
      } else {
        newSet.add(index);
      }
      return newSet;
    });
  };

  return (
    <div className="h-full flex flex-col overflow-hidden">
      <div className="flex-1 overflow-y-auto">
        <div className="space-y-4 p-4">
          {history.length === 0 ? (
            <div className="text-sm text-muted-foreground text-center py-4">
              No execution history available
            </div>
          ) : (
            history.map((item, index) => (
              <div
                key={index}
                className="border rounded-lg p-3 space-y-2 hover:bg-muted/30 transition-colors"
              >
                <div
                  className="flex items-center justify-between cursor-pointer"
                  onClick={() => toggleExpand(index)}
                >
                  <div className="flex items-center gap-2">
                    <span className="text-muted-foreground">
                      {expandedItems.has(index) ? (
                        <ChevronDown className="h-3.5 w-3.5" />
                      ) : (
                        <ChevronRight className="h-3.5 w-3.5" />
                      )}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      {new Date(item.timestamp).toLocaleString()}
                    </span>
                    {item.latencyMs !== undefined && (
                      <span className="text-xs text-muted-foreground">
                        {item.latencyMs}ms
                      </span>
                    )}
                    {item.status === 'success' && (
                      <span className="text-xs text-muted-foreground">
                        ({(item.rowCount ?? 0).toLocaleString()} row{(item.rowCount ?? 0) !== 1 ? 's' : ''})
                      </span>
                    )}
                  </div>
                  <div className="flex items-center gap-3">
                    {item.rowsAffected !== undefined && item.status === 'success' && (
                      <span className="text-xs text-muted-foreground">
                        {item.rowsAffected.toLocaleString()} row{item.rowsAffected !== 1 ? 's' : ''}
                      </span>
                    )}
                    <div className="flex items-center gap-2">
                      <span className={`text-xs ${item.status === 'success' ? 'text-green-500' : 'text-red-500'}`}>
                        {item.status === 'success' ? 'Success' : 'Error'}
                      </span>
                    </div>
                  </div>
                </div>
                {expandedItems.has(index) && (
                  <>
                    <p className={`text-sm ${item.status === 'success' ? 'text-green-500' : 'text-red-500'}`}>
                      {item.message}
                    </p>
                    <pre className="text-sm bg-muted/30 p-2 rounded overflow-x-auto">
                      <code>{item.query}</code>
                    </pre>

                  </>
                )}
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
} 