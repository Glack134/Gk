import React from 'react';
import type { Metadata } from "next"
import Image from "next/image"
import Link from "next/link"
import { SignupForm } from "@/components/auth/signup-form"

export const metadata: Metadata = {
  title: "Регистрация | Secure Chat",
  description: "Создайте новый аккаунт",
}

export default function SignupPage() {
  return (
    <div className="auth-layout flex min-h-screen items-center justify-center">
      <div className="w-full max-w-md p-6 bg-card rounded-lg shadow-lg dark:bg-card">
        <div className="flex justify-center mb-8">
          <Image src="/placeholder.svg?height=50&width=50" width={50} height={50} alt="Logo" className="rounded-full" />
        </div>

        <div className="flex justify-center space-x-2 mb-6">
          <Link
            href="/login"
            className="px-6 py-2 bg-secondary text-secondary-foreground rounded-l-full rounded-r-full"
          >
            Войти
          </Link>
          <Link href="/signup" className="px-6 py-2 bg-primary text-primary-foreground rounded-l-full rounded-r-full">
            Создать
          </Link>
        </div>

        <SignupForm />

        <div className="mt-6 flex justify-between">
          <Link href="/reset-password" className="text-sm text-muted-foreground hover:text-primary">
            Восстановить
          </Link>
          <Link href="/support" className="text-sm text-muted-foreground hover:text-primary">
            Поддержка
          </Link>
        </div>
      </div>

      <div className="hidden md:block absolute right-8 bottom-8 text-black dark:text-white">
        <h1 className="text-4xl font-light mb-2">Добро пожаловать!</h1>
        <p className="text-xl font-light">Давайте зарегистрируемся</p>
      </div>
    </div>
  )
}

