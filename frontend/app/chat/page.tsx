import React from 'react';
import { Metadata } from "next";
import { ChatLayout } from "@/components/chat/chat-layout";
import { ChatList } from "@/components/chat/chat-list";
import { ChatWindow } from "@/components/chat/chat-window";

export const metadata: Metadata = {
  title: "Чат | Secure Chat",
  description: "Общайтесь безопасно",
};

export default function ChatPage() {
  return (
    <ChatLayout>
      <ChatList />
      <ChatWindow />
    </ChatLayout>
  );
}
