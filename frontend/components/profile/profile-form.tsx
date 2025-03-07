"use client"

import React from 'react';
import { useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import * as z from "zod"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { useToast } from "@/components/ui/use-toast"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

const formSchema = z.object({
  username: z.string().min(3, "Имя пользователя должно содержать минимум 3 символа"),
  email: z.string().email("Введите корректный email"),
  fullName: z.string().min(2, "Имя должно содержать минимум 2 символа"),
  phone: z.string().optional(),
})

export function ProfileForm() {
  const [isLoading, setIsLoading] = useState(false)
  const { toast } = useToast()

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "username",
      email: "user@example.com",
      fullName: "User Name",
      phone: "+7 (999) 123-45-67",
    },
  })

  async function onSubmit(values: z.infer<typeof formSchema>) {
    setIsLoading(true)

    try {
      // In a real app, this would be an API call to your backend
      await new Promise((resolve) => setTimeout(resolve, 1000))

      toast({
        title: "Профиль обновлен",
        description: "Ваши данные успешно сохранены",
      })
    } catch (error) {
      toast({
        title: "Ошибка",
        description: "Не удалось обновить профиль",
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Фото профиля</CardTitle>
          <CardDescription>Обновите свое фото профиля</CardDescription>
        </CardHeader>
        <CardContent className="flex items-center space-x-4">
          <Avatar className="h-24 w-24">
            <AvatarImage src="/placeholder.svg?height=96&width=96" />
            <AvatarFallback>UN</AvatarFallback>
          </Avatar>
          <div>
            <Button variant="outline" className="mr-2">
              Загрузить новое фото
            </Button>
            <Button variant="ghost">Удалить</Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Личная информация</CardTitle>
          <CardDescription>Обновите свои личные данные</CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="username"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Имя пользователя</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="fullName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Полное имя</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="phone"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Телефон</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button type="submit" disabled={isLoading}>
                {isLoading ? "Сохранение..." : "Сохранить изменения"}
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  )
}

