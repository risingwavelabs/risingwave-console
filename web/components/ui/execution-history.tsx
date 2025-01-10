import React from 'react';
import { ScrollArea } from "./scroll-area";

interface ExecutionHistoryProps {
  history?: Array<{
    query: string;
    timestamp: string;
    status: 'success' | 'error';
    message: string;
    rowsAffected?: number;
  }>;
}

export function ExecutionHistory({ history = [] }: ExecutionHistoryProps) {
  return (
    <ScrollArea className="h-full">
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
              <div className="flex items-center justify-between">
                <span className="text-xs text-muted-foreground">
                  {new Date(item.timestamp).toLocaleString()}
                </span>
                <span className={`text-xs ${
                  item.status === 'success' ? 'text-green-500' : 'text-red-500'
                }`}>
                  {item.status === 'success' ? 'Success' : 'Error'}
                  {item.rowsAffected !== undefined && item.status === 'success' && 
                    ` (${item.rowsAffected} rows affected)`
                  }
                </span>
              </div>
              <pre className="text-sm bg-muted/30 p-2 rounded overflow-x-auto">
                <code>{item.query}</code>
              </pre>
              <p className={`text-sm ${
                item.status === 'success' ? 'text-green-500' : 'text-red-500'
              }`}>
                {item.message}
              </p>
            </div>
          ))
        )}
      </div>
    </ScrollArea>
  );
} 