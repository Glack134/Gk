"use client"
import React from 'react';
import type React from "react"

import { useState } from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { Home, Users, MessageSquare, Settings, LogOut } from "lucide-react"

interface ChatLayoutProps {
  children: React.ReactNode
}

export function ChatLayout({ children }: ChatLayoutProps) {
  const pathname = usePathname()
  const [darkMode, setDarkMode] = useState(true)

  const toggleDarkMode = () => {
    setDarkMode(!darkMode)
    document.documentElement.classList.toggle("dark")
  }

  return (
    <div className={`flex h-screen ${darkMode ? "dark" : ""}`}>
      {/* Sidebar */}
      <div className="w-16 bg-black flex flex-col items-center py-6">
        <div className="w-10 h-10 bg-gray-700 rounded-full mb-10"></div>

        <nav className="flex flex-col items-center space-y-8 flex-1">
          <Link href="/chat" className={`p-2 rounded-lg ${pathname === "/chat" ? "bg-gray-800" : ""}`}>
            <MessageSquare className="w-6 h-6 text-white" />
          </Link>
          <Link href="/contacts" className={`p-2 rounded-lg ${pathname === "/contacts" ? "bg-gray-800" : ""}`}>
            <Users className="w-6 h-6 text-white" />
          </Link>
          <Link href="/" className={`p-2 rounded-lg ${pathname === "/" ? "bg-gray-800" : ""}`}>
            <Home className="w-6 h-6 text-white" />
          </Link>
        </nav>

        <div className="mt-auto flex flex-col items-center space-y-6">
          <Link href="/profile" className={`p-2 rounded-lg ${pathname === "/profile" ? "bg-gray-800" : ""}`}>
            <Settings className="w-6 h-6 text-white" />
          </Link>
          <button className="p-2 rounded-lg">
            <LogOut className="w-6 h-6 text-white" />
          </button>
        </div>
      </div>

      {/* Main content */}
      <div className="flex-1 bg-background dark:bg-background flex">{children}</div>
    </div>
  )
}

