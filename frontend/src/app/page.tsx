'use client';

import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function Home() {
  return (
    <div className="container mx-auto p-8 flex flex-col items-center justify-center min-h-[calc(100vh-80px)]">
      <div className="max-w-3xl w-full">
        <h1 className="text-3xl font-bold text-center mb-8">欢迎使用 SAT 题目审核系统</h1>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card>
            <CardHeader>
              <CardTitle>审核页面</CardTitle>
              <CardDescription>使用传统的表格视图审核题目</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-500">
                此页面提供标准的列表视图，适用于快速浏览和批量审核题目。
                支持多种筛选和排序选项。
              </p>
            </CardContent>
            <CardFooter>
              <Link href="/review" className="w-full">
                <Button className="w-full">进入审核页面</Button>
              </Link>
            </CardFooter>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>题目审核</CardTitle>
              <CardDescription>使用详细视图审核题目</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-500">
                此页面提供完整题目内容和详细信息，适用于深入审核题目内容。
                默认显示待审核题目，支持实时审核提交。
              </p>
            </CardContent>
            <CardFooter>
              <Link href="/questions" className="w-full">
                <Button className="w-full">进入题目审核</Button>
              </Link>
            </CardFooter>
          </Card>
        </div>
      </div>
    </div>
  );
}
