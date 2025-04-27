'use client';

import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function Home() {
  return (
    <div className="container mx-auto p-8 flex flex-col items-center justify-center min-h-[calc(100vh-80px)]">
      <div className="max-w-3xl w-full">
        <h1 className="text-3xl font-bold text-center mb-8">Welcome to the Starter Template</h1>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card>
            <CardHeader>
              <CardTitle>Item Management</CardTitle>
              <CardDescription>Manage your items with a full CRUD interface</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-500">
                This page provides a standard list view with filtering, sorting, and pagination.
                You can create, view, edit, and delete items.
              </p>
            </CardContent>
            <CardFooter>
              <Link href="/items" className="w-full">
                <Button className="w-full">Go to Items</Button>
              </Link>
            </CardFooter>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Dashboard</CardTitle>
              <CardDescription>View analytics and statistics</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-500">
                This page provides an overview of your data with charts and statistics.
                Monitor key metrics and track performance.
              </p>
            </CardContent>
            <CardFooter>
              <Link href="/dashboard" className="w-full">
                <Button className="w-full">Go to Dashboard</Button>
              </Link>
            </CardFooter>
          </Card>
        </div>
      </div>
    </div>
  );
}
