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
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Switch } from "@/components/ui/switch"

const passwordFormSchema = z
  .object({
    currentPassword: z.string().min(6, "Пароль должен содержать минимум 6 символов"),
    newPassword: z.string().min(6, "Пароль должен содержать минимум 6 символов"),
    confirmPassword: z.string().min(6, "Пароль должен содержать минимум 6 символов"),
  })
  .refine((data) => data.newPassword === data.confirmPassword, {
    message: "Пароли не совпадают",
    path: ["confirmPassword"],
  })

export function SecuritySettings() {
  const [isLoading, setIsLoading] = useState(false)
  const [twoFAEnabled, setTwoFAEnabled] = useState(false)
  const [showTwoFASetup, setShowTwoFASetup] = useState(false)
  const { toast } = useToast()

  const form = useForm<z.infer<typeof passwordFormSchema>>({
    resolver: zodResolver(passwordFormSchema),
    defaultValues: {
      currentPassword: "",
      newPassword: "",
      confirmPassword: "",
    },
  })

  async function onSubmit(values: z.infer<typeof passwordFormSchema>) {
    setIsLoading(true)

    try {
      // In a real app, this would be an API call to your backend
      await new Promise((resolve) => setTimeout(resolve, 1000))

      toast({
        title: "Пароль обновлен",
        description: "Ваш пароль успешно изменен",
      })

      form.reset()
    } catch (error) {
      toast({
        title: "Ошибка",
        description: "Не удалось обновить пароль",
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  const handleTwoFAToggle = async () => {
    if (!twoFAEnabled) {
      setShowTwoFASetup(true)
    } else {
      // In a real app, this would be an API call to disable 2FA
      await new Promise((resolve) => setTimeout(resolve, 1000))

      setTwoFAEnabled(false)
      toast({
        title: "2FA отключена",
        description: "Двухфакторная аутентификация отключена",
      })
    }
  }

  const setupTwoFA = async () => {
    // In a real app, this would be an API call to enable 2FA
    await new Promise((resolve) => setTimeout(resolve, 1000))

    setTwoFAEnabled(true)
    setShowTwoFASetup(false)
    toast({
      title: "2FA включена",
      description: "Двухфакторная аутентификация успешно включена",
    })
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Изменить пароль</CardTitle>
          <CardDescription>Обновите свой пароль для повышения безопасности</CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="currentPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Текущий пароль</FormLabel>
                    <FormControl>
                      <Input type="password" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="newPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Новый пароль</FormLabel>
                    <FormControl>
                      <Input type="password" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="confirmPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Подтвердите новый пароль</FormLabel>
                    <FormControl>
                      <Input type="password" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button type="submit" disabled={isLoading}>
                {isLoading ? "Обновление..." : "Обновить пароль"}
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Двухфакторная аутентификация (2FA)</CardTitle>
          <CardDescription>Повысьте безопасность вашего аккаунта с помощью 2FA</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div>
              <h4 className="font-medium">Статус 2FA</h4>
              <p className="text-sm text-muted-foreground">{twoFAEnabled ? "Включена" : "Отключена"}</p>
            </div>
            <Switch checked={twoFAEnabled} onCheckedChange={handleTwoFAToggle} />
          </div>

          {showTwoFASetup && (
            <div className="mt-6 p-4 border rounded-lg">
              <h4 className="font-medium mb-2">Настройка 2FA</h4>
              <p className="text-sm text-muted-foreground mb-4">
                Отсканируйте QR-код с помощью приложения аутентификации (Google Authenticator, Authy и т.д.)
              </p>

              <div className="flex justify-center mb-4">
                <div className="w-48 h-48 bg-gray-200 flex items-center justify-center">
                  <p className="text-xs text-center">QR-код для сканирования</p>
                </div>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="text-sm font-medium">Введите код подтверждения</label>
                  <Input className="mt-1" placeholder="Введите 6-значный код" />
                </div>

                <div className="flex space-x-2">
                  <Button onClick={setupTwoFA}>Подтвердить</Button>
                  <Button variant="outline" onClick={() => setShowTwoFASetup(false)}>
                    Отмена
                  </Button>
                </div>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Сеансы</CardTitle>
          <CardDescription>Управляйте активными сеансами вашего аккаунта</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="p-4 border rounded-lg">
              <div className="flex justify-between">
                <div>
                  <h4 className="font-medium">Текущий сеанс</h4>
                  <p className="text-sm text-muted-foreground">Браузер: Chrome, ОС: Windows, IP: 192.168.1.1</p>
                </div>
                <div className="text-sm text-green-500">Активен</div>
              </div>
            </div>

            <Button variant="destructive">Завершить все другие сеансы</Button>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}