'use client'

import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
} from "@/components/ui/sidebar"
import {
  Server,
  Settings,
  Code,
  Moon,
  Sun,
} from "lucide-react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { cn } from "@/lib/utils"
import { useTheme } from "next-themes"
import { useEffect, useState } from "react"

const mainMenuItems = [
  {
    title: "Clusters",
    href: "/clusters",
    icon: Server,
  },
  {
    title: "SQL Console",
    href: "/sqlconsole",
    icon: Code,
  },
  {
    title: "Settings",
    href: "/settings",
    icon: Settings,
  },
]

export function AppSidebar() {
  const pathname = usePathname()
  const { theme, setTheme } = useTheme()
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  return (
    <Sidebar collapsible="icon" className="border-r">
      <div className="h-4" />
      <SidebarContent>
        <SidebarGroup>
          {mainMenuItems.map((item) => {
            const isSelected = pathname === item.href || pathname?.startsWith(item.href + '/')
            return (
              <SidebarMenu key={item.title}>
                <SidebarMenuButton asChild tooltip={item.title}>
                  <Link
                    href={item.href}
                    className="flex items-center justify-center"
                  >
                    <item.icon className={cn("h-4 w-4", isSelected ? "text-black dark:text-white" : "text-muted-foreground")} />
                  </Link>
                </SidebarMenuButton>
              </SidebarMenu>
            )
          })}
          <SidebarMenu>
            <SidebarMenuButton
              tooltip="Toggle theme"
              onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
            >
              <div className="flex items-center justify-center">
                {mounted && (theme === 'dark' ? (
                  <Sun className="h-4 w-4 text-muted-foreground" />
                ) : (
                  <Moon className="h-4 w-4 text-muted-foreground" />
                ))}
              </div>
            </SidebarMenuButton>
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  )
}
