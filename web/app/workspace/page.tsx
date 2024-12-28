'use client'

import { DatabaseList } from "@/components/ui/database-list"
import { useState, useCallback, useEffect, useRef } from 'react'
import Editor, { useMonaco } from '@monaco-editor/react'
import type { languages } from 'monaco-editor'
import { Button } from "@/components/ui/button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Play, PlayCircle, Save, X, Settings, Plus } from 'lucide-react'
import { DatabaseManagement } from "@/components/ui/database-management"
import { GenerateQuery } from "@/components/ui/generate-query"

// Sample data - replace with real data from your backend
const sampleDatabases = [
  {
    id: "db1",
    name: "Main Database",
    tables: [
      { id: "t1", name: "users" },
      { id: "t2", name: "products" },
      { id: "t3", name: "orders" },
    ],
  },
  {
    id: "db2",
    name: "Analytics DB",
    tables: [
      { id: "t4", name: "events" },
      { id: "t5", name: "metrics" },
    ],
  },
]

// Sample saved queries - replace with real data from your backend
const savedQueries = [
  { id: "q1", name: "Get Active Users" },
  { id: "q2", name: "Monthly Revenue" },
  { id: "q3", name: "Product Inventory" },
]

// Add sample completions data
const SQL_COMPLETIONS = {
  keywords: ['SELECT', 'FROM', 'WHERE', 'GROUP BY', 'ORDER BY', 'HAVING', 'JOIN', 'LEFT JOIN', 'RIGHT JOIN', 'INNER JOIN', 'LIMIT', 'OFFSET'],
  functions: ['COUNT', 'SUM', 'AVG', 'MAX', 'MIN', 'COALESCE', 'NULLIF', 'CAST', 'DATE', 'EXTRACT'],
  operators: ['AND', 'OR', 'NOT', 'IN', 'BETWEEN', 'LIKE', 'IS NULL', 'IS NOT NULL'],
}

const MIN_WIDTH = 150 // 9.375rem
const MAX_WIDTH = 400 // 25rem
const THEME_STORAGE_KEY = 'editor-theme'

// Sample result datasets
const SAMPLE_RESULTS = [
  {
    message: '3 row(s) affected',
    rows: [
      { id: 1, name: 'John Doe', email: 'john@example.com', status: 'active' },
      { id: 2, name: 'Jane Smith', email: 'jane@example.com', status: 'inactive' },
      { id: 3, name: 'Bob Wilson', email: 'bob@example.com', status: 'active' },
    ]
  },
  {
    message: '4 row(s) affected',
    rows: [
      { order_id: 'ORD001', product: 'Laptop', quantity: 2, total: 2400 },
      { order_id: 'ORD002', product: 'Mouse', quantity: 5, total: 100 },
      { order_id: 'ORD003', product: 'Keyboard', quantity: 3, total: 300 },
      { order_id: 'ORD004', product: 'Monitor', quantity: 1, total: 500 },
    ]
  },
  {
    message: '5 row(s) affected',
    rows: [
      { event_date: '2024-01-01', event_type: 'login', user_id: 'U123', platform: 'web' },
      { event_date: '2024-01-01', event_type: 'purchase', user_id: 'U124', platform: 'mobile' },
      { event_date: '2024-01-02', event_type: 'login', user_id: 'U125', platform: 'web' },
      { event_date: '2024-01-02', event_type: 'logout', user_id: 'U123', platform: 'web' },
      { event_date: '2024-01-03', event_type: 'purchase', user_id: 'U126', platform: 'mobile' },
    ]
  }
]

// Sample error messages
const ERROR_MESSAGES = [
  'Syntax error in SQL statement: Unexpected token near "FROM"',
  'Table "nonexistent_table" does not exist',
  'Column "unknown_column" not found in any table',
  'Invalid operator in WHERE clause',
  'Subquery returned more than 1 row'
]

// Sample AI-generated queries
const SAMPLE_AI_QUERIES = [
  `-- Find active users with recent purchases
SELECT DISTINCT u.id, u.name, u.email
FROM users u
JOIN orders o ON u.id = o.user_id
WHERE u.status = 'active'
  AND o.created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)
ORDER BY o.created_at DESC;`,

  `-- Calculate monthly revenue by product
SELECT 
  p.name as product_name,
  DATE_FORMAT(o.created_at, '%Y-%m') as month,
  COUNT(*) as total_orders,
  SUM(o.quantity) as units_sold,
  SUM(o.quantity * p.price) as revenue
FROM orders o
JOIN products p ON o.product_id = p.id
GROUP BY p.name, DATE_FORMAT(o.created_at, '%Y-%m')
ORDER BY month DESC, revenue DESC;`,

  `-- Track user engagement metrics
SELECT 
  e.user_id,
  u.name,
  COUNT(DISTINCT DATE(e.event_date)) as active_days,
  COUNT(*) as total_events,
  STRING_AGG(DISTINCT e.event_type) as event_types
FROM events e
JOIN users u ON e.user_id = u.id
WHERE e.event_date >= DATE_SUB(NOW(), INTERVAL 7 DAY)
GROUP BY e.user_id, u.name
HAVING COUNT(*) > 5
ORDER BY active_days DESC;`
]

interface EditorTab {
  id: string
  name: string
  content: string
  isDirty?: boolean
}

export default function WorkspacePage() {
  const [isResizing, setIsResizing] = useState(false)
  const [panelWidth, setPanelWidth] = useState(192) // 12rem default
  const [editorWidth, setEditorWidth] = useState(0)
  const [selectedQuery, setSelectedQuery] = useState("")
  const [queryResult, setQueryResult] = useState<{ type: 'success' | 'error', message: string, rows?: any[] }>()
  const [editorHeight, setEditorHeight] = useState('60%')
  const [isResizingHeight, setIsResizingHeight] = useState(false)
  const [selectedTheme, setSelectedTheme] = useState(() => {
    if (typeof window !== 'undefined') {
      return localStorage.getItem(THEME_STORAGE_KEY) || 'light'
    }
    return 'light'
  })
  const containerRef = useRef<HTMLDivElement>(null)
  const panelRef = useRef<HTMLDivElement>(null)
  const monaco = useMonaco()
  const [popupPosition, setPopupPosition] = useState({ x: 200, y: 200 })
  const [isDragging, setIsDragging] = useState(false)
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 })
  const [aiPrompt, setAIPrompt] = useState('')
  const [isGenerating, setIsGenerating] = useState(false)
  const [generateError, setGenerateError] = useState<string>("")
  const [loadingDots, setLoadingDots] = useState('')
  const editorRef = useRef<any>(null)
  const [showDatabaseManagement, setShowDatabaseManagement] = useState(false)
  const [tabs, setTabs] = useState<EditorTab[]>([
    { id: '1', name: 'Query 1', content: '-- Write your SQL query here', isDirty: false }
  ])
  const [activeTab, setActiveTab] = useState('1')
  const [editingTabId, setEditingTabId] = useState<string | null>(null)
  const [editingName, setEditingName] = useState('')
  const editInputRef = useRef<HTMLInputElement>(null)

  // Load tabs from localStorage on client-side only
  useEffect(() => {
    const savedTabs = localStorage.getItem('editor-tabs')
    if (savedTabs) {
      const parsedTabs = JSON.parse(savedTabs)
      setTabs(parsedTabs)
      // Set active tab to the first tab if the current active tab doesn't exist
      if (!parsedTabs.find((tab: EditorTab) => tab.id === activeTab)) {
        setActiveTab(parsedTabs[0].id)
      }
    }
  }, [])

  // Save tabs to local storage when they change
  useEffect(() => {
    localStorage.setItem('editor-tabs', JSON.stringify(tabs))
  }, [tabs])

  // Handle Ctrl+S
  useEffect(() => {
    const handleSave = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === 's') {
        e.preventDefault()
        setTabs(prev => prev.map(tab => 
          tab.id === activeTab ? { ...tab, isDirty: false } : tab
        ))
      }
    }

    window.addEventListener('keydown', handleSave)
    return () => window.removeEventListener('keydown', handleSave)
  }, [activeTab])

  const handleMouseDown = useCallback((e: React.MouseEvent) => {
    e.preventDefault()
    setIsResizing(true)
    document.body.style.cursor = 'ew-resize'
  }, [])

  const handleMouseUp = useCallback(() => {
    setIsResizing(false)
    document.body.style.cursor = 'default'
  }, [])

  const handleMouseMove = useCallback((e: React.MouseEvent<HTMLDivElement> | MouseEvent) => {
    if (!isResizing || !containerRef.current) return

    const containerRect = containerRef.current.getBoundingClientRect()
    const newWidth = e.clientX - containerRect.left

    if (newWidth >= MIN_WIDTH && newWidth <= MAX_WIDTH) {
      setPanelWidth(newWidth)
    }
  }, [isResizing])

  // Add and remove event listeners
  const handleMouseMoveGlobal = useCallback((e: MouseEvent) => {
    handleMouseMove(e)
  }, [handleMouseMove])

  useEffect(() => {
    if (isResizing) {
      window.addEventListener('mousemove', handleMouseMoveGlobal)
      window.addEventListener('mouseup', handleMouseUp)
    }

    return () => {
      window.removeEventListener('mousemove', handleMouseMoveGlobal)
      window.removeEventListener('mouseup', handleMouseUp)
    }
  }, [isResizing, handleMouseMoveGlobal, handleMouseUp])

  useEffect(() => {
    localStorage.setItem(THEME_STORAGE_KEY, selectedTheme)
  }, [selectedTheme])

  const calculateEditorWidth = useCallback(() => {
    if (containerRef.current) {
      const totalWidth = containerRef.current.getBoundingClientRect().width
      setEditorWidth(totalWidth - panelWidth - 5)
    }
  }, [panelWidth])

  useEffect(() => {
    calculateEditorWidth()
    window.addEventListener('resize', calculateEditorWidth)

    return () => {
      window.removeEventListener('resize', calculateEditorWidth)
    }
  }, [calculateEditorWidth])

  // Configure Monaco Editor with SQL completions
  useEffect(() => {
    if (monaco) {
      // Register SQL completion provider
      const disposable = monaco.languages.registerCompletionItemProvider('sql', {
        provideCompletionItems: (model, position) => {
          const word = model.getWordUntilPosition(position)
          const range = {
            startLineNumber: position.lineNumber,
            endLineNumber: position.lineNumber,
            startColumn: word.startColumn,
            endColumn: word.endColumn,
          }

          // Get the current line text up to the cursor
          const lineContent = model.getLineContent(position.lineNumber)
          const textUntilPosition = lineContent.substring(0, position.column - 1).toUpperCase()

          const suggestions: languages.CompletionItem[] = []

          // Add keyword suggestions
          SQL_COMPLETIONS.keywords.forEach(keyword => {
            if (keyword.startsWith(word.word.toUpperCase())) {
              suggestions.push({
                label: keyword,
                kind: monaco!.languages.CompletionItemKind.Keyword,
                insertText: keyword,
                range,
                preselect: true,
              })
            }
          })

          // Add function suggestions
          SQL_COMPLETIONS.functions.forEach(func => {
            if (func.startsWith(word.word.toUpperCase())) {
              suggestions.push({
                label: func,
                kind: monaco!.languages.CompletionItemKind.Function,
                insertText: func + '()',
                insertTextRules: monaco!.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                range,
              })
            }
          })

          // Add operator suggestions
          SQL_COMPLETIONS.operators.forEach(operator => {
            if (operator.startsWith(word.word.toUpperCase())) {
              suggestions.push({
                label: operator,
                kind: monaco!.languages.CompletionItemKind.Operator,
                insertText: operator,
                range,
              })
            }
          })

          // Add table suggestions from sampleDatabases
          sampleDatabases.forEach(db => {
            db.tables.forEach(table => {
              if (table.name.toLowerCase().startsWith(word.word.toLowerCase())) {
                suggestions.push({
                  label: table.name,
                  kind: monaco!.languages.CompletionItemKind.Class,
                  insertText: table.name,
                  range,
                  detail: `Table from ${db.name}`,
                })
              }
            })
          })

          return {
            suggestions,
          }
        },
        triggerCharacters: [' ', '.', '('],
      })

      return () => {
        disposable.dispose()
      }
    }
  }, [monaco])

  const handleTableSelect = (dbId: string, tableId: string) => {
    console.log(`Selected table ${tableId} from database ${dbId}`)
  }

  // Handle vertical resizing
  const handleMouseDownVertical = useCallback((e: React.MouseEvent) => {
    e.preventDefault()
    if (!containerRef.current) return
    
    const startY = e.clientY
    const startHeight = containerRef.current.getBoundingClientRect().height
    const startEditorHeight = parseFloat(editorHeight)
    
    const handleMouseMove = (e: MouseEvent) => {
      const deltaY = e.clientY - startY
      const newPercentage = (startEditorHeight + (deltaY / startHeight) * 100)
      
      // Limit the editor height between 30% and 80%
      if (newPercentage >= 30 && newPercentage <= 80) {
        setEditorHeight(`${newPercentage}%`)
      }
    }
    
    const handleMouseUp = () => {
      document.removeEventListener('mousemove', handleMouseMove)
      document.removeEventListener('mouseup', handleMouseUp)
      document.body.style.cursor = 'default'
    }
    
    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
    document.body.style.cursor = 'ns-resize'
  }, [editorHeight])

  // Remove the old vertical resize handlers since we're using the new approach
  const handleMouseUpVertical = useCallback(() => {
    setIsResizingHeight(false)
    document.body.style.cursor = 'default'
  }, [])

  const handleMouseMoveVertical = useCallback((e: React.MouseEvent<HTMLDivElement> | MouseEvent) => {
    if (!isResizingHeight || !containerRef.current) return

    const containerRect = containerRef.current.getBoundingClientRect()
    const containerHeight = containerRect.height
    const mouseY = e.clientY - containerRect.top
    const percentage = (mouseY / containerHeight) * 100

    // Limit the editor height between 30% and 80%
    if (percentage >= 30 && percentage <= 80) {
      setEditorHeight(`${percentage}%`)
    }
  }, [isResizingHeight])

  useEffect(() => {
    if (isResizingHeight) {
      window.addEventListener('mousemove', handleMouseMoveVertical)
      window.addEventListener('mouseup', handleMouseUpVertical)
    }

    return () => {
      window.removeEventListener('mousemove', handleMouseMoveVertical)
      window.removeEventListener('mouseup', handleMouseUpVertical)
    }
  }, [isResizingHeight, handleMouseMoveVertical, handleMouseUpVertical])

  // Sample function to handle query execution
  const handleRunQuery = useCallback(() => {
    // Simulate random success/error (80% success rate)
    const isSuccess = Math.random() < 0.8

    if (isSuccess) {
      // Pick a random successful result
      const randomResult = SAMPLE_RESULTS[Math.floor(Math.random() * SAMPLE_RESULTS.length)]
      setQueryResult({
        type: 'success',
        message: randomResult.message,
        rows: randomResult.rows
      })
    } else {
      // Pick a random error message
      const randomError = ERROR_MESSAGES[Math.floor(Math.random() * ERROR_MESSAGES.length)]
      setQueryResult({
        type: 'error',
        message: randomError
      })
    }
  }, [])

  // Handle popup dragging
  const handlePopupMouseDown = useCallback((e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      setIsDragging(true)
      setDragStart({
        x: e.clientX - popupPosition.x,
        y: e.clientY - popupPosition.y
      })
    }
  }, [popupPosition])

  const handlePopupMouseMove = useCallback((e: MouseEvent) => {
    if (isDragging) {
      setPopupPosition({
        x: e.clientX - dragStart.x,
        y: e.clientY - dragStart.y
      })
    }
  }, [isDragging, dragStart])

  const handlePopupMouseUp = useCallback(() => {
    setIsDragging(false)
  }, [])

  useEffect(() => {
    if (isDragging) {
      window.addEventListener('mousemove', handlePopupMouseMove)
      window.addEventListener('mouseup', handlePopupMouseUp)
    }

    return () => {
      window.removeEventListener('mousemove', handlePopupMouseMove)
      window.removeEventListener('mouseup', handlePopupMouseUp)
    }
  }, [isDragging, handlePopupMouseMove, handlePopupMouseUp])

  // Add editor reference handler
  const handleEditorDidMount = (editor: any) => {
    editorRef.current = editor
  }

  // Loading dots animation
  useEffect(() => {
    let interval: NodeJS.Timeout
    if (isGenerating) {
      interval = setInterval(() => {
        setLoadingDots(dots => dots.length >= 3 ? '' : dots + '.')
      }, 200)
    }
    return () => clearInterval(interval)
  }, [isGenerating])

  const handleGenerateQuery = useCallback((prompt: string) => {
    if (!prompt.trim()) return

    setIsGenerating(true)
    setGenerateError("")
    console.log("Generating query for prompt:", prompt)

    // Simulate AI query generation with a delay
    setTimeout(() => {
      try {
        const randomQuery = SAMPLE_AI_QUERIES[Math.floor(Math.random() * SAMPLE_AI_QUERIES.length)]
        if (editorRef.current) {
          const currentValue = editorRef.current.getValue()
          const newValue = currentValue
            ? `${currentValue}\n\n${randomQuery}`  // Add two newlines before new query
            : randomQuery
          editorRef.current.setValue(newValue)

          // Move cursor to the end of the editor
          const lineCount = editorRef.current.getModel().getLineCount()
          editorRef.current.setPosition({
            lineNumber: lineCount,
            column: editorRef.current.getModel().getLineMaxColumn(lineCount)
          })
          editorRef.current.focus()
        } else {
          throw new Error("Editor not initialized")
        }
      } catch (error) {
        setGenerateError(error instanceof Error ? error.message : "Failed to generate query")
      } finally {
        setIsGenerating(false)
      }
    }, 1000)
  }, [])

  const handleNewTab = () => {
    const newId = `${Date.now()}`
    setTabs(prev => [...prev, {
      id: newId,
      name: `Query ${prev.length + 1}`,
      content: '-- Write your SQL query here',
      isDirty: false
    }])
    setActiveTab(newId)
  }

  const handleCloseTab = (tabId: string, e: React.MouseEvent) => {
    e.stopPropagation()
    if (tabs.length === 1) return // Don't close the last tab

    setTabs(prev => prev.filter(tab => tab.id !== tabId))
    if (activeTab === tabId) {
      // Set the previous tab as active, or the first tab if closing the first one
      const index = tabs.findIndex(tab => tab.id === tabId)
      const newActiveTab = tabs[index - 1]?.id || tabs[index + 1]?.id
      setActiveTab(newActiveTab)
    }
  }

  const handleEditorChange = (value: string | undefined) => {
    if (!value) return
    setTabs(prev => prev.map(tab => 
      tab.id === activeTab 
        ? { ...tab, content: value, isDirty: true }
        : tab
    ))
  }

  const handleTabDoubleClick = (tab: EditorTab) => {
    setEditingTabId(tab.id)
    setEditingName(tab.name)
  }

  const handleTabNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEditingName(e.target.value)
  }

  const handleTabNameSave = () => {
    if (editingTabId && editingName.trim()) {
      setTabs(prev => prev.map(tab =>
        tab.id === editingTabId ? { ...tab, name: editingName.trim() } : tab
      ))
    }
    setEditingTabId(null)
  }

  const handleTabNameKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleTabNameSave()
    } else if (e.key === 'Escape') {
      setEditingTabId(null)
    }
  }

  // Focus input when editing starts
  useEffect(() => {
    if (editingTabId && editInputRef.current) {
      editInputRef.current.focus()
      editInputRef.current.select()
    }
  }, [editingTabId])

  return (
    <div 
      ref={containerRef} 
      className="flex h-full w-full relative"
      onMouseMove={handleMouseMove}
    >
      <div
        ref={panelRef}
        style={{ width: panelWidth }}
        className="relative flex-shrink-0 border-r bg-background h-full"
      >
        <div className="p-4 flex items-center justify-between w-full border-b pb-2">
          <h2 className="text-sm">Databases</h2>
          <Button
            variant="ghost"
            size="sm"
            className="h-7 px-2 text-muted-foreground hover:text-foreground"
            onClick={() => setShowDatabaseManagement(true)}
          >
            <Settings className="w-4 h-4 mr-1" />
          </Button>
        </div>
        <div className="p-4">
          <DatabaseList
            databases={sampleDatabases}
            onSelectTable={handleTableSelect}
          />
        </div>
        <div
          className="absolute right-0 top-0 bottom-0 w-[2px] cursor-ew-resize hover:bg-blue-500/20 active:bg-blue-500/40 transition-colors"
          onMouseDown={handleMouseDown}
        />
      </div>
      <div className="flex flex-col h-full min-w-0 flex-1">
        <div className="flex items-center gap-2 p-2 border-b">
          <Button size="sm" variant="default" onClick={handleRunQuery}>
            <Play className="w-4 h-4 mr-1" />
            Run
          </Button>
          <Button size="sm" variant="secondary">
            <PlayCircle className="w-4 h-4 mr-1" />
            Run Selected
          </Button>
          <Select value={selectedQuery} onValueChange={setSelectedQuery}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="My SQL Queries" />
            </SelectTrigger>
            <SelectContent>
              {savedQueries.map(query => (
                <SelectItem key={query.id} value={query.id}>
                  {query.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Button size="sm" variant="outline">
            <Save className="w-4 h-4 mr-1" />
            Save Query
          </Button>
          <Select value={selectedTheme} onValueChange={setSelectedTheme}>
            <SelectTrigger className="w-[120px]">
              <SelectValue placeholder="Select theme" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="light">Light</SelectItem>
              <SelectItem value="vs-dark">Dark</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="flex flex-col flex-1 min-h-0">
          <div className="border-b bg-muted/30 px-2">
            <div className="flex items-center">
              {tabs.map(tab => (
                <div
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`group flex items-center gap-2 px-4 py-2 border-r cursor-pointer hover:bg-muted/50 ${
                    activeTab === tab.id ? 'bg-background border-b-2 border-b-primary' : 'text-muted-foreground'
                  }`}
                >
                  <div className="flex items-center gap-1">
                    {editingTabId === tab.id ? (
                      <input
                        ref={editInputRef}
                        value={editingName}
                        onChange={handleTabNameChange}
                        onBlur={handleTabNameSave}
                        onKeyDown={handleTabNameKeyDown}
                        className="bg-transparent border-none outline-none focus:ring-1 focus:ring-primary text-sm px-0 w-[96px] h-5 -my-[1px]"
                        onClick={e => e.stopPropagation()}
                      />
                    ) : (
                      <span 
                        className="text-sm select-none"
                        onDoubleClick={() => handleTabDoubleClick(tab)}
                      >
                        {tab.name}
                      </span>
                    )}
                    {tab.isDirty && (
                      <span className="h-1.5 w-1.5 rounded-full bg-primary/70" />
                    )}
                  </div>
                  {tabs.length > 1 && (
                    <button
                      onClick={(e) => handleCloseTab(tab.id, e)}
                      className="opacity-0 group-hover:opacity-100 hover:text-red-500"
                    >
                      <X className="h-3 w-3" />
                    </button>
                  )}
                </div>
              ))}
              <button
                onClick={handleNewTab}
                className="p-2 hover:bg-muted/50"
              >
                <Plus className="h-4 w-4 text-muted-foreground" />
              </button>
            </div>
          </div>
          <div
            style={{ height: editorHeight }}
            className="relative"
          >
            <Editor
              defaultLanguage="sql"
              value={tabs.find(tab => tab.id === activeTab)?.content}
              onChange={handleEditorChange}
              theme={selectedTheme}
              onMount={handleEditorDidMount}
              options={{
                minimap: { enabled: false },
                fontSize: 14,
                scrollBeyondLastLine: false,
                automaticLayout: true,
                padding: { top: 16 },
                quickSuggestions: {
                  other: true,
                  comments: true,
                  strings: true
                },
                suggestOnTriggerCharacters: true,
                acceptSuggestionOnEnter: "on",
                tabCompletion: "on",
                suggest: {
                  preview: true,
                  showIcons: true,
                  showStatusBar: true,
                  showInlineDetails: true,
                }
              }}
              className="h-full"
              width={editorWidth}
            />
          </div>
          <div className="relative h-0">
            <div
              className="absolute inset-x-0 h-4 -top-2 cursor-ns-resize group"
              onMouseDown={handleMouseDownVertical}
            >
              <div className="absolute inset-x-0 top-[7px] h-[2px] bg-border group-hover:bg-blue-500/20 group-active:bg-blue-500/40" />
            </div>
          </div>
          <div className="flex-1 min-h-0 flex flex-col">
            <GenerateQuery
              onGenerate={handleGenerateQuery}
              isGenerating={isGenerating}
              error={generateError}
              className="border-b bg-background/95"
            />
            <div className="flex-1 overflow-auto p-4 min-h-0 bg-muted/30">
              {queryResult && (
                <div>
                  <div className={`mb-2 text-sm ${queryResult.type === 'error' ? 'text-red-500' : 'text-green-500'}`}>
                    {queryResult.message}
                  </div>
                  {queryResult.rows && (
                    <div className="overflow-auto">
                      <table className="w-full border-collapse">
                        <thead>
                          <tr>
                            {Object.keys(queryResult.rows[0]).map(key => (
                              <th key={key} className="text-left p-2 border bg-muted font-medium text-sm">
                                {key}
                              </th>
                            ))}
                          </tr>
                        </thead>
                        <tbody>
                          {queryResult.rows.map((row, i) => (
                            <tr key={i}>
                              {Object.values(row).map((value, j) => (
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
            </div>
          </div>
        </div>
      </div>
      <DatabaseManagement
        isOpen={showDatabaseManagement}
        onClose={() => setShowDatabaseManagement(false)}
      />
    </div>
  )
} 
