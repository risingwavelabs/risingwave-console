'use client'

import * as React from "react"
import { Sparkles } from "lucide-react"
import { Button } from "./button"
import { Input } from "./input"

interface GenerateQueryProps {
  onGenerate: (prompt: string) => void
  isGenerating?: boolean
  error?: string
  className?: string
}

export function GenerateQuery({ onGenerate, isGenerating, error, className }: GenerateQueryProps) {
  const [prompt, setPrompt] = React.useState("")

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (prompt.trim() && !isGenerating) {
      onGenerate(prompt.trim())
      setPrompt("") // Clear the input after generating
    }
  }

  return (
    <div className={className}>
      <form onSubmit={handleSubmit}>
        <div className="flex gap-2 p-2 border-t">
          <Input
            placeholder="Describe the query you want to generate..."
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            disabled={isGenerating}
            className="flex-1"
          />
          <Button 
            type="submit" 
            disabled={!prompt.trim() || isGenerating}
            variant="outline"
            size="sm"
          >
            <Sparkles className="w-4 h-4 mr-2" />
            {isGenerating ? "Generating..." : "Generate"}
          </Button>
        </div>
      </form>
      {(isGenerating || error) && (
        <div className="px-2 pb-2">
          <p className={`text-sm ${error ? 'text-red-500' : 'text-muted-foreground'}`}>
            {error ? error : 'Generating your SQL query...'}
          </p>
        </div>
      )}
    </div>
  )
}