# softtrans

柔性事务实现(分布式事务)，包含基于TCC(Try/Confirm/Cancel)的最终一致事务...

## 特性

1. 支持TCC事务机制

    Try: 尝试执行业务
    ```markdown
    完成所有业务检查（一致性）
    预留必须业务资源（准隔离性）
    ```
    Confirm: 确认执行业务
    ```markdown
    真正执行业务
    不作任何业务检查
    只使用Try阶段预留的业务资源
    Confirm操作满足幂等性
    ```
    Cancel: 取消执行业务
    ```markdown
    释放Try阶段预留的业务资源
    Cancel操作满足幂等性
    ```
