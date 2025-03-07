import React from 'react';
import type { Metadata } from "next"
import { ProfileForm } from "@/components/profile/profile-form"
import { SecuritySettings } from "@/components/profile/security-settings"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

export const metadata: Metadata = {
  title: "Профиль | Secure Chat",
  description: "Управление профилем и настройками безопасности",
}

export default function ProfilePage() {
  return (
    <div className="container mx-auto py-10">
      <h1 className="text-3xl font-bold mb-6">Профиль</h1>

      <Tabs defaultValue="profile" className="w-full">
        <TabsList className="mb-6">
          <TabsTrigger value="profile">Личные данные</TabsTrigger>
          <TabsTrigger value="security">Безопасность</TabsTrigger>
        </TabsList>

        <TabsContent value="profile">
          <ProfileForm />
        </TabsContent>

        <TabsContent value="security">
          <SecuritySettings />
        </TabsContent>
      </Tabs>
    </div>
  )
}

