## monstcache依赖问题

- monstcache和so插件必须基于相同的依赖库编译出来，否则无法兼容运行，所有依赖都在vendor;
- monstcache社区的一些特性和BUG修复不受控，故此单独在vendor中维护私有版本(vendor/github.com/rwynn/monstache/);
- 系统部分package引入了blog，其中blog参数和monstcache冲突造成启动异常，在私有版本中修改了monstcache命令行参数v;
