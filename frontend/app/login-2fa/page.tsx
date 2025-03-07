import React from 'react';
import { Metadata } from "next";
import Image from "next/image";
import { TwoFactorForm } from "@/components/auth/two-factor-form";

export const metadata: Metadata = {
  title: "2FA Verification | Secure Chat",
  description: "Подтвердите вход с помощью двухфакторной аутентификации",
};

export default function TwoFactorPage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-black text-white">
      <div className="absolute top-4 left-4">
        <Image
          src="/placeholder.svg?height=50&width=50"
          width={50}
          height={50}
          alt="Logo"
          className="rounded-full"
        />
      </div>
      
      <h1 className="text-4xl font-light mb-16">Нужен код 2FA</h1>
      
      <div className="w-full max-w-md">
        <TwoFactorForm />
      </div>
    </div>
  );
}
