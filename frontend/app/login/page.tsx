import React from 'react';
import { Metadata } from "next";
import Image from "next/image";
import Link from "next/link";
import { LoginForm } from "@/components/auth/login-form";

export const metadata: Metadata = {
  title: "Войти | Secure Chat",
  description: "Войдите в свой аккаунт",
};

export default function LoginPage() {
  return (
    <div className="auth-layout flex min-h-screen items-center justify-center">
      <div className="w-full max-w-md p-6 bg-card rounded-lg shadow-lg dark:bg-card">
        <div className="flex justify-center mb-8">
          <Image
            src="/placeholder.svg?height=50&width=50"
            width={50}
            height={50}
            alt="Logo"
            className="rounded-full"
          />
        </div>
        
        <div className="flex justify-center space-x-2 mb-6">
          <Link 
            href="/login" 
            className="px-6 py-2 bg-primary text-primary-foreground rounded-l-full rounded-r-full"
          >
            Войти
          </Link>
          <Link 
            href="/signup" 
            className="px-6 py-2 bg-secondary text-secondary-foreground rounded-l-full rounded-r-full"
          >
            Создать
          </Link>
        </div>
        
        <LoginForm />
        
        <div className="mt-6 flex justify-between">
          <Link 
            href="/reset-password" 
            className="text-sm text-muted-foreground hover:text-primary"
          >
            Восстановить
          </Link>
          <Link 
            href="/support" 
            className="text-sm text-muted-foreground hover:text-primary"
          >
            Поддержка
          </Link>
        </div>
      </div>
      
      <div className="hidden md:block absolute left-8 bottom-8 text-white">
        <h1 className="text-4xl font-light mb-2">Добро пожаловать!</h1>
        <p className="text-xl font-light">Мы рады вас видеть</p>
      </div>
    </div>
  );
}
