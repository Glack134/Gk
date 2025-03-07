"use client"
import React from 'react';
import type React from "react"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { useToast } from "@/components/ui/use-toast"

export function TwoFactorForm() {
  const [code, setCode] = useState(["", "", "", "", "", ""])
  const [isLoading, setIsLoading] = useState(false)
  const router = useRouter()
  const { toast } = useToast()

  const handleChange = (index: number, value: string) => {
    if (value.length > 1) {
      value = value.slice(0, 1)
    }

    const newCode = [...code]
    newCode[index] = value
    setCode(newCode)

    // Auto-focus next input
    if (value && index < 5) {
      const nextInput = document.getElementById(`code-${index + 1}`)
      if (nextInput) {
        nextInput.focus()
      }
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)

    try {
      // In a real app, this would be an API call to verify the 2FA code
      await new Promise((resolve) => setTimeout(resolve, 1000))

      // For demo purposes, we'll accept any code
      router.push("/chat")
    } catch (error) {
      toast({
        title: "Ошибка проверки",
        description: "Неверный код 2FA",
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-8">
      <div className="flex justify-center space-x-2">
        {code.map((digit, index) => (
          <input
            key={index}
            id={`code-${index}`}
            type="text"
            inputMode="numeric"
            pattern="[0-9]*"
            maxLength={1}
            value={digit}
            onChange={(e) => handleChange(index, e.target.value)}
            className="w-12 h-16 text-center text-xl border border-gray-600 rounded-md bg-black text-white focus:border-white focus:outline-none"
          />
        ))}
      </div>

      <Button type="submit" className="w-full" disabled={isLoading || code.some((digit) => !digit)}>
        {isLoading ? "Проверка..." : "Проверить"}
      </Button>
    </form>
  )
}

