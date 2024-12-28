'use client'

import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
} from "@/components/ui/sidebar"
import {
  LayoutDashboard,
  Settings,
} from "lucide-react"
import Link from "next/link"

const mainMenuItems = [
  {
    title: "Workspace",
    href: "/workspace",
    icon: LayoutDashboard,
  },
  {
    title: "Settings",
    href: "/settings",
    icon: Settings,
  },
]

export function AppSidebar() {
  return (
    <Sidebar collapsible="icon" className="border-r">
      <div className="h-4"/>
      <SidebarContent>
        <SidebarGroup>
          {mainMenuItems.map((item) => (
            <SidebarMenu key={item.title}>
              <SidebarMenuButton asChild tooltip={item.title}>
                <Link
                  href={item.href}
                  className="flex items-center justify-center"
                >
                  <item.icon className="h-4 w-4" />
                </Link>
              </SidebarMenuButton>
            </SidebarMenu>
          ))}
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  )
}
