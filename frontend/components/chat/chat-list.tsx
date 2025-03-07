"use client"
import React from 'react';
import { useState } from "react"
import { Search } from "lucide-react"
import { Input } from "@/components/ui/input"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

// Mock data for chat list
const mockChats = [
  { id: 1, name: "User", message: "Message", time: "20:30", unread: false },
  { id: 2, name: "User", message: "Message", time: "20:30", unread: false },
  { id: 3, name: "User", message: "Message", time: "20:30", unread: false },
  { id: 4, name: "User", message: "Message", time: "20:30", unread: false },
  { id: 5, name: "User", message: "Message", time: "20:30", unread: false },
  { id: 6, name: "User", message: "Message", time: "20:30", unread: false },
  { id: 7, name: "User", message: "Message", time: "20:30", unread: false },
]

export function ChatList() {
  const [searchQuery, setSearchQuery] = useState("")
  const [activeChat, setActiveChat] = useState<number | null>(1)

  const filteredChats = mockChats.filter(
    (chat) =>
      chat.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      chat.message.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  return (
    <div className="w-80 border-r border-border flex flex-col">
      <div className="p-4 border-b border-border">
        <div className="relative">
          <Search className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Поиск"
            className="pl-10"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
      </div>

      <Tabs defaultValue="all" className="flex-1 flex flex-col">
        <div className="px-2">
          <TabsList className="w-full">
            <TabsTrigger value="all" className="flex-1">
              All in
            </TabsTrigger>
            <TabsTrigger value="favorite" className="flex-1">
              Favorite
            </TabsTrigger>
            <TabsTrigger value="chanel" className="flex-1">
              Chanel
            </TabsTrigger>
            <TabsTrigger value="wor" className="flex-1">
              Wor
            </TabsTrigger>
          </TabsList>
        </div>

        <TabsContent value="all" className="flex-1 overflow-y-auto">
          <div className="space-y-1 p-2">
            {filteredChats.map((chat) => (
              <button
                key={chat.id}
                className={`w-full flex items-center p-3 rounded-lg ${
                  activeChat === chat.id ? "bg-secondary" : "hover:bg-secondary/50"
                }`}
                onClick={() => setActiveChat(chat.id)}
              >
                <div className="w-10 h-10 bg-gray-300 rounded-full mr-3"></div>
                <div className="flex-1 text-left">
                  <div className="flex justify-between">
                    <span className="font-medium">{chat.name}</span>
                    <span className="text-xs text-muted-foreground">{chat.time}</span>
                  </div>
                  <p className="text-sm text-muted-foreground truncate">{chat.message}</p>
                </div>
              </button>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="favorite" className="flex-1">
          <div className="p-4 text-center text-muted-foreground">Нет избранных чатов</div>
        </TabsContent>

        <TabsContent value="chanel" className="flex-1">
          <div className="p-4 text-center text-muted-foreground">Нет каналов</div>
        </TabsContent>

        <TabsContent value="wor" className="flex-1">
          <div className="p-4 text-center text-muted-foreground">Нет рабочих чатов</div>
        </TabsContent>
      </Tabs>
    </div>
  )
}

