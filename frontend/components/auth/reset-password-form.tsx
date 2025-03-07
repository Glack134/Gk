"use client"
import React from 'react';
import { useState } from "react"
import { useRouter } from "next/navigation"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import * as z from "zod"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Form, FormControl, FormField, FormItem, FormMessage } from "@/components/ui/form"
import { useToast } from "@/components/ui/use-toast"

const formSchema = z.object({
  email: z.string().email("Введите корректный email"),
})

export function ResetPasswordForm() {
  const [isLoading, setIsLoading] = useState(false)
  const router = useRouter()
  const { toast } = useToast()

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
    },
  })

  async function onSubmit(values: z.infer<typeof formSchema>) {
    setIsLoading(true)

    try {
      // In a real app, this would be an API call to your backend
      await new Promise((resolve) => setTimeout(resolve, 1000))

      toast({
        title: "Запрос отправлен",
        description: "Инструкции по восстановлению пароля отправлены на ваш email",
      })

      router.push("/login")
    } catch (error) {
      toast({
        title: "Ошибка",
        description: "Не удалось отправить запрос на восстановление пароля",
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <Input placeholder="Введите почту" {...field} className="bg-black border-gray-600 text-white" />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type="submit" className="w-full" disabled={isLoading}>
          {isLoading ? "Отправка..." : "Проверить"}
        </Button>
      </form>
    </Form>
  )
}

