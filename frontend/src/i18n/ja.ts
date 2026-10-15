export const ja: Record<string, string> = {
  // タイトルバー
  'title.minimize': '最小化',
  'title.maximize': '最大化 / 元に戻す',
  'title.close': '閉じる',

  // トップバー
  'sdk.label': 'Android SDK',
  'sdk.notSet': '未設定',
  'sdk.pick': 'フォルダを選択',
  'sdk.manual': 'パスを入力',
  'sdk.placeholder': '例: C:\\Users\\<ユーザー>\\AppData\\Local\\Android\\Sdk',
  'sdk.apply': '適用',
  'sdk.ready': '✓ adb 準備完了（{version}）',
  'sdk.notReady': '✕ adb 利用不可',
  'sdk.toggleTheme': 'テーマ切替',

  // サイドバー
  'sidebar.title': 'デバイス設定',
  'sidebar.create': 'スマホ設定を新規作成',
  'sidebar.delete': 'この設定を削除',
  'sidebar.serialAuto': 'serial 空 · 自動マッチ',
  'sidebar.empty': 'スマホ設定がありません',
  'sidebar.emptyHint': '右上の ＋ から作成',
  'sidebar.connUnknown': '接続状態不明 —「接続チェック」を実行',

  // デバイスパネル
  'devices.title': 'adb デバイス',
  'devices.refreshing': '更新中…',
  'devices.refresh': '更新',
  'state.device': 'オンライン·許可済み',
  'state.unauthorized': '未許可',
  'state.offline': 'オフライン',
  'state.recovery': 'Recovery',
  'state.sideload': 'Sideload',
  'devices.none': 'デバイスが見つかりません。スマホを接続して「更新」をクリック',
  'devices.adbNotReady': 'adb 利用不可 — まず上部で Android SDK ディレクトリを選択',

  // フォーム
  'form.basic': '基本情報',
  'form.package': 'APK パッケージ',
  'form.options': 'インストールオプション',
  'form.name': 'スマホ名',
  'form.nameHint': '（= yml ファイル名）',
  'form.namePlaceholder': '例：社内テスト機',
  'form.serial': 'デバイス Serial',
  'form.serialHint': '（adb devices の1列目。空欄=唯一のオンライン機器に自動マッチ）',
  'form.serialPlaceholder': '例：emulator-5554 / 1A2B3C4D',
  'form.apkPath': 'APK パス',
  'form.apkPathHint': '（絶対パス、.apk）',
  'form.apkPlaceholder': 'D:\\apk\\demo-release.apk',
  'form.browse': '参照…',
  'form.reinstall': '上書きインストール',
  'form.reinstallHint': '-r データを保持して更新',
  'form.downgrade': 'ダウングレード許可',
  'form.grantPerms': '実行時権限を付与',
  'form.autoLaunch': 'インストール後に自動起動',
  'form.launchPackagePlaceholder': 'パッケージ名、例：com.example.app',
  'form.notesPlaceholder': 'メモ（任意）',
  'form.hint': 'インストール前に2つのチェックを実行：① adb の利用可否 → ② デバイスへの接続。いずれか失敗で中止します。',
  'form.empty': '左でスマホを選択または新規作成',
  'check.ok': '✓ 接続 OK：',
  'check.fail': '✕ 失敗：',

  // 操作
  'action.save': '保存',
  'action.check': '接続チェック',
  'action.checking': 'チェック中…',
  'action.install': 'チェックしてインストール',
  'action.installing': 'インストール中…',

  // ログ
  'log.title': '実行ログ',
  'log.openConfigs': 'configs フォルダを開く',
  'log.clear': 'ログをクリア',
  'log.empty': 'ログはまだありません',

  // ダイアログ
  'common.cancel': 'キャンセル',
  'common.confirm': 'OK',
  'modal.createTitle': '新規スマホ設定',
  'modal.createBody': 'configs/ に「{name}.yml」を作成します。1台=1ファイル。',
  'modal.namePlaceholder': 'スマホ名、例：社内テスト機',
  'modal.deleteTitle': 'スマホ設定を削除',
  'modal.deleteBody': '「{name}」を削除しますか？',
  'modal.deleteBodyHint': 'configs/{name}.yml を削除します。元に戻せません。',

  // トースト / ログ文言
  'toast.saved': 'configs/{name}.yml を保存しました',
  'toast.created': 'configs/{name}.yml を作成しました',
  'toast.deleted': 'configs/{name}.yml を削除しました',
  'toast.needSdkDir': '先に SDK ディレクトリを入力してください',
  'toast.needPhoneName': 'スマホ名を入力してください',
  'log.devicesRefreshed': 'デバイス一覧を更新：{n} 台',
  'log.refreshFailed': 'デバイス更新失敗：{msg}',
  'log.initDone': '初期化完了：設定 {n} 件；adb 準備完了（{version}）',
  'log.initDoneNotReady': '初期化完了：設定 {n} 件；adb 未接続 — 上部で SDK ディレクトリを選択',
  'log.initFailed': '初期化失敗：{msg}',
}
