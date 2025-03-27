# 项目文档

本目录包含项目相关的技术文档、设计方案和指南。

## 文档列表

- [ELK 日志管理系统设计方案](elk_logging_system_design.md) - 详细的ELK栈日志管理系统设计方案，包括架构设计、配置示例和实施路线图。
- [ELK 日志收集最佳实践：基于文件的日志收集](elk_file_logging_best_practices.md) - 关于使用ELK栈收集文件日志的最佳实践指南，包括详细的配置、部署策略和故障排除方法。
- [ELK 日志收集路径选择：直连 vs. Logstash 处理](elk_logging_path_comparison.md) - 比较Filebeat直接连接Elasticsearch与使用Logstash处理的两种方法，帮助选择最适合您需求的日志收集路径。
- [Filebeat 配置与应用指南](filebeat_configuration_guide.md) - 全面的Filebeat配置指南，包括安装部署、基本和高级配置、性能优化、安全性配置和常见问题排查。
- [日志实现指南：从应用到 ELK](logging_implementation_guide.md) - 详细说明应用程序的日志格式规范、收集流程和最佳实践的完整指南。
- [Filebeat 配置过程记录](filebeat_setup_process.md) - 记录了 Filebeat 的配置过程，包括遇到的问题、解决方案以及最终配置，可作为实施参考。
- [Filebeat 安装配置与问题解决全记录](filebeat_setup_summary.md) - 全面记录了Filebeat的安装、配置和问题解决过程，包括Docker部署、索引模板配置、时区处理和常见问题解决方案。
- [cc-starship 应用日志架构设计](application_logging_architecture.md) - 详细说明了项目的日志架构设计，包括日志模型、处理器、格式规范和与Filebeat的集成方案。
- [API开发与日志系统集成最佳实践](api_logging_best_practices.md) - 提供了在API开发中正确集成日志系统的最佳实践指南，包括日志记录原则、级别使用、分布式追踪和性能考量。

## 文档使用指南

这些文档旨在提供项目的技术细节和最佳实践。在实施新功能或修改现有功能时，请参考相关文档。如需更新文档，请确保遵循相同的格式和结构。

## 贡献指南

添加新文档时，请遵循以下规则：

1. 使用Markdown格式
2. 文件名使用小写并以下划线分隔单词
3. 在文档开头添加清晰的目录结构
4. 更新此README.md文件，添加对新文档的引用 