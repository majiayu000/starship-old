'use client';

import React from 'react';
import { ContentItem, ContentRenderOptions, ContentType } from '../../components/content/types';

// 组件注册类型
export interface ComponentRegistration {
  component: React.ComponentType<any>;
  priority: number;
}

/**
 * 内容组件注册中心类
 * 负责管理和提供内容渲染组件
 */
class ContentComponentRegistry {
  private registry: Map<string, ComponentRegistration[]> = new Map();
  private defaultComponent: React.ComponentType<any> | null = null;
  
  /**
   * 注册组件
   * @param type 内容类型
   * @param component 组件
   * @param priority 优先级（较高的优先级会覆盖较低的优先级）
   */
  register(type: string, component: React.ComponentType<any>, priority: number = 0): void {
    if (!this.registry.has(type)) {
      this.registry.set(type, []);
    }
    
    const registrations = this.registry.get(type)!;
    registrations.push({ component, priority });
    
    // 按优先级排序
    registrations.sort((a, b) => b.priority - a.priority);
  }
  
  /**
   * 注册默认组件
   * @param component 默认组件，用于处理未注册类型
   */
  registerDefaultComponent(component: React.ComponentType<any>): void {
    this.defaultComponent = component;
  }
  
  /**
   * 获取组件
   * @param type 内容类型
   * @returns 对应的组件或默认组件
   */
  getComponent(type: string): React.ComponentType<any> | null {
    const registrations = this.registry.get(type);
    if (!registrations || registrations.length === 0) {
      return this.defaultComponent;
    }
    return registrations[0].component;
  }
  
  /**
   * 卸载组件
   * @param type 内容类型
   * @param component 可选，如果提供则只卸载特定组件
   */
  unregister(type: string, component?: React.ComponentType<any>): void {
    if (!this.registry.has(type)) return;
    
    if (component) {
      const registrations = this.registry.get(type)!;
      const newRegistrations = registrations.filter(reg => reg.component !== component);
      this.registry.set(type, newRegistrations);
    } else {
      this.registry.delete(type);
    }
  }
  
  /**
   * 检查类型是否已注册
   * @param type 内容类型
   */
  hasType(type: string): boolean {
    return this.registry.has(type) && this.registry.get(type)!.length > 0;
  }
  
  /**
   * 获取所有已注册类型
   */
  getRegisteredTypes(): string[] {
    return Array.from(this.registry.keys());
  }
}

// 创建全局内容组件注册中心实例
export const contentRegistry = new ContentComponentRegistry();

/**
 * 内容项工厂组件属性
 */
export interface ContentItemFactoryProps {
  item: ContentItem;
  options?: ContentRenderOptions;
}

/**
 * 内容项工厂组件
 * 工厂模式实现，根据内容类型选择适当的组件进行渲染
 */
export const ContentItemFactory: React.FC<ContentItemFactoryProps> = ({ 
  item, 
  options = {} 
}) => {
  if (!item) return null;
  
  const { type, content, metadata } = item;
  
  // 从注册中心获取组件
  const Component = contentRegistry.getComponent(type);
  
  if (!Component) {
    console.warn(`未找到内容类型 ${type} 的渲染组件`);
    return (
      <div className="unknown-content-type p-2 border border-yellow-300 bg-yellow-50 text-yellow-700 rounded text-xs">
        无法渲染 {type} 类型的内容
      </div>
    );
  }
  
  // 渲染组件
  try {
    return (
      <Component 
        content={content} 
        metadata={metadata} 
        options={options}
      />
    );
  } catch (error) {
    console.error(`渲染 ${type} 类型内容时出错:`, error);
    return (
      <div className="content-render-error p-2 border border-red-300 bg-red-50 text-red-700 rounded text-xs">
        渲染错误: {error instanceof Error ? error.message : '未知错误'}
      </div>
    );
  }
};

export default ContentItemFactory; 