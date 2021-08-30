提供两种同步策略

- notify_hexo:监听本地目录变化,按需生成文章.
- ticker_hexo:定时将本地目录重新生成一遍



为了防止多次重复提交,notify_hexo会将默认10min内的消息聚合后才提交.具体逻辑在ticker_channel.go中.ticker_hexo默认每24h执行一次.

因为notify_hexo使用的第三方库fsnotify行为在不同平台的执行效果不一致且难以预测,因此有一定概率生成错误文章或历史文章没有删除.

(eg:对于操作`mv test/test1.txt test/test2.txt`,在unix会按序触发create test2.txt; remove test1.txt; rename test1.txt三个事件. 但是在windows平台触发create test2.txt; remove test1.txt; write test; rename test1.txt;事件)

所以建议使用时同时开启notify_hexo和ticker_hexo,既保证监听的响应速度,又能在闲暇时纠错.

当然,因为错误只会局限在文章,相关资源的重新生成文章的时候不会重新下载,复制.且由于文章数据都是纯文本,因此ticker_hexo重新生成文章并部署能在秒级完成,因此只使用ticker_hexo也是可以的.