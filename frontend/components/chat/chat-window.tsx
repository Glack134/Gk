"use client"

import { useState } from "react"
import { Send, Paperclip, Smile } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

// Mock data for messages
const mockMessages = [
  { id: 1, sender: "User", content: "Привет! Как дела?", time: "10:30", isMine: false },
  { id: 2, sender: "Me", content: "Привет! Все хорошо, спасибо. Как у тебя?", time: "10:32", isMine: true },
  { id: 3, sender: "User", content: "Тоже неплохо. Что нового?", time: "10:33", isMine: false },
  { id: 4, sender: "Me", content: "Работаю над новым проектом. Очень интересно!", time: "10:35", isMine: true },
]

export function ChatWindow() {
  const [message, setMessage] = useState("")
  const [messages, setMessages] = useState(mockMessages)

  const handleSendMessage = () => {
    if (!message.trim()) return

    const newMessage = {
      id: messages.length + 1,
      sender: "Me",
      content: message,
      time: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
      isMine: true,
    }

    setMessages([...messages, newMessage])
    setMessage("")
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault()
      handleSendMessage()
    }
  }

  return (
    <div className="flex-1 flex flex-col">
      {/* Chat header */}
      <div className="p-4 border-b border-border flex items-center">
        <div className="w-10 h-10 bg-gray-300 rounded-full mr-3"></div>
        <div>
          <h2 className="font-medium">User</h2>
          <p className="text-xs text-muted-foreground">Online</p>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((msg) => (
          <div key={msg.id} className={`flex ${msg.isMine ? "justify-end" : "justify-start"}`}>
            <div
              className={`max-w-[70%] rounded-lg p-3 ${
                msg.isMine ? "bg-primary text-primary-foreground" : "bg-secondary text-secondary-foreground"
              }`}
            >
              <p>{msg.content}</p>
              <p className="text-xs mt-1 opacity-70">{msg.time}</p>
            </div>
          </div>
        ))}
      </div>

      {/* Message input */}
      <div className="p-4 border-t border-border">
        <div className="flex items-center space-x-2">
          <Button variant="ghost" size="icon">
            <Paperclip className="h-5 w-5" />
          </Button>

          <Input
            placeholder="Введите сообщение..."
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            onKeyDown={handleKeyDown}
            className="flex-1"
          />

          <Button variant="ghost" size="icon">
            <Smile className="h-5 w-5" />
          </Button>

          <Button variant="primary" size="icon" onClick={handleSendMessage} disabled={!message.trim()}>
            <Send className="h-5 w-5" />
          </Button>
        </div>
      </div>
    </div>
  )
}

