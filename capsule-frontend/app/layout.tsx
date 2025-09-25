
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Capsule - Time Capsule Blockchain',
  description: 'Create and manage time-locked digital capsules on the blockchain',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
          <header className="bg-white shadow-sm border-b">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
              <div className="flex justify-between items-center py-4">
                <div className="flex items-center">
                  <h1 className="text-2xl font-bold text-gray-900">Capsule</h1>
                  <span className="ml-2 text-sm text-gray-500">Time Capsule Blockchain</span>
                </div>
                <nav className="flex space-x-4">
                  <a href="/" className="text-gray-600 hover:text-gray-900">Home</a>
                  <a href="/create" className="text-gray-600 hover:text-gray-900">Create Capsule</a>
                  <a href="/explore" className="text-gray-600 hover:text-gray-900">Explore</a>
                </nav>
              </div>
            </div>
          </header>
          <main>{children}</main>
        </div>
      </body>
    </html>
  )
}