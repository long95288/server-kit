### Description

some simple tool kits for server

### Design

后端接口API定义,所有的API接口定义为`$program/api/version/action`。为什么需要程序作为前缀呢?为了使用nginx进行转发方便。所有的`$program`的前缀均由该程序进行处理。

