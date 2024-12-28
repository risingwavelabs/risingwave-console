'use client'

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarCollapseButton,
  useSidebar,
  SidebarProvider,
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
  const { state } = useSidebar()
  const isCollapsed = state === "collapsed"

  return (
    <SidebarProvider defaultOpen={false}>
      <Sidebar collapsible="icon">
        <div className={`flex items-center py-2 ${isCollapsed ? 'justify-center' : 'justify-end mr-2'}`}>
          <SidebarCollapseButton />
        </div>
        <SidebarContent>
          <SidebarGroup>
            {mainMenuItems.map((item) => (
              <SidebarMenu key={item.title}>
                <SidebarMenuButton asChild tooltip={item.title}>
                  <Link
                    href={item.href}
                    className={`flex items-center ${isCollapsed ? 'justify-center' : 'px-4'}`}
                  >
                    <item.icon className={`h-4 w-4 ${!isCollapsed && 'mr-2'}`} />
                    {!isCollapsed && <span>{item.title}</span>}
                  </Link>
                </SidebarMenuButton>
              </SidebarMenu>
            ))}
          </SidebarGroup>
        </SidebarContent>
      </Sidebar>
    </SidebarProvider>
  )
}
