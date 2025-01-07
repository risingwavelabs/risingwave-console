"use client"

import "./globals.css";
import { ThemeProvider } from "next-themes"
import { SidebarProvider } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/ui/app-sidebar"
import { usePathname } from "next/navigation"
import { Toaster } from "react-hot-toast"
import { initService } from "@/lib/init";

interface Props {
  children: React.ReactNode;
}

declare global {
  interface Window {
    hasInit: boolean;
  }
}

const ClientSideProvider = ({ children }: Props) => {
  const init = async () => {
    if (typeof window == "undefined" || window.hasInit) {
      return;
    }
    await initService();
    window.hasInit = true;
  };

  init();

  return (
    <>
      {children}
    </>
  );
};


export default function Layout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const isLoginPage = pathname === "/login"

  return (
    <html lang="en" className="h-full" suppressHydrationWarning>
      <body className="h-full w-full">
        <Toaster />
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <ClientSideProvider>
            {isLoginPage ? (
              children
            ) : (
              <SidebarProvider defaultOpen={false} className="h-full w-full">
                <div className="flex h-full w-full">
                  <AppSidebar />
                  <div className="flex-1 h-full w-full">
                    {children}
                  </div>
                </div>
              </SidebarProvider>
            )}
          </ClientSideProvider>
        </ThemeProvider>
      </body>
    </html>
  )
}
