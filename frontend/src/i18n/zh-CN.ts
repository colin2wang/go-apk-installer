export const zhCN: Record<string, string> = {
  // 标题栏
  'title.minimize': '最小化',
  'title.maximize': '最大化 / 还原',
  'title.close': '关闭',

  // 顶栏
  'sdk.label': 'Android SDK',
  'sdk.notSet': '未设置',
  'sdk.pick': '选择目录',
  'sdk.manual': '手动填写',
  'sdk.placeholder': '例如 C:\\Users\\<你>\\AppData\\Local\\Android\\Sdk',
  'sdk.apply': '应用',
  'sdk.ready': '✓ adb 就绪（{version}）',
  'sdk.notReady': '✕ adb 未就绪',
  'sdk.toggleTheme': '切换主题',

  // 侧栏
  'sidebar.title': '设备配置',
  'sidebar.create': '新建手机配置',
  'sidebar.delete': '删除该手机配置',
  'sidebar.serialAuto': 'serial 留空 · 自动匹配',
  'sidebar.empty': '暂无手机配置',
  'sidebar.emptyHint': '点击右上角 ＋ 创建',
  'sidebar.connUnknown': '连接状态未知，可点「检查连接」',

  // 设备面板
  'devices.title': '当前 adb 设备',
  'devices.refreshing': '刷新中…',
  'devices.refresh': '刷新设备',
  'state.device': '在线·已授权',
  'state.unauthorized': '未授权调试',
  'state.offline': '离线',
  'state.recovery': 'Recovery',
  'state.sideload': 'Sideload',
  'devices.none': '未检测到设备，插上手机后点「刷新设备」',
  'devices.adbNotReady': 'adb 未就绪，先在顶部选择 Android SDK 目录',

  // 表单
  'form.basic': '基本信息',
  'form.package': '安装包',
  'form.options': '安装选项',
  'form.name': '手机名称',
  'form.nameHint': '（= yml 文件名）',
  'form.namePlaceholder': '如：办公测试机',
  'form.serial': '设备 Serial',
  'form.serialHint': '（adb devices 第一列，留空=自动匹配唯一在线设备）',
  'form.serialPlaceholder': '如 emulator-5554 / 1A2B3C4D',
  'form.apkPath': 'APK 路径',
  'form.apkPathHint': '（绝对路径，.apk）',
  'form.apkPlaceholder': 'D:\\apk\\demo-release.apk',
  'form.browse': '浏览…',
  'form.reinstall': '覆盖安装',
  'form.reinstallHint': '-r 保留数据升级',
  'form.downgrade': '允许降级',
  'form.grantPerms': '授予运行时权限',
  'form.autoLaunch': '安装后自动启动',
  'form.launchPackagePlaceholder': '包名，如 com.example.app',
  'form.notesPlaceholder': '备注（可选）',
  'form.hint': '安装前会自动执行两道检查：① adb 是否可用 → ② adb 能否连到该手机；任一不通过都不会执行安装。',
  'form.empty': '左侧选择或新建一台手机开始',
  'check.ok': '✓ 连接正常：',
  'check.fail': '✕ 未通过：',

  // 操作
  'action.save': '保存配置',
  'action.check': '检查连接',
  'action.checking': '检查中…',
  'action.install': '检查并安装',
  'action.installing': '检查并安装中…',

  // 日志
  'log.title': '运行日志',
  'log.openConfigs': '打开 configs 目录',
  'log.clear': '清空日志',
  'log.empty': '暂无日志',

  // 弹窗
  'common.cancel': '取消',
  'common.confirm': '确定',
  'modal.createTitle': '新建手机配置',
  'modal.createBody': '将在 configs/ 目录下创建「{name}.yml」，一台手机一个配置文件。',
  'modal.namePlaceholder': '手机名称，如：办公测试机',
  'modal.deleteTitle': '删除手机配置',
  'modal.deleteBody': '确定删除「{name}」吗？',
  'modal.deleteBodyHint': '将删除配置文件 configs/{name}.yml，此操作不可恢复。',

  // Toast / 日志文案
  'toast.saved': '已保存 configs/{name}.yml',
  'toast.created': '已创建手机配置 configs/{name}.yml',
  'toast.deleted': '已删除手机配置 configs/{name}.yml',
  'toast.needSdkDir': '请先填写 SDK 目录路径',
  'toast.needPhoneName': '请输入手机名称',
  'log.devicesRefreshed': '设备列表已刷新：共 {n} 台',
  'log.refreshFailed': '刷新设备失败：{msg}',
  'log.initDone': '初始化完成：手机配置 {n} 台；adb 就绪（{version}）',
  'log.initDoneNotReady': '初始化完成：手机配置 {n} 台；adb 未就绪，请在顶部选择 Android SDK 目录',
  'log.initFailed': '初始化失败：{msg}',
}
