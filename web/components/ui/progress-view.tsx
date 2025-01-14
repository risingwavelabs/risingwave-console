import React, { useState, useEffect } from 'react';
import { X, ChevronRight, ChevronDown } from 'lucide-react';
import { Button } from './button';
import Editor from '@monaco-editor/react';
import { useTheme } from 'next-themes';
import { DefaultService } from '@/api-gen';
import { number } from 'yup';

// Export the interface so it can be used by other components
export interface ProgressItem {
  name: string;
  status: 'running' | 'completed' | 'failed';
  progress: string;
  startTime: number;
  sql: string;
  ddlId: string;
}

export interface ProgressViewProps {
  databaseId?: string | null;
  onCancel?: (ddlId: string) => void;
}

export function ProgressView({ databaseId, onCancel }: ProgressViewProps) {
  const [items, setItems] = useState<ProgressItem[]>([]);
  const [cancelId, setCancelId] = useState<string | null>(null);
  const [now, setNow] = useState(Date.now());
  const [expandedItems, setExpandedItems] = useState<Set<string>>(new Set());
  const { theme } = useTheme();

  // Fetch progress items
  useEffect(() => {
    if (!databaseId) return;

    const fetchProgress = async () => {
      try {
        const progress = await DefaultService.getDdlProgress(Number(databaseId));
        setItems(progress.map(p => ({
          name: p.statement,
          status: parseProgress(p.progress) === 100 ? 'completed' : 'running',
          progress: p.progress,
          startTime: Date.parse(p.initializedAt),
          sql: p.statement,
          ddlId: String(p.ID)
        })));
      } catch (error) {
        console.error('Error fetching DDL progress:', error);
      }
    };

    fetchProgress();
    const interval = setInterval(fetchProgress, 1000);
    return () => clearInterval(interval);
  }, [databaseId]);

  // Update time every second for running items
  useEffect(() => {
    const hasRunningItems = items.some(
      (item: ProgressItem) => item.status === 'running'
    );

    if (hasRunningItems) {
      const interval = setInterval(() => {
        setNow(Date.now());
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [items]);

  const toggleExpand = (ddlId: string) => {
    setExpandedItems(prev => {
      const newSet = new Set(prev);
      if (newSet.has(ddlId)) {
        newSet.delete(ddlId);
      } else {
        newSet.add(ddlId);
      }
      return newSet;
    });
  };

  const getStatusColor = (status: ProgressItem['status']) => {
    switch (status) {
      case 'completed':
        return 'bg-green-500';
      case 'failed':
        return 'bg-red-500';
      default:
        return 'bg-gray-400';
    }
  };

  const getStatusIcon = (status: ProgressItem['status']) => {
    switch (status) {
      case 'completed':
        return '✓';
      case 'failed':
        return '✗';
      default:
        return '↻';
    }
  };

  const parseProgress = (progressStr: string): number => {
    const match = progressStr.match(/\d+\.?\d*/);
    if (match) {
      return Number(match[0]);
    }
    return 0
  };

  const formatDuration = (startTime: number) => {
    const duration = Math.floor((now - startTime) / 1000); // duration in seconds
    const hours = Math.floor(duration / 3600);
    const minutes = Math.floor((duration % 3600) / 60);
    const seconds = duration % 60;

    if (hours > 0) {
      return `${hours}h ${minutes}m ${seconds}s`;
    } else if (minutes > 0) {
      return `${minutes}m ${seconds}s`;
    } else {
      return `${seconds}s`;
    }
  };

  const handleConfirmCancel = (ddlId: string) => {
    if (onCancel) {
      onCancel(ddlId)
    }
    setCancelId(null)
  };

  return (
    <div className="h-full overflow-auto">
      <div className="p-4 space-y-2">
        {items.length === 0 ? (
          <div className="flex items-center justify-center text-muted-foreground text-sm">
            No operations in progress
          </div>
        ) : (
          items.map((item: ProgressItem, index: number) => (
            <div key={index} className="bg-background border rounded-md shadow-sm">
              <div
                className="p-2 cursor-pointer hover:bg-accent/50"
                onClick={() => toggleExpand(item.ddlId)}
              >
                <div className="flex items-center justify-between mb-1.5">
                  <div className="flex items-center gap-1.5">
                    <span className="text-muted-foreground">
                      {expandedItems.has(item.ddlId) ? (
                        <ChevronDown className="h-3.5 w-3.5" />
                      ) : (
                        <ChevronRight className="h-3.5 w-3.5" />
                      )}
                    </span>
                    <span className={`w-4 h-4 rounded-full flex items-center justify-center text-white text-[10px] ${getStatusColor(item.status)}`}>
                      {getStatusIcon(item.status)}
                    </span>
                    <span className="font-medium text-sm text-foreground">{item.progress}</span>
                    <span className="text-xs text-muted-foreground">
                      Elapased: {formatDuration(item.startTime)}
                    </span>
                  </div>
                  {item.status === 'running' && (
                    <div className="relative" onClick={e => e.stopPropagation()}>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 w-6 p-0 hover:bg-destructive/10 hover:text-destructive"
                        onClick={() => setCancelId(item.ddlId)}
                      >
                        <X className="h-3.5 w-3.5" />
                      </Button>
                      {cancelId === item.ddlId && (
                        <div className="absolute right-0 top-full mt-1 bg-background border rounded-lg shadow-lg py-1.5 px-2 z-50 whitespace-nowrap">
                          <p className="text-xs text-muted-foreground mb-1.5">Cancel this operation?</p>
                          <div className="flex gap-1.5">
                            <Button
                              size="sm"
                              variant="ghost"
                              className="h-6 px-2 text-xs"
                              onClick={() => setCancelId(null)}
                            >
                              No
                            </Button>
                            <Button
                              size="sm"
                              variant="destructive"
                              className="h-6 px-2 text-xs"
                              onClick={() => handleConfirmCancel(item.ddlId)}
                            >
                              Yes, cancel
                            </Button>
                          </div>
                        </div>
                      )}
                    </div>
                  )}
                </div>
                <div className="flex items-center gap-2">
                  <div className="flex-1 bg-accent rounded-full h-1.5">
                    <div
                      className={`h-1.5 rounded-full ${getStatusColor(item.status)} transition-all duration-500`}
                      style={{ width: `${parseProgress(item.progress)}%` }}
                    />
                  </div>
                </div>
              </div>
              {expandedItems.has(item.ddlId) && (
                <div className="px-2 pb-2 pt-0">
                  <div className="border rounded-md overflow-hidden">
                    <Editor
                      height="150px"
                      defaultLanguage="sql"
                      value={item.sql}
                      theme={theme === 'dark' ? 'vs-dark' : 'light'}
                      options={{
                        readOnly: true,
                        minimap: { enabled: false },
                        scrollBeyondLastLine: false,
                        lineNumbers: 'off',
                        folding: false,
                        tabSize: 2,
                      }}
                    />
                  </div>
                </div>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
} 