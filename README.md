# Qt 翻译工具

本质上是调用 lupdate 生成一个 ts 文件
从 zh_dict.json 中把英文替换成中文
最后调用 lrelease 把 ts 文件转换成 qm 文件

todo.txt 文件是待翻译列表，只需要把中文填上并 copy 到 zh_dict.json 中，再次执行 translator.exe 即可

<br/>
目前只处理了 windows 版本
具体用法请使用 translator.exe --help 查看

