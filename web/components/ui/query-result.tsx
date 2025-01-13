"use client"

import React, { useState, useMemo, useEffect } from 'react';
import { ArrowUpDown, Loader2 } from "lucide-react";

interface QueryResultProps {
  result: { 
    type: 'success' | 'error';
    message: string;
    rows?: Record<string, unknown>[];
    columns?: string[];
    latencyMs?: number;
    rowCount?: number;
  };
  isLoading?: boolean;
}

function LoadingSpinner() {
  const [elapsedTime, setElapsedTime] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      setElapsedTime(prev => prev + 1);
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="flex flex-col items-center justify-center gap-4 h-full">
      <Loader2 className="h-8 w-8 animate-spin text-primary" />
      <div className="text-sm text-muted-foreground">
        Elapsed time: {elapsedTime}s
      </div>
    </div>
  );
}

export function QueryResult({ result, isLoading = false }: QueryResultProps) {
  const [sortConfig, setSortConfig] = useState<{
    key: string;
    direction: 'asc' | 'desc' | null;
  }>({ key: '', direction: null });

  const columns = useMemo(() => {
    if (result.columns) return result.columns;
    if (!result.rows?.length) return [];
    return Object.keys(result.rows[0]);
  }, [result.columns, result.rows]);

  const sortData = (data: Record<string, unknown>[], key: string) => {
    if (!data.length) return data;
    
    const sorted = [...data].sort((a, b) => {
      const aVal = a[key];
      const bVal = b[key];
      
      if (aVal === null || aVal === undefined) return 1;
      if (bVal === null || bVal === undefined) return -1;
      
      if (typeof aVal === 'number' && typeof bVal === 'number') {
        return aVal - bVal;
      }
      
      return String(aVal).localeCompare(String(bVal));
    });

    return sortConfig.direction === 'desc' ? sorted.reverse() : sorted;
  };

  const handleSort = (key: string) => {
    setSortConfig(current => ({
      key,
      direction:
        current.key === key
          ? current.direction === null
            ? 'asc'
            : current.direction === 'asc'
              ? 'desc'
              : null
          : 'asc'
    }));
  };

  if (isLoading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="h-full flex flex-col">
      <div className="flex-1 overflow-hidden">
        <div className="h-full flex flex-col">
          <div className="flex items-center justify-between p-4">
            <div className="flex items-center gap-4">
              <div className={`text-xs ${result.type === 'error' ? 'text-red-500' : 'text-green-500'}`}>
                {result.message}
              </div>
              <div className="text-xs text-muted-foreground">
                {(result.rowCount ?? result.rows?.length ?? 0).toLocaleString()} row{(result.rowCount ?? result.rows?.length ?? 0) !== 1 ? 's' : ''}
              </div>
            </div>
            {result.latencyMs !== undefined && (
              <div className="text-xs text-muted-foreground">
                {result.latencyMs}ms
              </div>
            )}
          </div>
          {result.rows && result.rows.length > 0 && result.columns && (
            <div className="flex-1 overflow-auto">
              <table className="w-full border-collapse">
                <thead>
                  <tr>
                    {columns.map((key: string) => (
                      <th 
                        key={key} 
                        className="text-left p-2 border bg-muted font-medium text-sm cursor-pointer hover:bg-muted/80 group"
                        onClick={() => handleSort(key)}
                      >
                        <div className="flex items-center justify-between gap-2">
                          <span>{key}</span>
                          <div className="flex items-center text-muted-foreground">
                            <ArrowUpDown className="h-3.5 w-3.5 transition-opacity opacity-0 group-hover:opacity-100" />
                            {sortConfig.key === key && (
                              <span className="ml-0.5 text-xs font-semibold text-foreground select-none">
                                {sortConfig.direction === 'asc' ? '↑' : sortConfig.direction === 'desc' ? '↓' : ''}
                              </span>
                            )}
                          </div>
                        </div>
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {(sortConfig.key
                    ? sortData(result.rows, sortConfig.key)
                    : result.rows
                  ).map((row: Record<string, unknown>, i: number) => (
                    <tr key={`row-${i}-${Object.values(row).join('-')}`}>
                      {columns.map((key: string) => (
                        <td key={`${key}-${i}`} className="p-2 border text-sm">
                          {String(row[key])}
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  );
} 