import React from 'react';
import { Html, Head, Main, NextScript } from "next/document"

export default function Document() {
  return (
    <Html lang="ru">
      <Head>{/* Здесь можно добавить дополнительные мета-теги, скрипты или стили */}</Head>
      <body>
        <Main />
        <NextScript />
      </body>
    </Html>
  )
}

