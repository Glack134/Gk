import React from 'react';
import { redirect } from "next/navigation"

export default function Home() {
  // In a real app, we would check for authentication here
  // For now, we'll just redirect to the login page
  redirect("/login")
}

