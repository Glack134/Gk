import React from 'react';
import type { Metadata } from "next"
import Image from "next/image"
import { ResetPasswordForm } from "@/components/auth/reset-password-form"

export const metadata: Metadata = {
  title: "Восстановление пароля | Secure Chat",
  description: "Восстановите доступ к вашему аккаунту",
}

export default function ResetPasswordPage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-black text-white">
      <div className="absolute top-4 left-4">
        <Image src="/placeholder.svg?height=50&width=50" width={50} height={50} alt="Logo" className="rounded-full" />
      </div>

      <h1 className="text-4xl font-light mb-16">Восстановление</h1>

      <div className="w-full max-w-md">
        <ResetPasswordForm />
      </div>
    </div>
  )
}

