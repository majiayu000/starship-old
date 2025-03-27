'use client';

import React, { ErrorInfo, ReactNode } from 'react';

interface ErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error?: Error;
}

/**
 * 内容错误边界组件
 * 捕获子组件中的 JavaScript 错误，并显示回退 UI
 */
export class ContentErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    // 可以将错误日志发送到服务器
    console.error('内容渲染错误:', error, errorInfo);
  }

  render(): ReactNode {
    if (this.state.hasError) {
      // 自定义回退 UI
      return this.props.fallback || (
        <div className="content-error p-2 border border-red-300 bg-red-50 text-red-700 rounded">
          <p className="text-sm">内容渲染失败</p>
        </div>
      );
    }

    return this.props.children;
  }
}

/**
 * 函数式错误边界包装器
 * 为函数组件提供一个简单的错误边界包装
 */
export const withErrorBoundary = <P extends object>(
  Component: React.ComponentType<P>,
  fallback?: ReactNode
): React.FC<P> => {
  return (props: P) => (
    <ContentErrorBoundary fallback={fallback}>
      <Component {...props} />
    </ContentErrorBoundary>
  );
}; 